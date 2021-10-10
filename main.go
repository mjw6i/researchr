package main

import (
	"html/template"
	"log"
	"net/http"
)

var home = template.Must(template.ParseFiles("template/layout.htm", "template/home.htm"))
var submit = template.Must(template.ParseFiles("template/layout.htm", "template/submit.htm", "template/mosquito.htm"))
var results = template.Must(template.ParseFiles("template/layout.htm", "template/results.htm"))
var assets = template.Must(template.ParseFiles("template/layout.htm", "template/assets.htm"))

var static = http.StripPrefix("/static", http.FileServer(http.Dir("./static")))

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/submit", submitHandler)
	http.HandleFunc("/results", resultsHandler)
	http.HandleFunc("/assets", assetsHandler)
	http.Handle("/static/", static)
	log.Fatal(http.ListenAndServe(":9000", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
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
