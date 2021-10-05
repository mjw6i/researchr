package main

import (
	"html/template"
	"log"
	"net/http"
)

var layout = template.Must(template.ParseFiles("index.htm", "mosquito.htm"))

func main() {
	http.HandleFunc("/", homeHandler)

	log.Fatal(http.ListenAndServe(":9000", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	layout.ExecuteTemplate(w, "index.htm", "")
}
