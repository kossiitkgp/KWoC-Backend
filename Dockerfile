FROM golang:1.19-alpine as builder
WORKDIR /app
COPY ./ ./

RUN go mod tidy
RUN gofmt -s -w .
RUN go build

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/kwoc20-backend ./

ENTRYPOINT [ "./kwoc20-backend" ]