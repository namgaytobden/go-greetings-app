FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

# Build the Go application
RUN go build -o go-greetings

FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/go-greetings .

EXPOSE 8080

# Run the Go binary
CMD ["./go-greetings"]
