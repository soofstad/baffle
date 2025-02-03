FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . /app
RUN go build -o main .
ENTRYPOINT ["/app/main"]