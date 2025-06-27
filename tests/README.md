# E2E tests

These are the Go tests used in the end-to-end testing framework.

To use the tests, run the end-to-end test suite with a manifest:

To run the full e2e test suite:
```
make test-e2e-ci
```

To build all the necessary tools, but not run them:
```
make build-e2e
```

To run a specific test, but not build the tools:
```
make test-e2e-ci-no-build
```

You can replace `ci` with `single` for a single-node testnet or `simple` for a four-node
testnet.

For monitoring, you need to pre-create the `prometheus.yaml` before running the tests:
```
make build-e2e
bin/build/runner -f testing/networks/ci.toml setup
make start-monitoring
make test-e2e-ci-no-build
```

Monitoring is not too picky. If you already started a testnet run and the `prometheus.yml`
was created, you can simply run `make start-monitoring` and it will start collecting the data
from that point.
