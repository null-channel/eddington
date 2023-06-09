FROM golang:1.20-buster as builder
WORKDIR /app

COPY go.* ./
RUN go mod download
COPY . .
RUN go build -o api ./main.go

# The Real container
FROM debian:buster-slim
COPY --from=builder /app/api /app/api
CMD ["/app/api"]
