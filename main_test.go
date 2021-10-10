package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRandomRoute(t *testing.T) {
	recorder := makeRequest(t, homeHandler, "/random")

	assertStatus(t, recorder, 404)
	assertBody(t, recorder, "404 page not found\n")
}

func TestHomeRoute(t *testing.T) {
	recorder := makeRequest(t, homeHandler, "/")

	assertStatus(t, recorder, 200)
	assertBodyStartsWith(t, recorder, "<!DOCTYPE html>")
}

func TestSubmitRoute(t *testing.T) {
	recorder := makeRequest(t, submitHandler, "/submit")

	assertStatus(t, recorder, 200)
	assertBodyStartsWith(t, recorder, "<!DOCTYPE html>")
}

func TestResultsRoute(t *testing.T) {
	recorder := makeRequest(t, resultsHandler, "/results")

	assertStatus(t, recorder, 200)
	assertBodyStartsWith(t, recorder, "<!DOCTYPE html>")
}

func TestAssetsRoute(t *testing.T) {
	recorder := makeRequest(t, assetsHandler, "/assets")

	assertStatus(t, recorder, 200)
	assertBodyStartsWith(t, recorder, "<!DOCTYPE html>")
}

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

func assertBody(t *testing.T, r *httptest.ResponseRecorder, expected string) {
	if r.Body.String() != expected {
		t.Errorf("Expected response: '%v' got '%v'", expected, r.Body.String())
	}
}

func assertBodyStartsWith(t *testing.T, r *httptest.ResponseRecorder, expected string) {
	if !strings.HasPrefix(r.Body.String(), expected) {
		comp := r.Body.String() + strings.Repeat("#", len(expected))
		t.Errorf("Expected prefix: '%v' got '%v'", expected, comp[0:len(expected)])
	}
}

func assertStatus(t *testing.T, r *httptest.ResponseRecorder, expected int) {
	if r.Code != expected {
		t.Errorf("Expected status: '%v' got '%v'", expected, r.Code)
	}
}
