
FROM golang:rc-alpine3.13

RUN mkdir -p /go/src/pool
WORKDIR /go/src/pool
COPY go.mod /go/src/pool/
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go mod download
COPY . /go/src/pool/
RUN go build -o server.bin main.go

EXPOSE 20150
CMD /go/src/pool/server.bin