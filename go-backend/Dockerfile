# Use the official Golang image as the base image
FROM golang:1.23-alpine

# Install bash and air
RUN apk update && apk add --no-cache bash curl git && \
    curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

# Set the working directory inside the container
WORKDIR /app

# Install the dependencies and tools
RUN go install github.com/air-verse/air@latest
# Copy go.mod and go.sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .
# Expose the port on which the app will run
EXPOSE 8080

# Command to run the Air tool (hot-reload)
CMD ["air"]
