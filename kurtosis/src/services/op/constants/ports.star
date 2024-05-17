execution = import_module("../../../nodes/execution/execution.star")

# L1 Ports
L1_ETH_RPC = execution.RPC_PORT_NUM
L1_ETH_WS = execution.WS_PORT_NUM
L1_ENGINE_RPC = execution.ENGINE_RPC_PORT_NUM

# TODO: Make these configurable
GETH_RPC = "8545"
NODE_RPC = "8547"
BATCHER_RPC = "8548"
PROPOSER_RPC = "8560"
