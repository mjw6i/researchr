package internal

import (
	"net/http"
)

var home = LoadNestedTemplates("home.htm")

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	Render(w, home, nil)
}
