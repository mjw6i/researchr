package internal

import (
	"testing"

	"github.com/mjw6i/researchr/pkg"
)

func TestHomeRoute(t *testing.T) {
	recorder := pkg.MakeRequest(t, HomeHandler, "/")

	pkg.AssertStatus(t, recorder, 200)
	pkg.AssertBodyStartsWith(t, recorder, "<!DOCTYPE html>")
}
