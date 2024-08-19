shared_utils = import_module("github.com/ethpandaops/ethereum-package/src/shared_utils/shared_utils.star")
builtins = import_module("./builtins.star")

COMETBFT_RPC_PORT_NUM = 26657
COMETBFT_P2P_PORT_NUM = 26656
COMETBFT_REST_PORT_NUM = 1317
COMETBFT_PPROF_PORT_NUM = 6060
METRICS_PORT_NUM = 26660

# Port IDs
COMETBFT_RPC_PORT_ID = "cometbft-rpc"
COMETBFT_P2P_PORT_ID = "cometbft-p2p"
COMETBFT_GRPC_PORT_ID = "cometbft-grpc"
COMETBFT_REST_PORT_ID = "cometbft-rest"
COMETBFT_PPROF_PORT_ID = "cometbft-pprof"
ENGINE_RPC_PORT_ID = "engine-rpc"
METRICS_PORT_ID = "metrics"
METRICS_PATH = "/metrics"
RPC_PORT_NUM = 8545
ENGINE_RPC_PORT_NUM = 8551
PUBLIC_RPC_PORT_NUM = 8547
DEFAULT_PRIVATE_IP_ADDRESS_PLACEHOLDER = "KURTOSIS_IP_ADDR_PLACEHOLDER"
DEFAULT_MAX_CPU = 2000  # 2 cores
DEFAULT_MAX_MEMORY = 2048  # 2 GB


USED_PORTS = {
    COMETBFT_RPC_PORT_ID: shared_utils.new_port_spec(COMETBFT_RPC_PORT_NUM, shared_utils.TCP_PROTOCOL),
    COMETBFT_P2P_PORT_ID: shared_utils.new_port_spec(COMETBFT_P2P_PORT_NUM, shared_utils.TCP_PROTOCOL),
    COMETBFT_REST_PORT_ID: shared_utils.new_port_spec(COMETBFT_REST_PORT_NUM, shared_utils.TCP_PROTOCOL),
    COMETBFT_PPROF_PORT_ID: shared_utils.new_port_spec(COMETBFT_PPROF_PORT_NUM, shared_utils.TCP_PROTOCOL),
    # ENGINE_RPC_PORT_ID: shared_utils.new_port_spec(ENGINE_RPC_PORT_NUM, shared_utils.TCP_PROTOCOL),
    METRICS_PORT_ID: shared_utils.new_port_spec(METRICS_PORT_NUM, shared_utils.TCP_PROTOCOL, wait = None),
}

JWT_FILEPATH = "/testing/sync-network/network/jwt-secret.hex"
GENESIS_FILEPATH = "/testing/sync-network/network/80084/genesis.json"

def run(plan, network = {}, nodes =[], node_settings = {}):
    jwt_file = plan.upload_files(
        src = JWT_FILEPATH,
        name = "jwt_file",
    )
    genesis_file = plan.upload_files(
        src = GENESIS_FILEPATH,
        name = "genesis",
    )

    # start_cl_nodes(plan, node_settings, jwt_file, genesis_file)
    
    start_el_nodes(plan, nodes, node_settings, jwt_file, genesis_file)


def start_cl_nodes(plan, node_settings, jwt_file, genesis_file):
    plan.print("starting nodes")
    beacond_config = create_node_config(plan, node_settings, jwt_file, genesis_file)
    plan.print("beacond_config", str(beacond_config))
    plan.add_service(name = "beacond", config = beacond_config)

def create_node_config(plan, node_settings, jwt_file = None, genesis_file = None):
    # engine_dial_url = "http://{}:{}".format(paired_el_client_name, execution.ENGINE_RPC_PORT_NUM)
    cmd = start(plan)
    plan.print("cmd", cmd)

    beacond_config = get_config(
        plan,
        node_settings,
        entrypoint = ["bash", "-c"],
        cmd = [cmd],
        jwt_file = jwt_file,
        genesis_file = genesis_file,
    )
    plan.print(beacond_config)

    return beacond_config

def start(plan):
    mv_genesis = "mv root/.tmp_genesis/eth-genesis.json /root/.beacond/config/genesis.json"

    start_node = "CHAIN_SPEC=testnet /usr/bin/beacond start \
    --beacon-kit.engine.jwt-secret-path=/root/jwt/jwt-secret.hex \
    --beacon-kit.kzg.implementation={} \
    --api.enable --api.enabled-unsafe-cors --minimum-gas-prices {}".format(
        "crate-crypto/go-kzg-4844","$BEACOND_MINIMUM_GAS_PRICE"
        )

    # --beacon-kit.engine.rpc-dial-url {}  BEACOND_ENGINE_DIAL_URL\
    #    --rpc.laddr tcp://0.0.0.0:26657 --api.address tcp://0.0.0.0:1317 
    # return "{} && {}".format(mv_genesis, start_node)
    return "{}".format(start_node)

