global_constants = import_module("../../../constants.star")
defaults = import_module("./../config.star")

GLOBAL_LOG_LEVEL = global_constants.GLOBAL_LOG_LEVEL
KURTOSIS_IP_ADDRESS_PLACEHOLDER = global_constants.KURTOSIS_IP_ADDRESS_PLACEHOLDER

# The dirpath of the execution data directory on the client container
EXECUTION_DATA_DIRPATH_ON_CLIENT_CONTAINER = defaults.EXECUTION_DATA_DIRPATH_ON_CLIENT_CONTAINER

NODE_CONFIG_ARTIFACT_NAME = "reth-config"
CONFIG_FILENAME = "reth-config.toml"
GENESIS_FILENAME = "genesis.json"

# The files that only need to be uploaded once to be read by every node
# NOTE: THIS MUST REFERENCE THE FILEPATH RELATIVE TO execution.star
GLOBAL_FILES = [
    ("./reth/reth-config.toml", NODE_CONFIG_ARTIFACT_NAME),
]

WS_PORT_NUM = defaults.WS_PORT_NUM
ENGINE_RPC_PORT_NUM = defaults.ENGINE_RPC_PORT_NUM
METRICS_PORT_NUM = defaults.METRICS_PORT_NUM

METRICS_PATH = defaults.METRICS_PATH

ENTRYPOINT = ["sh", "-c"]
FILES = {
    "/root/.reth": NODE_CONFIG_ARTIFACT_NAME,
    "/root/genesis": "genesis_file",
    "/jwt": "jwt_file",
}
CMD = [
    "reth",
    "init",
    "--datadir",
    EXECUTION_DATA_DIRPATH_ON_CLIENT_CONTAINER,
    "--chain",
    "/root/genesis/{}".format(GENESIS_FILENAME),
    "&&",
    "reth",
    "node",
    "--datadir",
    EXECUTION_DATA_DIRPATH_ON_CLIENT_CONTAINER,
    "--chain",
    "/root/genesis/{}".format(GENESIS_FILENAME),
    "--http",
    "--http.addr",
    "0.0.0.0",
    "--http.corsdomain",
    "'*'",
    "--http.api",
    "admin,eth,net,web3,txpool,debug,trace",
    "--ws",
    "--ws.addr",
    "0.0.0.0",
    "--ws.port",
    str(WS_PORT_NUM),
    "--ws.api",
    "net,eth",
    "--ws.origins",
    "'*'",
    "--authrpc.port",
    str(ENGINE_RPC_PORT_NUM),
    "--authrpc.jwtsecret",
    global_constants.JWT_MOUNT_PATH_ON_CONTAINER,
    "--authrpc.addr",
    "0.0.0.0",
    "--metrics",
    "0.0.0.0:{0}".format(METRICS_PORT_NUM),
    # "--config", CONFIG_LOCATION,
    "--nat",
    "extip:" + KURTOSIS_IP_ADDRESS_PLACEHOLDER,
    "--builder.deadline",
    "2",
    "--builder.max-tasks",
    "20",
    "--txpool.pending-max-count",
    "100000",
    "--txpool.pending-max-size",
    "100",
    "--txpool.basefee-max-count",
    "100000",
    "--txpool.basefee-max-size",
    "100",
    "--txpool.queued-max-count",
    "100000",
    "--txpool.queued-max-size",
    "100",
    "--txpool.max-account-slots",
    "1000",
    "--txpool.max-cached-entries",
    "1000",
    "--rpc-cache.max-receipts",
    "10000",
]
BOOTNODE_CMD = "--bootnodes"
MAX_PEERS_OUTBOUND_CMD = "--max-outbound-peers"
MAX_PEERS_INBOUND_CMD = "--max-inbound-peers"

# Modify command flag --verbosity to change the verbosity level
VERBOSITY_LEVELS = {
    GLOBAL_LOG_LEVEL.error: "1",
    GLOBAL_LOG_LEVEL.warn: "2",
    GLOBAL_LOG_LEVEL.info: "3",
    GLOBAL_LOG_LEVEL.debug: "4",
    GLOBAL_LOG_LEVEL.trace: "5",
}

USED_PORTS = defaults.USED_PORTS
USED_PORTS_TEMPLATE = defaults.USED_PORTS_TEMPLATE

def set_max_peers(config, max_peers):
    cmdList = config["cmd"][:]
    cmdList.append(MAX_PEERS_OUTBOUND_CMD)
    cmdList.append(max_peers)
    cmdList.append(MAX_PEERS_INBOUND_CMD)
    cmdList.append(max_peers)
    config["cmd"] = cmdList
    return config
