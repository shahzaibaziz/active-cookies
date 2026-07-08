BIN_DIR   := bin
BINARY    := $(BIN_DIR)/main
CMD       := ./cmd/activecookies

IMAGE_NAME ?= activecookie
DOCKER_USER ?=
DOCKER_TAG ?= latest
IMAGE     := $(if $(DOCKER_USER),$(DOCKER_USER)/$(IMAGE_NAME):$(DOCKER_TAG),$(IMAGE_NAME):$(DOCKER_TAG))

GOOS      ?= linux
GOARCH    ?= amd64
LDFLAGS   := -s -w

DOCKERRUN = docker run --rm -v "$(PWD):/app" -w /app $(DEV_IMAGE)
DEV_IMAGE ?= activecookie-dev

.PHONY: build build-linux test cover run clean format check \
        docker-build docker-run docker-push docker-dev

build:
	@mkdir -p $(BIN_DIR)
	CGO_ENABLED=0 go build -trimpath -ldflags="$(LDFLAGS)" -o $(BINARY) $(CMD)

build-linux:
	@mkdir -p $(BIN_DIR)
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -trimpath -ldflags="$(LDFLAGS)" -o $(BINARY) $(CMD)

format: docker-dev
	$(DOCKERRUN) ./scripts/format.sh

check: docker-dev
	$(DOCKERRUN) ./scripts/check.sh

test: docker-dev
	$(DOCKERRUN) ./scripts/test.sh

cover: check
	go test ./... -cover

clean:
	./scripts/clean.sh

docker-build: build-linux
	docker build -t $(IMAGE) -f Dockerfile .

docker-run: docker-build
	docker run --rm -v "$(PWD):/data" $(IMAGE) -f /data/$(FILE) -d $(DATE)

docker-push: docker-build
	@if [ -z "$(DOCKER_USER)" ]; then \
		echo "DOCKER_USER is required, e.g. make docker-push DOCKER_USER=myuser" >&2; \
		exit 1; \
	fi
	docker push $(IMAGE)

docker-dev:
	docker build -t $(DEV_IMAGE) -f Dockerfile.dev .
