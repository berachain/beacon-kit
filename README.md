
<div align="center">
  <a href="https://github.com/berachain/beacon-kit">
    <picture>
      <source media="(prefers-color-scheme: dark)" srcset="https://res.cloudinary.com/duv0g402y/image/upload/v1718034312/BeaconKitBanner.png">
      <img alt="BeaconKit Banner" src="https://res.cloudinary.com/duv0g402y/image/upload/v1718034312/BeaconKitBanner.png" width="auto" height="auto">
    </picture>
  </a>
</div>

<h2>
  ⚡ BeaconKit – A Modular Framework for Building EVM Consensus Clients ⛵✨
</h2>

<div>

[![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/berachain/beacon-kit/pipeline.yml?label=CI&logo=github)](https://github.com/berachain/beacon-kit/actions/workflows/pipeline.yml)
[![Code Coverage](https://img.shields.io/codecov/c/github/berachain/beacon-kit?logo=codecov)](https://codecov.io/gh/berachain/beacon-kit)
[![GitHub Repo stars](https://img.shields.io/github/stars/berachain/beacon-kit?logo=github&color=yellow)](https://github.com/berachain/beacon-kit/stargazers)
[![GitHub forks](https://img.shields.io/github/forks/berachain/beacon-kit?logo=github&color=blue)](https://github.com/berachain/beacon-kit/network/members)
[![GitHub last commit](https://img.shields.io/github/last-commit/berachain/beacon-kit?logo=git)](https://github.com/berachain/beacon-kit/commits/main)
[![Discord](https://img.shields.io/discord/924442927399313448?logo=discord&color=5865F2)](https://discord.gg/berachain)
[![X Follow](https://img.shields.io/twitter/follow/berachain)](https://x.com/berachain)
[![Telegram Chat](https://img.shields.io/endpoint?color=neon&logo=telegram&label=chat&url=https%3A%2F%2Ftg.sumanjay.workers.dev%2Fbeacon_kit)](https://t.me/beacon_kit)
</div>

---

## 🔹 What is BeaconKit?

[BeaconKit](https://docs.berachain.com/learn/what-is-beaconkit) is a **modular framework** for building **EVM-based consensus clients**.  
It provides the **most user-friendly** way to build and operate an EVM blockchain while ensuring **full compatibility** with the Ethereum Mainnet.

### ✨ **Key Features**
- ⚙️ **Modular architecture** – highly customizable and flexible.
- 🚀 **Optimized performance** – designed for efficiency and scalability.
- 🔄 **Ethereum-compatible** – works seamlessly via the [Ethereum Engine API](https://github.com/ethereum/execution-apis/blob/main/src/engine).

---

## 🖥 **Supported Execution Clients**

BeaconKit integrates with Ethereum's Execution Layer via the **Ethereum Engine API**, supporting **six major execution clients**:

| 🚀 Client | 🌍 Description |
|-----------|--------------|
| [**Geth**](https://geth.ethereum.org/) | Official Go implementation of the Ethereum protocol. |
| [**Erigon**](https://erigon.tech/) | High-performance, feature-rich client forked from `go-ethereum`. |
| [**Nethermind**](https://www.nethermind.io/) | .NET-based client with full Ethereum protocol support. |
| [**Besu**](https://www.lfdecentralizedtrust.org/projects/besu) | Enterprise-grade client, Apache 2.0 licensed, written in Java. |
| [**Reth**](https://reth.rs/) | Rust-based client focusing on performance and reliability. |
| [**EthereumJS**](https://ethereumjs.readthedocs.io/en/latest/#) | JavaScript-based client managed by the Ethereum Foundation. |

---

## 🛠 **Running a Local Development Network**

### ✅ **Prerequisites**
📦 [Docker](https://docs.docker.com/engine/install/)  
🏗 [Golang 1.23.0+](https://go.dev/doc/install)  
🔥 [Foundry](https://book.getfoundry.sh/)  

### 🚀 **Run the Network in Two Terminals**

📌 **Terminal 1 – Start the BeaconKit Consensus Client**
```bash
make start
```


📌 Terminal 2 – Start an Ethereum Execution Client
```bash
make start-reth # or start-geth start-besu start-erigon start-nethermind start-ethereumjs
```
🔑 Preloaded Accounts
After starting the local network, you will have access to a preloaded account with native EVM tokens:

```bash
Private Key: 0xfffdbb37105441e14b0ee6330d855d8504ff39e705c3afa8f859ac9865f99306
Address: 0x20f33ce90a13a4b5e7697e3544c3083b8f8a51d4
```
You can use this account to deploy contracts and interact with the network.


## Multinode Local Devnet

Please refer to the [Kurtosis README](https://github.com/berachain/beacon-kit/blob/main/kurtosis/README.md) for more information on how to run a multinode local devnet.

## 💬 **Join the Community**

<p align="left">
  <a href="https://t.me/beacon_kit">
    <img src="https://img.shields.io/badge/Telegram-26A5E4?logo=telegram&logoColor=white&style=for-the-badge" alt="Telegram">
  </a>
  <a href="https://discord.gg/berachain">
    <img src="https://img.shields.io/badge/Discord-5865F2?logo=discord&logoColor=white&style=for-the-badge" alt="Discord">
  </a>
  <a href="https://x.com/berachain">
    <img src="https://img.shields.io/badge/Twitter-000000?logo=x&logoColor=white&style=for-the-badge" alt="Twitter (X)">
  </a>
</p>



