FROM golang:alpine
WORKDIR /go/src/github.com/gotoolkit/miner
RUN apk add --no-cache git
COPY . /go/src/github.com/gotoolkit/miner
RUN go build -o miner .

FROM alpine:latest
RUN apk add --no-cache ca-certificates
COPY --from=0 /go/src/github.com/gotoolkit/miner/miner /usr/local/bin/miner
ENTRYPOINT ["miner"]
CMD ["-h"]