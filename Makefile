REPO=malice-plugins/nsrl
ORG=malice
NAME=nsrl
VERSION=?= sha1

all: build size test avtest gotest

build:
	cd $(VERSION); docker build -t $(ORG)/$(NAME):$(VERSION) .

base:
	cd $(VERSION); docker build -f Dockerfile.base -t $(ORG)/$(NAME):base .

dev:
	cd $(VERSION); docker build -f Dockerfile.dev -t $(ORG)/$(NAME):$(VERSION) .

size:
	sed -i.bu 's/docker%20image-.*-blue/docker%20image-$(shell docker images --format "{{.Size}}" $(ORG)/$(NAME):$(VERSION)| cut -d' ' -f1)-blue/' README.md

tags:
	docker images --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}" $(ORG)/$(NAME)

ssh:
	@docker run --init -it --rm --entrypoint=bash $(ORG)/$(NAME):$(VERSION)

tar:
	docker save $(ORG)/$(NAME):$(VERSION) -o $(NAME).tar

test:
	docker run --init -d --name elasticsearch -p 9200:9200 blacktop/elasticsearch
	sleep 10; docker run --init --rm $(ORG)/$(NAME):$(VERSION)
	docker run --rm $(REPO)/$(NAME):$(BUILD) -V lookup 6b82f126555e7644816df5d4e4614677ee0bda5c > results.json
	cat docs/results.json | jq .
	http localhost:9200/malice/_search | jq . > docs/elastic.json
	cat docs/elastic.json | jq -r '.hits.hits[] ._source.plugins.av.${NAME}.markdown' > docs/SAMPLE.md
	docker rm -f elasticsearch

circle: ci-size
	@sed -i.bu 's/docker%20image-.*-blue/docker%20image-$(shell cat .circleci/SIZE)-blue/' README.md
	@echo "===> Image size is: $(shell cat .circleci/SIZE)"

ci-build:
	@echo "===> Getting CircleCI build number"
	@http https://circleci.com/api/v1.1/project/github/${REPO} | jq '.[0].build_num' > .circleci/build_num

ci-size: ci-build
	@echo "===> Getting image build size from CircleCI"
	@http "$(shell http https://circleci.com/api/v1.1/project/github/${REPO}/$(shell cat .circleci/build_num)/artifacts${CIRCLE_TOKEN} | jq '.[].url')" > .circleci/SIZE
	@echo "===> Getting NSRL DB size from CircleCI"
	@http "$(shell http https://circleci.com/api/v1.1/project/github/${REPO}/$(shell cat .circleci/build_num)/artifacts${CIRCLE_TOKEN} | jq '.[].url')" > .circleci/SIZE
	@echo "===> Getting bloomfilter size from CircleCI"
	@http "$(shell http https://circleci.com/api/v1.1/project/github/${REPO}/$(shell cat .circleci/build_num)/artifacts${CIRCLE_TOKEN} | jq '.[].url')" > .circleci/SIZE

clean:
	docker-clean stop
	docker rmi $(ORG)/$(NAME):$(VERSION)
	docker rmi $(ORG)/$(NAME):base

.PHONY: build dev size tags test gotest clean circle
