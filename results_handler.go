package main

import (
	"log"
	"net/http"
)

func (env *Env) resultsHandler(w http.ResponseWriter, r *http.Request) {
	res, err := env.store.getResult()

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	render(w, results, res)
}
