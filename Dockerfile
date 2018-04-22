FROM alpine:3.7
ENTRYPOINT ["/bin/kube-universe"]

COPY . /go/src/github.com/afritzler/kube-universe
RUN apk --no-cache add -t build-deps build-base go git \
	&& apk --no-cache add ca-certificates \
	&& cd /go/src/github.com/afritzler/kube-universe \
	&& export GOPATH=/go \
	&& go build -o /bin/kube-universe \
	&& rm -rf /go \
	&& apk del --purge build-deps