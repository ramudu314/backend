# Use the official Golang image
FROM golang:1.20-alpine

# Set the working directory
WORKDIR /app

# Copy the Go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o main .

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./main"]