FROM golang:1.27.0-alpine3.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o edge-portal main.go

FROM alpine:3.23

WORKDIR /app

COPY --from=builder /app/edge-portal .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static

EXPOSE 8080

CMD ["./edge-portal"]
