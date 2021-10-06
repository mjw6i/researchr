package main

import (
	"html/template"
	"log"
	"net/http"
)

var layout = template.Must(template.ParseFiles("template/layout.htm", "template/home.htm", "template/submit.htm", "template/mosquito.htm"))

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/submit", submitHandler)
	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./static"))))
	log.Fatal(http.ListenAndServe(":9000", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	err := layout.ExecuteTemplate(w, "home.htm", "")
	if err != nil {
		log.Fatal(err)
	}
	err = layout.ExecuteTemplate(w, "layout.htm", "")
	if err != nil {
		log.Fatal(err)
	}
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
	err := layout.ExecuteTemplate(w, "submit.htm", "")
	if err != nil {
		log.Fatal(err)
	}
	err = layout.ExecuteTemplate(w, "layout.htm", "")
	if err != nil {
		log.Fatal(err)
	}
}
