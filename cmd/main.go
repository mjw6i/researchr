package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mjw6i/researchr/internal"
)

var static = http.StripPrefix("/static", http.FileServer(http.Dir("../web/static")))

func main() {
	var addr string
	flag.StringVar(&addr, "addr", ":9000", "listen address")
	flag.Parse()

	config, err := pgxpool.ParseConfig(os.Getenv("COCKROACH_URL"))
	if err != nil {
		log.Fatal("error configuring the database: ", err)
	}

	pool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()

	ds := internal.NewDatabaseStore(pool)
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
