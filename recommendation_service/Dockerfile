FROM golang:1.22 AS builder

WORKDIR /app

COPY recommendation_service/go.mod recommendation_service/go.sum ./
RUN go mod download

COPY recommendation_service/ .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/main.go

FROM scratch

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/config/config.json /app/config/config.json

EXPOSE 8001

CMD ["./main"]