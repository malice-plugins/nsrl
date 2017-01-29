![NSRL logo](https://raw.githubusercontent.com/maliceio/malice-nsrl/master/logo.png)

malice-nsrl
===========

[![Circle CI](https://circleci.com/gh/maliceio/malice-nsrl.png?style=shield)](https://circleci.com/gh/maliceio/malice-nsrl) [![License](http://img.shields.io/:license-mit-blue.svg)](http://doge.mit-license.org) [![Docker Stars](https://img.shields.io/docker/stars/malice/nsrl.svg)](https://hub.docker.com/r/malice/nsrl/) [![Docker Pulls](https://img.shields.io/docker/pulls/malice/nsrl.svg)](https://hub.docker.com/r/malice/nsrl/) [![Docker Image](https://img.shields.io/badge/docker image-104 MB-blue.svg)](https://hub.docker.com/r/malice/nsrl/)

Malice NSRL Plugin - This takes the **6.49GB** [NSRL](http://www.nsrl.nist.gov/Downloads.htm) minimal set and converts it into a **86M** [bloom filter](https://en.wikipedia.org/wiki/Bloom_filter) with an Estimate False Positive Rate of `0.001`

This repository contains a **Dockerfile** of the [NSRL](http://www.nsrl.nist.gov) lookup malice plugin **malice/nsrl**.

### Dependencies

-	[malice/alpine:tini](https://hub.docker.com/r/malice/alpine/)

### Installation

1.	Install [Docker](https://www.docker.io/).
2.	Download [trusted build](https://hub.docker.com/r/malice/nsrl/) from public [DockerHub](https://hub.docker.com): `docker pull malice/nsrl`

### Usage

```
docker run --rm malice/nsrl:md5 lookup MD5
docker run --rm malice/nsrl:sha1 lookup SHA1
```

```bash
Usage: nsrl [OPTIONS] COMMAND [arg...]

Malice nsrl Plugin

Version: v0.1.0, BuildTime: 20161119

Author:
  blacktop - <https://github.com/blacktop>

Options:
  --verbose, -V		verbose output
  --post, -p		POST results to Malice webhook [$MALICE_ENDPOINT]
  --proxy, -x		proxy settings for Malice webhook endpoint [$MALICE_PROXY]
  --table, -t		output as Markdown table
  --timeout value       malice plugin timeout (in seconds) (default: 10) [$MALICE_TIMEOUT]    
  --elasitcsearch value	elasitcsearch address for Malice to store results [$MALICE_ELASTICSEARCH]
  --help, -h		show help
  --version, -v		print the version

Commands:
  web		Create a NSRL lookup web service
  build		Build bloomfilter from NSRL database
  lookup	Query NSRL for hash
  help		Shows a list of commands or help for one command

Run 'nsrl COMMAND --help' for more information on a command.
```

Sample Output
-------------

### JSON:

---

```json
{
  "nsrl": {
    "found": true
  }
}
```

---

### Markdown Table:

---

#### NSRL Database

-	Found :white_check_mark:

---

Documentation
-------------

-	[To write results to ElasticSearch](https://github.com/maliceio/malice-nsrl/blob/master/docs/elasticsearch.md)
-	[To create a nsrl lookup micro-service](https://github.com/maliceio/malice-nsrl/blob/master/docs/web.md)
-	[To post results to a webhook](https://github.com/maliceio/malice-nsrl/blob/master/docs/callback.md)

### Issues

Find a bug? Want more features? Find something missing in the documentation? Let me know! Please don't hesitate to [file an issue](https://github.com/maliceio/malice-nsrl/issues/new)

### CHANGELOG

See [`CHANGELOG.md`](https://github.com/maliceio/malice-nsrl/blob/master/CHANGELOG.md)

### Contributing

[See all contributors on GitHub](https://github.com/maliceio/malice-nsrl/graphs/contributors).

Please update the [CHANGELOG.md](https://github.com/maliceio/malice-nsrl/blob/master/CHANGELOG.md) and submit a [Pull Request on GitHub](https://help.github.com/articles/using-pull-requests/).

### License

MIT Copyright (c) 2016 **blacktop**
