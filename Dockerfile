FROM golang:1.22.6-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /app/service ./cmd/service

FROM alpine:3.18.3

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/service .

EXPOSE 8080

ENTRYPOINT ["/app/service"]