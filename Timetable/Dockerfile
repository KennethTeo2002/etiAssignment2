# syntax=docker/dockerfile:1
FROM golang:1.17

# Set destination for COPY
WORKDIR /timetable

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
COPY *.json ./

RUN go build -o /timetable-docker

EXPOSE 8073

CMD ["/timetable-docker"]