FROM golang:1-alpine as build

WORKDIR /app
COPY src src
RUN go build src/web-req-info.go

FROM alpine:3.11.3

WORKDIR /app
COPY --from=build /app/web-req-info /app/web-req-info

EXPOSE 80
ENTRYPOINT ["./web-req-info"]
