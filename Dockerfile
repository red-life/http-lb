FROM golang:1.21-bookworm AS builder
WORKDIR /src
COPY . ./
RUN go mod download
RUN go build -o http-lb ./cmd/main.go

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /src/http-lb ./http-lb
CMD ["./http-lb"]
