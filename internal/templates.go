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
	t := append([]string{"layout.htm"}, filenames...)
	t = templatePath(t)
	return template.Must(template.ParseFiles(t...))
}

func templatePath(templates []string) []string {
	for i, t := range templates {
		templates[i] = "../web/template/" + t
	}

	return templates
}
