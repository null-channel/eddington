#simple dockerfile to run sleep in alpine
FROM golang:alpine as builder

WORKDIR /app
COPY ./go.mod .
RUN go mod download
RUN apk add --no-cache gcc musl-dev
COPY . .
RUN CGO_ENABLED=1 GOOS=linux  go build -o container-builder

FROM alpine

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/container-builder /container-builder

EXPOSE 8080

ENTRYPOINT ["/container-builder"]