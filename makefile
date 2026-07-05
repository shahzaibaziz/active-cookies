BINARY    := activecookie
CMD       := ./cmd/activecookie
IMAGE     ?= activecookie
DEV_IMAGE ?= activecookie-dev

DOCKERRUN = docker run --rm -v "$(PWD):/app" -w /app $(DEV_IMAGE)

.PHONY: build test cover run clean format check \
        docker-build docker-push docker-dev

build:
	go build -o $(BINARY) $(CMD)

format: docker-dev
	$(DOCKERRUN) ./scripts/format.sh

check: docker-dev
	$(DOCKERRUN) ./scripts/check.sh

test: docker-dev
	$(DOCKERRUN) ./scripts/test.sh

cover:
	go test ./... -cover

run: build
	./$(BINARY) -f $(FILE) -d $(DATE)

clean:
	./scripts/clean.sh

docker-build:
	docker build -t $(IMAGE) -f Dockerfile .

docker-push: docker-build
	docker push $(IMAGE)

docker-dev:
	docker build -t $(DEV_IMAGE) -f Dockerfile.dev .
