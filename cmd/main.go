package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/mjw6i/researchr/internal"
)

var static = http.StripPrefix("/static", http.FileServer(http.Dir("./static")))

func main() {
	var addr string
	flag.StringVar(&addr, "addr", ":9000", "listen address")
	flag.Parse()
	db, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxIdleConns(80)
	db.SetMaxOpenConns(80)
	defer db.Close()
	ds := internal.NewDatabaseStore(db)
	env := &internal.Env{Store: ds}

	http.HandleFunc("/", internal.BaseHandler)
	http.HandleFunc("/submit", internal.SubmitHandler)
	http.HandleFunc("/receive", env.ReceiveHandler)
	http.HandleFunc("/results", env.ResultsHandler)
	http.HandleFunc("/assets", internal.AssetsHandler)
	http.Handle("/static/", static)
	log.Println("Listening on: ", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
