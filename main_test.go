package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRandomRoute(t *testing.T) {
	recorder := makeRequest(t, homeHandler, "/random")

	expectedBody := "404 page not found\n"

	assertStatus(t, recorder, 404)

	if recorder.Body.String() != expectedBody {
		t.Errorf("Expected response: '%v' got '%v'", expectedBody, recorder.Body.String())
	}
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

func assertStatus(t *testing.T, r *httptest.ResponseRecorder, expected int) {
	if r.Code != expected {
		t.Errorf("Expected status: '%v' got '%v'", expected, r.Code)
	}
}
