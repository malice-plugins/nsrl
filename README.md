![NSRL logo](https://raw.githubusercontent.com/malice-plugins/nsrl/master/docs/logo.png)

# malice-nsrl

[![Circle CI](https://circleci.com/gh/malice-plugins/nsrl.png?style=shield)](https://circleci.com/gh/malice-plugins/nsrl) [![License](http://img.shields.io/:license-mit-blue.svg)](http://doge.mit-license.org) [![Docker Stars](https://img.shields.io/docker/stars/malice/nsrl.svg)](https://hub.docker.com/r/malice/nsrl/) [![Docker Pulls](https://img.shields.io/docker/pulls/malice/nsrl.svg)](https://hub.docker.com/r/malice/nsrl/) [![Docker Image](https://img.shields.io/badge/docker%20image-117MB-blue.svg)](https://hub.docker.com/r/malice/nsrl/)

Malice NSRL Plugin - This takes the **5.5 GB** [NSRL](http://www.nsrl.nist.gov/Downloads.htm) minimal set and converts it into a **77.4 MB** [bloom filter](https://en.wikipedia.org/wiki/Bloom_filter) with an Estimate False Positive Rate of `0.001`

This repository contains a **Dockerfile** of the [NSRL](http://www.nsrl.nist.gov) lookup malice plugin **malice/nsrl**.

### Dependencies

- [malice/alpine](https://hub.docker.com/r/malice/alpine/)

## Image Tags

```
REPOSITORY          TAG                 SIZE
malice/nsrl         latest              51.9MB
malice/nsrl         0.1.0               51.9MB
malice/nsrl         sha1                51.3MB
malice/nsrl         md5                 51.3MB
```

> **NOTE:** tag `sha1` can query by sha1 hash and tag `md5` can query by md5 hash

## Installation

1. Install [Docker](https://www.docker.io/).
2. Download [trusted build](https://hub.docker.com/r/malice/nsrl/) from public [DockerHub](https://hub.docker.com): `docker pull malice/nsrl`

## Usage

```bash
docker run --rm malice/nsrl --help

Usage: nsrl [OPTIONS] COMMAND [arg...]

Malice nsrl Plugin

Version: v0.1.0, BuildTime: 20161119

Author:
  blacktop - <https://github.com/blacktop>

Options:
  --verbose, -V  verbose output
  --help, -h     show help
  --version, -v  print the version

Commands:
  web     Create a NSRL lookup web service
  build   Build bloomfilter from NSRL database
  lookup  Query NSRL for hash
  help    Shows a list of commands or help for one command

Run 'nsrl COMMAND --help' for more information on a command.
```

### Lookup By Hash `md5|sha1`

```
docker run --rm malice/nsrl:md5 lookup MD5
docker run --rm malice/nsrl:sha1 lookup SHA1
```

```
NAME:
   nsrl lookup - Query NSRL for hash

USAGE:
   nsrl lookup [command options] SHA1 to query NSRL with

OPTIONS:
   --elasticsearch value  elasticsearch url for Malice to store results [$MALICE_ELASTICSEARCH_URL]
   --post, -p             POST results to Malice webhook [$MALICE_ENDPOINT]
   --proxy, -x            proxy settings for Malice webhook endpoint [$MALICE_PROXY]
   --timeout value        malice plugin timeout (in seconds) (default: 10) [$MALICE_TIMEOUT]
   --table, -t            output as Markdown table
```

## Sample Output

### [JSON](https://github.com/malice-plugins/nsrl/blob/master/docs/results.json)

---

```json
{
  "nsrl": {
    "found": true,
    "hash": "5A272B7441328E09704B6D7EABDBD51B8858FDE4"
  }
}
```

---

### [Markdown](https://github.com/malice-plugins/nsrl/blob/master/docs/SAMPLE.md)

---

#### NSRL Database

- Found :white_check_mark:

---

## Documentation

- [To write results to ElasticSearch](https://github.com/malice-plugins/nsrl/blob/master/docs/elasticsearch.md)
- [To create a nsrl lookup micro-service](https://github.com/malice-plugins/nsrl/blob/master/docs/web.md)
- [To post results to a webhook](https://github.com/malice-plugins/nsrl/blob/master/docs/callback.md)

## Issues

Find a bug? Want more features? Find something missing in the documentation? Let me know! Please don't hesitate to [file an issue](https://github.com/malice-plugins/nsrl/issues/new)

## CHANGELOG

See [`CHANGELOG.md`](https://github.com/malice-plugins/nsrl/blob/master/CHANGELOG.md)

## Contributing

[See all contributors on GitHub](https://github.com/malice-plugins/nsrl/graphs/contributors).

Please update the [CHANGELOG.md](https://github.com/malice-plugins/nsrl/blob/master/CHANGELOG.md) and submit a [Pull Request on GitHub](https://help.github.com/articles/using-pull-requests/).

## License

MIT Copyright (c) 2015 **blacktop**
