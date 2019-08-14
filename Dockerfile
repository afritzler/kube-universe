FROM golang:1.12.8
WORKDIR /go/src/github.com/afritzler/kube-universe
RUN go get github.com/rakyll/statik
COPY . .
RUN make

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /go/src/github.com/afritzler/kube-universe/kube-universe .
CMD ["./kube-universe"]