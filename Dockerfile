# syntax=docker/dockerfile:1

FROM golang:1.21-alpine AS builder
WORKDIR /app

RUN apk add --no-cache git build-base
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/app ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/migrate ./cmd/migrate

FROM gcr.io/distroless/base-debian11
WORKDIR /app
ENV APP_NAME=matchee-services \
    HTTP_HOST=0.0.0.0 \
    HTTP_PORT=8080 \
    SHUTDOWN_TIMEOUT=10s

COPY --from=builder /out/app /usr/local/bin/app
COPY --from=builder /out/migrate /usr/local/bin/migrate
COPY migrations ./migrations

EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/app"]


