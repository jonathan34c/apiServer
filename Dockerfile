# Use the official Golang image as the base image
FROM golang:1.20-bullseye

# Set the working directory to /app
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . .

# install curl 
RUN apt-get -y update; apt-get -y install curl
# Build the Go app
RUN go mod tidy

# Expose port 8080 for the container
EXPOSE 8080

# Command to run the executable
CMD ["go", "run", "."]
