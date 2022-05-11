package internal

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mjw6i/researchr/pkg"
)

const absoluteQuery = "SELECT .* FROM experiments$"
const headlessQuery = "SELECT .* WHERE head = FALSE$"
const missingQuery = "SELECT .* GROUP BY extremities$"

func TestResult(t *testing.T) {
	db, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	store := NewDatabaseStore(db)
	_, err = store.getResult()

	if err != nil {
		t.Fatal(err)
	}
}

func BenchmarkResult(b *testing.B) {
	db, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		b.Fatal(err)
	}
	defer db.Close()

	store := NewDatabaseStore(db)

	for i := 0; i < b.N; i++ {
		_, err = store.getResult()

		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestGetResultErrorAbsolute(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectQuery(absoluteQuery).WillReturnError(fmt.Errorf("SQL Error"))

	store := DatabaseStore{db: db}
	_, err = store.getResult()

	pkg.AssertError(t, "DB error", err)

	err = mock.ExpectationsWereMet()

	if err != nil {
		t.Fatal(err)
	}
}

func TestGetResultErrorHeadless(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"count", "responsive", "extremity"}).AddRow(0, 0, 0)
	mock.ExpectQuery(absoluteQuery).WillReturnRows(rows)
	mock.ExpectQuery(headlessQuery).WillReturnError(fmt.Errorf("SQL Error"))

	store := DatabaseStore{db: db}
	_, err = store.getResult()

	pkg.AssertError(t, "DB error", err)

	err = mock.ExpectationsWereMet()

	if err != nil {
		t.Fatal(err)
	}
}

func TestGetResultErrorMissing(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"count", "responsive", "extremity"}).AddRow(0, 0, 0)
	mock.ExpectQuery(absoluteQuery).WillReturnRows(rows)
	rows = sqlmock.NewRows([]string{"count", "responsive"}).AddRow(0, 0)
	mock.ExpectQuery(headlessQuery).WillReturnRows(rows)
	mock.ExpectQuery(missingQuery).WillReturnError(fmt.Errorf("SQL Error"))

	store := DatabaseStore{db: db}
	_, err = store.getResult()

	pkg.AssertError(t, "DB error", err)

	err = mock.ExpectationsWereMet()

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

	mock.ExpectQuery(absoluteQuery).WillReturnError(fmt.Errorf("SQL Error"))

	store := DatabaseStore{db: db}
	_, _, err = store.getAbsoluteData()

	pkg.AssertError(t, "DB error", err)

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
	mock.ExpectQuery(absoluteQuery).WillReturnRows(rows)

	store := DatabaseStore{db: db}
	remainedPercent, averageRemoved, err := store.getAbsoluteData()

	if err != nil {
		t.Fatal(err)
	}

	pkg.AssertFloat(t, 0, remainedPercent)
	pkg.AssertFloat(t, 0, averageRemoved)

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
	mock.ExpectQuery(absoluteQuery).WillReturnRows(rows)

	store := DatabaseStore{db: db}
	remainedPercent, averageRemoved, err := store.getAbsoluteData()

	if err != nil {
		t.Fatal(err)
	}

	pkg.AssertFloat(t, 100*float64(responsive)/float64(count), remainedPercent)
	pkg.AssertFloat(t, float64(extremitiesRemoved)/float64(count), averageRemoved)

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

	mock.ExpectQuery(headlessQuery).WillReturnError(fmt.Errorf("SQL Error"))

	store := DatabaseStore{db: db}
	_, err = store.getHeadlessData()

	pkg.AssertError(t, "DB error", err)

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
	mock.ExpectQuery(headlessQuery).WillReturnRows(rows)

	store := DatabaseStore{db: db}
	remainedPercent, err := store.getHeadlessData()

	if err != nil {
		t.Fatal(err)
	}

	pkg.AssertFloat(t, 0, remainedPercent)

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
	mock.ExpectQuery(headlessQuery).WillReturnRows(rows)

	store := DatabaseStore{db: db}
	remainedPercent, err := store.getHeadlessData()

	if err != nil {
		t.Fatal(err)
	}

	pkg.AssertFloat(t, 100*float64(responsive)/float64(count), remainedPercent)

	err = mock.ExpectationsWereMet()

	if err != nil {
		t.Fatal(err)
	}
}

func TestGetExtremitiesMissingDataQueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectQuery(missingQuery).WillReturnError(fmt.Errorf("SQL Error"))

	store := DatabaseStore{db: db}

	_, err = store.getExtremitiesMissingData()

	pkg.AssertError(t, "DB error", err)

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetExtremitiesMissingDataScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"extremities", "count", "responsive"}).
		AddRow(0, 0, 0).AddRow(nil, nil, nil)
	mock.ExpectQuery(missingQuery).WillReturnRows(rows)

	store := DatabaseStore{db: db}

	_, err = store.getExtremitiesMissingData()

	pkg.AssertError(t, "DB error", err)

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetExtremitiesMissingDataErrError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"extremities", "count", "responsive"}).
		AddRow(0, 0, 0).RowError(0, fmt.Errorf("SQL Error"))
	mock.ExpectQuery(missingQuery).WillReturnRows(rows)

	store := DatabaseStore{db: db}

	_, err = store.getExtremitiesMissingData()

	pkg.AssertError(t, "DB error", err)

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

	rows := sqlmock.NewRows([]string{"extremities", "count", "responsive"})
	mock.ExpectQuery(missingQuery).WillReturnRows(rows)

	store := DatabaseStore{db: db}
	remainedPercentages, err := store.getExtremitiesMissingData()

	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i <= 8; i++ {
		percent := remainedPercentages[i]
		pkg.AssertFloat(t, 0, percent)
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

	rows := sqlmock.NewRows([]string{"extremities", "count", "responsive"})
	for i := 0; i <= 8; i++ {
		rows.AddRow(i, total, responsive)
	}
	mock.ExpectQuery(missingQuery).WillReturnRows(rows)

	store := DatabaseStore{db: db}
	remainedPercentages, err := store.getExtremitiesMissingData()

	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i <= 8; i++ {
		percent := remainedPercentages[i]
		pkg.AssertFloat(t, 100*float64(responsive)/float64(total), percent)
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

	pkg.AssertError(t, "Store error", err)

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

	store := NewDatabaseStore(db)

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

func BenchmarkStore(b *testing.B) {
	db, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		b.Fatal(err)
	}
	defer db.Close()

	store := NewDatabaseStore(db)

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

	for i := 0; i < b.N; i++ {
		err = store.storeExperiment(e)

		if err != nil {
			b.Fatal(err)
		}
	}
}
