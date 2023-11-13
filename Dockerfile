FROM golang:1.19-alpine

WORKDIR /app

COPY . /app

RUN apk add build-base

RUN go build -o app ./cmd/backend.go

EXPOSE 8080

CMD ["./app"]
