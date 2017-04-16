REPO=malice
NAME=nsrl
BUILD ?= sha1

all: build size test

build:
	cd $(BUILD); docker build -t $(REPO)/$(NAME):$(BUILD) .

size:
	sed -i.bu 's/docker%20image-.*-blue/docker%20image-$(shell docker images --format "{{.Size}}" $(REPO)/$(NAME):$(BUILD)| cut -d' ' -f1)%20MB-blue/' README.md

tags:
	docker images --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}" $(REPO)/$(NAME)

test:
	docker run --rm $(REPO)/$(NAME):$(BUILD) --help
	docker run --rm $(REPO)/$(NAME):$(BUILD) -V lookup 6b82f126555e7644816df5d4e4614677ee0bda5c > results.json
	cat results.json | jq .

.PHONY: build size tags test
