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
	var responsive, head, leg1, leg2, leg3, leg4, leg5, leg6, wing1, wing2 bool
	if r.FormValue("responsive") == "yes" {
		responsive = true
	} else if r.FormValue("responsive") == "no" {
		responsive = false
	} else {
		return Experiment{}, errors.New("Incorrect responsive value")
	}
	if r.FormValue("head") == "on" {
		head = false
	} else {
		head = true
	}
	if r.FormValue("leg1") == "on" {
		leg1 = false
	} else {
		leg1 = true
	}
	if r.FormValue("leg2") == "on" {
		leg2 = false
	} else {
		leg2 = true
	}
	if r.FormValue("leg3") == "on" {
		leg3 = false
	} else {
		leg3 = true
	}
	if r.FormValue("leg4") == "on" {
		leg4 = false
	} else {
		leg4 = true
	}
	if r.FormValue("leg5") == "on" {
		leg5 = false
	} else {
		leg5 = true
	}
	if r.FormValue("leg6") == "on" {
		leg6 = false
	} else {
		leg6 = true
	}
	if r.FormValue("wing1") == "on" {
		wing1 = false
	} else {
		wing1 = true
	}
	if r.FormValue("wing2") == "on" {
		wing2 = false
	} else {
		wing2 = true
	}

	return Experiment{
		Responsive: responsive,
		Head:       head,
		Leg1:       leg1,
		Leg2:       leg2,
		Leg3:       leg3,
		Leg4:       leg4,
		Leg5:       leg5,
		Leg6:       leg6,
		Wing1:      wing1,
		Wing2:      wing2,
	}, nil
}
