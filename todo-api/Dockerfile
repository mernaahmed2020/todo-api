# Use Go base image
FROM golang:1.24.4

# Set working directory inside container
WORKDIR /app

# Copy go.mod and go.sum first to cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the app
COPY . .

# Build the Go app
RUN go build -o todo-api

# Set environment variables (or override using .env later)
ENV PORT=8080

# Expose port
EXPOSE 8080

# Run the app
CMD ["./todo-api"]
