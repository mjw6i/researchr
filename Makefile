export DATABASE_URL = postgresql://postgres:secret@127.0.0.1:5432/rsc

run:
	go run .

test:
	go test

race:
	go test -race

cover:
	go test -coverprofile=coverage.out && go tool cover -html=coverage.out

db-up:
	podman create --rm \
		--pod new:pg-rsc \
		--name pg-rsc-container \
		-p 5432:5432 \
		-e POSTGRES_DB=rsc \
		-e POSTGRES_USER=postgres \
		-e POSTGRES_PASSWORD=secret \
		docker.io/library/postgres:alpine && \
	podman cp migration.sql pg-rsc-container:/docker-entrypoint-initdb.d/ && \
	podman start pg-rsc-container

db-down:
	podman pod stop pg-rsc && podman pod rm pg-rsc

db-conn:
	podman exec -it pg-rsc-container psql "${DATABASE_URL}"

stat:
	watch -n 1 -t "ss -t | grep '127.0.0.1:postgres' | wc -l"
