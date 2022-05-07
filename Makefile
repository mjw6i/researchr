export DATABASE_URL = postgresql://postgres:secret@127.0.0.1:5432/rsc
export COCKROACH_URL = postgres://root:@localhost:26257/defaultdb?pool_max_conns=64
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

coc-up:
	cockroach start-single-node \
		--background \
		--store path=/mnt/scrapyard,size=20GB \
		--cache=2GB \
		--listen-addr=localhost:26257 \
		--http-addr=localhost:8080 \
		--insecure

coc-down:
	cockroach quit \
		--host=localhost:26257 \
		--insecure

db-down:
	podman pod stop pg-rsc && podman pod rm pg-rsc

docker-db-down:
	docker stop rsc && docker rm rsc

db-conn:
	podman exec -it pg-rsc-container psql "${DATABASE_URL}"

stat:
	watch -n 1 -t "ss -t | grep '127.0.0.1:cockroach' | wc -l"

ab:
	ab -c 20 -n 20000 http://localhost:9000/results

ab-post:
	ab -c 200 -n 20000 -p ab-post-data -T 'application/x-www-form-urlencoded' http://localhost:9000/receive
