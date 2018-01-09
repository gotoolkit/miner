FROM containerize/dep AS builder

WORKDIR /go/src/github.com/gotoolkit/miner

COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure -vendor-only

COPY . .

RUN go install .

FROM alpine:latest
RUN apk add --no-cache ca-certificates
COPY --from=builder /go/bin/miner /usr/local/bin/miner
ENTRYPOINT ["miner"]
CMD ["-h"]