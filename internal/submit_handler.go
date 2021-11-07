package internal

import (
	"net/http"
)

var submit = LoadNestedTemplates("template/submit.htm", "template/mosquito.htm")

func SubmitHandler(w http.ResponseWriter, r *http.Request) {
	Render(w, submit, nil)
}
