package internal

import (
	"net/http"
	"strings"
	"testing"

	"github.com/mjw6i/researchr/pkg"
)

func TestReceiveRedirect(t *testing.T) {
	ds := SuccessStore{}
	env := &Env{Store: ds}

	recorder := pkg.MakeFormRequest(t, env.ReceiveHandler, "/receive", "responsive=yes")

	pkg.AssertStatus(t, recorder, 303)
	pkg.AssertHeader(t, recorder, "Location", "/results")
}

func TestReceiveError(t *testing.T) {
	ds := SuccessStore{}
	env := &Env{Store: ds}

	recorder := pkg.MakeRequest(t, env.ReceiveHandler, "/receive")

	pkg.AssertStatus(t, recorder, 400)
	pkg.AssertBody(t, recorder, "")
}

func TestParseExperimentFormAllChecked(t *testing.T) {
	body := "leg1=removed&leg2=removed&leg3=removed&leg4=removed&leg5=removed&leg6=removed&wing1=removed&wing2=removed&head=removed&responsive=yes"

	req, err := http.NewRequest("POST", "", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	expected := Experiment{
		Responsive: true,
		Head:       false,
		Leg1:       false,
		Leg2:       false,
		Leg3:       false,
		Leg4:       false,
		Leg5:       false,
		Leg6:       false,
		Wing1:      false,
		Wing2:      false,
	}

	result, err := parseExperimentFormData(req)

	if err != nil {
		t.Fatal(err)
	}

	if expected != result {
		t.Fatal("Structs are not equal")
	}
}

func TestParseExperimentFormAllUnchecked(t *testing.T) {
	body := "responsive=no"

	req, err := http.NewRequest("POST", "", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	expected := Experiment{
		Responsive: false,
		Head:       true,
		Leg1:       true,
		Leg2:       true,
		Leg3:       true,
		Leg4:       true,
		Leg5:       true,
		Leg6:       true,
		Wing1:      true,
		Wing2:      true,
	}

	result, err := parseExperimentFormData(req)

	if err != nil {
		t.Fatal(err)
	}

	if expected != result {
		t.Fatal("Structs are not equal")
	}
}

func TestReceiveStore(t *testing.T) {
	ds := FailureStore{}
	env := &Env{Store: ds}

	recorder := pkg.MakeFormRequest(t, env.ReceiveHandler, "/receive", "responsive=yes")

	pkg.AssertStatus(t, recorder, 500)
	pkg.AssertBody(t, recorder, "")
}
