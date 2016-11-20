FROM malice/alpine:tini

MAINTAINER blacktop, https://github.com/blacktop

COPY shrink_nsrl.sh /nsrl/shrink_nsrl.sh
COPY . /go/src/github.com/maliceio/malice-nsrl
RUN apk-install -t .build-deps \
                    build-base \
                    mercurial \
                    musl-dev \
                    openssl \
                    p7zip \
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
  && go build -ldflags "-X main.Version=$(cat VERSION) -X main.ErrorRate=$(cat ERROR) -X main.BuildTime=$(date -u +%Y%m%d)" -o /bin/nsrl \
  && echo "[INFO] Creating bloomfilter from NSRL database..." \ 
  && /nsrl/shrink_nsrl.sh \
  && rm -rf /go /usr/local/go /usr/lib/go /tmp/* \
  && apk del --purge .build-deps

VOLUME ["/malware"]

WORKDIR /malware

ENTRYPOINT ["gosu","malice","/sbin/tini","--","nsrl"]

CMD ["--help"]
