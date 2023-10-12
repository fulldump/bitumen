PROJECT=bitumen

COMMIT = $(shell git log -1 --format="%h" 2>/dev/null || echo "0")
VERSION = $(shell git describe --tags --always)
BUILD_DATE = $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
COMPILER = $(shell go version)
FLAGS = -ldflags "\
  -X $(PROJECT)/server.COMMIT=$(COMMIT) \
  -X $(PROJECT)/server.VERSION=$(VERSION) \
  -X $(PROJECT)/server.BUILD_DATE=$(BUILD_DATE) \
  -X '$(PROJECT)/server.COMPILER=$(COMPILER)' \
  "

GO = GOFLAGS=-mod=vendor go

PORTS ?= -p 8080:8080


.PHONY: build
build:
	$(GO) build $(FLAGS) -mod vendor -o bin/$(PROJECT) ./cmd/bitumen
	$(GO) build $(FLAGS) -mod vendor -o bin/sftpexample ./cmd/sftpexample

.PHONY: run
run:
	$(GO) run $(FLAGS) ./cmd/bitumen

.PHONY: apitree
apitree:
	APITREE=1 $(GO) run $(FLAGS) ./main.go

.PHONY: deps
deps:
	go mod tidy
	go mod download
	go mod vendor

.PHONY: test
test:
	$(GO) test $(PROJECT)/... --cover

.PHONY: cloc
cloc:
	cloc --exclude-dir=vendor,lib .

.PHONY: version
version:
	@echo $(VERSION)

.PHONY: docker-%
docker-%:
	docker-compose run $(PORTS) --use-aliases --rm app make $*

