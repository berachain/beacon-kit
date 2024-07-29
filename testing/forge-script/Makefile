#!/usr/bin/make -f

###############################################################################
###                                Kurtosis                                 ###
###############################################################################

# Starts a Kurtosis enclave containing a local devnet
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
.PHONY: start-forge-script stop-forge-script reset-forge-script rm-forge-script star-lint-forge-script star-fix-forge-script
