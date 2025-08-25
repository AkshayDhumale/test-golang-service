package main

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"strings"
	"github.com/stretchr/testify/assert"
)

func TestSearchUsers_EmptyParams(t *testing.T) {
	s := setupTestService()
	req := httptest.NewRequest("GET", "/api/v1/users/search", nil)
	w := httptest.NewRecorder()
	s.searchUsers(w, req)
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestSearchUsers_ByName(t *testing.T) {
	s := setupTestService()
	req := httptest.NewRequest("GET", "/api/v1/users/search?name=John", nil)
	w := httptest.NewRecorder()
	s.searchUsers(w, req)
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func setupTestService() *Service {
	// You may want to use a mock DB and Redis for real unit tests
	return NewService()
}
