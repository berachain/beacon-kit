# beacon-kit 

<!-- [![CI status](https://github.com/itsdevbear/bolaris/workflows/ci/badge.svg)][gh-ci] -->
<!-- [![cargo-deny status](https://github.com/paradigmxyz/reth/workflows/deny/badge.svg)][gh-deny]
[![Codecov](https://img.shields.io/codecov/c/github/paradigmxyz/reth?token=c24SDcMImE)][codecov] -->
<!-- [![Telegram Chat][tg-badge]][tg-url] -->

**A modular and customizable consensus layer for Ethereum based blockchains**

![](.github/assets/banner.png)


## What is BeaconKit?

BeaconKit introduces an innovative framework that utilizes the Cosmos-SDK to create a flexible, customizable consensus layer tailored for Ethereum-based blockchains. The framework offers the most user-friendly way to build and operate an EVM blockchain, while ensuring a functionally identical execution environment to that of the Ethereum Mainnet.

First there was EVM Compatibility; next, EVM Equivalence; and now with BeaconKit, **EVM Identicality**.

## Why BeaconKit? 

TODO: Talk about Polaris / Ethermint compatibility issue.

## Supported Execution Clients

Through utilizing the [Ethereum Engine API](https://github.com/ethereum/execution-apis/blob/main/src/engine) BeaconKit is able to support all 5 major Ethereum execution clients:

- **Geth**: Official Go implementation of the Ethereum protocol.
- **Erigon**: Formerly known as Turbo-Geth, it is a more performant and feature-rich Ethereum client forked from `go-ethereum`.
- **Nethermind**: .NET based Ethereum client with full support for Ethereum and other blockchain protocols.
- **Besu**: An enterprise-grade Ethereum client developed under the Apache 2.0 license and written in Java.
- **Reth**: A Rust-based Ethereum client, focusing on performance and reliability.

## Documentation
BeaconKit leverages `godoc` for it's core documentation, you can run `godoc` locally and run a web-ui of the 
latest documentation:

```bash
make godoc 
```

## Running a Local Development Network

**Prerequisites:**
- [Docker](https://docs.docker.com/engine/install/)
- [Golang 1.21.6+](https://go.dev/doc/install)
- [Foundry](https://book.getfoundry.sh/getting-started/installation)

Start by opening two terminals side-by-side:

**Terminal 1:**
```bash
# Start the sample BeaconKit Consensus Client:
make start
```

**Terminal 2:**
```bash
# Start an Ethereum Execution Client:
make start-reth # or start-geth start-besu start-erigon start-nethermind
```

The account with `private-key=0xfffdbb37105441e14b0ee6330d855d8504ff39e705c3afa8f859ac9865f99306` corresponding 
with `address=0x20f33ce90a13a4b5e7697e3544c3083b8f8a51d4` is preloaded with the native EVM token.