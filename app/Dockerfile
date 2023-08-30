FROM golang:1.20 AS build
WORKDIR /app
COPY ./src /app
RUN ls /app
RUN go mod download
RUN go build -o app
RUN ls /app

FROM debian:buster-slim
WORKDIR /
COPY --from=build /app/app .
RUN ls /
CMD ["/app"]