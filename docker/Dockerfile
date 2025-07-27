FROM golang:1.24 AS builder

WORKDIR /app

COPY ../go.mod ../go.sum ./
RUN go mod download

COPY ../ ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o transaction-test ./cmd/server/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/transaction-test .

COPY ../config/ ./config/

COPY ./migrations ./migrations

EXPOSE 8080

CMD ["./transaction-test"]