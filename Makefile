.PHONY: build test unittest lint clean prepare update docker

MICROSERVICES=cmd/device-gpio

.PHONY: $(MICROSERVICES)

ARCH=$(shell uname -m)

DOCKERS=docker_device_gpio_go
.PHONY: $(DOCKERS)

VERSION=$(shell cat ./VERSION 2>/dev/null || echo 0.0.0)
GIT_SHA=$(shell git rev-parse HEAD)

SDKVERSION=$(shell cat ./go.mod | grep 'github.com/edgexfoundry/device-sdk-go/v4 v' | awk '{print $$2}')
GOFLAGS=-ldflags "-X github.com/edgexfoundry/device-gpio.Version=$(VERSION) \
                  -X github.com/edgexfoundry/device-sdk-go/v4/internal/common.SDKVersion=$(SDKVERSION)" \
                   -trimpath -mod=readonly

build: $(MICROSERVICES)

build-nats:
	make -e ADD_BUILD_TAGS=include_nats_messaging build

tidy:
	go mod tidy

cmd/device-gpio:
	CGO_ENABLED=0 go build -tags "$(ADD_BUILD_TAGS)" $(GOFLAGS) -o $@ ./cmd

unittest:
	go test ./... -coverprofile=coverage.out ./...

lint:
	@which golangci-lint >/dev/null || echo "WARNING: go linter not installed. To install, run make install-lint"
	@if [ "z${ARCH}" = "zx86_64" ] && which golangci-lint >/dev/null ; then golangci-lint run --config .golangci.yml ; else echo "WARNING: Linting skipped (not on x86_64 or linter not installed)"; fi

install-lint:
	sudo curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.61.0

test: unittest lint
	go vet ./...
	gofmt -l $$(find . -type f -name '*.go'| grep -v "/vendor/")
	[ "`gofmt -l $$(find . -type f -name '*.go'| grep -v "/vendor/")`" = "" ]
	./bin/test-attribution-txt.sh

clean:
	rm -f $(MICROSERVICES)

docker: $(DOCKERS)

docker_device_gpio_go:
	docker build \
		--build-arg ADD_BUILD_TAGS=$(ADD_BUILD_TAGS) \
		--label "git_sha=$(GIT_SHA)" \
		-t edgexfoundry/device-gpio:$(GIT_SHA) \
		-t edgexfoundry/device-gpio:$(VERSION)-dev \
		.

docker-nats:
	make -e ADD_BUILD_TAGS=include_nats_messaging docker

vendor:
	CGO_ENABLED=0 go mod vendor
