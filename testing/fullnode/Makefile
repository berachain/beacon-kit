start-full-node:
	kurtosis run --enclave full-node-enclave ./testing/fullnode

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

TESTAPP_FILES_DIR_full = testing/networks/80084
JWT_PATH = ${TESTAPP_FILES_DIR_full}/jwt.hex
ETH_GENESIS_PATH = ${TESTAPP_FILES_DIR_full}/eth-genesis.json
ETH_DATA_DIR = .tmp/eth-home

start-geth-full-node: ## start an ephemeral `geth` node with docker
	rm -rf ${ETH_DATA_DIR}
	docker run \
	--rm -v $(PWD)/${TESTAPP_FILES_DIR_full}:/${TESTAPP_FILES_DIR_full} \
	-v $(PWD)/.tmp:/.tmp \
	ethereum/client-go init \
	--datadir ${ETH_DATA_DIR} \
	${ETH_GENESIS_PATH}

	docker run \
	-p 30303:30303 \
	-p 8545:8545 \
	-p 8551:8551 \
	--rm -v $(PWD)/${TESTAPP_FILES_DIR_full}:/${TESTAPP_FILES_DIR_full} \
	-v $(PWD)/.tmp:/.tmp \
	ethereum/client-go \
	--networkid 80084 \
	--http \
	--http.addr 0.0.0.0 \
	--http.corsdomain "*" \
	--http.api eth,net,web3 \
	--authrpc.jwtsecret $(JWT_PATH) \
	--authrpc.vhosts "*" \
	--datadir ${ETH_DATA_DIR} \
	--ipcpath ${IPC_PATH} \
	--syncmode=snap \
	--verbosity=3 \
	--bootnodes="enode://0401e494dbd0c84c5c0f72adac5985d2f2525e08b68d448958aae218f5ac8198a80d1498e0ebec2ce38b1b18d6750f6e61a56b4614c5a6c6cf0981c39aed47dc@34.159.32.127:30303,enode://7a2f67d22b12e10c6ba9cd951866dda6471604be5fbd5102217dbad1cc56e590befd2009ecc99958a468a5b8e0dc28e14d9b6822491719c93199be6aa0319077@34.124.220.31:30303,enode://e31aa249638083d34817eed2b499ccd4b0718a332f0ea530e3062e13f624cb03a7d6b6e0460193ee87b5fc12e73a726070a3126ef53492ffbdc5e6c102f6dfb3@34.64.198.56:30303,enode://f24b54da77cf604e92aeb5ee5e79401fd3e66111563ca630e72330ccab6f385ccbbde5eba4577ee7bfb5e83347263d0e4cad042fd4c10468d0e38906fc82ba31@bera-testnet-seeds.nodeinfra.com:30303" \


	# --authrpc.addr 0.0.0.0 \

## Testing:
start-beaconkit: ## start an ephemeral `beacond` node
	@JWT_SECRET_PATH=$(JWT_PATH) RPC_DIAL_URL=https://bartio.rpc.berachain.com/ testing/fullnode/entrypoint.sh

.PHONY: start-full-node stop-full-node reset-full-node rm-full-node start-geth-full-node start-beaconkit

