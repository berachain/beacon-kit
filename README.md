# beacon-kit

[![Telegram Chat](https://img.shields.io/endpoint?color=neon&logo=telegram&label=chat&url=https%3A%2F%2Ftg.sumanjay.workers.dev%2Fbeacon_kit)](https://t.me/beacon_kit)

## A modular consensus framework for building layer 1/2 evm blockchains ⛵️✨

![banner](.github/assets/banner.png)

## 🚧 WARNING: UNDER CONSTRUCTION 🚧

This project is work in progress and subject to frequent changes as we are still working on wiring up the final system. It has not been audited for security purposes and should not be used in production yet.

## What is BeaconKit

BeaconKit introduces an innovative framework that utilizes the Cosmos-SDK to
create a flexible, customizable consensus layer tailored for Ethereum-based
blockchains. The framework offers the most user-friendly way to build and
operate an EVM blockchain, while ensuring a functionally identical execution
environment to that of the Ethereum Mainnet.

First there was EVM Compatibility; next, EVM Equivalence; and now with
BeaconKit, **EVM Identicality**.

## Supported Execution Clients

Through utilizing the [Ethereum Engine API](https://github.com/ethereum/execution-apis/blob/main/src/engine)
BeaconKit is able to support all 6 major Ethereum execution clients:

- **Geth**: Official Go implementation of the Ethereum protocol.
- **Erigon**: More performant, feature-rich client forked from `go-ethereum`.
- **Nethermind**: .NET based client with full support for Ethereum protocols.
- **Besu**: Enterprise-grade client, Apache 2.0 licensed, written in Java.
- **Reth**: Rust-based client focusing on performance and reliability.
- **Ethereumjs**: Javascript based client managed by the Ethereum Foundation.

## Running a Local Development Network

**Prerequisites:**

- [Docker](https://docs.docker.com/engine/install/)
- [Golang 1.22.0+](https://go.dev/doc/install)
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
make start-reth # or start-geth start-besu start-erigon start-nethermind start-ethereumjs
```

The account with
`private-key=0xfffdbb37105441e14b0ee6330d855d8504ff39e705c3afa8f859ac9865f99306`
corresponding with `address=0x20f33ce90a13a4b5e7697e3544c3083b8f8a51d4` is
preloaded with the native EVM token.
