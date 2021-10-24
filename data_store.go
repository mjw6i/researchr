package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

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
	RemainedResponsive0MissingPercent string
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

type DatabaseStore struct {
	db *sql.DB
}

func (store *DatabaseStore) getResult() (Result, error) {
	remainedResponsivePercent, averageExtremitiesRemoved, err := store.getAbsoluteData()

	if err != nil {
		return Result{}, errors.New("DB error")
	}

	remainedResponsiveHeadlessPercent, err := store.getHeadlessData()

	if err != nil {
		return Result{}, errors.New("DB error")
	}

	remainedResponsiveMissingPercent, err := store.getExtremitiesMissingData()

	if err != nil {
		return Result{}, errors.New("DB error")
	}

	return Result{
		RemainedResponsivePercent:         fmt.Sprintf("%.2f", remainedResponsivePercent),
		RemainedResponsiveHeadlessPercent: fmt.Sprintf("%.2f", remainedResponsiveHeadlessPercent),
		AverageExtremitiesRemoved:         fmt.Sprintf("%.2f", averageExtremitiesRemoved),
		RemainedResponsive0MissingPercent: fmt.Sprintf("%.2f", remainedResponsiveMissingPercent[0]),
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

func (store *DatabaseStore) storeExperiment(e Experiment) error {
	_, err := store.db.Exec(`
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

func (store *DatabaseStore) getAbsoluteData() (float64, float64, error) {
	var count, responsive, extremity int

	row := store.db.QueryRow(`
	SELECT
		COUNT(*),
		COALESCE(SUM(responsive::int), 0) AS responsive,
		COALESCE(SUM(leg1::int + leg2::int + leg3::int + leg4::int + leg5::int + leg6::int + wing1::int + wing2::int), 0) AS extremity
	FROM experiments`)

	err := row.Scan(&count, &responsive, &extremity)
	if err != nil {
		log.Println(err)
		return 0, 0, errors.New("DB error")
	}

	if count == 0 {
		return 0, 0, nil
	}

	remainedResponsivePercent := 100 * float64(responsive) / float64(count)
	averageExtremitiesRemoved := float64(8*count-extremity) / float64(count)

	return remainedResponsivePercent, averageExtremitiesRemoved, nil
}

func (store *DatabaseStore) getHeadlessData() (float64, error) {
	var headless, headlessResponsive int

	row := store.db.QueryRow(`
	SELECT
		COUNT(*),
		COALESCE(SUM(responsive::int), 0) AS responsive
	FROM experiments
	WHERE head = FALSE
	`)

	err := row.Scan(&headless, &headlessResponsive)
	if err != nil {
		log.Println(err)
		return 0, errors.New("DB error")
	}

	if headless == 0 {
		return 0, nil
	}

	remainedResponsiveHeadlessPercent := 100 * float64(headlessResponsive) / float64(headless)

	return remainedResponsiveHeadlessPercent, nil
}

func (store *DatabaseStore) getExtremitiesMissingData() (map[int]float64, error) {
	remainedResponsiveMissingPercent := make(map[int]float64)

	for missing := 0; missing <= 8; missing++ {
		remaining := 8 - missing
		row := store.db.QueryRow(`
			SELECT
				COUNT(*),
				COALESCE(SUM(responsive::int), 0)
			FROM (
				SELECT responsive
				FROM experiments
				GROUP BY id
				HAVING SUM(leg1::int + leg2::int + leg3::int + leg4::int + leg5::int + leg6::int + wing1::int + wing2::int) = $1
			) as rows
		`, remaining)

		var responsive, total int

		err := row.Scan(&total, &responsive)
		if err != nil {
			log.Println(err)
			return make(map[int]float64), errors.New("DB error")
		}

		if total == 0 {
			remainedResponsiveMissingPercent[missing] = 0
		} else {
			remainedResponsiveMissingPercent[missing] = 100 * float64(responsive) / float64(total)
		}
	}

	return remainedResponsiveMissingPercent, nil
}
