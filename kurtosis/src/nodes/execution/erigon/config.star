global_constants = import_module("../../../constants.star")
defaults = import_module("./../config.star")

GLOBAL_LOG_LEVEL = global_constants.GLOBAL_LOG_LEVEL
KURTOSIS_IP_ADDRESS_PLACEHOLDER = global_constants.KURTOSIS_IP_ADDRESS_PLACEHOLDER

# The dirpath of the execution data directory on the client container
EXECUTION_DATA_DIRPATH_ON_CLIENT_CONTAINER = "/home/erigon{}".format(defaults.EXECUTION_DATA_DIRPATH_ON_CLIENT_CONTAINER)

# NODE_CONFIG_ARTIFACT_NAME = "erigon-config"
# CONFIG_FILENAME = "erigon-config.toml"
GENESIS_FILENAME = "genesis.json"

# The files that only need to be uploaded once to be read by every node
# NOTE: THIS MUST REFERENCE THE FILEPATH RELATIVE TO execution.star
GLOBAL_FILES = [
    # ("./erigon/erigon-config.toml", NODE_CONFIG_ARTIFACT_NAME)
]

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

GENESIS_FILEPATH = "/home/erigon/genesis"
ENTRYPOINT = ["sh", "-c"]

# CONFIG_LOCATION = "/root/.erigon/{}".format(CONFIG_FILENAME)
FILES = {
    # "/root/.erigon": NODE_CONFIG_ARTIFACT_NAME,
    GENESIS_FILEPATH: "genesis_file",
    "/jwt": "jwt_file",
}
CMD = [
    "erigon",
    "init",
    "--datadir",
    EXECUTION_DATA_DIRPATH_ON_CLIENT_CONTAINER,
    "{}/{}".format(GENESIS_FILEPATH, GENESIS_FILENAME),
    "&&",
    "erigon",
    # "--config", CONFIG_LOCATION,
    "--nat",
    "extip:" + KURTOSIS_IP_ADDRESS_PLACEHOLDER,
    "--metrics",
    "--datadir",
    EXECUTION_DATA_DIRPATH_ON_CLIENT_CONTAINER,
    "--port",
    str(DISCOVERY_PORT_NUM),
    "--http.api",
    "eth,erigon,engine,web3,net,debug,trace,txpool,admin,ots",
    "--http.vhosts",
    "'*'",
    "--ws",
    "--allow-insecure-unlock",
    "--nat",
    "extip:{}".format(KURTOSIS_IP_ADDRESS_PLACEHOLDER),
    "--http",
    "--http.addr",
    "0.0.0.0",
    "--http.corsdomain",
    "'*'",
    "--http.port",
    "{}".format(RPC_PORT_NUM),
    "--authrpc.jwtsecret",
    global_constants.JWT_MOUNT_PATH_ON_CONTAINER,
    "--authrpc.addr",
    "0.0.0.0",
    "--authrpc.port",
    "{}".format(ENGINE_RPC_PORT_NUM),
    "--authrpc.vhosts",
    "'*'",
    "--metrics",
    "--metrics.addr",
    "0.0.0.0",
    "--metrics.port",
    "{}".format(METRICS_PORT_NUM),
    "--networkid",
    "80087",
    "--db.size.limit={}MB".format(3000),
]
BOOTNODE_CMD = "--bootnodes"
MAX_PEERS_CMD = "--maxpeers"

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
    cmdList.append(MAX_PEERS_CMD)
    cmdList.append(max_peers)
    config["cmd"] = cmdList
    return config
