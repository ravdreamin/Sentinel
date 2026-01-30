# Build Stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

# Run Stage
FROM alpine:latest

WORKDIR /app

# Install basic dependencies
RUN apk --no-cache add ca-certificates curl

# Install golang-migrate
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz && \
    mv migrate /usr/bin/migrate

# Copy binary from builder
COPY --from=builder /app/main . 

# Copy migrations (Assumes migrations directory is in the build context root)
COPY migrations /migrations

# Expose the API port
EXPOSE 8081

# Run migrations and then start the application
# We use sh -c to allow environment variable expansion for DATABASE_URL
CMD ["sh", "-c", "/usr/bin/migrate -path /migrations -database \"$DATABASE_URL\" up && ./main"]