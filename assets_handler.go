package main

import (
	"net/http"
)

func assetsHandler(w http.ResponseWriter, r *http.Request) {
	render(w, assets, nil)
}
