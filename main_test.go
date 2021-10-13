package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRandomRoute(t *testing.T) {
	recorder := makeRequest(t, baseHandler, "/random")

	assertStatus(t, recorder, 404)
	assertBody(t, recorder, "404 page not found\n")
}

func TestBaseRoute(t *testing.T) {
	recorder := makeRequest(t, baseHandler, "/")

	assertStatus(t, recorder, 200)
	assertBodyStartsWith(t, recorder, "<!DOCTYPE html>")
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

func TestReceiveRoute(t *testing.T) {
	ds := SuccessStore{}
	env := &Env{store: ds}

	recorder := makeRequest(t, env.receiveHandler, "/receive")

	assertStatus(t, recorder, 303)
	assertHeader(t, recorder, "Location", "/results")
}

func TestReceiveAllChecked(t *testing.T) {
	ds := SuccessStore{}
	env := &Env{store: ds}

	body := "leg1=on&leg2=on&leg3=on&leg4=on&leg5=on&leg6=on&wing1=on&wing2=on&head=on&responsive=yes"
	recorder := makeFormRequest(t, env.receiveHandler, "/receive", body)

	assertStatus(t, recorder, 303)
	assertHeader(t, recorder, "Location", "/results")
}

func TestReceiveAllUnchecked(t *testing.T) {
	ds := SuccessStore{}
	env := &Env{store: ds}

	body := ""
	recorder := makeFormRequest(t, env.receiveHandler, "/receive", body)

	assertStatus(t, recorder, 303)
	assertHeader(t, recorder, "Location", "/results")
}

func TestResultsRoute(t *testing.T) {
	ds := SuccessStore{}
	env := &Env{store: ds}

	recorder := makeRequest(t, env.resultsHandler, "/results")

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

func assertHeader(t *testing.T, r *httptest.ResponseRecorder, key string, value string) {
	if r.Header().Get(key) != value {
		t.Errorf("Expected header: '%v' got '%v'", value, r.Header().Get(key))
	}
}
