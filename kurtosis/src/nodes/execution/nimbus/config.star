global_constants = import_module("../../../constants.star")
defaults = import_module("./../config.star")
port_spec_lib = import_module("../../../lib/port_spec.star")
shared_utils = import_module("github.com/ethpandaops/ethereum-package/src/shared_utils/shared_utils.star")

GLOBAL_LOG_LEVEL = global_constants.GLOBAL_LOG_LEVEL
KURTOSIS_IP_ADDRESS_PLACEHOLDER = global_constants.KURTOSIS_IP_ADDRESS_PLACEHOLDER

GLOBAL_FILES = []
ENTRYPOINT = ["nimbus"]
GENESIS_FILENAME = "genesis.json"

# USED_PORTS_TEMPLATE = defaults.USED_PORTS_TEMPLATE
WS_RPC_PORT_NUM = 8545
DISCOVERY_PORT_NUM = 30303
ENGINE_RPC_PORT_NUM = 8551
METRICS_PORT_NUM = 9001

# The min/max CPU/memory that the execution node can use
EXECUTION_MIN_CPU = 100
EXECUTION_MIN_MEMORY = 256

# Port IDs
WS_RPC_PORT_ID = "ws-rpc"
TCP_DISCOVERY_PORT_ID = "tcp-discovery"
UDP_DISCOVERY_PORT_ID = "udp-discovery"
ENGINE_RPC_PORT_ID = "engine-rpc"
METRICS_PORT_ID = "metrics"
RPC_PORT_ID = "eth-json-rpc"

VERBOSITY_LEVELS = {
    GLOBAL_LOG_LEVEL.error: "1",
    GLOBAL_LOG_LEVEL.warn: "2",
    GLOBAL_LOG_LEVEL.info: "3",
    GLOBAL_LOG_LEVEL.debug: "4",
    GLOBAL_LOG_LEVEL.trace: "5",
}

USED_PORTS_TEMPLATE = {
    RPC_PORT_ID: port_spec_lib.get_port_spec_template(defaults.RPC_PORT_NUM, shared_utils.TCP_PROTOCOL),
    # WS_PORT_ID: port_spec_lib.get_port_spec_template(defaults.WS_PORT_NUM, shared_utils.TCP_PROTOCOL),
    # WS_PORT_ENGINE_ID: port_spec_lib.get_port_spec_template(
    #     WS_PORT_ENGINE_NUM,
    #     shared_utils.TCP_PROTOCOL,
    # ),
    TCP_DISCOVERY_PORT_ID: port_spec_lib.get_port_spec_template(
        defaults.DISCOVERY_PORT_NUM,
        shared_utils.TCP_PROTOCOL,
    ),
    UDP_DISCOVERY_PORT_ID: port_spec_lib.get_port_spec_template(
        defaults.DISCOVERY_PORT_NUM,
        shared_utils.UDP_PROTOCOL,
    ),
    # ENGINE_RPC_PORT_ID: port_spec_lib.get_port_spec_template(
    #     defaults.ENGINE_RPC_PORT_NUM,
    #     shared_utils.TCP_PROTOCOL,
    # ),
}

# Paths
METRICS_PATH = "/metrics"
GENESIS_FILEPATH = "/app/genesis"
FILES = {
    "/app/genesis": "genesis_file",
    "/jwt": "jwt_file",
}

# FILES = {
#         "/root/.geth": NODE_CONFIG_ARTIFACT_NAME,
#         "/root/genesis": "genesis_file",
#         "/jwt": "jwt_file",
#     }
# The dirpath of the execution data directory on the client container
EXECUTION_DATA_DIRPATH_ON_CLIENT_CONTAINER = "/data/nimbus/execution-data"

BOOTNODE_CMD = "--bootstrap-node"
CMD = [
    "--data-dir=" + EXECUTION_DATA_DIRPATH_ON_CLIENT_CONTAINER,
    "--http-port={0}".format(WS_RPC_PORT_NUM),
    "--rpc",
    "--rpc-api=eth,debug,exp",
    "--ws",
    "--ws-api=eth,debug,exp",
    "--engine-api",
    "--engine-api-address=0.0.0.0",
    "--engine-api-port={0}".format(ENGINE_RPC_PORT_NUM),
    "--jwt-secret={0}".format(global_constants.JWT_MOUNT_PATH_ON_CONTAINER),
    "--metrics",
    "--metrics-address=0.0.0.0",
    "--metrics-port={0}".format(METRICS_PORT_NUM),
    "--nat=extip:{0}".format(KURTOSIS_IP_ADDRESS_PLACEHOLDER),
    "--tcp-port={0}".format(DISCOVERY_PORT_NUM),
    "--log-level={0}".format("DEBUG"),
    "--custom-network={}/{}".format(GENESIS_FILEPATH, "genesis.json"),
    "--engine-api-ws",
    "--allowed-origins=*",
    "--listen-address=0.0.0.0",
    "--http-address=0.0.0.0",
]

# "nimbus_eth1_max_mem": 16384,  # 16GB
# "nimbus_eth1_max_cpu": 4000,  # 4 cores
