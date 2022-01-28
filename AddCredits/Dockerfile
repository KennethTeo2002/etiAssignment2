# syntax=docker/dockerfile:1
FROM golang:1.17

# Set destination for COPY
WORKDIR /addcredits

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /addcredits-docker

EXPOSE 8071

CMD ["/addcredits-docker"]