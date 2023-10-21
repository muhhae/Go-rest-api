FROM golang:1.21.3

WORKDIR /app

COPY . /app

# Install dependencies
RUN go get -d -v ./...
RUN go install -v ./...

RUN go build -o main .

EXPOSE 8080

ENV GIN_MODE=release

# Run the application
CMD ["./main"]