package internal

import (
	"errors"
	"log"
	"net/http"
)

func (env *Env) ReceiveHandler(w http.ResponseWriter, r *http.Request) {
	experiment, err := parseExperimentFormData(r)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = env.Store.storeExperiment(experiment)

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
		Head:       isPresent(r, "head"),
		Leg1:       isPresent(r, "leg1"),
		Leg2:       isPresent(r, "leg2"),
		Leg3:       isPresent(r, "leg3"),
		Leg4:       isPresent(r, "leg4"),
		Leg5:       isPresent(r, "leg5"),
		Leg6:       isPresent(r, "leg6"),
		Wing1:      isPresent(r, "wing1"),
		Wing2:      isPresent(r, "wing2"),
	}, nil
}

func isPresent(r *http.Request, field string) bool {
	return r.FormValue(field) != "removed"
}
