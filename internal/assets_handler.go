package internal

import (
	"net/http"
)

var assets = loadNestedTemplates("assets.htm")

func AssetsHandler(w http.ResponseWriter, r *http.Request) {
	render(w, assets, nil)
}
