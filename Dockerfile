# ---------- Builder Stage ----------
FROM golang:1.24-alpine AS builder

# Install build tools
RUN apk add --no-cache git bash make

# Install goose (migration tool)
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

WORKDIR /app
COPY . .

# Download deps and build binary
RUN go mod tidy
RUN go build -o okies_core ./cmd/server

# ---------- Runtime Stage ----------
FROM alpine:3.19

WORKDIR /root/

# Copy binary
COPY --from=builder /app/okies_core .
# Copy goose binary too
COPY --from=builder /go/bin/goose /usr/local/bin/goose
# Copy migrations directory
COPY --from=builder /app/migrations ./migrations

# Default command: run the app
CMD ["./okies_core"]
