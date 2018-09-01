REPO=malice-plugins/nsrl
ORG=malice
NAME=nsrl
CATEGORY=intel
VERSION?=sha1

FOUND_HASH=5a272b7441328e09704b6d7eabdbd51b8858fde4
MISSING_HASH=6b82f126555e7644816df5d4e4614677ee0bdacc

all: build size tag test_all

.PHONY: build
build:
	docker build -t $(ORG)/$(NAME):$(VERSION) .

.PHONY: build.md5
build.md5:
	docker build --build-arg HASH=md5 -t $(ORG)/$(NAME):md5 .

.PHONY: size
size:
	sed -i.bu 's/docker%20image-.*-blue/docker%20image-$(shell docker images --format "{{.Size}}" $(ORG)/$(NAME):$(VERSION)| cut -d' ' -f1)-blue/' README.md

.PHONY: tag
tag:
	docker tag $(ORG)/$(NAME):$(VERSION) $(ORG)/$(NAME):latest

.PHONY: tags
tags:
	docker images --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}" $(ORG)/$(NAME)

.PHONY: ssh
ssh:
	@docker run --init -it --rm --entrypoint=bash $(ORG)/$(NAME):$(VERSION)

.PHONY: tar
tar:
	docker save $(ORG)/$(NAME):$(VERSION) -o $(NAME).tar

.PHONY: start_elasticsearch
start_elasticsearch:
ifeq ("$(shell docker inspect -f {{.State.Running}} elasticsearch)", "true")
	@echo "===> elasticsearch already running.  Stopping now..."
	@docker rm -f elasticsearch || true
endif
	@echo "===> Starting elasticsearch"
	@docker run --init -d --name elasticsearch -p 9200:9200 malice/elasticsearch:6.3; sleep 15

.PHONY: test_all
test_all: test test_elastic test_markdown test_web

.PHONY: test
test:
	@echo "===> ${NAME} --help"
	@docker run --rm $(ORG)/$(NAME):$(VERSION)
	@echo "===> ${NAME} test"
	docker run --rm $(ORG)/$(NAME):$(VERSION) -V lookup $(FOUND_HASH) > docs/results.json
	cat docs/results.json | jq .
	@echo "===> Test lookup NOT found"
	@docker run --rm $(ORG)/$(NAME):$(VERSION) -V lookup $(MISSING_HASH) | jq . > docs/no_results.json
	cat docs/no_results.json | jq .

.PHONY: test_elastic
test_elastic: start_elasticsearch
	@echo "===> ${NAME} test_elastic"
	docker run --rm --link elasticsearch -e MALICE_ELASTICSEARCH_URL=http://elasticsearch:9200 $(ORG)/$(NAME):$(VERSION) -V lookup $(FOUND_HASH)
	@echo "===> ${NAME} test_elastic NOT found"
	docker run --rm --link elasticsearch -e MALICE_ELASTICSEARCH_URL=http://elasticsearch:9200 $(ORG)/$(NAME):$(VERSION) -V lookup $(MISSING_HASH)
	http localhost:9200/malice/_search | jq . > docs/elastic.json

.PHONY: test_elastic_remote
test_elastic_remote:
	@echo "===> ${NAME} test_elastic"
	docker run --rm \
	-e MALICE_ELASTICSEARCH_URL=${MALICE_ELASTICSEARCH_URL} \
	-e MALICE_ELASTICSEARCH_USERNAME=${MALICE_ELASTICSEARCH_USERNAME} \
	-e MALICE_ELASTICSEARCH_PASSWORD=${MALICE_ELASTICSEARCH_PASSWORD} \
	-e MALICE_ELASTICSEARCH_INDEX="test" \
	$(ORG)/$(NAME):$(VERSION) -V lookup $(FOUND_HASH)

.PHONY: test_markdown
test_markdown:
	@echo "===> ${NAME} test_markdown"
	http localhost:9200/malice/_search | jq . > docs/elastic.json
	cat docs/elastic.json | jq -r '.hits.hits[] ._source.plugins.${CATEGORY}.${NAME}.markdown' > docs/SAMPLE.md

.PHONY: test_web
test_web: stop
	@echo "===> ${NAME} web service"
	@docker run --init -d --name $(NAME) -p 3993:3993 $(ORG)/$(NAME):$(VERSION) -V web
	http -f localhost:3993/lookup/$(FOUND_HASH)
	http -f localhost:3993/lookup/$(MISSING_HASH)

.PHONY: stop
stop:
	@echo "===> Stopping container ${NAME}"
	@docker container rm -f $(NAME) || true

.PHONY: circle
circle: ci-size
	@sed -i.bu 's/docker%20image-.*-blue/docker%20image-$(shell cat .circleci/size)-blue/' README.md
	@echo "===> Image size is: $(shell cat .circleci/size)"

ci-build:
	@echo "===> Getting CircleCI build number"
	@http https://circleci.com/api/v1.1/project/github/${REPO} | jq '.[0].build_num' > .circleci/build_num

ci-size: ci-build
	@echo "===> Getting artifact sizes from CircleCI"
	@cd .circleci; rm size nsrl bloom || true
	@http https://circleci.com/api/v1.1/project/github/${REPO}/$(shell cat .circleci/build_num)/artifacts${CIRCLE_TOKEN} | jq -r ".[] | .url" | xargs wget -q -P .circleci

clean:
	docker-clean stop
	docker rmi $(ORG)/$(NAME):$(VERSION)
	docker rmi $(ORG)/$(NAME):base

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := all
