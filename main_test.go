package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestStyleFile(t *testing.T) {
	recorder := makeRequest(t, static.ServeHTTP, "/static/style.css")

	assertStatus(t, recorder, 200)
	assertBodyStartsWith(t, recorder, ":root {")
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

func makeFormRequest(t *testing.T, h func(http.ResponseWriter, *http.Request), url string, body string) *httptest.ResponseRecorder {
	req, err := http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(h)

	handler.ServeHTTP(recorder, req)

	return recorder
}

func assertBody(t *testing.T, r *httptest.ResponseRecorder, expected string) {
	if r.Body.String() != expected {
		t.Fatalf("Expected response: '%v' got '%v'", expected, r.Body.String())
	}
}

func assertBodyStartsWith(t *testing.T, r *httptest.ResponseRecorder, expected string) {
	if !strings.HasPrefix(r.Body.String(), expected) {
		comp := r.Body.String() + strings.Repeat("#", len(expected))
		t.Fatalf("Expected prefix: '%v' got '%v'", expected, comp[0:len(expected)])
	}
}

func assertStatus(t *testing.T, r *httptest.ResponseRecorder, expected int) {
	if r.Code != expected {
		t.Fatalf("Expected status: '%v' got '%v'", expected, r.Code)
	}
}

func assertHeader(t *testing.T, r *httptest.ResponseRecorder, key string, value string) {
	if r.Header().Get(key) != value {
		t.Fatalf("Expected header: '%v' got '%v'", value, r.Header().Get(key))
	}
}

func assertFloat(t *testing.T, expected float64, actual float64) {
	if expected != actual {
		t.Fatalf("Expected value: '%v' got '%v'", expected, actual)
	}
}

func assertError(t *testing.T, expected string, err error) {
	if err == nil {
		t.Fatal("Expected an error")
	}

	if err.Error() != expected {
		t.Fatalf("Expected error: '%v' got '%v'", expected, err.Error())
	}
}
