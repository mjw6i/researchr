package internal

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type Store interface {
	getResult() (Result, error)
	storeExperiment(Experiment) error
}

type Result struct {
	RemainedResponsivePercent         string
	RemainedResponsiveHeadlessPercent string
	AverageExtremitiesRemoved         string
	RemainedResponsiveMissingPercent  [9]string
}

type DatabaseStore struct {
	db *sql.DB
}

func NewDatabaseStore(db *sql.DB) *DatabaseStore {
	store := DatabaseStore{db: db}
	return &store
}

type absoluteData struct {
	remainedResponsivePercent float64
	averageExtremitiesRemoved float64
	err                       error
}

type headlessData struct {
	remainedResponsivePercent float64
	err                       error
}

type missingData struct {
	remainedResponsivePercent [9]float64
	err                       error
}

func (store *DatabaseStore) getResult() (Result, error) {
	abs := make(chan absoluteData, 1)
	head := make(chan headlessData, 1)
	miss := make(chan missingData, 1)

	go func() {
		remainedResponsivePercent, averageExtremitiesRemoved, err := store.getAbsoluteData()
		abs <- absoluteData{remainedResponsivePercent, averageExtremitiesRemoved, err}
	}()

	go func() {
		remainedResponsivePercent, err := store.getHeadlessData()
		head <- headlessData{remainedResponsivePercent, err}
	}()

	go func() {
		remainedResponsivePercent, err := store.getExtremitiesMissingData()
		miss <- missingData{remainedResponsivePercent, err}
	}()

	absolute := <-abs
	headless := <-head
	missing := <-miss

	if absolute.err != nil {
		return Result{}, errors.New("DB error")
	}

	if headless.err != nil {
		return Result{}, errors.New("DB error")
	}

	if missing.err != nil {
		return Result{}, errors.New("DB error")
	}

	return Result{
		RemainedResponsivePercent:         fmt.Sprintf("%.2f", absolute.remainedResponsivePercent),
		RemainedResponsiveHeadlessPercent: fmt.Sprintf("%.2f", headless.remainedResponsivePercent),
		AverageExtremitiesRemoved:         fmt.Sprintf("%.2f", absolute.averageExtremitiesRemoved),
		RemainedResponsiveMissingPercent: [9]string{
			fmt.Sprintf("%.2f", missing.remainedResponsivePercent[0]),
			fmt.Sprintf("%.2f", missing.remainedResponsivePercent[1]),
			fmt.Sprintf("%.2f", missing.remainedResponsivePercent[2]),
			fmt.Sprintf("%.2f", missing.remainedResponsivePercent[3]),
			fmt.Sprintf("%.2f", missing.remainedResponsivePercent[4]),
			fmt.Sprintf("%.2f", missing.remainedResponsivePercent[5]),
			fmt.Sprintf("%.2f", missing.remainedResponsivePercent[6]),
			fmt.Sprintf("%.2f", missing.remainedResponsivePercent[7]),
			fmt.Sprintf("%.2f", missing.remainedResponsivePercent[8]),
		},
	}, nil
}

func (store *DatabaseStore) storeExperiment(e Experiment) error {
	_, err := store.db.Exec(`
	INSERT INTO experiments (
		responsive, head, leg1, leg2, leg3, leg4, leg5, leg6, wing1, wing2, extremities
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
	)`, e.Responsive, e.Head, e.Leg1, e.Leg2, e.Leg3, e.Leg4, e.Leg5, e.Leg6, e.Wing1, e.Wing2, e.Extremities())

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
		COALESCE(SUM(extremities), 0) AS extremity
	FROM experiments
	`)

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
	var count, responsive int

	row := store.db.QueryRow(`
	SELECT
		COUNT(*),
		COALESCE(SUM(responsive::int), 0) AS responsive
	FROM experiments
	WHERE head = FALSE
	`)

	err := row.Scan(&count, &responsive)
	if err != nil {
		log.Println(err)
		return 0, errors.New("DB error")
	}

	if count == 0 {
		return 0, nil
	}

	percent := 100 * float64(responsive) / float64(count)

	return percent, nil
}

func (store *DatabaseStore) getExtremitiesMissingData() ([9]float64, error) {
	var remainedResponsive [9]float64

	rows, err := store.db.Query(`
	SELECT
		extremities,
		COUNT(*),
		COALESCE(SUM(responsive::int), 0) as responsive
	FROM experiments
	GROUP BY extremities
	`)
	if err != nil {
		log.Println(err)
		return [9]float64{}, errors.New("DB error")
	}
	defer rows.Close()

	for rows.Next() {
		var extremities, total, responsive int

		err := rows.Scan(&extremities, &total, &responsive)
		if err != nil {
			log.Println(err)
			return [9]float64{}, errors.New("DB error")
		}

		missing := 8 - extremities
		remainedResponsive[missing] = 100 * float64(responsive) / float64(total)
	}

	err = rows.Err()
	if err != nil {
		log.Println(err)
		return [9]float64{}, errors.New("DB error")
	}

	return remainedResponsive, nil
}
