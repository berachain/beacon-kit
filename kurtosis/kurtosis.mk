#!/usr/bin/make -f

###############################################################################
###                                Kurtosis                                 ###
###############################################################################

# Starts a Kurtosis enclave containing a local devnet.
start-devnet:
	kurtosis run ./kurtosis --enclave my-local-devnet

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
.PHONY: start-kurtosis stop-kurtosis rm-kurtosis buildifier-install star-lint star-fix
