GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=kube-universe
BINARY_UNIX=$(BINARY_NAME)_linux
IMAGE=afritzler/kube-universe
TAG=latest

all: test build
build: deps
		statik -f -src=$(PWD)/web/
		$(GOBUILD) -o $(BINARY_NAME) -v
test:
		$(GOTEST) -v ./...
clean:
		$(GOCLEAN)
		rm -f $(BINARY_NAME)
		rm -f $(BINARY_UNIX)
run:
		$(GOBUILD) -o $(BINARY_NAME) -v ./...
		./$(BINARY_NAME)
deps:
		$(GOGET) -d github.com/rakyll/statik
fmt: ## Run go fmt against code.
		go fmt ./...
vet: ## Run go vet against code.
		go vet ./...
lint:
		golangci-lint run ./...
# Cross compilation
build-linux:
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v
docker-build: build
		docker build -t $(IMAGE):$(TAG) .
