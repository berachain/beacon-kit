start-nimbus:
	rm -rf ./tmp/nimbus
	docker run \
	--rm -v $(PWD)/${TESTAPP_DIR}:/${TESTAPP_DIR} \
	-v $(PWD)/.tmp:/.tmp \
	-p 30303:30303 \
	-p 8545:8545 \
	statusim/nimbus-fluffy:amd64-master-latest \
	--data-dir=.tmp/nimbus \
	--rpc \
	--rpc-port=8545 \
	--rpc-address=0.0.0.0

start-nimbus-1:
	rm -rf .tmp/nimbus-1
	docker run \
	--rm -v $(PWD)/${TESTAPP_DIR}:/${TESTAPP_DIR} \
	-v $(PWD)/.tmp:/.tmp \
	-p 30304:30303 \
	-p 8546:8545 \
	ethpandaops/nimbus-eth1:master \
	--data-dir=.tmp/nimbus-1 \
	--rpc

# start-nimbus-2:
# 	rm -rf .tmp/nimbus-2
# 	docker run \
# 	--rm -v $(PWD)/${TESTAPP_DIR}:/${TESTAPP_DIR} \
# 	-v $(PWD)/.tmp:/.tmp \
# 	-p 30303:30303 \
# 	-p 8545:8545 \
# 	statusim/nimbus-validator-client:multiarch-latest \
# 	--data-dir=.tmp/nimbus-2


start-ethereumjs:
	docker run \
	--rm -v $(PWD)/${TESTAPP_DIR}:/${TESTAPP_DIR} \
	-v $(PWD)/.tmp:/.tmp \
	-v $(PWD)/beacond/eth-genesis.json:/.tmp/beacond/eth-genesis.json \
	-p 30303:30303 \
	-p 8545:8545 \
	-p 8551:8551 \
	ethereumjs:nidhi \
	--gethGenesis /.tmp/beacond/eth-genesis.json \
	--rpcEngine \
	--rpc \
	--rpcEngineAddr 0.0.0.0 \
	--rpcEnginePort 8551 \
	--jwtSecret ../../$(JWT_PATH) \
	--logLevel debug \
	--rpcDebugVerbose \
	--rpcCors * 

# --startExecution \
# --startExecutionFrom 0


# --gethGenesis /usr/app/.tmp/beacond/eth-genesis.json \
#--dataDir .tmp/ethereumjs \
	ethpandaops/ethereumjs:master-46d09ca \
# start-ethereumjs:
# 	rm -rf .tmp/ethereumjs
# 	docker run \
# 	--rm -v $(PWD)/${TESTAPP_DIR}:/${TESTAPP_DIR} \
# 	-v $(PWD)/.tmp:/.tmp \
# 	-v $(PWD)/beacond/eth-genesis.json:/usr/app/.tmp/beacond/eth-genesis.json \
# 	-p 30303:30303 \
# 	-p 8545:8545 \
# 	-p 8551:8551 \
# 	ethpandaops/ethereumjs:stable \
# 	--dataDir .tmp/ethereumjs \
# 	--jwtSecret ../../$(JWT_PATH) \
# 	--gethGenesis /usr/app/.tmp/beacond/eth-genesis.json \
# 	--rpcEngine \
# 	--port 30303 \
# 	--ws \
# 	--rpc \
# 	--rpcApi eth,net,engine \
# 	--logLevel debug \
# 	--maxPeers 50 \
# 	--init /usr/app/.tmp/beacond/eth-genesis.json \
# 	--networkId 80087

	# -v $(PWD)/${ETH_GENESIS_PATH}:/.tmp/${ETH_GENESIS_PATH} \
	# --gethGenesis .tmp/${ETH_GENESIS_PATH} \
