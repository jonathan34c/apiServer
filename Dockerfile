# Use the official Golang image as the base image
FROM golang:1.20-alpine

# Set the working directory to /app
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . .

# Build the Go app
RUN go mod tidy

# Expose port 8080 for the container
EXPOSE 8888

# Command to run the executable
CMD ["go", "run", "."]
