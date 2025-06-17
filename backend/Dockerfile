FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
COPY static/ .
RUN go mod download

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY . .

RUN go build -o main ./cmd/server/main.go

COPY ./migrations ./migrations

FROM golang:1.24

WORKDIR /app

RUN apt-get update && apt-get install -y netcat-openbsd

COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY --from=builder /app/static ./static

COPY ./internal/config ./internal/config

EXPOSE 8080

CMD ["./main"]
