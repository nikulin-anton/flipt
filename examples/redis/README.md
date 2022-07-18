# Redis Example

This example shows how you can run Flipt with a Redis cache.

This works by setting the following (or similar) configuration options:

```yaml
cache:
  enabled: true
  backend: redis
  ttl: 60s
  redis:
    host: redis
    port: 6379
```

## Requirements

To run this example application you'll need:

* [Docker](https://docs.docker.com/install/)
* [docker-compose](https://docs.docker.com/compose/install/)

## Running the Example

1. Run `docker-compose up` from this directory
1. Open the Flipt UI (default: [http://localhost:8080](http://localhost:8080))
