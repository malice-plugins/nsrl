# To write results to [ElasticSearch](https://www.elastic.co/products/elasticsearch)

## Write to a `elasticsearch` docker container

```bash
$ docker volume create --name malice
$ docker run -d --name elastic \
                -p 9200:9200 \
                -v malice:/usr/share/elasticsearch/data \
                 blacktop/elasticsearch
$ docker run --rm --link elastic malice/nsrl HASH
```

## Write to an external `elasticsearch` database

```bash
$ docker run --rm \
             -e MALICE_ELASTICSEARCH_URL=$MALICE_ELASTICSEARCH_URL \
             -e MALICE_ELASTICSEARCH_USERNAME=$MALICE_ELASTICSEARCH_USERNAME \
             -e MALICE_ELASTICSEARCH_PASSORD=$MALICE_ELASTICSEARCH_PASSORD \
              malice/nsrl -V lookup SHA1_HASH
```
