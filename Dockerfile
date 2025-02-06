FROM golang:1.23-alpine AS builder
WORKDIR /app/src
COPY src /app/src
RUN go build -o goatbff .

FROM alpine:latest AS prod
COPY --from=builder /app/src/goatbff /app/goatbff
ENTRYPOINT ["/app/goatbff"]
