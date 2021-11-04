package main

import (
	"errors"
	"log"
	"net/http"
)

func (env *Env) receiveHandler(w http.ResponseWriter, r *http.Request) {
	experiment, err := parseExperimentFormData(r)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = env.store.storeExperiment(experiment)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/results", http.StatusSeeOther)
}

func parseExperimentFormData(r *http.Request) (Experiment, error) {
	var responsive bool
	if r.FormValue("responsive") == "yes" {
		responsive = true
	} else if r.FormValue("responsive") == "no" {
		responsive = false
	} else {
		return Experiment{}, errors.New("incorrect responsive value")
	}

	return Experiment{
		Responsive: responsive,
		Head:       formBool(r, "head"),
		Leg1:       formBool(r, "leg1"),
		Leg2:       formBool(r, "leg2"),
		Leg3:       formBool(r, "leg3"),
		Leg4:       formBool(r, "leg4"),
		Leg5:       formBool(r, "leg5"),
		Leg6:       formBool(r, "leg6"),
		Wing1:      formBool(r, "wing1"),
		Wing2:      formBool(r, "wing2"),
	}, nil
}

func formBool(r *http.Request, field string) bool {
	if r.FormValue(field) == "on" {
		return false
	}

	return true
}
