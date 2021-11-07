package internal

import (
	"net/http"
	"strings"
	"testing"
)

func TestReceiveRedirect(t *testing.T) {
	ds := SuccessStore{}
	env := &Env{store: ds}

	recorder := makeFormRequest(t, env.receiveHandler, "/receive", "responsive=yes")

	assertStatus(t, recorder, 303)
	assertHeader(t, recorder, "Location", "/results")
}

func TestReceiveError(t *testing.T) {
	ds := SuccessStore{}
	env := &Env{store: ds}

	recorder := makeRequest(t, env.receiveHandler, "/receive")

	assertStatus(t, recorder, 400)
	assertBody(t, recorder, "")
}

func TestParseExperimentFormAllChecked(t *testing.T) {
	body := "leg1=on&leg2=on&leg3=on&leg4=on&leg5=on&leg6=on&wing1=on&wing2=on&head=on&responsive=yes"

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
	env := &Env{store: ds}

	recorder := makeFormRequest(t, env.receiveHandler, "/receive", "responsive=yes")

	assertStatus(t, recorder, 500)
	assertBody(t, recorder, "")
}
