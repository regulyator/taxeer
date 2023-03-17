# Use the official Go image as the base image
FROM golang:1.20.1-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to the workspace
COPY go.mod go.sum ./

# Download all the dependencies
RUN go mod download

# Copy the entire source code from the current directory to the workspace
COPY . .

# Build the Go app
RUN go build -o main .

# Run the app
CMD ["./main"]

