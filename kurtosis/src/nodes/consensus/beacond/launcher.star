shared_utils = import_module("github.com/kurtosis-tech/ethereum-package/src/shared_utils/shared_utils.star")
execution = import_module("../../execution/execution.star")

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

def get_config(image, engine_dial_url, cl_service_name, entrypoint = [], cmd = [], persistent_peers = "", expose_ports = True, jwt_file = None, kzg_trusted_setup_file = None):
    exposed_ports = {}
    if expose_ports:
        exposed_ports = USED_PORTS

    files = {}
    if jwt_file:
        files["/root/jwt"] = jwt_file
    if kzg_trusted_setup_file:
        files["/root/kzg"] = kzg_trusted_setup_file

    config = ServiceConfig(
        image = image,
        files = files,
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

def perform_genesis_ceremony(plan, validators, jwt_file):
    num_validators = len(validators)

    node_peering_info = []
    for n, validator in enumerate(validators):
        cl_service_name = "cl-validator-beaconkit-{}".format(n)
        engine_dial_url = ""  # not needed for this step
        beacond_config = get_config(
            validator.cl_image,
            engine_dial_url,
            cl_service_name,
            expose_ports = False,
            jwt_file = jwt_file,
        )

        if n > 0:
            beacond_config.files["/root/.beacond/config"] = Directory(
                artifact_names = ["cosmos-genesis-{}".format(n - 1)],
            )

        if n == num_validators - 1 and n != 0:
            collected_gentx = []
            for other_validator_id in range(num_validators - 1):
                collected_gentx.append("cosmos-gentx-{}".format(other_validator_id))

            beacond_config.files["/root/.beacond/config/gentx"] = Directory(
                artifact_names = collected_gentx,
            )

        plan.add_service(
            name = cl_service_name,
            config = beacond_config,
        )

        exec_recipe = None
        if n == 0:
            exec_recipe = ExecRecipe(
                # Initialize the Cosmos genesis file
                command = ["/usr/bin/init.sh", "-f"],
            )
        else:
            exec_recipe = ExecRecipe(
                # Initialize the Cosmos genesis file
                command = ["/usr/bin/init.sh"],
            )

        plan.exec(
            service_name = cl_service_name,
            recipe = exec_recipe,
        )

        peer_result = plan.exec(
            service_name = cl_service_name,
            recipe = ExecRecipe(
                command = ["bash", "-c", "/usr/bin/beacond comet show-node-id --home $BEACOND_HOME | tr -d '\n'"],
            ),
        )

        node_peering_info.append(peer_result["output"])

        file_suffix = "{}".format(n)
        if n == num_validators - 1:
            finalize_recipe = ExecRecipe(
                # Initialize the Cosmos genesis file
                command = ["/usr/bin/finalize.sh"],
            )
            result = plan.exec(
                service_name = cl_service_name,
                recipe = finalize_recipe,
            )
            file_suffix = "final"

        node_beacond_config = plan.store_service_files(
            service_name = cl_service_name,
            src = "/root/.beacond",
            name = "node-beacond-config-{}".format(n),
        )

        genesis_artifact = plan.store_service_files(
            # The service name of a preexisting service from which the file will be copied.
            service_name = cl_service_name,

            # The path on the service's container that will be copied into a files artifact.
            # MANDATORY
            src = "/root/.beacond/config/genesis.json",

            # The name to give the files artifact that will be produced.
            # If not specified, it will be auto-generated.
            # OPTIONAL
            name = "cosmos-genesis-{}".format(file_suffix),
        )

        gentx_artifact = plan.store_service_files(
            service_name = cl_service_name,
            src = "/root/.beacond/config/gentx/*",
            name = "cosmos-gentx-{}".format(n),
        )

        # Node has completed its genesis step. We will add it back later once genesis is complete
        plan.remove_service(
            cl_service_name,
        )
    return node_peering_info

def get_persistent_peers(plan, peers):
    persistent_peers = peers[:]
    for i in range(len(persistent_peers)):
        peer_cl_service_name = "cl-validator-beaconkit-{}".format(i)
        peer_service = plan.get_service(peer_cl_service_name)
        persistent_peers[i] = persistent_peers[i] + "@" + peer_service.ip_address + ":26656"
    return ",".join(persistent_peers)

def create_node(plan, cl_image, peers, paired_el_client_name, jwt_file = None, kzg_trusted_setup_file = None, index = 0):
    cl_service_name = "cl-validator-beaconkit-{}".format(index)
    engine_dial_url = "http://{}:{}".format(paired_el_client_name, execution.ENGINE_RPC_PORT_NUM)

    # Get peers for the cl node
    persistent_peers = get_persistent_peers(plan, peers)

    beacond_config = get_config(
        cl_image,
        engine_dial_url,
        cl_service_name,
        entrypoint = ["bash"],
        cmd = ["-c", "/usr/bin/start.sh"],
        persistent_peers = persistent_peers,
        jwt_file = jwt_file,
        kzg_trusted_setup_file = kzg_trusted_setup_file,
    )

    # Add back in the node's config data and overwrite genesis.json with final genesis file
    beacond_config.files["/root"] = Directory(
        artifact_names = ["node-beacond-config-{}".format(index)],
    )
    beacond_config.files["/root/.tmp_genesis"] = Directory(artifact_names = ["cosmos-genesis-final"])

    plan.print(beacond_config)

    plan.add_service(
        name = cl_service_name,
        config = beacond_config,
    )

def create_full_node_config(plan, cl_image, peers, paired_el_client_name, jwt_file = None, kzg_trusted_setup_file = None, index = 0):
    cl_service_name = "cl-full-beaconkit-{}".format(index)
    engine_dial_url = "http://{}:{}".format(paired_el_client_name, execution.ENGINE_RPC_PORT_NUM)

    persistent_peers = get_persistent_peers(plan, peers)

    beacond_config = get_config(
        cl_image,
        engine_dial_url,
        cl_service_name,
        entrypoint = ["bash", "-c"],
        cmd = ["/usr/bin/init.sh && /usr/bin/start.sh"],
        persistent_peers = persistent_peers,
        jwt_file = jwt_file,
        kzg_trusted_setup_file = kzg_trusted_setup_file,
    )

    beacond_config.files["/root/.tmp_genesis"] = Directory(artifact_names = ["cosmos-genesis-final"])

    plan.print(beacond_config)

    return beacond_config
