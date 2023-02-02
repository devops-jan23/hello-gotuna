FROM golang:bullseye
WORKDIR /go/app
COPY . .
CMD go run examples/fullapp/cmd/main.go
EXPOSE 8888

