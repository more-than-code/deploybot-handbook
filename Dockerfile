FROM golang:1.19-alpine3.16 AS build
RUN apk --no-cache add gcc g++ make ca-certificates
WORKDIR /go/src/deploybot

COPY go.mod go.sum ./
COPY . .

RUN go build -o /go/bin/app ./main.go

FROM alpine:3.16
WORKDIR /usr/bin
COPY --from=build /go/bin .
CMD ["app"]