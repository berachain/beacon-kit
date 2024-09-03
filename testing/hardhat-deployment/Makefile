#!/usr/bin/make -f

###############################################################################
###                                Kurtosis                                 ###
###############################################################################

# Starts a Kurtosis enclave containing a local devnet
start-hardhat-deploy: install-kurtosis get-ip-address
	kurtosis run ./testing/hardhat-deployment --args-file ./testing/hardhat-deployment/config.yaml \
	--enclave hardhat-devnet

# Stops the running Kurtosis enclave
stop-hardhat-deploy:
	kurtosis enclave stop hardhat-devnet

# Stops and removes the specified Kurtosis enclave
reset-hardhat-deploy:
	$(MAKE) stop-hardhat-deploy
	kurtosis enclave rm hardhat-devnet 

# Removes the specified Kurtosis enclave
rm-hardhat-deploy:
	kurtosis enclave rm hardhat-devnet --force

# Lints Starlark (.star) files in the Kurtosis directory using buildifier
star-lint-hardhat-deploy:
	@$(MAKE) buildifier-install
	@echo "--> Running buildifier to format starlark files..."
	find ./testing/hardhat-deployment -name "*.star" -exec buildifier -mode=check {} +

# Automatically fixes formatting issues in Starlark (.star) files using buildifier
star-fix-hardhat-deploy:
	@$(MAKE) buildifier-install
	@echo "--> Running buildifier to format starlark files..."
	find ./testing/hardhat-deployment -name "*.star" -exec buildifier --mode=fix {} +

# Marks targets as not being associated with files
.PHONY: start-hardhat-deploy stop-hardhat-deploy reset-hardhat-deploy rm-hardhat-deploy star-lint-hardhat-deploy star-fix-hardhat-deploy
