package main

import (
	"encoding/json"
	"net/http"
	"fmt"
	"strings"
	"time"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
)

// UserSearchParams holds possible search parameters
// Extend as needed
func (s *Service) searchUsers(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	name := r.URL.Query().Get("name")
	email := r.URL.Query().Get("email")

	// Build cache key from params
	cacheKey := fmt.Sprintf("usersearch:name=%s:email=%s", name, email)
	cached, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil && cached != "" {
		w.Header().Set("X-Cache", "HIT")
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cached))
		return
	}

	// Build SQL query
	var queryBuilder strings.Builder
	var args []interface{}
	queryBuilder.WriteString("SELECT id, name, email FROM users WHERE 1=1")
	if name != "" {
		queryBuilder.WriteString(" AND name ILIKE $1")
		args = append(args, "%"+name+"%")
	}
	if email != "" {
		queryBuilder.WriteString(" AND email ILIKE $2")
		args = append(args, "%"+email+"%")
	}
	query := queryBuilder.String()

	rows, err := s.db.Query(query, args...)
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

	resultJSON, _ := json.Marshal(users)
	// Cache result for 5 minutes
	s.redis.Set(ctx, cacheKey, resultJSON, 5*time.Minute)

	w.Header().Set("X-Cache", "MISS")
	w.Header().Set("Content-Type", "application/json")
	w.Write(resultJSON)
}
