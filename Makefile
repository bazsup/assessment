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

# require database
test_cover_html:
	go clean -testcache && DATABASE_URL=postgres://root:root@localhost/assessment?sslmode=disable go test -cover -coverprofile=c.out --tags=unit,integration ./... && go tool cover -html=c.out -o coverage.html
