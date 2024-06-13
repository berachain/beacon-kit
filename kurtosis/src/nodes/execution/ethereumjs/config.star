defaults = import_module("./../config.star")
global_constants = import_module("../../../constants.star")
shared_utils = import_module("github.com/ethpandaops/ethereum-package/src/shared_utils/shared_utils.star")
port_spec_lib = import_module("../../../lib/port_spec.star")

GLOBAL_LOG_LEVEL = global_constants.GLOBAL_LOG_LEVEL
VERBOSITY_LEVELS = {
    GLOBAL_LOG_LEVEL.error: "error",
    GLOBAL_LOG_LEVEL.warn: "warn",
    GLOBAL_LOG_LEVEL.info: "info",
    GLOBAL_LOG_LEVEL.debug: "debug",
}

PRIVATE_IP_ADDRESS_PLACEHOLDER = "KURTOSIS_IP_ADDR_PLACEHOLDER"

GENESIS_FILEPATH = "/app/genesis"
FILES = {
    "/app/genesis": "genesis_file",
    "/jwt": "jwt_file",
}

ENTRYPOINT = ["sh", "-c"]

BOOTNODE_CMD = "--bootnodes"
GLOBAL_FILES = []

# Port IDs
RPC_PORT_ID = "eth-json-rpc"
WS_PORT_ID = "eth-json-rpc-ws"
TCP_DISCOVERY_PORT_ID = "tcp-discovery"
UDP_DISCOVERY_PORT_ID = "udp-discovery"
ENGINE_RPC_PORT_ID = "engine-rpc"
ENGINE_WS_PORT_ID = "engineWs"
WS_PORT_ENGINE_ID = "ws-engine"

WS_PORT_ENGINE_NUM = 8547

# Redefining it here as we don't want metrics ports in here
USED_PORTS_TEMPLATE = {
    RPC_PORT_ID: port_spec_lib.get_port_spec_template(defaults.RPC_PORT_NUM, shared_utils.TCP_PROTOCOL),
    WS_PORT_ID: port_spec_lib.get_port_spec_template(defaults.WS_PORT_NUM, shared_utils.TCP_PROTOCOL),
    WS_PORT_ENGINE_ID: port_spec_lib.get_port_spec_template(
        WS_PORT_ENGINE_NUM,
        shared_utils.TCP_PROTOCOL,
    ),
    TCP_DISCOVERY_PORT_ID: port_spec_lib.get_port_spec_template(
        defaults.DISCOVERY_PORT_NUM,
        shared_utils.TCP_PROTOCOL,
    ),
    UDP_DISCOVERY_PORT_ID: port_spec_lib.get_port_spec_template(
        defaults.DISCOVERY_PORT_NUM,
        shared_utils.UDP_PROTOCOL,
    ),
    ENGINE_RPC_PORT_ID: port_spec_lib.get_port_spec_template(
        defaults.ENGINE_RPC_PORT_NUM,
        shared_utils.TCP_PROTOCOL,
    ),
}

CMD = [
    "dumb-init",
    "node /usr/app/node_modules/.bin/ethereumjs",
    "--gethGenesis={}/{}".format(GENESIS_FILEPATH, "genesis.json"),
    "--dataDir=" + defaults.EXECUTION_DATA_DIRPATH_ON_CLIENT_CONTAINER,
    "--port={0}".format(defaults.DISCOVERY_PORT_NUM),
    "--rpc",
    "--rpcAddr=0.0.0.0",
    "--rpcPort={0}".format(defaults.RPC_PORT_NUM),
    "--rpcCors=*",
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
    "--sync=full",
]
