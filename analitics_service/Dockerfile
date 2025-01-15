FROM golang:1.22 AS builder

WORKDIR /app

COPY analitics_service/go.mod analitics_service/go.sum ./
RUN go mod download

COPY analitics_service/ .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/main.go

FROM scratch

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/config/config.json /app/config/config.json

EXPOSE 8003

CMD ["./main"]
