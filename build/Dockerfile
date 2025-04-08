# Stage 1: Build the Go app
FROM golang:1.23.8-alpine AS build

WORKDIR /app

# Copy dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the Go app
RUN go build -o /event ./cmd/server  # Ensure correct path

# Stage 2: Run the Go app
FROM alpine:latest

# Install necessary libraries
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary correctly
COPY --from=build /event /event

# Ensure execution permissions
RUN chmod +x /event

# Expose port
EXPOSE 8080

# Run the application
CMD ["/event"]
