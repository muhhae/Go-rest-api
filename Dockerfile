FROM golang:1.21.3

WORKDIR /app

COPY . /app

# Install dependencies
RUN go get -d -v ./...
RUN go install -v ./...

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]

# Load .env file and set environment variables
ENV GO_ENV=production
ADD .env /app/.env

# Set the correct port
ENV PORT=8080

# Expose the correct port
EXPOSE $PORT

# Run the application
CMD ["./main"]