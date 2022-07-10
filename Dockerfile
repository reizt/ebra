# Dockerfile
FROM golang:1.18

ENV TZ /usr/share/zoneinfo/Asia/Tokyo

ENV ROOT=/app
WORKDIR ${ROOT}

ENV GO111MODULE=on

COPY . .
EXPOSE 8080

RUN go get github.com/labstack/echo/v4
RUN go install github.com/cosmtrek/air@latest
