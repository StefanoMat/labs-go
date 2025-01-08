# Build Stage
FROM golang:1.23 AS builder
WORKDIR /app

# Copy Go dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy all project files
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app .

# Final Stage
FROM alpine:3.18
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/app .

# Expose the application port
EXPOSE 8080

# Command to run the application
CMD ["/app/app"]
