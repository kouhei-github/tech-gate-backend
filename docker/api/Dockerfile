FROM golang:1.19-alpine3.15

RUN mkdir /app
WORKDIR /app

RUN go install github.com/cosmtrek/air@v1.29.0

ADD ./src/go.mod /app
RUN go mod download
RUN go mod tidy
