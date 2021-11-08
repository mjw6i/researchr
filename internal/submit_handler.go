package internal

import (
	"net/http"
)

var submit = loadNestedTemplates("submit.htm", "mosquito.htm")

func SubmitHandler(w http.ResponseWriter, r *http.Request) {
	render(w, submit, nil)
}
