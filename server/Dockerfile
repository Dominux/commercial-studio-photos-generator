FROM golang:1.22rc1-bookworm AS build
RUN mkdir /app 
ADD . /app/ 
WORKDIR /app 
RUN go mod download
RUN go build -o /main ./cmd/main.go
CMD ["/main"]
