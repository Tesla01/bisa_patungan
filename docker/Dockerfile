##
## Build
##
FROM golang:1.18-alpine AS build

WORKDIR /app

COPY go.* ./

RUN apk add git
RUN go mod download

COPY . ./

RUN go build -o /go-api

##
## Deploy
##
FROM alpine

WORKDIR /

COPY --from=build /go-api /go-api

EXPOSE 9000

ENTRYPOINT ["/go-api"]