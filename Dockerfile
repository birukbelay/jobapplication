# Use the official Golang image as the base image
FROM golang:alpine as builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go app with optimizations
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Use a minimal base image for the final stage
FROM alpine:latest

# Set the working directory
WORKDIR /root/
# this is needed for image upload, golangs multipart uses the /tmp folder
COPY --from=alpine /tmp /tmp 
# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Expose the port your app runs on (change 8080 if needed)
EXPOSE 8080

# Command to run the app
CMD ["./main"]