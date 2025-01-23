# Running BeaconKit with Kurtosis

## What is Kurtosis

[Kurtosis](https://www.kurtosis.com/) is a platform for running distributed systems on Docker / Kubernetes.
It provides a simple, powerful framework for spinning up and tearing down distributed systems programmatically.

## How to Use

To use BeaconKit with Kurtosis, you'll first need to install the Kurtosis CLI and its dependencies.
You can find instructions for doing so [here](https://docs.kurtosis.com/install).

### Docker/local environment

Once you've installed the Kurtosis CLI, you can use it to spin up a Beacon network with the following command from within the root directory of the `beacon-kit` repository:

```sh
make start-devnet
```

This will automatically build your `beacond` Docker image from the local source code and spin up a Kurtosis network based on the config file in `kurtosis/beaconkit-local.yaml`. Once complete, this will output all the network information for your nodes, like so:

![Example Network](./img/example-network.png)

When you want to tear down your network, you can do so with the following commands:

```sh
make stop-devnet
make rm-devnet
```

And that's it!

## Deploy Devnet to Kubernetes Networks

### Deploy to a Google Cloud Network

This will allow you to deploy a network orchestrated in the same way as `make start-devnet`, but on a cloud environment. A similar approach can be taken for a local Kubernetes environment (alternative commands will be commented with [Docker Desktop K8s]).

1. First, open your Kurtosis config:

   ```sh
   kurtosis config path

   # The command will output a path which you need to open in an editor
   /Users/.../kurtosis-config.yml
   ```

2. Update the Kurtosis config with the following, replacing the entire file:

   ```yaml
   config-version: 2
   should-send-metrics: true
   kurtosis-clusters:
     docker:
       type: "docker"
     docker-desktop:
       type: "kubernetes"
       config:
         kubernetes-cluster-name: "docker-desktop"
         storage-class: "hostpath"
     cloud:
       type: "kubernetes"
       config:
         kubernetes-cluster-name: "cloud"
         storage-class: "premium-rwo"
   ```

3. Next, ensure Kurtosis is using the correct config so it deploys to the cloud instead of local Docker Desktop:

   ```sh
   kurtosis cluster set cloud

   # [Docker Desktop K8s]: kurtosis config use-context docker-desktop
   ```

4. Now ensure your Kubernetes config is using the correct context, i.e., the one context you wish to deploy to:

   ```sh
   kubectl config use-context gke_prj-.....

   # [Docker Desktop K8s]: kubectl config use-context docker-desktop
   ```

5. Run Kurtosis Gateway. This command will start a local "gateway" to connect your local machine to your remote Kubernetes cluster. Run this in a separate shell:

   ```sh
   kurtosis gateway
   ```

6. Cloud-based deployments require a Docker image, as local Docker images cannot be pulled from the remote instance. If you want to update the image, edit:

   ```yaml
   # Found in beacon-kit/kurtosis/beaconkit-cloud.yaml
   images:
     beaconkit: ghcr.io/berachain/beacon-kit:main
   ```

7. Deploy. Note that re-executing the same command twice will start the network from zero again unless you change the enclave name in the `Makefile`:

   ```sh
   make start-devnet-cloud
   ```

8. View your deployment in K9s, navigating to the relevant namespace. It should be named `kt-my-cloud-devnet-${whoami}`.

## Helper Commands

If you want to start from a clean state and remove all existing pods:

```sh
# Everything is wrecked
kurtosis clean -a
kurtosis engine restart
```

If you manually kill a pod and want to restart it:

```sh
# Note that for the namespace, you should remove the "kt-" prefix
kurtosis service start {namespace} {podname}
