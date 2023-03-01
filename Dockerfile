# syntax=docker/dockerfile:1

FROM golang:1.20.1-alpine3.17

WORKDIR /app

COPY server/go.mod ./
COPY server/go.sum ./

RUN go mod download

COPY server/*.go ./
COPY server/pb ./pb

RUN go build -o main

EXPOSE 9000

CMD [ "./main" ]