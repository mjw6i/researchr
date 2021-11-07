package internal

import (
	"net/http"
)

func BaseHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	homeHandler(w, r)
}
