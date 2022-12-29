FROM golang:1.18-alpine

RUN mkdir /app
ADD . /app
WORKDIR /app

RUN go mod download

RUN go build -o /app

EXPOSE 3000

CMD ["/app"]