def get_config(plan, node_settings, entrypoint = [], cmd = [], jwt_file = None, genesis_file = None):
    app_file = plan.upload_files(
        src = "/testing/sync-network/network/80084/app.toml",
        name = "app_file",
    )

    client_file = plan.upload_files(
        src = "/testing/sync-network/network/80084/client.toml",
        name = "client_file",
    )

    config_file = plan.upload_files(
        src = "/testing/sync-network/network/80084/config.toml",
        name = "config_file",
    )

    genesis_file = plan.upload_files(
        src = "/testing/sync-network/network/80084/genesis.json",
        name = "genesis_file",
    )

    # network_files = plan.upload_files(
    #     src = "/testing/sync-network/network/80084/",
    #     name = "files",
    # )

    files = {}
    if jwt_file:
        plan.print("jwt file present")
        files["/root/jwt"] = jwt_file
    if genesis_file:
        plan.print("genesis file present")
        files["/root/.beacond/config"] = genesis_file

    files["/root/.beacond/config"] = app_file
    files["/root/.beacond/config"] = config_file
    files["/root/.beacond/config"] = client_file

    settings = node_settings["consensus_settings"]

    config = ServiceConfig(
        image = settings["images"]["beaconkit"],
        files = files,
        entrypoint = entrypoint,
        cmd = cmd,
        min_cpu = settings["specs"]["min_cpu"],
        max_cpu = settings["specs"]["max_cpu"],
        min_memory = settings["specs"]["min_memory"],
        max_memory = settings["specs"]["max_memory"],
        env_vars = {
            "BEACOND_MONIKER": "test",
            "BEACOND_NET": "VALUE_2",
            "BEACOND_HOME": "/root/.beacond",
            "BEACOND_CHAIN_ID": "bartio-beacon-80084",
            "BEACOND_DEBUG": "false",
            "BEACOND_KEYRING_BACKEND": "test",
            "BEACOND_MINIMUM_GAS_PRICE": "0abgt",
            "BEACOND_ETH_CHAIN_ID": "80084",
            "BEACOND_ENABLE_PROMETHEUS": "true",
            "BEACOND_CONSENSUS_KEY_ALGO": "bls12_381",
        },
        ports = USED_PORTS,
        #     labels = node_labels,
        #     node_selectors = settings.node_selectors,
    )

    return config

# def get_peer_info(plan, cl_service_name):
#     peer_result = exec_on_service(plan, cl_service_name, "/usr/bin/beacond comet show-node-id --home $BEACOND_HOME | tr -d '\n'")
#     return peer_result["output"]

# def exec_on_service(plan, service_name, command):
#     return plan.exec(
#         service_name = service_name,
#         recipe = ExecRecipe(
#             command = ["bash", "-c", command],
#         ),
#     )


def start_el_nodes(plan, nodes, node_settings, jwt_file, genesis_file):
    node_modules = {}
    full_node_el_client_configs = []
    for node in nodes:
            eth_type = node["el_type"]
            plan.print("node", str(eth_type))
            node_path = "./nodes/execution/{}/config.star".format(eth_type)
            plan.print("node_path", node_path)
            node_module = import_module(node_path)
            node_modules[eth_type] = node_module

    for n, full in enumerate(nodes):
        el_client_config = generate_node_config(plan, node_modules, full,node_settings)
        full_node_el_client_configs.append(el_client_config)

    if full_node_el_client_configs != []:
        full_node_el_clients = deploy_nodes(plan, full_node_el_client_configs, True)

#### execution method helpers
def generate_node_config(plan, node_modules, node_struct,node_settings):
    node_module = node_modules[node_struct["el_type"]]

    # 4a. Launch EL
    el_service_config_dict = get_default_service_config(plan, node_struct, node_settings, node_module)

    # el_service_config_dict = add_bootnodes(node_module, el_service_config_dict, bootnode_enode_addrs)

    return el_service_config_dict


