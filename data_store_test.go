package main

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
)

func TestResult(t *testing.T) {
	db, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	store := DatabaseStore{db: db}
	res, err := store.getResult()

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(res)
}

func TestStore(t *testing.T) {
	db, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	store := DatabaseStore{db: db}

	e1 := Experiment{
		Responsive: true,
		Head:       false,
		Leg1:       true,
		Leg2:       false,
		Leg3:       true,
		Leg4:       false,
		Leg5:       true,
		Leg6:       false,
		Wing1:      true,
		Wing2:      false,
	}

	e2 := Experiment{
		Responsive: false,
		Head:       true,
		Leg1:       false,
		Leg2:       true,
		Leg3:       false,
		Leg4:       true,
		Leg5:       false,
		Leg6:       true,
		Wing1:      false,
		Wing2:      true,
	}

	err = store.storeExperiment(e1)

	if err != nil {
		t.Fatal(err)
	}

	err = store.storeExperiment(e2)

	if err != nil {
		t.Fatal(err)
	}
}
