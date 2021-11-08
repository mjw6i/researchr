module github.com/mjw6i/researchr/internal

go 1.17

require (
	github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/jackc/pgx/v4 v4.13.0
	github.com/mjw6i/researchr/pkg v0.0.0-00010101000000-000000000000
)

require (
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.10.0 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.1.1 // indirect
	github.com/jackc/pgservicefile v0.0.0-20200714003250-2b9c44734f2b // indirect
	github.com/jackc/pgtype v1.8.1 // indirect
	golang.org/x/crypto v0.0.0-20210711020723-a769d52b0f97 // indirect
	golang.org/x/text v0.3.6 // indirect
)

replace github.com/mjw6i/researchr/pkg => ../pkg
