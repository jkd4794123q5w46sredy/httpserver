# Use the official Go image as a base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules manifests
COPY . .

RUN go mod download

# Build the Go application
RUN go build -o main .

# Command to run when the container starts
CMD ["/app/main"]
