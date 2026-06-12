include ../../Makefile

MAKEFILE_DIR := $(dir $(abspath $(firstword $(MAKEFILE_LIST))))

MOCKERY_SOURCE := github.com/vektra/mockery/v3@v3.7.0

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: tidy
lint:
	go vet ./...
	go fmt ./...

.PHONY: unit-test
unit-test:
	go test ./...

.PHONY: build
build:
	../scripts/build.sh

.PHONY: mocks
mocks:
	@go install "$(MOCKERY_SOURCE)" || exit 1; \
	\
	cd "$(MAKEFILE_DIR)"; \
	"$$(go env GOPATH)/bin/mockery" --config .mockery.yml
