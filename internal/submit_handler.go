package internal

import (
	"net/http"
)

var submit = LoadNestedTemplates("submit.htm", "mosquito.htm")

func SubmitHandler(w http.ResponseWriter, r *http.Request) {
	Render(w, submit, nil)
}
