init:
	go mod download
	go install github.com/swaggo/swag/cmd/swag@v1.16.3
	go get -u github.com/swaggo/swag
	swag init -g ./cmd/main.go
	go mod tidy

test:
	go test ./internal/handler ./internal/repo -v

up-build: init test
	docker-compose up --build

up:
	docker-compose up

down:
	docker-compose down