FROM golang:alpine3.16
WORKDIR /go/app
COPY . .
CMD go run examples/fullapp/cmd/main.go
EXPOSE 8888

