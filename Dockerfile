FROM golang:1-alpine as build

WORKDIR /app
COPY src src
RUN go build src/web-req-info.go

FROM alpine:latest

WORKDIR /app
COPY --from=build /app/web-req-info /app/web-req-info

EXPOSE 80
ENTRYPOINT ["./web-req-info"]
