# http-route

Reverse-proxy HTTP server that simply maps urls to other urls.

Config by JSON:

```json
{
  "http://www.mydomain.com/subpath/": "http://www.example.com/",
  "http://www.mydomain.com/otherpath/": "http://www.example.com/someotherpath/"
}
```

This means a request to `http://www.mydomain.com/subpath/foobar` will route to `http://www.example.com/foobar`, and a request to `http://www.mydomain.com/otherpath/foobar` will route to `http://www.example.com/someotherpath/foobar`. It's that simple.

## How to run

```
docker run -p 80:80 teamjobot/http-route http-route -json '{
  "http://www.mydomain.com/subpath/": "http://www.example.com/",
  "http://www.mydomain.com/otherpath/": "http://www.example.com/someotherpath/"
}'
```

or with docker-compose:

```yaml
version: '3.0'
services:
  http-proxy:
    image: teamjobot/http-route
    command: >
      http-route -json
      '{
        "http://www.mydomain.com/subpath/": "http://www.example.com/",
        "http://www.mydomain.com/otherpath/": "http://www.example.com/someotherpath/"
      }'
```