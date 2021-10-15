package main

import (
	"testing"
)

func TestSubmitRoute(t *testing.T) {
	recorder := makeRequest(t, submitHandler, "/submit")

	assertStatus(t, recorder, 200)
	assertBodyStartsWith(t, recorder, "<!DOCTYPE html>")
}
