FROM golang:1.20-buster as builder
WORKDIR /app
ADD . .
RUN go mod download
RUN go mod tidy
RUN go build -o ./container-builder ./main.go


# The Real container
FROM ubuntu:20.04
RUN apt-get update && apt-get install -y openssl ca-certificates curl docker.io
RUN update-ca-certificates
RUN apt-get install -y libssl-dev
RUN rm -rf /var/lib/apt/lists/*
COPY --from=builder /app/container-builder app/container-builder
CMD ["/app/container-builder"]
