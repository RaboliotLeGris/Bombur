clean:
	rm -rf build

build: clean
	go build -o build/bombur

build_image: clean
	docker build . -t bombur:latest

run:
	go run .

lint:
	docker run --rm -v $(PWD):/app -w /app golangci/golangci-lint:v1.38.0 golangci-lint run -v

fmt:
	docker run --rm -v $(PWD):/data cytopia/gofmt -s -w .

dev_start:
	BOMBUR_DB_URI="postgresql://localhost/bombur?user=bombur&password=bombur" go run .

dev_start_pg13:
	docker run --rm --name bombur_pg -p 5432:5432 -e POSTGRES_DB=bombur -e POSTGRES_USER=bombur -e POSTGRES_PASSWORD=bombur postgres:13.3-buster

dev_start_pg13_daemon:
	docker run -d --rm --name bombur_pg -p 5432:5432 -e POSTGRES_DB=bombur -e POSTGRES_USER=bombur -e POSTGRES_PASSWORD=bombur postgres:13.3-buster

dev_connect_pg13:
	docker exec -ti `docker ps -aqf "name=bombur_pg"` psql --user bombur