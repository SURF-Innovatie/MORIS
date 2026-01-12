# Builder stage
FROM golang:1.25.4-alpine AS builder

# Install core build dependencies
RUN apk add --no-cache git

WORKDIR /app

# Copy dependency files first for better caching
# Note: This assumes build context is the root of the monorepo or at least contains apps/backend
# If building from apps/backend directly, adjust paths accordingly. 
# Based on plan, we assume root context context, so we copy specifically.
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
# Targeting the main entry point defined in package.json dev script: cmd/dev/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/dev/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies if needed (e.g. ca-certificates for external requests)
RUN apk add --no-cache ca-certificates

# Copy binary from builder
COPY --from=builder /server /server

# Expose port (default for Go servers is often 8080 or depends on env, checking docker-compose usually helps but 8080 is safe default)
EXPOSE 8080

CMD ["/server"]
