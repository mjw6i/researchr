package internal

import (
	"net/http"
)

var home = LoadNestedTemplates("template/home.htm")

func homeHandler(w http.ResponseWriter, r *http.Request) {
	Render(w, home, nil)
}
