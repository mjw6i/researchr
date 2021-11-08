package internal

import (
	"testing"

	"github.com/mjw6i/researchr/pkg"
)

func TestResultsRoute(t *testing.T) {
	ds := SuccessStore{}
	env := &Env{Store: ds}

	recorder := pkg.MakeRequest(t, env.ResultsHandler, "/results")

	pkg.AssertStatus(t, recorder, 200)
	pkg.AssertBodyStartsWith(t, recorder, "<!DOCTYPE html>")
}

func TestResultsStore(t *testing.T) {
	ds := FailureStore{}
	env := &Env{Store: ds}

	recorder := pkg.MakeRequest(t, env.ResultsHandler, "/results")

	pkg.AssertStatus(t, recorder, 500)
	pkg.AssertBody(t, recorder, "")
}
