FROM golang:alpine

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go test ./handler/
RUN go build -o main .

WORKDIR /dist

COPY configs/config.yml .

RUN cp /build/main .

EXPOSE 1323

CMD ["/dist/main"]
