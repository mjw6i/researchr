package internal

import (
	"testing"

	"github.com/mjw6i/researchr/pkg"
)

func TestRandomRoute(t *testing.T) {
	recorder := pkg.MakeRequest(t, BaseHandler, "/random")

	pkg.AssertStatus(t, recorder, 404)
	pkg.AssertBody(t, recorder, "404 page not found\n")
}

func TestBaseRoute(t *testing.T) {
	recorder := pkg.MakeRequest(t, BaseHandler, "/")

	pkg.AssertStatus(t, recorder, 200)
	pkg.AssertBodyStartsWith(t, recorder, "<!DOCTYPE html>")
}
