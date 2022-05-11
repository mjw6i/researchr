package internal

import (
	"testing"

	"github.com/mjw6i/researchr/pkg"
)

func TestEmpty(t *testing.T) {
	e := Experiment{}

	pkg.AssertSmall(t, 0, e.Extremities())
}

func TestResponsiveAndHead(t *testing.T) {
	e := Experiment{Responsive: true, Head: true}

	pkg.AssertSmall(t, 0, e.Extremities())
}

func TestAll(t *testing.T) {
	e := Experiment{false, false, true, true, true, true, true, true, true, true}

	pkg.AssertSmall(t, 8, e.Extremities())
}
