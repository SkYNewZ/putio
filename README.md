# Put.io HTTP interface

![putio](https://put.io/images/nav-logo-black.png)

This tool is an index directory listing as Nginx format for your Put.io files. Its supported navigation like Nginx too.
For each file found, a download/streaming link is generated.

![index](img/index.png)

My instance is hosted on Google Cloud Run and I store users in Google Cloud Firestore.

## Features

- Use a Put.io OAuth token. You must set `PUT_IO_TOKEN=XXX`
- Ofuscation URI for a minimum of privacy. You must set `OFUSCATION_TOKEN=XXXX`
- Cache response (enable by default, can be disable with `NO_CACHE=1`)
- HTTP Basic auth (enable by default, can be disable with `NO_AUTH=1`)
- Store encrypted users passwords in Google Cloud Firestore (need `GOOGLE_CLOUD_PROJECT` of your Firestore project in authentication is enable)
- Support Google Stackdriver logging format using https://github.com/joonix/log

## Deploy with Docker

```bash
$ docker run -itd --name putio -e PUT_IO_TOKEN=XXX -e NO_AUTH=1 -e OFUSCATION_TOKEN=XXXX -p 127.0.0.1:8080:8080 skynewz/putio
```

Homepage URL is `http://localhost:8080/<OFUSCATION_TOKEN>/0`. `0` is the root folder ID.
