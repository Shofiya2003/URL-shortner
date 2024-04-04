package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/assert"
)

func TestPingServer(t *testing.T) {
	mockResponse := `{"message":"url shortner apis!"}`
	router := SetUpRouter()
	router.GET("/ping", PingServer)
	req, _ := http.NewRequest("GET", "/ping", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	response, _ := io.ReadAll(w.Body)
	assert.Equal(t, string(mockResponse), string(response))
	assert.Equal(t, http.StatusOK, w.Code)
}
