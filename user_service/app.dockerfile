FROM golang:1.23-alpine3.20 AS build
RUN apk --no-cache add gcc g++ make ca-certificates
WORKDIR /app
COPY go.mod go.sum ./
COPY user_service user_service
COPY books books
RUN go build -o app ./user_service/cmd

FROM alpine:3.20
WORKDIR /usr/bin
COPY --from=build /app/app .
EXPOSE 8080
CMD ["app"]