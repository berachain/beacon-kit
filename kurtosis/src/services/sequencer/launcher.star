# Sequencer service launcher for preconfirmation testing.
# This service runs bera-sequencer as the execution client for the sequencer node.

shared_utils = import_module("github.com/ethpandaops/ethereum-package/src/shared_utils/shared_utils.star")

SERVICE_NAME = "sequencer"

# Ports
AUTHRPC_PORT_ID = "authrpc"
AUTHRPC_PORT_NUMBER = 8551

HTTP_PORT_ID = "http"
HTTP_PORT_NUMBER = 8545

WS_PORT_ID = "ws"
WS_PORT_NUMBER = 8546

PRECONF_PORT_ID = "preconf"
PRECONF_PORT_NUMBER = 9090

# Default image for bera-sequencer
DEFAULT_IMAGE = "ghcr.io/berachain/bera-sequencer:latest"

# Resource limits
MIN_CPU = 1000
MAX_CPU = 4000
MIN_MEMORY = 2048
MAX_MEMORY = 8192

USED_PORTS = {
    AUTHRPC_PORT_ID: shared_utils.new_port_spec(
        AUTHRPC_PORT_NUMBER,
        shared_utils.TCP_PROTOCOL,
    ),
    HTTP_PORT_ID: shared_utils.new_port_spec(
        HTTP_PORT_NUMBER,
        shared_utils.TCP_PROTOCOL,
        shared_utils.HTTP_APPLICATION_PROTOCOL,
    ),
    WS_PORT_ID: shared_utils.new_port_spec(
        WS_PORT_NUMBER,
        shared_utils.TCP_PROTOCOL,
    ),
    PRECONF_PORT_ID: shared_utils.new_port_spec(
        PRECONF_PORT_NUMBER,
        shared_utils.TCP_PROTOCOL,
        shared_utils.HTTP_APPLICATION_PROTOCOL,
    ),
}

def launch_sequencer(
        plan,
        jwt_file,
        genesis_files,
        chain_id,
        image = DEFAULT_IMAGE,
        extra_args = []):
    """
    Launch the bera-sequencer service for preconfirmation testing.

    Args:
        plan: The Kurtosis plan object.
        jwt_file: The JWT secret file artifact for Engine API authentication.
        genesis_files: The genesis files artifact for chain initialization.
        chain_id: The chain ID for the network.
        image: The Docker image to use for the sequencer.
        extra_args: Additional command-line arguments for the sequencer.

    Returns:
        The service context for the launched sequencer.
    """
    config = get_config(
        jwt_file,
        genesis_files,
        chain_id,
        image,
        extra_args,
    )

    return plan.add_service(SERVICE_NAME, config)

def get_config(
        jwt_file,
        genesis_files,
        chain_id,
        image,
        extra_args):
    """
    Generate the service configuration for the sequencer.
    """
    cmd = [
        "node",
        "--datadir=/data",
        "--authrpc.addr=0.0.0.0",
        "--authrpc.port={}".format(AUTHRPC_PORT_NUMBER),
        "--authrpc.vhosts=*",
        "--authrpc.jwtsecret=/jwt/jwt-secret.hex",
        "--http",
        "--http.addr=0.0.0.0",
        "--http.port={}".format(HTTP_PORT_NUMBER),
        "--http.api=eth,net,web3,txpool",
        "--http.corsdomain=*",
        "--http.vhosts=*",
        "--ws",
        "--ws.addr=0.0.0.0",
        "--ws.port={}".format(WS_PORT_NUMBER),
        "--ws.api=eth,net,web3,txpool",
        "--ws.origins=*",
        "--chain.id={}".format(chain_id),
    ]

    # Add extra arguments
    cmd.extend(extra_args)

    return ServiceConfig(
        image = image,
        ports = USED_PORTS,
        files = {
            "/jwt": jwt_file,
            "/genesis": genesis_files,
        },
        cmd = cmd,
        min_cpu = MIN_CPU,
        max_cpu = MAX_CPU,
        min_memory = MIN_MEMORY,
        max_memory = MAX_MEMORY,
        ready_conditions = ReadyCondition(
            recipe = GetHttpRequestRecipe(
                port_id = HTTP_PORT_ID,
                endpoint = "/",
                body = '{"jsonrpc":"2.0","method":"eth_chainId","params":[],"id":1}',
                content_type = "application/json",
            ),
            field = "code",
            assertion = "==",
            target_value = 200,
        ),
    )

def get_sequencer_url(service_context):
    """
    Get the preconf API URL for the sequencer.

    Args:
        service_context: The service context returned from launch_sequencer.

    Returns:
        The URL for connecting to the sequencer's preconf API.
    """
    return "http://{}:{}".format(
        service_context.ip_address,
        service_context.ports[PRECONF_PORT_ID].number,
    )

def get_engine_url(service_context):
    """
    Get the Engine API URL for the sequencer.

    Args:
        service_context: The service context returned from launch_sequencer.

    Returns:
        The URL for connecting to the sequencer's Engine API.
    """
    return "http://{}:{}".format(
        service_context.ip_address,
        service_context.ports[AUTHRPC_PORT_ID].number,
    )
