shared_utils = import_module("github.com/kurtosis-tech/ethereum-package/src/shared_utils/shared_utils.star")

COMETBFT_RPC_PORT_NUM = 26657
COMETBFT_P2P_PORT_NUM = 26656
COMETBFT_GRPC_PORT_NUM = 9090
COMETBFT_REST_PORT_NUM = 1317
PROMETHEUS_PORT_NUM = 26660
ENGINE_RPC_PORT_NUM = 8551

# Port IDs
COMETBFT_RPC_PORT_ID = "cometbft-rpc"
COMETBFT_P2P_PORT_ID = "cometbft-p2p"
COMETBFT_GRPC_PORT_ID = "cometbft-grpc"
COMETBFT_REST_PORT_ID = "cometbft-rest"
ENGINE_RPC_PORT_ID = "engine-rpc"
PROMETHEUS_PORT_ID = "prometheus"

USED_PORTS = {
    COMETBFT_RPC_PORT_ID: shared_utils.new_port_spec(COMETBFT_RPC_PORT_NUM, shared_utils.TCP_PROTOCOL),
    COMETBFT_P2P_PORT_ID: shared_utils.new_port_spec(COMETBFT_P2P_PORT_NUM, shared_utils.TCP_PROTOCOL),
    COMETBFT_GRPC_PORT_ID: shared_utils.new_port_spec(COMETBFT_GRPC_PORT_NUM, shared_utils.TCP_PROTOCOL),
    COMETBFT_REST_PORT_ID: shared_utils.new_port_spec(COMETBFT_REST_PORT_NUM, shared_utils.TCP_PROTOCOL),
    # ENGINE_RPC_PORT_ID: shared_utils.new_port_spec(ENGINE_RPC_PORT_NUM, shared_utils.TCP_PROTOCOL),
    PROMETHEUS_PORT_ID: shared_utils.new_port_spec(PROMETHEUS_PORT_NUM, shared_utils.TCP_PROTOCOL, wait = None),
}

def get_config(image, jwt_file, engine_dial_url, cl_service_name, entrypoint = [], cmd = [], persistent_peers = "", expose_ports = True):
    exposed_ports = {}
    if expose_ports:
        exposed_ports = USED_PORTS

    config = ServiceConfig(
        image = image,
        files = {
            "/root/app": jwt_file,
        },
        entrypoint = entrypoint,
        cmd = cmd,
        env_vars = {
            "BEACOND_MONIKER": cl_service_name,
            "BEACOND_NET": "VALUE_2",
            "BEACOND_HOME": "/root/.beacond",
            "BEACOND_CHAIN_ID": "beacon-kurtosis-80087",
            "BEACOND_DEBUG": "false",
            "BEACOND_KEYRING_BACKEND": "test",
            "BEACOND_MINIMUM_GAS_PRICE": "0abgt",
            "BEACOND_ENGINE_DIAL_URL": engine_dial_url,
            "BEACOND_ETH_CHAIN_ID": "80087",
            "BEACOND_PERSISTENT_PEERS": persistent_peers,
            "BEACOND_ENABLE_PROMETHEUS": "true",
        },
        ports = exposed_ports,
    )

    return config
