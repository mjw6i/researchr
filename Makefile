export DATABASE_URL = postgresql://127.0.0.1/rsc?user=postgres&password=secret

run:
	go run .

test:
	go test

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
