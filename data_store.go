package main

import (
	"database/sql"
	"errors"
	"log"
	"os"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type DataStore interface {
	getResult() (Result, error)
	storeExperiment(Experiment) error
}

type Result struct {
	RemainedResponsivePercent         string
	RemainedResponsiveHeadlessPercent string
	AverageExtremitiesRemoved         string
	RemainedResponsive1MissingPercent string
	RemainedResponsive2MissingPercent string
	RemainedResponsive3MissingPercent string
	RemainedResponsive4MissingPercent string
	RemainedResponsive5MissingPercent string
	RemainedResponsive6MissingPercent string
	RemainedResponsive7MissingPercent string
	RemainedResponsive8MissingPercent string
}

type Experiment struct {
	Responsive bool
	Head       bool
	Leg1       bool
	Leg2       bool
	Leg3       bool
	Leg4       bool
	Leg5       bool
	Leg6       bool
	Wing1      bool
	Wing2      bool
}

type SuccessStore struct{}

func (store SuccessStore) getResult() (Result, error) {
	result := Result{
		RemainedResponsivePercent:         "57.03",
		RemainedResponsiveHeadlessPercent: "1.63",
		AverageExtremitiesRemoved:         "3.72",
		RemainedResponsive1MissingPercent: "95.88",
		RemainedResponsive2MissingPercent: "87.80",
		RemainedResponsive3MissingPercent: "75.61",
		RemainedResponsive4MissingPercent: "65.91",
		RemainedResponsive5MissingPercent: "34.82",
		RemainedResponsive6MissingPercent: "24.27",
		RemainedResponsive7MissingPercent: "11.56",
		RemainedResponsive8MissingPercent: "0.03",
	}

	return result, nil
}

func (store SuccessStore) storeExperiment(e Experiment) error {
	return nil
}

type FailureStore struct{}

func (store FailureStore) getResult() (Result, error) {
	return Result{}, errors.New("Store error")
}

func (store FailureStore) storeExperiment(e Experiment) error {
	return errors.New("Store error")
}

type DatabaseStore struct{}

func (store DatabaseStore) getResult() (Result, error) {
	db, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Println(err)
		return Result{}, errors.New("DB error")
	}
	defer db.Close()

	return Result{}, errors.New("Store error")
}

func (store DatabaseStore) storeExperiment(e Experiment) error {
	db, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Println(err)
		return errors.New("DB error")
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO experiments (responsive, head, leg1, leg2, leg3, leg4, leg5, leg6, wing1, wing2) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)", e.Responsive, e.Head, e.Leg1, e.Leg2, e.Leg3, e.Leg4, e.Leg5, e.Leg6, e.Wing1, e.Wing2)

	if err != nil {
		log.Println(err)
		return errors.New("Store error")
	}

	return nil
}
