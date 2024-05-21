# Start from the official Golang image so we have a good base.
FROM golang:1.22 as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Set necessary environmet variables needed for our image and build the application
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# Build the Go app
RUN go build -o AQI-predictor .

# Start a new stage from scratch for a smaller, final image
FROM alpine:latest

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/ .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./AQI-predictor"]
