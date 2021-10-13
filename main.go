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
	ds := MockDataStore{}
	env := &Env{store: ds}

	http.HandleFunc("/", baseHandler)
	http.HandleFunc("/submit", submitHandler)
	http.HandleFunc("/receive", receiveHandler)
	http.HandleFunc("/results", env.resultsHandler)
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
	render(w, home, nil)
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
	render(w, submit, nil)
}

func receiveHandler(w http.ResponseWriter, r *http.Request) {
	log.Print(r.FormValue("responsiveness"))
}

func (env *Env) resultsHandler(w http.ResponseWriter, r *http.Request) {
	res, _ := env.store.getResult()

	render(w, results, res)
}

func assetsHandler(w http.ResponseWriter, r *http.Request) {
	render(w, assets, nil)
}

func render(w http.ResponseWriter, t *template.Template, data interface{}) {
	err := t.ExecuteTemplate(w, "layout.htm", data)
	if err != nil {
		log.Fatal(err)
	}
}

func loadNestedTemplates(filenames ...string) *template.Template {
	t := append([]string{"template/layout.htm"}, filenames...)
	return template.Must(template.ParseFiles(t...))
}
