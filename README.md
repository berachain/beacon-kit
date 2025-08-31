</br>

<div align="center">
  <a href="https://github.com/berachain/beacon-kit">
    <picture>
      <source media="(prefers-color-scheme: dark)" srcset="https://res.cloudinary.com/duv0g402y/image/upload/v1718034312/BeaconKitBanner.png">
      <img alt="beacon-kit-banner" src="https://res.cloudinary.com/duv0g402y/image/upload/v1718034312/BeaconKitBanner.png" width="auto" height="auto">
    </picture>
  </a>
</div>
<h2 >
  A modular framework for building EVM consensus clients ⛵️✨
</h2>

<div>

[![CI status](https://github.com/berachain/beacon-kit/workflows/pipeline/badge.svg)](https://github.com/berachain/beacon-kit/actions/workflows/pipeline.yml)
[![CodeCov](https://codecov.io/gh/berachain/beacon-kit/graph/badge.svg?token=0l5iJ3ZbzV)](https://codecov.io/gh/berachain/beacon-kit)
[![Telegram Chat](https://img.shields.io/endpoint?color=neon&logo=telegram&label=chat&url=https%3A%2F%2Ftg.sumanjay.workers.dev%2Fbeacon_kit)](https://t.me/beacon_kit)
[![X Follow](https://img.shields.io/twitter/follow/berachain)](https://x.com/berachain)
[![Discord](https://img.shields.io/discord/924442927399313448?label=discord)](https://discord.gg/berachain)

</div>

## What is BeaconKit?

[BeaconKit](https://docs.berachain.com/learn/what-is-beaconkit) is a modular framework for building EVM based consensus clients.
The framework offers the most user-friendly way to build and operate an EVM blockchain, while ensuring a functionally identical execution environment to that of the Ethereum Mainnet.

## Supported Execution Clients

Through utilizing the [Ethereum Engine API](https://github.com/ethereum/execution-apis/blob/main/src/engine)
BeaconKit supports the following execution clients:

- [**Bera-Geth**](https://github.com/berachain/bera-geth): Official Go implementation of the Berachain protocol.
- [**Bera-Reth**](https://github.com/berachain/bera-reth): Rust-based client focusing on performance and reliability.

## Running a Local Development Network

**Prerequisites:**

- [Docker](https://docs.docker.com/engine/install/)
- [Golang 1.23.0+](https://go.dev/doc/install)
- [Foundry](https://book.getfoundry.sh/)

Start by opening two terminals side-by-side:

**Terminal 1:**

```bash
# Start the sample BeaconKit Consensus Client:
make start
```

**Terminal 2:**

**Note:** This must be run *after* the `beacond` node is started since `make start` will populate the
eth-genesis file used by the Execution Client.

```bash
# Start an Ethereum Execution Client:
make start-reth # or start-geth
```

The account with
`private-key=0xfffdbb37105441e14b0ee6330d855d8504ff39e705c3afa8f859ac9865f99306`
corresponding with `address=0x20f33ce90a13a4b5e7697e3544c3083b8f8a51d4` is
preloaded with the native EVM token.

## Multinode Local Devnet

Please refer to the [Kurtosis README](https://github.com/berachain/beacon-kit/blob/main/kurtosis/README.md) for more information on how to run a multinode local devnet.

## Important Commands and Options

`beacond help` lists available commands. Some commands have sub-commands.

`beacond init` creates a folder structure for beacond to operate in, along with initial configuration files, the important ones being `app.toml` and `config.toml`.

`beacond start` starts the chain client and begins the syncing process.

The key configuration value is the Beacon chain specification ("chainspec")  - analogous to an Eth genesis - is held in a spec.toml file. There are three well-known chainspecs (`mainnet|testnet|devnet`).  Or, a custom file can be provided to provide your own scenario's settings - see this [example using Docker](https://docs.berachain.com/nodes/guides/docker-devnet).

The chainspec is set with the `--beacon-kit.chain-spec` command line option, and if necessary for a custom chainspec, use `--beacon-kit.chain-spec-file`.  If used during `beacond init`, this value is written into app.toml and does not need to be specified anymore.

You can override the default operating directories for beacond with the `--home <path>` option.

The [Berachain Node Quickstart](https://docs.berachain.com/nodes/quickstart) provides a quick deployment of mainnet or testnet on your desk. 

For developing with beacon-kit, you have options:
1. see the Makefile for targets to start stand-alone processes
2. see [an example deploying a team with Docker on a custom chainspec](https://docs.berachain.com/nodes/guides/docker-devnet)
3. see [a deployment with Kurtosis](https://docs.berachain.com/nodes/guides/kurtosis)
4. see [an expect script that does a complete cycle of deposits, withdrawal, eviction, voluntary exit](https://github.com/berachain/guides/blob/main/apps/local-docker-devnet/devnet-automation.exp)
