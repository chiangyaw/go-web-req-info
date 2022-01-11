FROM golang:1-alpine as build

WORKDIR /app
COPY src src
RUN go build src/web-req-info.go

FROM alpine:3.13.2

LABEL owner=prismaclouddemo@gmail.com

WORKDIR /app
COPY --from=build /app/web-req-info /app/web-req-info
COPY startup.sh .

EXPOSE 80
ENTRYPOINT ["sh", "./startup.sh"]
