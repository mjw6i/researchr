package main

import (
	"html/template"
	"log"
	"net/http"
)

var home = loadNestedTemplates("template/home.htm")
var submit = loadNestedTemplates("template/submit.htm", "template/mosquito.htm")
var results = loadNestedTemplates("template/results.htm")
var assets = loadNestedTemplates("template/assets.htm")

var static = http.StripPrefix("/static", http.FileServer(http.Dir("./static")))

func main() {
	http.HandleFunc("/", baseHandler)
	http.HandleFunc("/submit", submitHandler)
	http.HandleFunc("/results", resultsHandler)
	http.HandleFunc("/assets", assetsHandler)
	http.Handle("/static/", static)
	log.Fatal(http.ListenAndServe(":9000", nil))
}

func baseHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	homeHandler(w, r)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	render(w, home)
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
	render(w, submit)
}

func resultsHandler(w http.ResponseWriter, r *http.Request) {
	render(w, results)
}

func assetsHandler(w http.ResponseWriter, r *http.Request) {
	render(w, assets)
}

func render(w http.ResponseWriter, t *template.Template) {
	err := t.ExecuteTemplate(w, "layout.htm", "")
	if err != nil {
		log.Fatal(err)
	}
}

func loadNestedTemplates(filenames ...string) *template.Template {
	t := append([]string{"template/layout.htm"}, filenames...)
	return template.Must(template.ParseFiles(t...))
}
