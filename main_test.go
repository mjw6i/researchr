package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRandomRoute(t *testing.T) {
	recorder := makeRequest(t, homeHandler, "/random")

	assertStatus(t, recorder, 404)
	assertBody(t, recorder, "404 page not found\n")
}

func makeRequest(t *testing.T, h func(http.ResponseWriter, *http.Request), url string) *httptest.ResponseRecorder {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(h)

	handler.ServeHTTP(recorder, req)

	return recorder
}

func assertBody(t *testing.T, r *httptest.ResponseRecorder, expected string) {
	if r.Body.String() != expected {
		t.Errorf("Expected response: '%v' got '%v'", expected, r.Body.String())
	}
}

func assertStatus(t *testing.T, r *httptest.ResponseRecorder, expected int) {
	if r.Code != expected {
		t.Errorf("Expected status: '%v' got '%v'", expected, r.Code)
	}
}
