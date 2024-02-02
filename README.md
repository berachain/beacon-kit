# beacon-kit 

<!-- [![CI status](https://github.com/itsdevbear/bolaris/workflows/ci/badge.svg)][gh-ci] -->
<!-- [![cargo-deny status](https://github.com/paradigmxyz/reth/workflows/deny/badge.svg)][gh-deny]
[![Codecov](https://img.shields.io/codecov/c/github/paradigmxyz/reth?token=c24SDcMImE)][codecov] -->
<!-- [![Telegram Chat][tg-badge]][tg-url] -->

**A modular and customizable consensus layer for Ethereum based blockchains**

![](.github/assets/banner.png)


## What is BeaconKit?

BeaconKit introduces an innovative framework that utilizes the Cosmos-SDK to create a flexible, customizable consensus layer tailored for Ethereum-based blockchains. The framework offers the most user-friendly way to build and operate an EVM blockchain, while ensuring a functionally identical execution environment to that of the Ethereum Mainnet.

You've heard of EVM Compatibility, then EVM Equivalence, and now with BeaconKit, we introduce *EVM Identicality*.

## Supported Execution Clients

Through utilizing the [Ethereum Engine API](https://github.com/ethereum/execution-apis/blob/main/src/engine) BeaconKit is able to support all 5 major Ethereum execution clients:

- **Geth**: Official Go implementation of the Ethereum protocol.
- **Erigon**: Formerly known as Turbo-Geth, it is a more performant and feature-rich Ethereum client forked from `go-ethereum`.
- **Nethermind**: .NET based Ethereum client with full support for Ethereum and other blockchain protocols.
- **Besu**: An enterprise-grade Ethereum client developed under the Apache 2.0 license and written in Java.
- **Reth**: A Rust-based Ethereum client, focusing on performance and reliability.


## Why BeaconKit? 

TODO: Talk about Polaris / Ethermint compatibility issue.