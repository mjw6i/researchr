package main

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestResult(t *testing.T) {
	db, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	store := DatabaseStore{db: db}
	_, err = store.getResult()

	if err != nil {
		t.Fatal(err)
	}
}

func TestGetAbsoluteDataError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT").WillReturnError(fmt.Errorf("SQL Error"))

	store := DatabaseStore{db: db}
	_, _, err = store.getAbsoluteData()

	assertError(t, "DB error", err)

	err = mock.ExpectationsWereMet()

	if err != nil {
		t.Fatal(err)
	}
}

func TestGetAbsoluteDataEmpty(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"count", "responsive", "extremity"}).AddRow(0, 0, 0)

	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	store := DatabaseStore{db: db}
	remainedPercent, averageRemoved, err := store.getAbsoluteData()

	if err != nil {
		t.Fatal(err)
	}

	assertFloat(t, 0, remainedPercent)
	assertFloat(t, 0, averageRemoved)

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetAbsoluteDataFilled(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	count := 3
	responsive := 2
	extremitiesRemaining := 11
	extremitiesRemoved := 13

	rows := sqlmock.NewRows([]string{"count", "responsive", "extremity"}).AddRow(count, responsive, extremitiesRemaining)

	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	store := DatabaseStore{db: db}
	remainedPercent, averageRemoved, err := store.getAbsoluteData()

	if err != nil {
		t.Fatal(err)
	}

	assertFloat(t, 100*float64(responsive)/float64(count), remainedPercent)
	assertFloat(t, float64(extremitiesRemoved)/float64(count), averageRemoved)

	err = mock.ExpectationsWereMet()

	if err != nil {
		t.Fatal(err)
	}
}

func TestGetHeadlessDataError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT").WillReturnError(fmt.Errorf("SQL Error"))

	store := DatabaseStore{db: db}
	_, err = store.getHeadlessData()

	assertError(t, "DB error", err)

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetHeadlessDataEmpty(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"count", "responsive"}).AddRow(0, 0)

	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	store := DatabaseStore{db: db}
	remainedPercent, err := store.getHeadlessData()

	if err != nil {
		t.Fatal(err)
	}

	assertFloat(t, 0, remainedPercent)

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetHeadlessDataFilled(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	count := 3
	responsive := 2

	rows := sqlmock.NewRows([]string{"count", "responsive"}).AddRow(count, responsive)

	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	store := DatabaseStore{db: db}
	remainedPercent, err := store.getHeadlessData()

	if err != nil {
		t.Fatal(err)
	}

	assertFloat(t, 100*float64(responsive)/float64(count), remainedPercent)

	err = mock.ExpectationsWereMet()

	if err != nil {
		t.Fatal(err)
	}
}

func TestGetExtremitiesMissingDataError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT").WillReturnError(fmt.Errorf("SQL Error"))

	store := DatabaseStore{db: db}

	_, err = store.getExtremitiesMissingData()

	assertError(t, "DB error", err)

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetExtremitiesMissingDataEmpty(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	for i := 1; i <= 9; i++ {
		rows := sqlmock.NewRows([]string{"count", "sum"}).AddRow(0, 0)

		mock.ExpectQuery("SELECT").WillReturnRows(rows)
	}

	store := DatabaseStore{db: db}
	remainedPercentages, err := store.getExtremitiesMissingData()

	if err != nil {
		t.Fatal(err)
	}

	for _, percent := range remainedPercentages {
		assertFloat(t, 0, percent)
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetExtremitiesMissingDataFilled(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	total := 11
	responsive := 3

	for i := 1; i <= 9; i++ {
		rows := sqlmock.NewRows([]string{"count", "sum"}).AddRow(total, responsive)

		mock.ExpectQuery("SELECT").WillReturnRows(rows)
	}

	store := DatabaseStore{db: db}
	remainedPercentages, err := store.getExtremitiesMissingData()

	if err != nil {
		t.Fatal(err)
	}

	for _, percent := range remainedPercentages {
		assertFloat(t, 100*float64(responsive)/float64(total), percent)
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Fatal(err)
	}
}

func TestStoreError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectExec("INSERT").WillReturnError(fmt.Errorf("SQL Error"))

	store := DatabaseStore{db: db}

	e := Experiment{
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

	err = store.storeExperiment(e)

	assertError(t, "Store error", err)

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Fatal(err)
	}
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
