# Build Stage
FROM golang:1.24 AS builder

ENV CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

# Install OS-level dependencies for C libraries
RUN apt-get update && apt-get install -y gcc musl-dev

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./
RUN go build -o app .


# Run Stage
FROM debian:12.10-slim

# Install any runtime C library dependencies (e.g. sqlite3)
RUN apt-get update && apt-get install -y libsqlite3-0 ca-certificates && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/app /app
ENTRYPOINT ["/app"]

