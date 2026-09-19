FROM golang:1.27-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o edge-portal main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/edge-portal .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static

EXPOSE 8080

CMD ["./edge-portal"]
