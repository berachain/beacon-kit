#!/usr/bin/make -f

###############################################################################
###                                Kurtosis                                 ###
###############################################################################

verify-kurtosis-installed:
	@command -v "kurtosis" > /dev/null || (echo "Please install Kurtosis, see kurtosis/README.md for instructions" && exit 1)

# Starts a Kurtosis enclave containing a local devnet.
start-devnet: verify-kurtosis-installed
	$(MAKE) build-docker VERSION=kurtosis-local start-devnet-no-build

# Starts a Kurtosis enclave containing a local devnet without building the image
start-devnet-no-build:
	kurtosis run ./kurtosis --args-file ./kurtosis/beaconkit-all.yaml \
		--enclave my-local-devnet --parallelism 200

# Starts a Kurtosis enclave containing a local devnet on GCP.
# --production flag is used to indicate that the enclave is
# running in production mode to allow pod restarts when doing chaos testing.
start-gcp-devnet-no-build:
	kurtosis run ./kurtosis --args-file ./kurtosis/beaconkit-base-gcp.yaml \
		--enclave my-gcp-devnet2 --parallelism 200 --production --image-download always

# Stops the running Kurtosis enclave on GCP
stop-gcp-devnet:
	kurtosis enclave stop my-gcp-devnet

# Stops and removes the specified Kurtosis enclave
reset-gcp-devnet:
	$(MAKE) stop-gcp-devnet
	kurtosis enclave rm my-gcp-devnet

# Removes the specified Kurtosis enclave
rm-gcp-devnet:
	kurtosis enclave rm my-gcp-devnet --force

# Stops the running Kurtosis enclave
stop-devnet:
	kurtosis enclave stop my-local-devnet

# Stops and removes the specified Kurtosis enclave
reset-devnet:
	$(MAKE) stop-devnet
	kurtosis enclave rm my-local-devnet

# Removes the specified Kurtosis enclave
rm-devnet:
	kurtosis enclave rm my-local-devnet --force

# Installs buildifier, a tool for linting and formatting starlark files.
buildifier-install:
	@echo "--> Installing buildifier"
	@go install github.com/bazelbuild/buildtools/buildifier

# Lints Starlark (.star) files in the Kurtosis directory using buildifier
star-lint:
	@$(MAKE) buildifier-install
	@echo "--> Running buildifier to format starlark files..."
	find ./kurtosis -name "*.star" -exec buildifier -mode=check {} +

# Automatically fixes formatting issues in Starlark (.star) files using buildifier
star-fix:
	@$(MAKE) buildifier-install
	@echo "--> Running buildifier to format starlark files..."
	find ./kurtosis -name "*.star" -exec buildifier --mode=fix {} +

# Marks targets as not being associated with files
.PHONY: start-devnet stop-devnet reset-devnet rm-devnet buildifier-install \
  star-lint star-fix start-gcp-devnet-no-build stop-gcp-devnet \
  reset-gcp-devnet rm-gcp-devnet install-kurtosis
