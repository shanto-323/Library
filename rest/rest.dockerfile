FROM golang:1.23-alpine3.20 AS build
RUN apk --no-cache add gcc g++ make ca-certificates
WORKDIR /app
COPY go.mod go.sum ./
COPY vendor vendor
COPY books books
COPY rest rest
RUN go build -o app ./rest

FROM alpine:3.20
WORKDIR /usr/bin
COPY --from=build /app/app .
EXPOSE 8080
CMD ["app"]

