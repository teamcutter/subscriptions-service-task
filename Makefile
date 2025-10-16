init:
	go mod download
	go install github.com/swaggo/swag/cmd/swag@v1.16.3
	go get -u github.com/swaggo/swag
	swag init -g ./cmd/main.go
	go mod tidy

test:
	go test ./internal/handler ./internal/repo -v

run-build:
	docker-compose up --build

run:
	docker-compose up