clean:
	rm -rf build

build: clean
	go build -o build/bombur

build_image: clean
	docker build . -t bombur:latest

run:
	go run .

test:
	# PGSQL must already be running
	BOMBUR_DB_URI="postgresql://localhost/bombur?user=bombur&password=bombur" go test ./... -v

lint:
	docker run --rm -v $(PWD):/app -w /app golangci/golangci-lint:v1.38.0 golangci-lint run -v

cover:
	BOMBUR_DB_URI="postgresql://localhost/bombur?user=bombur&password=bombur" go test 'cover ./...

fmt:
	docker run --rm -v $(PWD):/data cytopia/gofmt -s -w .

# Command wrapped into docker and used mostly for CI
ci-test: dev_stop_all_containers
	docker network create bombur-network
	docker run -d --rm --name bombur_pg --network=bombur-network -e POSTGRES_DB=bombur -e POSTGRES_USER=bombur -e POSTGRES_PASSWORD=bombur postgres:13.3-buster
	sleep 5 # To let PG boot
	docker run --rm --name bombur_app --network=bombur-network -v $(PWD):/go/src/app -w /go/src/app -e BOMBUR_DB_URI="postgresql://bombur_pg/bombur?user=bombur&password=bombur" golang:1.16.7-buster go test ./... -v

ci-cover: dev_stop_all_containers
	docker network create bombur-network
	docker run -d --rm --name bombur_pg --network=bombur-network -e POSTGRES_DB=bombur -e POSTGRES_USER=bombur -e POSTGRES_PASSWORD=bombur postgres:13.3-buster
	sleep 5 sleep 5 # To let PG boot
	docker run --rm --name bombur_app --network=bombur-network -v $(PWD):/go/src/app -w /go/src/app -e BOMBUR_DB_URI="postgresql://bombur_pg/bombur?user=bombur&password=bombur" golang:1.16.7-buster go test -cover ./...


# Command to ease development
dev_start:
	BOMBUR_DB_URI="postgresql://localhost/bombur?user=bombur&password=bombur" go run .

dev_start_container:
	docker build . -t bombur_app
	docker run --rm --network=bombur-network -p 7777:7777 -e BOMBUR_DB_URI="postgresql://bombur_pg/bombur?user=bombur&password=bombur" bombur_app

dev_start_pg13:
	docker run --rm --name bombur_pg -p 5432:5432 -e POSTGRES_DB=bombur -e POSTGRES_USER=bombur -e POSTGRES_PASSWORD=bombur postgres:13.3-buster

dev_start_pg13_daemon:
	docker run -d --rm --network=bombur-network --name bombur_pg -p 5432:5432 -e POSTGRES_DB=bombur -e POSTGRES_USER=bombur -e POSTGRES_PASSWORD=bombur postgres:13.3-buster

dev_connect_pg13:
	docker exec -ti `docker ps -aqf "name=bombur_pg"` psql --user bombur

dev_stop_all_containers:
	-docker stop $$(docker ps -q)
	-docker container rm $$(docker container ls -aq)
	-docker volume rm -f $$(docker volume ls -q)
	-docker network prune -f 