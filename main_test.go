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

func TestReceiveRedirect(t *testing.T) {
	ds := SuccessStore{}
	env := &Env{store: ds}

	recorder := makeFormRequest(t, env.receiveHandler, "/receive", "responsive=yes")

	assertStatus(t, recorder, 303)
	assertHeader(t, recorder, "Location", "/results")
}

func TestReceiveError(t *testing.T) {
	ds := SuccessStore{}
	env := &Env{store: ds}

	recorder := makeRequest(t, env.receiveHandler, "/receive")

	assertStatus(t, recorder, 400)
	assertBody(t, recorder, "")
}

func TestParseExperimentFormAllChecked(t *testing.T) {
	body := "leg1=on&leg2=on&leg3=on&leg4=on&leg5=on&leg6=on&wing1=on&wing2=on&head=on&responsive=yes"

	req, err := http.NewRequest("POST", "", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	expected := Experiment{
		Responsive: true,
		Head:       false,
		Leg1:       false,
		Leg2:       false,
		Leg3:       false,
		Leg4:       false,
		Leg5:       false,
		Leg6:       false,
		Wing1:      false,
		Wing2:      false,
	}

	result, err := parseExperimentFormData(req)

	if err != nil {
		t.Fatal(err)
	}

	if expected != result {
		t.Fatal("Structs are not equal")
	}
}

func TestParseExperimentFormAllUnchecked(t *testing.T) {
	body := "responsive=no"

	req, err := http.NewRequest("POST", "", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	expected := Experiment{
		Responsive: false,
		Head:       true,
		Leg1:       true,
		Leg2:       true,
		Leg3:       true,
		Leg4:       true,
		Leg5:       true,
		Leg6:       true,
		Wing1:      true,
		Wing2:      true,
	}

	result, err := parseExperimentFormData(req)

	if err != nil {
		t.Fatal(err)
	}

	if expected != result {
		t.Fatal("Structs are not equal")
	}
}

func TestReceiveStore(t *testing.T) {
	ds := FailureStore{}
	env := &Env{store: ds}

	recorder := makeFormRequest(t, env.receiveHandler, "/receive", "responsive=yes")

	assertStatus(t, recorder, 500)
	assertBody(t, recorder, "")
}

func TestResultsRoute(t *testing.T) {
	ds := SuccessStore{}
	env := &Env{store: ds}

	recorder := makeRequest(t, env.resultsHandler, "/results")

	assertStatus(t, recorder, 200)
	assertBodyStartsWith(t, recorder, "<!DOCTYPE html>")
}

func TestResultsStore(t *testing.T) {
	ds := FailureStore{}
	env := &Env{store: ds}

	recorder := makeRequest(t, env.resultsHandler, "/results")

	assertStatus(t, recorder, 500)
	assertBody(t, recorder, "")
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
