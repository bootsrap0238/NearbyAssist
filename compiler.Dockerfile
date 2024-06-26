# Start from the official Golang image
FROM golang:1.21-bookworm as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN GOOS=windows GOARCH=amd64 go build -o nearbyassist.exe .

# List files for debugging
RUN ls -la /app
