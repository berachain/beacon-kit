shared_utils = import_module("github.com/ethpandaops/ethereum-package/src/shared_utils/shared_utils.star")
execution = import_module("../../execution/execution.star")
node = import_module("./node.star")
bash = import_module("../../../lib/bash.star")

COMETBFT_RPC_PORT_NUM = 26657
COMETBFT_P2P_PORT_NUM = 26656
COMETBFT_REST_PORT_NUM = 1317
COMETBFT_PPROF_PORT_NUM = 6060
METRICS_PORT_NUM = 26660
ENGINE_RPC_PORT_NUM = 8551
NODE_API_PORT_NUM = 3500

# Port IDs
COMETBFT_RPC_PORT_ID = "cometbft-rpc"
COMETBFT_P2P_PORT_ID = "cometbft-p2p"
COMETBFT_GRPC_PORT_ID = "cometbft-grpc"
COMETBFT_REST_PORT_ID = "cometbft-rest"
COMETBFT_PPROF_PORT_ID = "cometbft-pprof"
ENGINE_RPC_PORT_ID = "engine-rpc"
METRICS_PORT_ID = "metrics"
METRICS_PATH = "/metrics"
NODE_API_PORT_ID = "node-api"

USED_PORTS = {
    COMETBFT_RPC_PORT_ID: shared_utils.new_port_spec(COMETBFT_RPC_PORT_NUM, shared_utils.TCP_PROTOCOL),
    COMETBFT_P2P_PORT_ID: shared_utils.new_port_spec(COMETBFT_P2P_PORT_NUM, shared_utils.TCP_PROTOCOL),
    COMETBFT_REST_PORT_ID: shared_utils.new_port_spec(COMETBFT_REST_PORT_NUM, shared_utils.TCP_PROTOCOL),
    COMETBFT_PPROF_PORT_ID: shared_utils.new_port_spec(COMETBFT_PPROF_PORT_NUM, shared_utils.TCP_PROTOCOL),
    # ENGINE_RPC_PORT_ID: shared_utils.new_port_spec(ENGINE_RPC_PORT_NUM, shared_utils.TCP_PROTOCOL),
    METRICS_PORT_ID: shared_utils.new_port_spec(METRICS_PORT_NUM, shared_utils.TCP_PROTOCOL, wait = None),
    NODE_API_PORT_ID: shared_utils.new_port_spec(NODE_API_PORT_NUM, shared_utils.TCP_PROTOCOL),
}

def get_config(node_struct, engine_dial_url, entrypoint = [], cmd = [], persistent_peers = "", expose_ports = True, jwt_file = None, kzg_trusted_setup_file = None):
    exposed_ports = {}
    if expose_ports:
        exposed_ports = USED_PORTS

    files = {}
    if jwt_file:
        files["/root/jwt"] = jwt_file
    if kzg_trusted_setup_file:
        files["/root/kzg"] = kzg_trusted_setup_file

    settings = node_struct.consensus_settings

    node_labels = dict(settings.labels)
    node_labels["node_type"] = "consensus"

    config = ServiceConfig(
        image = node_struct.cl_image,
        files = files,
        entrypoint = entrypoint,
        cmd = cmd,
        min_cpu = settings.specs.min_cpu,
        max_cpu = settings.specs.max_cpu,
        min_memory = settings.specs.min_memory,
        max_memory = settings.specs.max_memory,
        env_vars = {
            "BEACOND_MONIKER": node_struct.cl_service_name,
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
            "BEACOND_CONSENSUS_KEY_ALGO": "bls12_381",
        },
        ports = exposed_ports,
        labels = node_labels,
        node_selectors = settings.node_selectors,
    )

    return config

def perform_genesis_ceremony(plan, validators, jwt_file):
    num_validators = len(validators)

    node_peering_info = []
    beacond_configs = []
    stored_configs = []

    for n in range(num_validators):
        beacond_configs.append("node-beacond-config-{}".format(n))
        stored_configs.append(StoreSpec(src = "/tmp/config{}".format(n), name = beacond_configs[n]))

    stored_configs.append(StoreSpec(src = "/tmp/config_genesis/.beacond/config/genesis.json", name = "cosmos-genesis-final"))

    multiple_gentx_file = plan.upload_files(
        src = "./scripts/multiple-premined-deposits.sh",
        name = "multiple-premined-deposits",
        description = "Uploading multiple-premined-deposits script",
    )

    multiple_gentx_env_vars = node.get_genesis_env_vars("cl-validator-beaconkit-0")
    multiple_gentx_env_vars["NUM_VALS"] = str(num_validators)

    plan.print(multiple_gentx_env_vars)
    plan.print(stored_configs)

    plan.run_sh(
        run = "chmod +x /app/scripts/multiple-premined-deposits.sh && /app/scripts/multiple-premined-deposits.sh",
        image = validators[0].cl_image,
        files = {
            "/app/scripts": "multiple-premined-deposits",
            "/root/eth_genesis": "genesis_file",
        },
        env_vars = multiple_gentx_env_vars,
        store = stored_configs,
        description = "Collecting beacond genesis files",
    )

