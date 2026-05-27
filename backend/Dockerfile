# 🏗️ Stage 1: Build the Go binary using the official Golang image
FROM golang:latest AS builder

# Set the current working directory inside the container
WORKDIR /app

# Copy the dependency management files first to leverage Docker cache
COPY go.mod go.sum ./

# Download all required dependencies specified in the go.mod file
RUN go mod download

# Copy the entire source code into the container working directory
COPY . .

# Compile the application into a standalone static binary for Linux
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o queue-ticket-api ./cmd/main.go


# 📦 Stage 2: Create a minimal production image using Alpine Linux
FROM alpine:3.19

# Set the active working directory for the runtime container
WORKDIR /root/

# Install the required CA certificates to allow secure HTTPS networking connections
RUN apk add --no-cache ca-certificates

# Copy only the compiled executable binary from the builder stage
COPY --from=builder /app/queue-ticket-api .

# Copy the environment configuration blueprint file to initialize runtime variables
COPY .env.example .env

# Inform Docker that the container listens on port 3000 at runtime
EXPOSE 3000

# Specify the command to execute the backend server binary when the container starts
CMD ["./queue-ticket-api"]
