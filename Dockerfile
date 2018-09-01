FROM malice/alpine

LABEL maintainer "https://github.com/blacktop"

ARG HASH=sha1
ENV HASH=$HASH

LABEL malice.plugin.repository = "https://github.com/malice-plugins/nsrl.git"
LABEL malice.plugin.category="intel"
LABEL malice.plugin.mime="hash"
LABEL malice.plugin.docker.engine="*"

COPY shrink_nsrl.sh /nsrl/shrink_nsrl.sh
COPY . /go/src/github.com/malice-plugins/nsrl
RUN apk --update add --no-cache ca-certificates
RUN apk add --no-cache -t .build-deps \
  build-base \
  mercurial \
  musl-dev \
  openssl \
  p7zip \
  unzip \
  bash \
  wget \
  git \
  gcc \
  dep \
  go \
  && set -ex \
  && echo "===> Building info Go binary..." \
  && cd /go/src/github.com/malice-plugins/nsrl \
  && export GOPATH=/go \
  && export PATH=$GOPATH/bin:/usr/local/go/bin:$PATH \
  && export CGO_CFLAGS="-I/usr/local/include" \
  && export CGO_LDFLAGS="-L/usr/local/lib" \
  && go version \
  && dep ensure \
  && go build -ldflags "-s -w -X main.HashType=${HASH} \
  -X main.ErrorRate=$(cat ERROR) \
  -X main.Version=v$(cat VERSION) \
  -X main.BuildTime=$(date -u +%Y%m%d)" -o /bin/nsrl \
  && echo "===> Creating bloomfilter from NSRL database..." \
  && /nsrl/shrink_nsrl.sh \
  && echo "===> Clean up unnecessary files..." \
  && rm -rf /go /usr/local/go /usr/lib/go /tmp/* \
  && apk del --purge .build-deps

VOLUME ["/nsrl"]

WORKDIR /nsrl

ENTRYPOINT ["su-exec","malice","/sbin/tini","--","nsrl"]
CMD ["--help"]
