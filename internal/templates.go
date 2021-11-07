package internal

import (
	"html/template"
	"log"
	"net/http"
)

func Render(w http.ResponseWriter, t *template.Template, data interface{}) {
	err := t.ExecuteTemplate(w, "layout.htm", data)
	if err != nil {
		log.Fatal(err)
	}
}

func LoadNestedTemplates(filenames ...string) *template.Template {
	t := append([]string{"template/layout.htm"}, filenames...)
	return template.Must(template.ParseFiles(t...))
}
