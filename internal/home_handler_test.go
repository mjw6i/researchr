package internal

import (
	"testing"
)

func TestHomeRoute(t *testing.T) {
	recorder := makeRequest(t, homeHandler, "/")

	assertStatus(t, recorder, 200)
	assertBodyStartsWith(t, recorder, "<!DOCTYPE html>")
}
