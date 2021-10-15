package main

import (
	"net/http"
)

func submitHandler(w http.ResponseWriter, r *http.Request) {
	render(w, submit, nil)
}
