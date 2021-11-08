export DATABASE_URL = postgresql://postgres:secret@127.0.0.1:5432/rsc
export PACKAGES = ./cmd github.com/mjw6i/researchr/internal

run:
	cd cmd && go run .

test:
	go test ${PACKAGES}

race:
	go test ${PACKAGES} -race

cover:
	go vet ${PACKAGES} && \
	go test ${PACKAGES} -coverprofile=coverage.out && \
	go tool cover -html=coverage.out

st:
	staticcheck ${PACKAGES}

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

docker-db-up:
	docker create --rm --network host \
		--name rsc \
		-e POSTGRES_DB=rsc \
		-e POSTGRES_USER=postgres \
		-e POSTGRES_PASSWORD=secret \
		postgres:alpine && \
    docker cp migration.sql rsc:/docker-entrypoint-initdb.d/ && \
    docker start rsc

db-down:
	podman pod stop pg-rsc && podman pod rm pg-rsc

docker-db-down:
	docker stop rsc && docker rm rsc

db-conn:
	podman exec -it pg-rsc-container psql "${DATABASE_URL}"

stat:
	watch -n 1 -t "ss -t | grep '127.0.0.1:postgres' | wc -l"

ab:
	ab -c 200 -n 20000 http://localhost:9000/results
