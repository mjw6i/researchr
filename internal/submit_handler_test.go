package internal

import (
	"testing"

	"github.com/mjw6i/researchr/pkg"
)

func TestSubmitRoute(t *testing.T) {
	recorder := pkg.MakeRequest(t, SubmitHandler, "/submit")

	pkg.AssertStatus(t, recorder, 200)
	pkg.AssertBodyStartsWith(t, recorder, "<!DOCTYPE html>")
}
