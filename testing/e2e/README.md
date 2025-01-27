# Running BeaconKit with Kurtosis for testing

## What is Kurtosis

[Kurtosis](https://www.kurtosis.com/) is a platform for running distributed
systems on Docker / Kubernetes. It provides a simple, powerful framework for
spinning up and tearing down distributed systems programmatically.

## How to Use

To use BeaconKit with Kurtosis, you'll first need to install the Kurtosis CLI
and its dependencies. You can find instructions for doing so
[here](https://docs.kurtosis.com/install).

### Docker/Test environment

Once you've installed the Kurtosis CLI, you can use it to spin up a Beacon
network with the following command from within the root directory of the
beacon-kit repo:

```bash
make test-e2e
```
If required, add tests under testing/e2e folder.

This will automatically build your beacond docker image from the local source
code, and spin up a Kurtosis network based on the config file in
`testing/e2e/config/defaults.go`.

Currently, the e2e tests runs in different kurtosis enclaves. Check the default configuration in `TestBeaconKitE2ESuite()`.
Play around with the configuration to see how it works. You need to pass the chain ID and chain spec name to the `suite.WithChain()` function. Ensure that the chain ID and chain spec name are valid.


## Configuration
In case you want to configure(change) the validator set, consider doing changes in `defaultValidators`.
The user can specify the number of replicas they want per type.

All the default configuration are listed in `testing/e2e/config/defaults.go`

Note: The default chainID for this local network is 80087, which is our dev network configuration. To make changes to the 80087 chain spec used, modify parameters [here](https://github.com/berachain/beacon-kit/blob/main/config/spec/devnet.go#L40).

## Configure the default network configuration
To change the chainID, modify the `ChainID` field in the `NetworkConfiguration` struct in `defaultNetworkConfiguration` 
function in `testing/e2e/config/defaults.go`.

To change the chainSpec, modify the `ChainSpec` field in the `NetworkConfiguration` struct in `defaultNetworkConfiguration`
function in `testing/e2e/config/defaults.go`.

## Add your tests
Add your tests in here like how it is done in `TestBasicStartup()`



