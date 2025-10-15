FROM golang:1.24.6

WORKDIR /app

COPY go.mod go.sum  ./

RUN go mod download

COPY . .

COPY .env .env 

RUN go install github.com/swaggo/swag/cmd/swag@v1.16.3
RUN go get -u github.com/swaggo/swag
RUN go mod tidy
RUN swag init -g ./cmd/main.go
RUN go build -o subscriptions-service ./cmd/main.go

EXPOSE 8080

CMD ["./subscriptions-service"]