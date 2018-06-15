## cryptoexchange-dashboard

Project info TBD

### How to run

- copy `docker/env.template` to `docker/env` and change it especially in the section commented by `###change me`
- add your read only Bittrex API keys `EXCHANGE_API_KEY`, `EXCHANGE_API_SECRET` to `env` file or export them as environment variables
- the simpliest way is to run with `docker-compose`

	```bash
	make docker-compose-x86 DCO_ARGS="up -d"
	```

	Check that everything is Ok

	```bash
	make docker-compose-x86 DCO_ARGS="logs"
	```

-
### How to build docker image

```bash
make docker-image-build-x86
```