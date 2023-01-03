dev:
	DATABASE_URL=postgres://root:root@localhost/assessment?sslmode=disable PORT=:2565 go run server.go

fmt:
	gofmt -w .

it_test:
	docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit --exit-code-from it_tests

it_test_down:
	docker-compose -f docker-compose.test.yml down
