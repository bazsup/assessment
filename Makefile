dev:
	DATABASE_URL=postgres://root:root@localhost/assessment?sslmode=disable PORT=:2565 go run server.go

fmt:
	gofmt -w .