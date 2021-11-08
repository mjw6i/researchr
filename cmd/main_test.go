package main

import (
	"testing"

	"github.com/mjw6i/researchr/pkg"
)

func TestStyleFile(t *testing.T) {
	recorder := pkg.MakeRequest(t, static.ServeHTTP, "/static/style.css")

	pkg.AssertStatus(t, recorder, 200)
	pkg.AssertBodyStartsWith(t, recorder, ":root {")
}
