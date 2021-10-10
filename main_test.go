package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRandomRoute(t *testing.T) {
	req, err := http.NewRequest("GET", "/random", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(homeHandler)

	handler.ServeHTTP(recorder, req)

	expectedStatus := 404
	expectedBody := "404 page not found\n"

	if recorder.Code != expectedStatus {
		t.Errorf("Expected status: '%v' got '%v'", expectedStatus, recorder.Code)
	}

	if recorder.Body.String() != expectedBody {
		t.Errorf("Expected response: '%v' got '%v'", expectedBody, recorder.Body.String())
	}
}
