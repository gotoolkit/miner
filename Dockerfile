FROM golang:alpine
WORKDIR /go/src/github.com/gotoolkit/miner
RUN apk add --no-cache git && git clone https://github.com/gotoolkit/miner.git .
RUN go build -o miner .

FROM alpine:latest
COPY --from=0 /go/src/github.com/gotoolkit/miner/miner /usr/local/bin/miner
ENTRYPOINT ["miner"]
CMD ["-h"]