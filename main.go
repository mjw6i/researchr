package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"os"
)

var home = loadNestedTemplates("template/home.htm")
var submit = loadNestedTemplates("template/submit.htm", "template/mosquito.htm")
var results = loadNestedTemplates("template/results.htm")
var assets = loadNestedTemplates("template/assets.htm")

var static = http.StripPrefix("/static", http.FileServer(http.Dir("./static")))

func main() {
	db, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	ds := DatabaseStore{db: db}
	env := &Env{store: &ds}

	http.HandleFunc("/", baseHandler)
	http.HandleFunc("/submit", submitHandler)
	http.HandleFunc("/receive", env.receiveHandler)
	http.HandleFunc("/results", env.resultsHandler)
	http.HandleFunc("/assets", assetsHandler)
	http.Handle("/static/", static)
	log.Fatal(http.ListenAndServe(":9000", nil))
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
