# Use the official Go image as the base image
FROM golang:1.19

# Set the working directory inside the container
WORKDIR /app

# Define build arguments for environment variables
ARG BOT_API_KEY
ARG DB_HOST_KEY
ARG DB_USER_KEY
ARG DB_PASSWORD_KEY

# Set runtime environment variables from build arguments
ENV BOT_API_KEY=$BOT_API_KEY
ENV DB_HOST_KEY=$DB_HOST_KEY
ENV DB_USER_KEY=$DB_USER_KEY
ENV DB_PASSWORD_KEY=$DB_PASSWORD_KEY

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

