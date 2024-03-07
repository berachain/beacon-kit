shared_utils = import_module("github.com/kurtosis-tech/ethereum-package/src/shared_utils/shared_utils.star")
constants = import_module("github.com/kurtosis-tech/ethereum-package/src/package_io/constants.star")
el_client_context = import_module("github.com/kurtosis-tech/ethereum-package/src/el/el_client_context.star")

service_config_lib = import_module("../../../lib/service_config.star")
port_spec_lib = import_module("../../../lib/port_spec.star")
builtins = import_module("../../../lib/builtins.star")

# The dirpath of the execution data directory on the client container
EXECUTION_DATA_DIRPATH_ON_CLIENT_CONTAINER = "/data/nethermind/execution-data"

METRICS_PATH = "/metrics"

RPC_PORT_NUM = 8545
WS_PORT_NUM = 8546
DISCOVERY_PORT_NUM = 30303
ENGINE_RPC_PORT_NUM = 8551
METRICS_PORT_NUM = 9001

# Port IDs
RPC_PORT_ID = "rpc"
WS_PORT_ID = "ws"
TCP_DISCOVERY_PORT_ID = "tcp-discovery"
UDP_DISCOVERY_PORT_ID = "udp-discovery"
ENGINE_RPC_PORT_ID = "engine-rpc"
METRICS_PORT_ID = "metrics"

PRIVATE_IP_ADDRESS_PLACEHOLDER = "KURTOSIS_IP_ADDR_PLACEHOLDER"

USED_PORTS = {
    RPC_PORT_ID: shared_utils.new_port_spec(RPC_PORT_NUM, shared_utils.TCP_PROTOCOL),
    WS_PORT_ID: shared_utils.new_port_spec(WS_PORT_NUM, shared_utils.TCP_PROTOCOL),
    TCP_DISCOVERY_PORT_ID: shared_utils.new_port_spec(
        DISCOVERY_PORT_NUM,
        shared_utils.TCP_PROTOCOL,
    ),
    UDP_DISCOVERY_PORT_ID: shared_utils.new_port_spec(
        DISCOVERY_PORT_NUM,
        shared_utils.UDP_PROTOCOL,
    ),
    ENGINE_RPC_PORT_ID: shared_utils.new_port_spec(
        ENGINE_RPC_PORT_NUM,
        shared_utils.TCP_PROTOCOL,
    ),
    METRICS_PORT_ID: shared_utils.new_port_spec(
        METRICS_PORT_NUM,
        shared_utils.TCP_PROTOCOL,
    ),
}

USED_PORTS_TEMPLATE = {
    RPC_PORT_ID: port_spec_lib.get_port_spec_template(RPC_PORT_NUM, shared_utils.TCP_PROTOCOL),
    WS_PORT_ID: port_spec_lib.get_port_spec_template(WS_PORT_NUM, shared_utils.TCP_PROTOCOL),
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
    # METRICS_PORT_ID: port_spec_lib.get_port_spec_template(
    #     METRICS_PORT_NUM, shared_utils.TCP_PROTOCOL
    # ),
}

VERBOSITY_LEVELS = {
    constants.GLOBAL_CLIENT_LOG_LEVEL.error: "ERROR",
    constants.GLOBAL_CLIENT_LOG_LEVEL.warn: "WARN",
    constants.GLOBAL_CLIENT_LOG_LEVEL.info: "INFO",
    constants.GLOBAL_CLIENT_LOG_LEVEL.debug: "DEBUG",
    constants.GLOBAL_CLIENT_LOG_LEVEL.trace: "TRACE",
}

DEFAULT_IMAGE = "nethermindeth/nethermind:release-1.25.4"
DEFAULT_ENTRYPOINT_ARGS = ["sh", "-c"]
DEFAULT_CONFIG_LOCATION = "/root/.nethermind/nethermind-config.cfg"
DEFAULT_CMD = ["./Nethermind.Runner", "--config", DEFAULT_CONFIG_LOCATION, "--Network.ExternalIp", PRIVATE_IP_ADDRESS_PLACEHOLDER]
DEFAULT_FILES = {
    "/root/.nethermind": "nethermind-config",
    "/root/genesis": "nether_genesis_file",
    "/jwt": "jwt_file",
}

# Because structs are immutable, we pass around a map to allow full modification up until we create the final ServiceConfig
def get_default_service_config(service_name):
    sc = service_config_lib.get_service_config_template(service_name, DEFAULT_IMAGE, ports = USED_PORTS_TEMPLATE, entrypoint = DEFAULT_ENTRYPOINT_ARGS, cmd = DEFAULT_CMD, files = DEFAULT_FILES)

    return sc

# Uploads files that all geth nodes use
def upload_global_files(plan):
    artifact_names = []

    nethermind_genesis_file = plan.upload_files(
        src = "../../../networks/kurtosis-devnet/network-configs/nethermind-genesis.json",
        name = "nether_genesis_file",
    )
    nethermind_config_artifact = plan.upload_files(
        src = "./nethermind-config.cfg",
        name = "nethermind-config",
    )
    artifact_names.append(nethermind_config_artifact)

    return artifact_names

def add_bootnodes(config, bootnodes):
    if type(bootnodes) == builtins.types.list:
        if len(bootnodes) > 0:
            cmdList = config["cmd"][:]
            cmdList.append("--Network.Bootnodes")
            config["cmd"] = cmdList

            bootnodes_str = ",".join(bootnodes)
            config["cmd"].append(bootnodes_str)
    elif type(bootnodes) == builtins.types.str:
        if len(bootnodes) > 0:
            config["cmd"].append("--Network.Bootnodes")
            config["cmd"].append(bootnodes)
    else:
        fail("Bootnodes was not a list or string, but instead a {}", type(bootnodes))

    return config

def deploy_node(plan, config):
    service_config = service_config_lib.create_from_config(config)

    plan.add_service(
        name = config["name"],
        config = service_config,
    )
