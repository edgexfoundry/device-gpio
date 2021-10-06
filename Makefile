.PHONY: build test clean prepare update docker

GO = CGO_ENABLED=0 GO111MODULE=on go

MICROSERVICES=cmd/device-gpio

.PHONY: $(MICROSERVICES)

DOCKERS=docker_device_gpio_go
.PHONY: $(DOCKERS)

VERSION=$(shell cat ./VERSION 2>/dev/null || echo 0.0.0)
GIT_SHA=$(shell git rev-parse HEAD)
GOFLAGS=-ldflags "-X github.com/edgexfoundry/device-gpio.Version=$(VERSION)"

build: $(MICROSERVICES)

tidy:
	go mod tidy

cmd/device-gpio:
	$(GO) build $(GOFLAGS) -o $@ ./cmd

test:
	$(GO) test ./... -coverprofile=coverage.out
	$(GO) vet ./...
	gofmt -l $$(find . -type f -name '*.go'| grep -v "/vendor/")
	[ "`gofmt -l $$(find . -type f -name '*.go'| grep -v "/vendor/")`" = "" ]
	./bin/test-attribution-txt.sh


clean:
	rm -f $(MICROSERVICES)

docker: $(DOCKERS)

docker_device_gpio_go:
	docker build \
		--label "git_sha=$(GIT_SHA)" \
		--build-arg http_proxy \
		--build-arg https_proxy \
		-t edgexfoundry/device-gpio:$(GIT_SHA) \
		-t edgexfoundry/device-gpio:$(VERSION)-dev \
		.

vendor:
	$(GO) mod vendor
