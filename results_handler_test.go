package main

import (
	"testing"
)

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
