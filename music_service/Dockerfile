FROM golang:1.22 AS builder

WORKDIR /app

COPY music_service/go.mod music_service/go.sum ./
RUN go mod download

COPY music_service/ .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/main.go

FROM scratch

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/config/config.json /app/config/config.json

EXPOSE 8000

CMD ["./main"]