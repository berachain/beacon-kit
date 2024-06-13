shared_utils = import_module("github.com/ethpandaops/ethereum-package/src/shared_utils/shared_utils.star")
global_constants = import_module("../../../constants.star")
defaults = import_module("./../config.star")
port_spec_lib = import_module("../../../lib/port_spec.star")

GLOBAL_LOG_LEVEL = global_constants.GLOBAL_LOG_LEVEL
KURTOSIS_IP_ADDRESS_PLACEHOLDER = global_constants.KURTOSIS_IP_ADDRESS_PLACEHOLDER

# The dirpath of the execution data directory on the client container
EXECUTION_DATA_DIRPATH_ON_CLIENT_CONTAINER = defaults.EXECUTION_DATA_DIRPATH_ON_CLIENT_CONTAINER

NODE_TYPE = "nethermind"

NODE_CONFIG_ARTIFACT_NAME = "{}-config".format(NODE_TYPE)
CONFIG_FILENAME = "{}-config.cfg".format(NODE_TYPE)
GENESIS_FILENAME = "genesis.json"
CONFIG_FOLDER = "/root/.{}".format(NODE_TYPE)

# The files that only need to be uploaded once to be read by every node
# NOTE: THIS MUST REFERENCE THE FILEPATH RELATIVE TO execution.star
GLOBAL_FILES = [
    ("./{}/{}".format(NODE_TYPE, CONFIG_FILENAME), NODE_CONFIG_ARTIFACT_NAME),
    ("./{}/{}".format(NODE_TYPE, GENESIS_FILENAME), "nether_genesis_file"),
]

ENTRYPOINT = ["sh", "-c"]
CONFIG_LOCATION = "{}/{}".format(CONFIG_FOLDER, CONFIG_FILENAME)
FILES = {
    CONFIG_FOLDER: NODE_CONFIG_ARTIFACT_NAME,
    "/root/genesis": "nether_genesis_file",
    "/jwt": "jwt_file",
}
CMD = [
    "./nethermind",
    "--config",
    CONFIG_LOCATION,
    "--datadir",
    EXECUTION_DATA_DIRPATH_ON_CLIENT_CONTAINER,
    "--Network.ExternalIp",
    KURTOSIS_IP_ADDRESS_PLACEHOLDER,
    "--Merge.Enabled",
    "true",
]

BOOTNODE_CMD = "--Network.Bootnodes"
MAX_PEERS_CMD = "--Network.MaxActivePeers"

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

# Modify command flag --verbosity to change the verbosity level
VERBOSITY_LEVELS = {
    GLOBAL_LOG_LEVEL.error: "ERROR",
    GLOBAL_LOG_LEVEL.warn: "WARN",
    GLOBAL_LOG_LEVEL.info: "INFO",
    GLOBAL_LOG_LEVEL.debug: "DEBUG",
    GLOBAL_LOG_LEVEL.trace: "TRACE",
}

USED_PORTS = defaults.USED_PORTS
USED_PORTS_TEMPLATE = {
    RPC_PORT_ID: port_spec_lib.get_port_spec_template(RPC_PORT_NUM, shared_utils.TCP_PROTOCOL, wait = "30s"),
    WS_PORT_ID: port_spec_lib.get_port_spec_template(WS_PORT_NUM, shared_utils.TCP_PROTOCOL, wait = "30s"),
    TCP_DISCOVERY_PORT_ID: port_spec_lib.get_port_spec_template(
        DISCOVERY_PORT_NUM,
        shared_utils.TCP_PROTOCOL,
        wait = "30s",
    ),
    UDP_DISCOVERY_PORT_ID: port_spec_lib.get_port_spec_template(
        DISCOVERY_PORT_NUM,
        shared_utils.UDP_PROTOCOL,
        wait = "30s",
    ),
    ENGINE_RPC_PORT_ID: port_spec_lib.get_port_spec_template(
        ENGINE_RPC_PORT_NUM,
        shared_utils.TCP_PROTOCOL,
        wait = "30s",
    ),
    METRICS_PORT_ID: port_spec_lib.get_port_spec_template(
        METRICS_PORT_NUM,
        shared_utils.TCP_PROTOCOL,
    ),
}

def set_max_peers(config, max_peers):
    cmdList = config["cmd"][:]
    cmdList.append(MAX_PEERS_CMD)
    cmdList.append(max_peers)
    config["cmd"] = cmdList
    return config
