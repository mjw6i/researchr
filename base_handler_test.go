package main

import (
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
