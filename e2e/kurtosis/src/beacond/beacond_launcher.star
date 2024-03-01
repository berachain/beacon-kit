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
    # COMETBFT_RPC_PORT_ID: shared_utils.new_port_spec(COMETBFT_RPC_PORT_NUM, shared_utils.TCP_PROTOCOL),
    COMETBFT_P2P_PORT_ID: shared_utils.new_port_spec(COMETBFT_P2P_PORT_NUM, shared_utils.TCP_PROTOCOL),
    # COMETBFT_GRPC_PORT_ID: shared_utils.new_port_spec(COMETBFT_GRPC_PORT_NUM, shared_utils.TCP_PROTOCOL),
    # COMETBFT_REST_PORT_ID: shared_utils.new_port_spec(COMETBFT_REST_PORT_NUM, shared_utils.TCP_PROTOCOL),
    # ENGINE_RPC_PORT_ID: shared_utils.new_port_spec(ENGINE_RPC_PORT_NUM, shared_utils.TCP_PROTOCOL),
    # PROMETHEUS_PORT_ID: shared_utils.new_port_spec(PROMETHEUS_PORT_NUM, shared_utils.TCP_PROTOCOL),
}

def build_beacond_docker_image():
    # temp workaround to get it working
    # run `make docker-build` before running kurtosis
    # make sure the image name and tag match, otherwise edit below
    image = "beacond:0.0.3-alpha-19-g1c02a0c5"

    # image = ImageBuildSpec(
    #     image_name="berachain/beacond",
    #     build_context_dir="examples/beacond"
    # )

    return image


def get_config(jwt_file, engine_dial_url):
    config = ServiceConfig(
        image=build_beacond_docker_image(),
        files = {
            "/root/app": jwt_file
        },
        entrypoint = [
            "bash",
        ],
        cmd=[
            "-c",
            "/usr/bin/init.sh",
        ],
        env_vars = {
            "BEACOND_MONIKER": "kurtosis",
            "BEACOND_NET": "VALUE_2",
            "BEACOND_HOME": "/root/.beacond",
            "BEACOND_CHAIN_ID": "beacon-kurtosis-80087",
            "BEACOND_DEBUG": "false",
            "BEACOND_KEYRING_BACKEND": "test",
            "BEACOND_MINIMUM_GAS_PRICE": "0stake",
            "BEACOND_ENGINE_DIAL_URL": engine_dial_url,
            "BEACOND_ETH_CHAIN_ID": "80087",
        },
        ports=USED_PORTS,
    )

    return config

def new_beacond_launcher():
    print("Launching new beacond")