#!/usr/bin/make -f

###############################################################################
###                                Kurtosis                                 ###
###############################################################################

# Installs Kurtosis if not already installed
install-kurtosis:
	@echo "Checking for Kurtosis installation..."
	@if ! command -v kurtosis &> /dev/null; then \
		echo "Kurtosis could not be found, installing..."; \
		OS=$$(uname -s | tr A-Z a-z); \
		if [ "$$OS" = "darwin" ]; then \
			brew install kurtosis-tech/tap/kurtosis-cli; \
		elif [ "$$OS" = "linux" ]; then \
			ARCH=$$(uname -m); \
			if [ "$$ARCH" = "x86_64" ]; then ARCH="amd64"; \
			elif [ "$$ARCH" = "arm64" ]; then ARCH="arm64"; \
			else echo "Unsupported architecture $$ARCH for Kurtosis installation" && exit 1; fi; \
			TAG=`curl -s "https://api.github.com/repos/kurtosis-tech/kurtosis-cli-release-artifacts/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/'`; \
			curl -Lo kurtosis.tar.gz "https://github.com/kurtosis-tech/kurtosis-cli-release-artifacts/releases/download/$TAG/kurtosis-cli_${TAG}_${OS}_${ARCH}.tar.gz"; \
			tar -xzf kurtosis.tar.gz; \
			rm kurtosis.tar.gz; \
			chmod +x kurtosis; \
			sudo mv kurtosis /usr/local/bin/; \
		else \
			echo "Unsupported OS $$OS for Kurtosis installation" && exit 1; \
		fi; \
	else \
		echo "Kurtosis is already installed"; \
	fi

# Starts a Kurtosis enclave containing a local devnet.
start-devnet: install-kurtosis
	$(MAKE) build-docker VERSION=kurtosis-local start-devnet-no-build

# Starts a Kurtosis enclave containing a local devnet without building the image
start-devnet-no-build:
	kurtosis run ./kurtosis --args-file ./kurtosis/beaconkit-local.yaml \
		--enclave my-local-devnet --parallelism 200

# Starts a Kurtosis enclave in the cloud using the latest image tag
start-devnet-cloud: install-kurtosis
	kurtosis run ./kurtosis --args-file ./kurtosis/beaconkit-cloud.yaml \
		--enclave my-cloud-devnet-$(shell whoami) --parallelism 200 --production --image-download always

# Remove the running Kurtosis enclave on GCP
rm-devnet-cloud:
	kurtosis enclave rm my-cloud-devnet-$(shell whoami) --force

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
	@go install github.com/bazelbuild/buildtools/buildifier@latest

# Lints Starlark (.star) files in the Kurtosis directory using buildifier
star-lint: buildifier-install
	@echo "--> Running buildifier to format starlark files..."
	find ./kurtosis -name "*.star" -exec buildifier -mode=check {} +

# Automatically fixes formatting issues in Starlark (.star) files using buildifier
star-fix: buildifier-install
	@echo "--> Running buildifier to format starlark files..."
	find ./kurtosis -name "*.star" -exec buildifier --mode=fix {} +

# Marks targets as not being associated with files
.PHONY: start-devnet stop-devnet start-devnet-cloud rm-devnet-cloud reset-devnet rm-devnet buildifier-install \
  star-lint star-fix install-kurtosis