def get_persistent_peers(plan, peers):
    persistent_peers = peers[:]
    for i in range(len(persistent_peers)):
        peer_cl_service_name = "cl-seed-beaconkit-{}".format(i)
        peer_service = plan.get_service(peer_cl_service_name)
        persistent_peers[i] = persistent_peers[i] + "@" + peer_service.ip_address + ":26656"
    return ",".join(persistent_peers)

def init_consensus_nodes():
    genesis_file = "{}/config/genesis.json".format("$BEACOND_HOME")

    # Check if genesis file exists, if not then initialize the beacond
    init_node = "if [ ! -f {} ]; then /usr/bin/beacond init --chain-id {} {} --home {}; fi".format(genesis_file, "$BEACOND_CHAIN_ID", "$BEACOND_MONIKER", "$BEACOND_HOME")
    add_validator = "/usr/bin/beacond genesis add-premined-deposit --home {}".format("$BEACOND_HOME")
    collect_gentx = "/usr/bin/beacond genesis collect-premined-deposits --home {}".format("$BEACOND_HOME")
    return "{} && {} && {}".format(init_node, add_validator, collect_gentx)

def create_node_config(plan, node_struct, peers, paired_el_client_name, jwt_file = None, kzg_trusted_setup_file = None):
    engine_dial_url = "http://{}:{}".format(paired_el_client_name, execution.ENGINE_RPC_PORT_NUM)

    persistent_peers = get_persistent_peers(plan, peers)
    config_settings = node_struct.consensus_settings.config
    app_settings = node_struct.consensus_settings.app
    kzg_impl = node_struct.kzg_impl

    cmd = "{} && {}".format(init_consensus_nodes(), node.start(persistent_peers, False, 0, config_settings, app_settings, kzg_impl))
    if node_struct.node_type == "validator":
        cmd = node.start(persistent_peers, False, node_struct.index, config_settings, app_settings, kzg_impl)
    elif node_struct.node_type == "seed":
        cmd = "{} && {}".format(init_consensus_nodes(), node.start(persistent_peers, True, 0, config_settings, app_settings, kzg_impl))

    beacond_config = get_config(
        node_struct,
        engine_dial_url,
        entrypoint = ["bash", "-c"],
        cmd = [cmd],
        persistent_peers = persistent_peers,
        jwt_file = jwt_file,
        kzg_trusted_setup_file = kzg_trusted_setup_file,
    )

    if node_struct.node_type == "validator":
        # Add back in the node's config data and overwrite genesis.json with final genesis file
        beacond_config.files["/root"] = Directory(
            artifact_names = ["node-beacond-config-{}".format(node_struct.index)],
        )

    beacond_config.files["/root/.tmp_genesis"] = Directory(artifact_names = ["cosmos-genesis-final"])

    plan.print(beacond_config)

    return beacond_config

def get_peer_info(plan, cl_service_name):
    peer_result = bash.exec_on_service(plan, cl_service_name, "/usr/bin/beacond comet show-node-id --home $BEACOND_HOME | tr -d '\n'")
    return peer_result["output"]

def dial_unsafe_peers(plan, seed_service_name, peers):
    peers_list = []
    for cl_service_name, peer_info in peers.items():
        p2p_addr = "\"{}@{}:26656\"".format(peer_info, plan.get_service(cl_service_name).ip_address)
        peers_list.append(p2p_addr)

    # Split peers_list into groups of 20
    peer_groups = [peers_list[i:i + 20] for i in range(0, len(peers_list), 20)]
    for group in peer_groups:
        peer_string = ",".join(group)
        endpoint = "/dial_peers?peers=%5B{}%5D&persistent=false".format(peer_string)
        curl_command = ["curl", "-X", "GET", "http://localhost:{}{}".format(COMETBFT_RPC_PORT_NUM, endpoint)]
        exec_recipe = ExecRecipe(
            command = curl_command,
        )
        plan.exec(
            service_name = seed_service_name,
            recipe = exec_recipe,
            description = "Adding peers to seed node",
        )
