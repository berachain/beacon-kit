<div align="center">
  <a href="https://github.com/berachain/beacon-kit">
    <picture>
      <source media="(prefers-color-scheme: dark)" srcset="https://res.cloudinary.com/duv0g402y/image/upload/v1718034312/BeaconKitBanner.png">
      <img alt="BeaconKit Banner" src="https://res.cloudinary.com/duv0g402y/image/upload/v1718034312/BeaconKitBanner.png">
    </picture>
  </a>
</div>

<h2 align="center">A Modular Framework for Building EVM Consensus Clients ⛵️✨</h2>

<p align="center">
  <em>A highly customizable framework for creating and operating EVM-based blockchains.</em><br>
  <strong>⚠️ Project under active development — see <a href="#status">status</a> for details.</strong>
</p>

<div align="center">
  <a href="https://github.com/berachain/beacon-kit/actions/workflows/pipeline.yml">
    <img src="https://github.com/berachain/beacon-kit/workflows/pipeline/badge.svg" alt="CI Status">
  </a>
  <a href="https://codecov.io/gh/berachain/beacon-kit">
    <img src="https://codecov.io/gh/berachain/beacon-kit/graph/badge.svg?token=0l5iJ3ZbzV" alt="Code Coverage">
  </a>
  <a href="https://t.me/beacon_kit">
    <img src="https://img.shields.io/endpoint?color=neon&logo=telegram&label=chat&url=https%3A%2F%2Ftg.sumanjay.workers.dev%2Fbeacon_kit" alt="Telegram Chat">
  </a>
  <a href="https://x.com/berachain">
    <img src="https://img.shields.io/twitter/follow/berachain" alt="Follow on X">
  </a>
  <a href="https://discord.gg/berachain">
    <img src="https://img.shields.io/discord/924442927399313448?label=discord" alt="Join on Discord">
  </a>
</div>

---

## 🚀 What is BeaconKit?

[BeaconKit](https://docs.berachain.com/learn/what-is-beaconkit) is a **modular framework** designed to simplify the process of building Ethereum Virtual Machine (EVM) based consensus clients. It ensures a functionally identical execution environment to Ethereum Mainnet, offering a user-friendly approach for blockchain developers.

---

## 💡 Supported Execution Clients

BeaconKit supports all major Ethereum execution clients through the [Ethereum Engine API](https://github.com/ethereum/execution-apis/blob/main/src/engine):

- **[Geth](https://geth.ethereum.org/):** Official Go implementation of Ethereum.
- **[Erigon](https://erigon.tech/):** High-performance, feature-rich fork of Geth.
- **[Nethermind](https://www.nethermind.io/):** .NET-based client with full Ethereum protocol support.
- **[Besu](https://www.lfdecentralizedtrust.org/projects/besu):** Java-based, enterprise-grade client.
- **[Reth](https://reth.rs/):** Rust-based client focused on performance.
- **[Ethereumjs](https://ethereumjs.readthedocs.io/en/latest/#):** JavaScript-based implementation managed by the Ethereum Foundation.

---

## 🛠️ Running a Local Development Network

### Prerequisites

Ensure the following tools are installed:

- [Docker](https://docs.docker.com/engine/install/)
- [Golang 1.23.0+](https://go.dev/doc/install)
- [Foundry](https://book.getfoundry.sh/)

### Steps

1. **Open two terminals side-by-side.**

   **Terminal 1:**
   ```bash
   # Start the sample BeaconKit Consensus Client:
   make start

# Start an Ethereum Execution Client:
make start-reth # or start-geth, start-besu, start-erigon, start-nethermind, start-ethereumjs


🔗 Multinode Local Devnet
Refer to the Kurtosis README for instructions on running a multinode local development network.

⚠️ Status
This project is a work in progress. Frequent changes are expected as we finalize the system. Audits are ongoing, and BeaconKit is not yet recommended for production environments.

📄 License
BeaconKit is licensed under MIT License.

🛡️ Support & Community
Join us on Telegram or Discord.
Follow updates on X (formerly Twitter).
