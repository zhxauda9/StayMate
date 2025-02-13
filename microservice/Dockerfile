FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the binary and name it "microservice_bin" to avoid conflicts
RUN go build -o microservice_bin ./cmd/microservice/microservice.go

EXPOSE 9090

# Execute the binary (make sure it's named microservice_bin)
CMD ["./microservice_bin"]
