package internal

import (
	"testing"

	"github.com/mjw6i/researchr/pkg"
)

func TestAssetsRoute(t *testing.T) {
	recorder := pkg.MakeRequest(t, AssetsHandler, "/assets")

	pkg.AssertStatus(t, recorder, 200)
	pkg.AssertBodyStartsWith(t, recorder, "<!DOCTYPE html>")
}
