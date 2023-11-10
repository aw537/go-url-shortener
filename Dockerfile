FROM golang:1.21.4

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod file
COPY go.mod ./

# Download all dependencies. 
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o shortener .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./shortener"]
