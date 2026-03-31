# Stage 1: Build
FROM golang:1.23-alpine AS builder

# Install git (needed for go mod to fetch from VCS)
RUN apk add --no-cache git

WORKDIR /app

# Copy dependency files first — Docker caches this layer
# Only re-downloads modules if go.mod or go.sum changes
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build a statically linked binary — no external dependencies at runtime
RUN CGO_ENABLED=0 GOOS=linux go build -o gateway ./main.go

# Stage 2: Run
FROM alpine:3.19

WORKDIR /app

# Copy only the compiled binary from the builder stage
# Final image has NO Go toolchain, no source code — just the binary
COPY --from=builder /app/gateway .

EXPOSE 8080

CMD ["./gateway"]