package main

import (
	"testing"
)

func TestAssetsRoute(t *testing.T) {
	recorder := makeRequest(t, assetsHandler, "/assets")

	assertStatus(t, recorder, 200)
	assertBodyStartsWith(t, recorder, "<!DOCTYPE html>")
}
