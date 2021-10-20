package main

import (
	"database/sql"
	"errors"
	"fmt"
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

	var count, responsive, extremity int

	row := db.QueryRow(`
	SELECT
		COUNT(*),
		COALESCE(SUM(responsive::int), 0) AS responsive,
		COALESCE(SUM(leg1::int + leg2::int + leg3::int + leg4::int + leg5::int + leg6::int + wing1::int + wing2::int), 0) AS extremity
	FROM experiments`)

	err = row.Scan(&count, &responsive, &extremity)
	if err != nil {
		log.Println(err)
		return Result{}, errors.New("DB error")
	}

	var remainedResponsivePercent, averageExtremitiesRemoved float64

	if count == 0 {
		remainedResponsivePercent = 0
		averageExtremitiesRemoved = 0
	} else {
		remainedResponsivePercent = float64(responsive) / float64(count)
		averageExtremitiesRemoved = float64(8*count-extremity) / float64(count)
	}

	var headless int

	row = db.QueryRow(`
	SELECT COUNT(*)
	FROM experiments
	WHERE responsive = TRUE AND head = FALSE
	`)

	err = row.Scan(&headless)
	if err != nil {
		log.Println(err)
		return Result{}, errors.New("DB error")
	}

	var remainedResponsiveHeadlessPercent float64

	if count == 0 {
		remainedResponsiveHeadlessPercent = 0
	} else {
		remainedResponsiveHeadlessPercent = float64(headless) / float64(count)
	}

	remainedResponsiveMissingPercent := make(map[int]float64)

	for missing := 1; missing <= 8; missing++ {
		remaining := 8 - missing
		row = db.QueryRow(`
			SELECT COUNT(*)
			FROM (
				SELECT id
				FROM experiments
				WHERE responsive = TRUE
				GROUP BY id
				HAVING SUM(leg1::int + leg2::int + leg3::int + leg4::int + leg5::int + leg6::int + wing1::int + wing2::int) = $1
			) as rows
		`, remaining)

		var temp int

		err = row.Scan(&temp)
		if err != nil {
			log.Println(err)
			return Result{}, errors.New("DB error")
		}

		if count == 0 {
			remainedResponsiveMissingPercent[missing] = 0
		} else {
			remainedResponsiveMissingPercent[missing] = float64(temp) / float64(count)
		}
	}

	log.Println(remainedResponsiveMissingPercent)

	return Result{
		RemainedResponsivePercent:         fmt.Sprintf("%.2f", remainedResponsivePercent),
		RemainedResponsiveHeadlessPercent: fmt.Sprintf("%.2f", remainedResponsiveHeadlessPercent),
		AverageExtremitiesRemoved:         fmt.Sprintf("%.2f", averageExtremitiesRemoved),
		RemainedResponsive1MissingPercent: fmt.Sprintf("%.2f", remainedResponsiveMissingPercent[1]),
		RemainedResponsive2MissingPercent: fmt.Sprintf("%.2f", remainedResponsiveMissingPercent[2]),
		RemainedResponsive3MissingPercent: fmt.Sprintf("%.2f", remainedResponsiveMissingPercent[3]),
		RemainedResponsive4MissingPercent: fmt.Sprintf("%.2f", remainedResponsiveMissingPercent[4]),
		RemainedResponsive5MissingPercent: fmt.Sprintf("%.2f", remainedResponsiveMissingPercent[5]),
		RemainedResponsive6MissingPercent: fmt.Sprintf("%.2f", remainedResponsiveMissingPercent[6]),
		RemainedResponsive7MissingPercent: fmt.Sprintf("%.2f", remainedResponsiveMissingPercent[7]),
		RemainedResponsive8MissingPercent: fmt.Sprintf("%.2f", remainedResponsiveMissingPercent[8]),
	}, nil

	return Result{}, errors.New("Store error")
}

func (store DatabaseStore) storeExperiment(e Experiment) error {
	db, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Println(err)
		return errors.New("DB error")
	}
	defer db.Close()

	_, err = db.Exec(`
	INSERT INTO experiments (
		responsive, head, leg1, leg2, leg3, leg4, leg5, leg6, wing1, wing2
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
	)`, e.Responsive, e.Head, e.Leg1, e.Leg2, e.Leg3, e.Leg4, e.Leg5, e.Leg6, e.Wing1, e.Wing2)

	if err != nil {
		log.Println(err)
		return errors.New("Store error")
	}

	return nil
}
