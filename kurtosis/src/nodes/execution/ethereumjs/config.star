defaults = import_module("./../config.star")
global_constants = import_module("../../../constants.star")
shared_utils = import_module("github.com/kurtosis-tech/ethereum-package/src/shared_utils/shared_utils.star")
port_spec_lib = import_module("../../../lib/port_spec.star")

GLOBAL_LOG_LEVEL = global_constants.GLOBAL_LOG_LEVEL
VERBOSITY_LEVELS = {
    GLOBAL_LOG_LEVEL.error: "error",
    GLOBAL_LOG_LEVEL.warn: "warn",
    GLOBAL_LOG_LEVEL.info: "info",
    GLOBAL_LOG_LEVEL.debug: "debug",
    GLOBAL_LOG_LEVEL.trace: "trace",
}

PRIVATE_IP_ADDRESS_PLACEHOLDER = "KURTOSIS_IP_ADDR_PLACEHOLDER"


GENESIS_FILEPATH = "/app/genesis"
GENESIS_DATA_MOUNTPOINT_ON_CLIENTS = "/network-configs"
GENESIS_CONFIG_MOUNT_PATH_ON_CONTAINER = (
    GENESIS_DATA_MOUNTPOINT_ON_CLIENTS + "/network-configs"
)
IMAGE = "ethpandaops/ethereumjs:master"
FILES = {
    "/app/genesis": "genesis_file",
    "/jwt": "jwt_file",
}

EXECUTION_DATA_DIRPATH_ON_CLIENT_CONTAINER = defaults.EXECUTION_DATA_DIRPATH_ON_CLIENT_CONTAINER

ENTRYPOINT = ["sh", "-c"]

BOOTNODE_CMD = "--bootnodes"
GLOBAL_FILES = []

METRICS_PATH = defaults.METRICS_PATH

# Port IDs
RPC_PORT_ID = "eth-json-rpc"
WS_PORT_ID = "eth-json-rpc-ws"
TCP_DISCOVERY_PORT_ID = "tcp-discovery"
UDP_DISCOVERY_PORT_ID = "udp-discovery"
ENGINE_RPC_PORT_ID = "engine-rpc"
ENGINE_WS_PORT_ID = "engineWs"
METRICS_PORT_ID = "metrics"
WS_PORT_ENGINE_ID = "ws-engine"

# ETH Execution constants

RPC_PORT_NUM = 8545
WS_PORT_NUM = 8546
DISCOVERY_PORT_NUM = 30303
ENGINE_RPC_PORT_NUM = 8551
WS_PORT_ENGINE_NUM = 8547


USED_PORTS_TEMPLATE = {
    RPC_PORT_ID: port_spec_lib.get_port_spec_template(RPC_PORT_NUM, shared_utils.TCP_PROTOCOL, shared_utils.HTTP_APPLICATION_PROTOCOL),
    WS_PORT_ID: port_spec_lib.get_port_spec_template(WS_PORT_NUM, shared_utils.TCP_PROTOCOL),
    WS_PORT_ENGINE_ID: port_spec_lib.get_port_spec_template(
        WS_PORT_ENGINE_NUM,
        shared_utils.TCP_PROTOCOL,
    ),
    TCP_DISCOVERY_PORT_ID: port_spec_lib.get_port_spec_template(
        DISCOVERY_PORT_NUM,
        shared_utils.TCP_PROTOCOL,
    ),
    UDP_DISCOVERY_PORT_ID: port_spec_lib.get_port_spec_template(
        DISCOVERY_PORT_NUM,
        shared_utils.UDP_PROTOCOL,
    ),
    ENGINE_RPC_PORT_ID: port_spec_lib.get_port_spec_template(
        ENGINE_RPC_PORT_NUM,
        shared_utils.TCP_PROTOCOL,
    ),
}

CMD = [
    "dumb-init",
    "node /usr/app/node_modules/.bin/ethereumjs",
    "--gethGenesis=/app/genesis/genesis.json",
    "--dataDir=" + EXECUTION_DATA_DIRPATH_ON_CLIENT_CONTAINER,
    "--port={0}".format(defaults.DISCOVERY_PORT_NUM),
    "--rpc",
    "--rpcAddr=0.0.0.0",
    "--rpcPort={0}".format(defaults.RPC_PORT_NUM),
    "--rpcEngine",
    "--rpcEngineAddr=0.0.0.0",
    "--rpcEnginePort={0}".format(defaults.ENGINE_RPC_PORT_NUM),
    "--ws",
    "--wsAddr=0.0.0.0",
    "--wsPort={0}".format(defaults.WS_PORT_NUM),
    "--wsEngineAddr=0.0.0.0",
    "--wsEnginePort={0}".format(WS_PORT_ENGINE_NUM),
    "--jwt-secret=" + global_constants.JWT_MOUNT_PATH_ON_CONTAINER,
    "--extIP={0}".format(PRIVATE_IP_ADDRESS_PLACEHOLDER),
    "--isSingleNode=true",
    "--logLevel={0}".format("debug"),
    "--rpcCors=*"
]