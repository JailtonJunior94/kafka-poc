FROM golang:1.21.3 AS builder
WORKDIR /go/src/api

RUN apt-get update
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN make build-api

FROM alpine:3.16
WORKDIR /app/

COPY --from=builder /go/src/api/bin/api .

EXPOSE 80
EXPOSE 443
CMD ["./api"]