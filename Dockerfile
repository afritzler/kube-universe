FROM golang:1.18.1 AS builder
WORKDIR /workspace
RUN go install github.com/rakyll/statik@latest

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY cmd/ cmd/
COPY pkg/ pkg/
COPY statik/ statik/
COPY web/ web/
COPY main.go main.go

ARG TARGETOS TARGETARCH

RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH GO111MODULE=on go build -ldflags="-s -w" -a -o kube-universe main.go

FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/kube-universe .
USER 65532:65532
ENTRYPOINT ["./kube-universe"]
