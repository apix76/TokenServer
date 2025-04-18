FROM golang:1.24-alpine

WORKDIR /TokenServer
COPY . .

RUN go build ./

CMD ./TokenServer