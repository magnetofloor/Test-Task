FROM golang:1.20-bullseye AS build
WORKDIR /app
COPY ./src /app
RUN go build -o app
CMD ["/app/app"]
