FROM golang:1.15-alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go test ./internal/api/
RUN go build -o istekbin-api ./cmd/istekbin-api

FROM alpine:3.11

WORKDIR /app

COPY --from=builder build/istekbin-api .
COPY --from=builder build/configs/config.yml configs/config.yml

EXPOSE 1323

CMD ["/app/istekbin-api"]
