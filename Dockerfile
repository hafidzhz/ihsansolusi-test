# Use official Golang image as the base
FROM golang:1.20-alpine

# Install bash and netcat-openbsd (for nc command to check port)
RUN apk add --no-cache bash netcat-openbsd

# Set the working directory inside the container
WORKDIR /app

# Copy wait-for-it script into the container
COPY ./scripts/wait-for-it.sh ./scripts/wait-for-it.sh
RUN chmod +x ./scripts/wait-for-it.sh

# Copy go module files and install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application (replace 'myapp' with your actual Go app binary name)
RUN go build -o myapp .

# Expose the port the app will run on
EXPOSE ${APP_PORT}

# Wait for PostgreSQL and then run the Go app
CMD ./scripts/wait-for-it.sh postgres:5432 -- ./myapp
