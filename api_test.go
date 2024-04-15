package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostUrlHandler(t *testing.T) {
	router := SetUpRouter()

	err := InitRedis()
	if err != nil {
		fmt.Println(err)
	}
	router.POST("/", PostURLHandler)

	url := Url{
		LongUrl: "Dummy URL",
	}

	jsonValue, _ := json.Marshal(url)
	req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(jsonValue))

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestUrlDeleteHandler(t *testing.T) {
	router := SetUpRouter()

	err := InitRedis()
	if err != nil {
		fmt.Println(err)
	}

	url := Url{
		LongUrl: "Dummy URL To Delete",
	}

	router.DELETE("/:key", DeleteUrlHandler)
	router.POST("/", PostURLHandler)

	jsonValue, _ := json.Marshal(url)
	req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(jsonValue))

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Result().StatusCode == http.StatusCreated {
		response := w.Body.Bytes()
		var newUrl Url

		json.Unmarshal(response, &newUrl)
		fmt.Print(newUrl.Key)
		req, _ = http.NewRequest("DELETE", "/"+newUrl.Key, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	}

}

func TestGetUrlsHandler(t *testing.T) {
	router := SetUpRouter()

	err := InitRedis()
	if err != nil {
		fmt.Println(err)
	}

	router.POST("/all", ListAllUrlHandler)

	req, _ := http.NewRequest("POST", "/all", nil)

	w := httptest.NewRecorder()

	var urls []Url

	router.ServeHTTP(w, req)
	json.Unmarshal(w.Body.Bytes(), &urls)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, urls)

}
