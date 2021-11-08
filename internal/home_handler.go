package internal

import (
	"net/http"
)

var home = loadNestedTemplates("home.htm")

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	render(w, home, nil)
}
