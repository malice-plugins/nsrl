REPO=malice
NAME=nsrl
BUILD ?= sha1

all: build size test

build:
	cd $(BUILD); docker build -t $(REPO)/$(NAME):$(BUILD) .

size:
	sed -i.bu 's/docker image-.*-blue/docker image-$(shell docker images --format "{{.Size}}" $(REPO)/$(NAME):$(BUILD))-blue/' README.md

tags:
	docker images --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}" $(REPO)/$(NAME)

test:
	docker run --init --rm $(REPO)/$(NAME):$(BUILD) -V 60B7C0FEAD45F2066E5B805A91F4F0FC > results.json
	cat results.json | jq .

.PHONY: build size tags test
