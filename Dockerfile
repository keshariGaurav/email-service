# Start from the official golang image
FROM golang:1.24-alpine

# Add git for go mod download
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Use a smaller image for the final build
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=0 /app/main .
COPY --from=0 /app/templates ./templates

# Expose the port the service runs on
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
