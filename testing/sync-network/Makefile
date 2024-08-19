#!/usr/bin/make -f
###############################################################################
###                                Kurtosis                                 ###
###############################################################################
start-full-node:
	kurtosis run ./testing/sync-network --enclave full-node-enclave --args-file ./testing/sync-network/config.yaml

# Stops the running Kurtosis enclave
stop-full-node:
	kurtosis enclave stop full-node-enclave

# Stops and removes the specified Kurtosis enclave
reset-full-node:
	$(MAKE) stop-devnet
	kurtosis enclave rm full-node-enclave

# Removes the specified Kurtosis enclave
rm-full-node:
	kurtosis enclave rm full-node-enclave --force

# Lints Starlark (.star) files in the Kurtosis directory using buildifier
star-lint-full-node:
	@$(MAKE) buildifier-install
	@echo "--> Running buildifier to format starlark files..."
	find ./testing/sync-network -name "*.star" -exec buildifier -mode=check {} +

# Automatically fixes formatting issues in Starlark (.star) files using buildifier
star-fix-full-node:
	@$(MAKE) buildifier-install
	@echo "--> Running buildifier to format starlark files..."
	find ./testing/sync-network -name "*.star" -exec buildifier --mode=fix {} +
	
.PHONY: start-full-node stop-full-node reset-full-node rm-full-node star-lint-full-node star-fix-full-node
