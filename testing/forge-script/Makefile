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
			curl -Lo kurtosis.tar.gz "https://github.com/kurtosis-tech/kurtosis-cli/releases/latest/download/kurtosis-$$OS-$$ARCH.tar.gz"; \
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


# Starts a Kurtosis enclave containing a local devnet without building the image
start-forge-script: install-kurtosis
	kurtosis run ./testing/forge-script --args-file ./testing/forge-script/forge-config.yaml \
	--enclave forge-devnet

# Stops the running Kurtosis enclave
stop-forge-script:
	kurtosis enclave stop forge-devnet

# Stops and removes the specified Kurtosis enclave
reset-forge-script:
	$(MAKE) stop-forge-script
	kurtosis enclave rm forge-devnet 

# Removes the specified Kurtosis enclave
rm-forge-script:
	kurtosis enclave rm forge-devnet --force

# Installs buildifier, a tool for linting and formatting starlark files.
buildifier-install:
	@echo "--> Installing buildifier"
	@go install github.com/bazelbuild/buildtools/buildifier

# Lints Starlark (.star) files in the Kurtosis directory using buildifier
star-lint-forge-script:
	@$(MAKE) buildifier-install
	@echo "--> Running buildifier to format starlark files..."
	find ./testing/forge-script -name "*.star" -exec buildifier -mode=check {} +

# Automatically fixes formatting issues in Starlark (.star) files using buildifier
star-fix-forge-script:
	@$(MAKE) buildifier-install
	@echo "--> Running buildifier to format starlark files..."
	find ./testing/forge-script -name "*.star" -exec buildifier --mode=fix {} +

# Marks targets as not being associated with files
.PHONY: start-forge-script stop-forge-script reset-forge-script rm-forge-script buildifier-install star-lint-forge-script star-fix-forge-script
