FROM golang:1.16-alpine

WORKDIR /app

# Copy go.mod first for better caching
COPY go.mod ./
RUN go mod download

# Copy the source code
COPY main.go ./

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]