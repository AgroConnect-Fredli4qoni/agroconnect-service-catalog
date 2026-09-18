# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod ./
RUN go mod download || true

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/service-catalog main.go

# Production stage
FROM alpine:3.18

RUN apk add --no-cache ca-certificates curl

WORKDIR /root/
COPY --from=builder /app/service-catalog .

EXPOSE 8081

HEALTHCHECK --interval=10s --timeout=5s --retries=3 \
  CMD curl -f http://localhost:8081/health || exit 1

CMD ["./service-catalog"]
