FROM golang:1.19.3-alpine3.16
RUN mkdir post
COPY . /post
WORKDIR /post
RUN go mod vendor
RUN go build -o main cmd/main.go
CMD ./main
EXPOSE 8000