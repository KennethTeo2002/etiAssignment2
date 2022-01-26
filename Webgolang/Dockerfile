# syntax=docker/dockerfile:1
FROM golang:1.16-alpine

# Set destination for COPY
WORKDIR /web

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY main.go ./
COPY ./Pages ./Pages

RUN go build -o /webapp-docker

EXPOSE 8070

CMD ["/webapp-docker"]