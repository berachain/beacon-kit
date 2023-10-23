<h1 align="center"> 🅱️olaris Monorepo ❄️🔭 </h1>

![](./docs/web/public/bear_banner.png)

*The project is still work in progress, see the [disclaimer below](#-warning-under-construction-).*

<div>
  <a href="https://codecov.io/gh/berachain/polaris" target="_blank">
    <img src="https://codecov.io/gh/berachain/polaris/branch/main/graph/badge.svg?token=5SYYGUS8GW"/> 
  </a>
  <a href="https://pkg.go.dev/github.com/itsdevbear/bolaris" target="_blank">
    <img src="https://pkg.go.dev/badge/github.com/itsdevbear/bolaris.svg" alt="Go Reference">
  </a>
  <a href="https://t.me/polaris_devs" target="_blank">
    <img alt="Telegram Chat" src="https://img.shields.io/endpoint?color=neon&logo=telegram&label=chat&url=https%3A%2F%2Ftg.sumanjay.workers.dev%2Fpolaris_devs">
  </a>
  <a href="https://twitter.com/berachain" target="_blank">
    <img alt="Twitter Follow" src="https://img.shields.io/twitter/follow/berachain">
  <a href="https://discord.gg/berachain">
   <img src="https://img.shields.io/discord/984015101017346058?color=%235865F2&label=Discord&logo=discord&logoColor=%23fff" alt="Discord">
  </a>
</div>

## Build & Test

[Golang 1.21+](https://go.dev/doc/install) and [Foundry](https://book.getfoundry.sh/getting-started/installation) are required for Polaris.

1. Install [go 1.21+ from the official site](https://go.dev/dl/) or the method of your choice. Ensure that your `GOPATH` and `GOBIN` environment variables are properly set up by using the following commands:

   For Ubuntu:

   ```sh
   cd $HOME
   sudo apt-get install golang jq -y
   export PATH=$PATH:/usr/local/go/bin
   export PATH=$PATH:$(go env GOPATH)/bin
   ```

   For Mac:

   ```sh
   cd $HOME
   brew install go jq
   export PATH=$PATH:/opt/homebrew/bin/go
   export PATH=$PATH:$(go env GOPATH)/bin
   ```

2. Install Foundry:

   ```sh
   curl -L https://foundry.paradigm.xyz | bash
   ```

3. Clone, Setup and Test:

   ```sh
   cd $HOME
   git clone https://github.com/berachain/polaris
   cd polaris
   git checkout main
   make test-unit
   ```

4. Start a local development network:

   ```sh
   make start
   ```

## 🚧 WARNING: UNDER CONSTRUCTION 🚧

This project is work in progress and subject to frequent changes as we are still working on wiring up the final system.
It has not been audited for security purposes and should not be used in production yet.

The network will have an Ethereum JSON-RPC server running at `http://localhost:8545` and a Tendermint RPC server running at `http://localhost:26657`.
