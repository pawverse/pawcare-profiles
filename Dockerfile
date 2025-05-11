FROM golang:1.24 AS builder

WORKDIR /src

COPY go.mod go.sum .

RUN go mod download

COPY . .

RUN make build

FROM alpine:latest

WORKDIR /app

COPY --from=builder /src/bin /app

EXPOSE 8000
EXPOSE 9000
VOLUME /data/conf

CMD ["./server", "-conf", "/data/conf"]
