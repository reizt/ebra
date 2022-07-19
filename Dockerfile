# Dockerfile
FROM golang:1.18

ENV TZ /usr/share/zoneinfo/Asia/Tokyo

ENV WORKDIR=/app
WORKDIR ${WORKDIR}

ENV GO111MODULE=on

COPY . .
EXPOSE 8080

RUN go get \
  github.com/aws/aws-sdk-go-v2/config \
  github.com/aws/aws-sdk-go-v2/service/sesv2 \
  github.com/golang-jwt/jwt \
  github.com/google/uuid \
  github.com/joho/godotenv \
  github.com/labstack/echo/v4 \
  github.com/stretchr/testify \
  gorm.io/driver/mysql \
  gorm.io/gorm
RUN go install github.com/cosmtrek/air@latest
