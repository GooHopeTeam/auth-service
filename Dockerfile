FROM golang:1.18-alpine

WORKDIR /auth-service

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY internal/ ./internal
COPY cmd/ ./cmd
COPY migrations/ ./migrations

RUN go build ./cmd/main/main.go

EXPOSE 8080

CMD ["./main"]