# Because structs are immutable, we pass around a map to allow full modification up until we create the final ServiceConfig
def get_default_service_config(plan, node_struct, node_settings, node_module):
    settings = node_settings["execution_settings"]

    plan.print("settings", str(settings))
    plan.print("settings", str(settings["images"][node_struct["el_type"]]))
    # Define common parameters
    common_params = {
        "name": node_struct["el_type"],
        "image": settings["images"][node_struct["el_type"]],
        "ports": node_module.USED_PORTS_TEMPLATE,
        "entrypoint": node_module.ENTRYPOINT,
        "cmd": node_module.CMD,
        "files": node_module.FILES,
        "min_cpu": settings["specs"]["min_cpu"],
        "max_cpu": settings["specs"]["max_cpu"],
        "min_memory": settings["specs"]["min_memory"],
        "max_memory": settings["specs"]["max_memory"],
    }

    # # Get the service config template
    sc = get_service_config_template(**common_params)

    return sc

def deploy_nodes(plan, configs, is_full_node = True):
    service_configs = {}
    for config in configs:
        service_configs[config["name"]] = create_from_config(config, is_full_node)

    return plan.add_services(
        configs = service_configs,
    )

def validate_port_spec(port_spec):
    if type(port_spec) != builtins.types.dict:
        fail("Port spec must be a dict, not {0}".format(type(port_spec)))

    if type(port_spec["number"]) != builtins.types.int:
        fail("Port spec number must be an int, not {0}".format(type(port_spec["number"])))

    if port_spec["transport_protocol"] != None:
        if type(port_spec["transport_protocol"]) != builtins.types.string:
            fail("Port spec transport_protocol must be a string, not {0}".format(type(port_spec["transport_protocol"])))

    if port_spec["application_protocol"] != None:
        if type(port_spec["application_protocol"]) != builtins.types.string:
            fail("Port spec application_protocol must be a string, not {0}".format(type(port_spec["application_protocol"])))

    if port_spec["wait"] != None:
        if type(port_spec["wait"]) != builtins.types.string:
            fail("Port spec wait must be a bool, not {0}".format(type(port_spec["wait"])))

def create_port_specs_from_config(config):
    ports = {}
    for port_key, port_spec in config["ports"].items():
        ports[port_key] = create_port_spec(port_spec)

    return ports

def create_port_spec(port_spec_dict):
    return PortSpec(
        number = port_spec_dict["number"],
        transport_protocol = port_spec_dict["transport_protocol"],
        application_protocol = port_spec_dict["application_protocol"],
        wait = port_spec_dict["wait"],
    )

def create_from_config(config, is_full_node = False):
    validate_service_config_types(config)

    return ServiceConfig(
        image = config["image"],
        ports = create_port_specs_from_config(config),
        public_ports =  {},
        files = config["files"] if config["files"] else {},
        entrypoint = config["entrypoint"] if config["entrypoint"] else [],
        cmd = [" ".join(config["cmd"])] if config["cmd"] else [],
        env_vars = config["env_vars"] if config["env_vars"] else {},
        private_ip_address_placeholder = config["private_ip_address_placeholder"] if config["private_ip_address_placeholder"] else DEFAULT_PRIVATE_IP_ADDRESS_PLACEHOLDER,
        max_cpu = config["max_cpu"] if config["max_cpu"] else DEFAULT_MAX_CPU,  # Needs a default, as 0 does not flag as optional
        min_cpu = config["min_cpu"] if config["min_cpu"] else 0,
        max_memory = config["max_memory"] if config["max_memory"] else DEFAULT_MAX_MEMORY,  # Needs a default, as 0 does not flag as optional
        min_memory = config["min_memory"] if config["min_memory"] else 0,
        #ready_conditions=config['ready_conditions'], Ready conditions not yet supported
        labels = config["labels"] if config["labels"] else {},
        #user=config['user'], User config not yet supported
        tolerations = config["tolerations"] if config["tolerations"] else [],
        node_selectors = config["node_selectors"] if config["node_selectors"] else {},
    )

