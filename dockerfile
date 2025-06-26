
FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN ls -a

RUN go build -o main ./cmd/server/main.go

FROM debian:bookworm-slim   

WORKDIR /app

COPY --from=builder /app/main .

# Copy config files
COPY --from=builder /app/config.yaml ./
COPY  --from=builder /app/.env ./



CMD ["./main"] 