# Use golang:1.19 as the base image
FROM golang:1.19

# Set the working directory in the container
WORKDIR /app

# Copy the Go file to the container
COPY main.go .
COPY go.mod .
COPY go.sum .

# Compile the Go file
RUN go build -o main .

Expose 10000

CMD ["./main", "--port", "10000"]
