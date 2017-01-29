# Create a NSRL lookup micro-service :new:

```bash
$ docker run malice/nsrl web

INFO[0000] web service listening on port :3993
```

Now you can perform queries like so

```bash
$ http localhost:3993/lookup/60B7C0FEAD45F2066E5B805A91F4F0FC
```

```bash
HTTP/1.1 200 OK
Content-Length: 24
Content-Type: application/json; charset=UTF-8
Date: Sun, 20 Nov 2016 21:43:30 GMT

{
    "nsrl": {
        "found": true
    }
}
```