def validate_service_config_types(service_config):
    if type(service_config) != builtins.types.dict:
        fail("Service config must be a dict, not {0}".format(type(service_config)))

    if type(service_config["name"]) != builtins.types.string:
        fail("Service config name must be a string, not {0}".format(type(service_config["name"])))

    if type(service_config["image"]) != builtins.types.string:
        fail("Service config image must be a string, not {0}".format(type(service_config["image"])))

    if service_config["ports"] != None:
        if type(service_config["ports"]) != builtins.types.dict:
            fail("Service config ports must be a dict, not {0}".format(type(service_config["ports"])))
        for port_key, port_spec in service_config["ports"].items():
            if type(port_key) != builtins.types.string:
                fail("Service config port key must be an int, not {0}".format(type(port_key)))
            validate_port_spec(port_spec)

    if service_config["files"] != None:
        if type(service_config["files"]) != builtins.types.dict:
            fail("Service config files must be a dict, not {0}".format(type(service_config["files"])))
        for path, content in service_config["files"].items():
            if type(path) != builtins.types.string:
                fail("Service config file path must be a string, not {0}".format(type(path)))
            if type(content) not in [builtins.types.string, builtins.types.directory]:
                fail("Service config file content must be a string or a Directory object, not {0}".format(type(content)))

    if service_config["entrypoint"] != None:
        if type(service_config["entrypoint"]) != builtins.types.list:
            fail("Service config entrypoint must be a list, not {0}".format(type(service_config["entrypoint"])))
        for entrypoint in service_config["entrypoint"]:
            if type(entrypoint) != builtins.types.string:
                fail("Service config entrypoint must be a string, not {0}".format(type(entrypoint)))

    if service_config["cmd"] != None:
        if type(service_config["cmd"]) != builtins.types.list:
            fail("Service config cmd must be a list, not {0}".format(type(service_config["cmd"])))
        for cmd in service_config["cmd"]:
            if type(cmd) != builtins.types.string:
                fail("Service config cmd must be a string, not {0}".format(type(cmd)))

    if service_config["env_vars"] != None:
        if type(service_config["env_vars"]) != builtins.types.dict:
            fail("Service config env_vars must be a dict, not {0}".format(type(service_config["env_vars"])))
        for env_var_key, env_var_value in service_config["env_vars"].items():
            if type(env_var_key) != builtins.types.string:
                fail("Service config env_var key must be a string, not {0}".format(type(env_var_key)))
            if type(env_var_value) != builtins.types.string:
                fail("Service config env_var value must be a string, not {0}".format(type(env_var_value)))

    if service_config["private_ip_address_placeholder"] != None:
        if type(service_config["private_ip_address_placeholder"]) != builtins.types.string:
            fail("Service config private_ip_address_placeholder must be a string, not {0}".format(type(service_config["private_ip_address_placeholder"])))

    if service_config["max_cpu"] != None:
        if type(service_config["max_cpu"]) != builtins.types.int:
            fail("Service config max_cpu must be a int, not {0}".format(type(service_config["max_cpu"])))
    if service_config["min_cpu"] != None:
        if type(service_config["min_cpu"]) != builtins.types.int:
            fail("Service config min_cpu must be a int, not {0}".format(type(service_config["min_cpu"])))
    if service_config["max_memory"] != None:
        if type(service_config["max_memory"]) != builtins.types.int:
            fail("Service config max_memory must be a int, not {0}".format(type(service_config["max_memory"])))
    if service_config["min_memory"] != None:
        if type(service_config["min_memory"]) != builtins.types.int:
            fail("Service config min_memory must be a int, not {0}".format(type(service_config["min_memory"])))

    # TODO(validation): Implement validation for ready_conditions
    # TODO(validation): Implement validation for labels
    # TODO(validation): Implement validation for user
    # TODO(validation): Implement validation for tolerations
    # TODO(validation): Implement validation for node_selectors

def get_service_config_template(
        name,
        image,
        ports = None,
        public_ports = None,
        files = None,
        entrypoint = None,
        cmd = None,
        env_vars = None,
        private_ip_address_placeholder = None,
        max_cpu = None,
        min_cpu = None,
        max_memory = None,
        min_memory = None,
        ready_conditions = None,
        labels = None,
        user = None,
        tolerations = None,
        node_selectors = None):
    service_config = {
        "name": name,
        "image": image,
        "ports": ports,
        "public_ports": public_ports,
        "files": files,
        "entrypoint": entrypoint,
        "cmd": cmd,
        "env_vars": env_vars,
        "private_ip_address_placeholder": private_ip_address_placeholder,
        "max_cpu": max_cpu,
        "min_cpu": min_cpu,
        "max_memory": max_memory,
        "min_memory": min_memory,
        "ready_conditions": ready_conditions,
        "labels": labels,
        "user": user,
        "tolerations": tolerations,
        "node_selectors": node_selectors,
    }

    # validate_service_config_types(service_config)
    return service_config    