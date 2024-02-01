# beacon-kit 

<!-- [![CI status](https://github.com/itsdevbear/bolaris/workflows/ci/badge.svg)][gh-ci] -->
<!-- [![cargo-deny status](https://github.com/paradigmxyz/reth/workflows/deny/badge.svg)][gh-deny]
[![Codecov](https://img.shields.io/codecov/c/github/paradigmxyz/reth?token=c24SDcMImE)][codecov] -->
<!-- [![Telegram Chat][tg-badge]][tg-url] -->

**A modular and customizable consensus layer framework for Ethereum based blockchains**

![](.github/assets/banner.png)


## What is BeaconKit?

BeaconKit is a novel framework that leverages the Cosmos-SDK to build a modular, customizable consensus layer for Ethereum based blockchains. It focuses on providing an extremely developer friendly way to spin up Ethereum Virtual Machine based blockchains, control / add consensus rules, while ensuring 100% compatibilty with existing Ethereum tooling. 


## Supported Execution Clients

Through utilizing the [Ethereum Engine API](https://github.com/ethereum/execution-apis/blob/main/src/engine) BeaconKit is able to support all 5 major Ethereum execution clients:

- **Geth**: Official Go implementation of the Ethereum protocol.
- **Erigon**: Formerly known as Turbo-Geth, it is a more performant and feature-rich Ethereum client forked from `go-ethereum`.
- **Nethermind**: .NET based Ethereum client with full support for Ethereum and other blockchain protocols.
- **Besu**: An enterprise-grade Ethereum client developed under the Apache 2.0 license and written in Java.
- **Reth**: A Rust-based Ethereum client, focusing on performance and reliability.


## Why BeaconKit? 

TODO: Talk about Polaris / Ethermint compatibility issue.