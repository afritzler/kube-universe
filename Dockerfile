FROM golang:1.17.5 AS builder
WORKDIR /go/src/github.com/afritzler/kube-universe
RUN go get github.com/rakyll/statik
COPY . .
RUN make

FROM alpine:3.15.0
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/afritzler/kube-universe/kube-universe .
CMD ["./kube-universe"]