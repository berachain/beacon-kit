# runtime

The runtime package contains the core components of beacon-kit:

- *abci*: This package connects the core beacon-chain business logic to the ABCI lifecycle.
- *modules*: These are the cosmos-sdk modules required for the beacon-kit chain to run.
- *service*: This defines the base service that all core beacon-chain embed.

The `runtime` package itself contains the `BeaconKitRuntime`, which is 
the main entrypoint for all of beacon-kit.
