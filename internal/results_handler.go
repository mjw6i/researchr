package internal

import (
	"log"
	"net/http"
)

var results = LoadNestedTemplates("template/results.htm")

func (env *Env) ResultsHandler(w http.ResponseWriter, r *http.Request) {
	res, err := env.Store.getResult()

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	Render(w, results, res)
}
