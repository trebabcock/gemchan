# syntax=docker/dockerfile:1

FROM golang:1.19-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY app/ ./app
COPY assets ./assets
COPY db ./db
COPY keys ./keys
COPY main.go ./

RUN apk add build-base
RUN mkdir -p /etc/gemchan/

RUN go build -o /gemchan

EXPOSE 1965

CMD [ "/gemchan" ]