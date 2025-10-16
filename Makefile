init:
	go mod download
	go install github.com/swaggo/swag/cmd/swag@v1.16.3
	go get -u github.com/swaggo/swag
	go mod tidy
	swag init -g ./cmd/main.go

run-build:
	docker-compose up --build

run:
	docker-compose up