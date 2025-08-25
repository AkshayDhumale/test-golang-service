// main.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"database/sql"
	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
	_ "github.com/lib/pq"
	"golang.org/x/net/context"
)

type User struct {
	ID    int    `json:"id" db:"id"`
	Name  string `json:"name" db:"name"`
	Email string `json:"email" db:"email"`
}

type Service struct {
	db    *sql.DB
	redis *redis.Client
}

func NewService() *Service {
	// PostgreSQL connection
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "userdb")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	// Redis connection
	redisHost := getEnv("REDIS_HOST", "localhost")
	redisPort := getEnv("REDIS_PORT", "6379")
	redisPassword := getEnv("REDIS_PASSWORD", "")

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisHost + ":" + redisPort,
		Password: redisPassword,
		DB:       0,
	})

	// Test Redis connection
	ctx := context.Background()
	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}

	// Initialize database tables
	initDB(db)

	log.Println("Connected to PostgreSQL and Redis successfully!")

	return &Service{
		db:    db,
		redis: rdb,
	}
}

func initDB(db *sql.DB) {
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		email VARCHAR(100) UNIQUE NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := db.Exec(createTableQuery)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}
}

func (s *Service) createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	query := "INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id"
	err := s.db.QueryRow(query, user.Name, user.Email).Scan(&user.ID)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Cache user in Redis for 1 hour
	ctx := context.Background()
	userJSON, _ := json.Marshal(user)
	cacheKey := fmt.Sprintf("user:%d", user.ID)
	s.redis.Set(ctx, cacheKey, userJSON, time.Hour)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (s *Service) getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	cacheKey := fmt.Sprintf("user:%d", userID)

	// Try to get from Redis cache first
	cached, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var user User
		if json.Unmarshal([]byte(cached), &user) == nil {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("X-Cache", "HIT")
			json.NewEncoder(w).Encode(user)
			return
		}
	}

	// If not in cache, get from database
	var user User
	query := "SELECT id, name, email FROM users WHERE id = $1"
	err = s.db.QueryRow(query, userID).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	// Cache the result
	userJSON, _ := json.Marshal(user)
	s.redis.Set(ctx, cacheKey, userJSON, time.Hour)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Cache", "MISS")
	json.NewEncoder(w).Encode(user)
}

func (s *Service) listUsers(w http.ResponseWriter, r *http.Request) {
	query := "SELECT id, name, email FROM users ORDER BY id"
	rows, err := s.db.Query(query)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			continue
		}
		users = append(users, user)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (s *Service) updateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	user.ID = userID

	query := "UPDATE users SET name = $1, email = $2 WHERE id = $3"
	result, err := s.db.Exec(query, user.Name, user.Email, userID)
	if err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Update cache
	ctx := context.Background()
	cacheKey := fmt.Sprintf("user:%d", userID)
	userJSON, _ := json.Marshal(user)
	s.redis.Set(ctx, cacheKey, userJSON, time.Hour)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (s *Service) deleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	query := "DELETE FROM users WHERE id = $1"
	result, err := s.db.Exec(query, userID)
	if err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Remove from cache
	ctx := context.Background()
	cacheKey := fmt.Sprintf("user:%d", userID)
	s.redis.Del(ctx, cacheKey)

	w.WriteHeader(http.StatusNoContent)
}

func (s *Service) healthCheck(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	
	// Check PostgreSQL
	if err := s.db.Ping(); err != nil {
		http.Error(w, "Database unhealthy", http.StatusServiceUnavailable)
		return
	}

	// Check Redis
	if _, err := s.redis.Ping(ctx).Result(); err != nil {
		http.Error(w, "Redis unhealthy", http.StatusServiceUnavailable)
		return
	}

	response := map[string]string{
		"status":     "healthy",
		"database":   "connected",
		"redis":      "connected",
		"timestamp":  time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func main() {
	service := NewService()
	defer service.db.Close()
	defer service.redis.Close()

	router := mux.NewRouter()

	// API routes
	api := router.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/users", service.createUser).Methods("POST")
	api.HandleFunc("/users", service.listUsers).Methods("GET")
	api.HandleFunc("/users/{id}", service.getUser).Methods("GET")
	api.HandleFunc("/users/{id}", service.updateUser).Methods("PUT")
	api.HandleFunc("/users/{id}", service.deleteUser).Methods("DELETE")

	// Health check
	router.HandleFunc("/health", service.healthCheck).Methods("GET")

	port := getEnv("PORT", "8080")
	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
