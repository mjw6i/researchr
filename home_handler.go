package main

import (
	"net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	render(w, home, nil)
}
