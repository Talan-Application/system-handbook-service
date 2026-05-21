FROM golang:1.25.5-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /system-handbook-service ./cmd/server

FROM alpine:3.21
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /system-handbook-service .
COPY config/ config/

EXPOSE 50053

CMD ["./system-handbook-service"]
