#!/usr/bin/make -f
###############################################################################
###                                Kurtosis                                 ###
###############################################################################
start-full-node:
	kurtosis run --enclave full-node-enclave ./testing/sync-network

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

.PHONY: start-full-node stop-full-node reset-full-node rm-full-node
