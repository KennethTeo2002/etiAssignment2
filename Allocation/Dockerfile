# syntax=docker/dockerfile:1
FROM golang:1.17

# Set destination for COPY
WORKDIR /allocation

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /allocation-docker

EXPOSE 8073

CMD ["/allocation-docker"]