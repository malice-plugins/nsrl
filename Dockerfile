FROM malice/alpine:tini

MAINTAINER blacktop, https://github.com/blacktop

COPY nsrl /nsrl
COPY . /go/src/github.com/maliceio/malice-nsrl
RUN apk-install -t .build-deps \
                    build-base \
                    mercurial \
                    musl-dev \
                    openssl \
                    bash \
                    wget \
                    git \
                    gcc \
                    go \                    
  && cd /tmp \
  && wget https://raw.githubusercontent.com/maliceio/go-plugin-utils/master/scripts/upgrade-alpine-go.sh \
  && chmod +x upgrade-alpine-go.sh \
  && ./upgrade-alpine-go.sh \
  && echo "Building info Go binary..." \
  && cd /go/src/github.com/maliceio/malice-nsrl \
  && export GOPATH=/go \
  && export PATH=$GOPATH/bin:/usr/local/go/bin:$PATH \
  && export CGO_CFLAGS="-I/usr/local/include" \
  && export CGO_LDFLAGS="-L/usr/local/lib" \
  && go version \
  && go get \
  && go build -ldflags "-X main.Version=$(cat VERSION) -X main.BuildTime=$(date -u +%Y%m%d)" -o /bin/nsrl \
  && rm -rf /go /usr/local/go /usr/lib/go /tmp/* \
  && apk del --purge .build-deps

RUN buildDeps='gcc libc-dev python-dev py-pip p7zip' \
  && set -x \
  && apk --update add python $buildDeps \
  && rm -f /var/cache/apk/* \
  && pip install pybloom \
  && /nsrl/shrink_nsrl.sh \
  && apk del --purge $buildDeps \
  && rm -rf /tmp/* /root/.cache /var/cache/apk/* /nsrl/shrink_nsrl.sh

VOLUME ["/malware"]

WORKDIR /malware

ENTRYPOINT ["gosu","malice","/sbin/tini","--","nsrl"]

CMD ["--help"]

# ENTRYPOINT ["/nsrl/search.py"]

# CMD ["-h"]
