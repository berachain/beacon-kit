global_constants = import_module("../../../constants.star")
defaults = import_module("./../config.star")

GLOBAL_LOG_LEVEL = global_constants.GLOBAL_LOG_LEVEL
KURTOSIS_IP_ADDRESS_PLACEHOLDER = global_constants.KURTOSIS_IP_ADDRESS_PLACEHOLDER

# The dirpath of the execution data directory on the client container
EXECUTION_DATA_DIRPATH_ON_CLIENT_CONTAINER = defaults.EXECUTION_DATA_DIRPATH_ON_CLIENT_CONTAINER

NODE_CONFIG_ARTIFACT_NAME = "besu-config"
CONFIG_FILENAME = "besu-config.toml"
GENESIS_FILENAME = "genesis.json"

# The files that only need to be uploaded once to be read by every node
# NOTE: THIS MUST REFERENCE THE FILEPATH RELATIVE TO execution.star
GLOBAL_FILES = []

RPC_PORT_NUM = defaults.RPC_PORT_NUM
WS_PORT_NUM = defaults.WS_PORT_NUM
DISCOVERY_PORT_NUM = defaults.DISCOVERY_PORT_NUM
ENGINE_RPC_PORT_NUM = defaults.ENGINE_RPC_PORT_NUM
METRICS_PORT_NUM = defaults.METRICS_PORT_NUM

# Port IDs
RPC_PORT_ID = defaults.RPC_PORT_ID
WS_PORT_ID = defaults.WS_PORT_ID
TCP_DISCOVERY_PORT_ID = defaults.TCP_DISCOVERY_PORT_ID
UDP_DISCOVERY_PORT_ID = defaults.UDP_DISCOVERY_PORT_ID
ENGINE_RPC_PORT_ID = defaults.ENGINE_RPC_PORT_ID
ENGINE_WS_PORT_ID = defaults.ENGINE_WS_PORT_ID
METRICS_PORT_ID = defaults.METRICS_PORT_ID

METRICS_PATH = defaults.METRICS_PATH

ENTRYPOINT = ["sh", "-c"]

# CONFIG_LOCATION = "/root/.geth/{}".format(CONFIG_FILENAME)
FILES = {
    "/jwt": "jwt_file",
}
CMD = [
    "besu",
    "--genesis-file={0}".format("/app/genesis/{}".format(GENESIS_FILENAME)),
    "--host-allowlist=*",
    "--rpc-http-enabled=true",
    "--rpc-http-host=0.0.0.0",
    "--rpc-http-port={0}".format(RPC_PORT_NUM),
    "--rpc-http-api=ADMIN,CLIQUE,ETH,NET,DEBUG,TXPOOL,ENGINE,TRACE,WEB3",
    "--rpc-http-cors-origins=*",
    "--rpc-ws-enabled=true",
    "--rpc-ws-host=0.0.0.0",
    "--rpc-ws-port={0}".format(WS_PORT_NUM),
    "--rpc-ws-api=ADMIN,CLIQUE,ETH,NET,DEBUG,TXPOOL,ENGINE,TRACE,WEB3",
    "--p2p-enabled=true",
    "--p2p-host=" + KURTOSIS_IP_ADDRESS_PLACEHOLDER,
    "--p2p-port={0}".format(DISCOVERY_PORT_NUM),
    "--engine-rpc-enabled=true",
    "--engine-jwt-secret=" + global_constants.JWT_MOUNT_PATH_ON_CONTAINER,
    "--engine-host-allowlist=*",
    "--engine-rpc-port={0}".format(ENGINE_RPC_PORT_NUM),
    "--sync-mode=FULL",
    "--data-storage-format=BONSAI",
    "--metrics-enabled=true",
    "--metrics-host=0.0.0.0",
    "--metrics-port={0}".format(METRICS_PORT_NUM),
]
BOOTNODE_CMD = "--bootnodes"
MAX_PEERS_CMD = "--max-peers"

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
    cmdList.append("{}={}".format(MAX_PEERS_CMD, max_peers))
    config["cmd"] = cmdList
    return config
