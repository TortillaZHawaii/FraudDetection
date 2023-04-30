FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o producer .

# Stage 2: Runtime stage
FROM golang:1.20-alpine
WORKDIR /app
COPY --from=builder /app/producer .

ENV KAFKA_BROKERS=kafka:9092
ENV KAFKA_TOPIC=test

CMD ["./producer"]
