global_constants = import_module("../../../constants.star")
defaults = import_module('./../constants.star')

GLOBAL_CLIENT_LOG_LEVEL = global_constants.GLOBAL_CLIENT_LOG_LEVEL
KURTOSIS_IP_ADDRESS_PLACEHOLDER = global_constants.KURTOSIS_IP_ADDRESS_PLACEHOLDER

# The dirpath of the execution data directory on the client container
EXECUTION_DATA_DIRPATH_ON_CLIENT_CONTAINER = defaults.EXECUTION_DATA_DIRPATH_ON_CLIENT_CONTAINER

# NODE_CONFIG_ARTIFACT_NAME = "nimbus-config"
# CONFIG_FILENAME = "nimbus-config.toml"
GENESIS_FILENAME = "genesis.json"
# The files that only need to be uploaded once to be read by every node
# NOTE: THIS MUST REFERENCE THE FILEPATH RELATIVE TO execution.star
GLOBAL_FILES = [
    # ("./nimbus/nimbus-config.toml", NODE_CONFIG_ARTIFACT_NAME)
]

RPC_PORT_NUM = defaults.RPC_PORT_NUM
WS_PORT_NUM = RPC_PORT_NUM
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

GENESIS_FILEPATH = "/home/nimbus/genesis"
IMAGE = "ethpandaops/nimbus-eth1:master"
ENTRYPOINT = ["sh", "-c"]
# CONFIG_LOCATION = "/root/.nimbus/{}".format(CONFIG_FILENAME)
FILES = {
    # "/root/.nimbus": NODE_CONFIG_ARTIFACT_NAME,
    GENESIS_FILEPATH: "genesis_file",
    "/jwt": "jwt_file",
}
CMD = [
    "nimbus",
    "--log-level=TRACE",
    "--data-dir={}".format(EXECUTION_DATA_DIRPATH_ON_CLIENT_CONTAINER),
    "--custom-network={}/{}".format(GENESIS_FILEPATH, GENESIS_FILENAME),
    # "--config", CONFIG_LOCATION,
    "--nat=extip:{}".format(KURTOSIS_IP_ADDRESS_PLACEHOLDER),
    "--network=80087",
    "--http-port={}".format(str(RPC_PORT_NUM)),
    "--http-address={}".format("0.0.0.0"),
    "--rpc",
    "--rpc-api=eth,debug,exp",
    # "--ws",
    "--engine-api",
    "--engine-api-port={}".format(str(ENGINE_RPC_PORT_NUM)),
    "--engine-api-address=0.0.0.0",
    "--jwt-secret={}".format(global_constants.JWT_MOUNT_PATH_ON_CONTAINER),
    "--metrics",
    "--metrics-port={}".format(str(METRICS_PORT_NUM)),
    "--metrics-address=0.0.0.0",
    # "--sync-mode", "full"

]
BOOTNODE_CMD = "--bootstrap-node="



# Modify command flag --verbosity to change the verbosity level
VERBOSITY_LEVELS = {
    GLOBAL_CLIENT_LOG_LEVEL.error: "1",
    GLOBAL_CLIENT_LOG_LEVEL.warn: "2",
    GLOBAL_CLIENT_LOG_LEVEL.info: "3",
    GLOBAL_CLIENT_LOG_LEVEL.debug: "4",
    GLOBAL_CLIENT_LOG_LEVEL.trace: "5",
}

USED_PORTS = defaults.USED_PORTS
USED_PORTS_TEMPLATE = {
    RPC_PORT_ID: defaults.USED_PORTS_TEMPLATE[RPC_PORT_ID],
    # WS_PORT_ID: defaults.USED_PORTS_TEMPLATE[RPC_PORT_ID],
    TCP_DISCOVERY_PORT_ID: defaults.USED_PORTS_TEMPLATE[TCP_DISCOVERY_PORT_ID],
    UDP_DISCOVERY_PORT_ID: defaults.USED_PORTS_TEMPLATE[UDP_DISCOVERY_PORT_ID],
    ENGINE_RPC_PORT_ID: defaults.USED_PORTS_TEMPLATE[ENGINE_RPC_PORT_ID],
    # METRICS_PORT_ID: defaults.USED_PORTS_TEMPLATE[METRICS_PORT_ID],
}