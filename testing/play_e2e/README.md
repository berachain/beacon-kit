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
make test-play
```
If required, add tests under testing/play_e2e folder.

This will automatically build your beacond docker image from the local source
code, and spin up a Kurtosis network based on the config file in
`testing/config/defaults.go`. 

When you want to tear down your network, you can do so
with the following commands:

```bash
make stop-playnet
make rm-playnet
```
or 

```bash
make reset-playnet
```

## Configuration 
In case you want to configure(change) the validator set, consider doing changes in `defaultValidators`.
The user can specify the number of replicas they want per type.

All the default configuration are listed in `testing/config/defaults.go`



