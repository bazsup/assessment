export PRIVATE_IP = $(shell ipconfig getifaddr en0)

dev:
	DATABASE_URL=postgres://root:root@localhost/assessment?sslmode=disable PORT=:2565 go run server.go

fmt:
	gofmt -w .

it_test:
	docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit --exit-code-from it_tests

it_test_down:
	docker-compose -f docker-compose.test.yml down

unit_test:
	go clean -testcache && go test -v -race ./... --tags=unit

unit_test_cover:
	go clean -testcache && go test -v -cover ./... --tags=unit

clear_testcache:
	go clean -testcache

# require database
test_cover_html:
	make clear_testcache && \
	PORT=2565 DATABASE_URL=postgres://root:root@localhost/assessment?sslmode=disable go test -cover -coverprofile=c.out --tags=unit,integration ./... && \
	go tool cover -html=c.out -o coverage.html

env:
	sed 's/localhost/$(PRIVATE_IP)/' .env.example > .env.local

build:
	docker build -t ghcr.io/bazsup/assessment:v1 .

start:
	docker run --name assessment -d --env-file .env.local -p 8080:2565 assessment:latest
