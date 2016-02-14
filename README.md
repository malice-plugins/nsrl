![NSRL logo](https://raw.githubusercontent.com/maliceio/malice-nsrl/master/logo.png)
# malice-nsrl

[![License](http://img.shields.io/:license-mit-blue.svg)](http://doge.mit-license.org)
[![Docker Stars](https://img.shields.io/docker/stars/malice/nsrl.svg)][hub]
[![Docker Pulls](https://img.shields.io/docker/pulls/malice/nsrl.svg)][hub]
[![Image Size](https://img.shields.io/imagelayers/image-size/malice/nsrl/latest.svg)](https://imagelayers.io/?images=malice/nsrl:latest)
[![Image Layers](https://img.shields.io/imagelayers/layers/malice/nsrl/latest.svg)](https://imagelayers.io/?images=malice/nsrl:latest)

Malice NSRL Plugin - This takes the **5.43GB** [NSRL](http://www.nsrl.nist.gov/Downloads.htm) minimal set and converts it into a **96M** [bloom filter](https://en.wikipedia.org/wiki/Bloom_filter).

This repository contains a **Dockerfile** of the **Malice NSRL Plugin** for [Docker](https://www.docker.io/)'s [trusted build](https://index.docker.io/u/malice/nsrl/) published to the public [Docker Registry](https://index.docker.io/).

### Dependencies
* [alpine](https://registry.hub.docker.com/_/alpine/)

### Image Tags
```bash
$ docker images

REPOSITORY          TAG                 VIRTUAL SIZE
malice/nsrl       latest              142 MB
malice/nsrl       sha1                142 MB
malice/nsrl       name                142 MB
malice/nsrl       error_0.001         192 MB
```
> **NOTE:** There are **3** other versions of this image:
 - [sha1](https://github.com/malice/docker-nsrl/tree/sha1) tag allows you to search the NSRL DB by **sha-1** hash
 - [name](https://github.com/malice/docker-nsrl/tree/name) tag allows you to search the NSRL DB by **filename**
 - [error_0.001](https://github.com/malice/docker-nsrl/tree/error_0.001) tag searches by **md5** hash and has a much **lower error_rate** threshold. It does, however, grow the size of the bloomfilter by 50MB.

### Installation

1. Install [Docker](https://www.docker.io/).

2. Download [trusted build](https://index.docker.io/u/malice/nsrl/) from public [Docker Registry](https://index.docker.io/): `docker pull malice/nsrl`

#### Alternatively, build an image from Dockerfile
`docker build -t malice/nsrl github.com/maliceio/docker-nsrl`

### Usage
```bash
$ docker run --rm malice/nsrl
```
#### Output:

    usage: malice/nsrl [-h] [-v] MD5 [MD5 ...]

    positional arguments:
      MD5            a md5 hash to search for.

    optional arguments:
      -h, --help     show this help message and exit
      -v, --verbose  Display verbose output message

#### Example (with `-v` option):
```bash
$ docker run --rm malice/nsrl -v 60B7C0FEAD45F2066E5B805A91F4F0FC
```
#### Output:
```bash
Hash 60B7C0FEAD45F2066E5B805A91F4F0FC found in NSRL Database.
```

#### To read from a **hash-list** file:
```bash
$ cat hash-list.txt
60B7C0FEAD45F2066E5B805A91F4F0FC
AABCA0896728846A9D5B841617EBE746
AABCA0896728846A9D5B841617EBE745

$ cat hash-list.txt | xargs docker run --rm malice/nsrl
True
True
False
```

### To Run on OSX
 - Install [Homebrew](http://brew.sh)

```bash
$ brew install caskroom/cask/brew-cask
$ brew cask install virtualbox
$ brew install docker
$ brew install docker-machine
$ docker-machine create --driver virtualbox malice
$ eval $(docker-machine env malice)
```
Add the following to your bash or zsh profile

```bash
alias nsrl='docker run --rm malice/nsrl $@'
```
#### Usage
```bash
$ nsrl -v 60B7C0FEAD45F2066E5B805A91F4F0FC AABCA0896728846A9D5B841617EBE746
```

### Optional Build Options
You can use different **NSRL databases** or **error-rates** for the bloomfilter (*which will increase it's accuracy*)

1. To use your own [NSRL](http://www.nsrl.nist.gov/Downloads.htm) database simply download the ZIP and place it in the `nsrl` folder and build the image like so: `docker build -t my_nsrl .`
2. To decrease the error-rate of the bloomfilter simply change the value of `ERROR_RATE` in the file `nsrl/shrink_nsrl.sh` and build as above.

### Issues
Find a bug? Want more features? Find something missing in the documentation? Let me know! Please don't hesitate to [file an issue](https://github.com/maliceio/malice-nsrl/issues/new) and I'll get right on it.

### Credits
Inspired by https://github.com/bigsnarfdude/Malware-Probabilistic-Data-Structres

### License
MIT Copyright (c) 2016 **blacktop**

[hub]: https://hub.docker.com/r/malice/nsrl/
