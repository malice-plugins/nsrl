# To Query By Filename

I personally don't see a point in querying by filename as I am primarly focused on identifying malware.

However, there might be reasons not apparent to me so this will let you query by any of the fields you want in the NSRL database.

## Getting Started

You will need to build the docker image:

```bash
$ git clone https://github.com/malice-plugins/nsrl.git
$ cd nsrl
$ docker build --build-arg HASH=filename -t malice/nsrl:filename .
```

This will take a while depending on your internet speed.

Once completed you can now do:

```bash
$ docker run --rm malice/nsrl:filename lookup calc.exe
```
