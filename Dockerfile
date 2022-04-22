# Build the manager binary
FROM --platform=$BUILDPLATFORM golang:1.18 as builder

ARG GOARCH=''
ARG GITHUB_PAT=''

WORKDIR /workspace
RUN go install github.com/rakyll/statik@latest

# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum

# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    go mod download

# Copy the go source
COPY cmd/ cmd/
COPY pkg/ pkg/
COPY statik/ statik/
COPY web/ web/
COPY main.go main.go

ARG TARGETOS
ARG TARGETARCH

# Build
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH GO111MODULE=on go build -ldflags="-s -w" -a -o kube-universe main.go

FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/kube-universe .
USER 65532:65532

ENTRYPOINT ["./kube-universe"]
