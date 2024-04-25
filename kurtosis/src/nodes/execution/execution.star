constants = import_module("../../constants.star")
service_config_lib = import_module("../../lib/service_config.star")
builtins = import_module("../../lib/builtins.star")

RPC_PORT_NUM = 8545
WS_PORT_NUM = 8546
DISCOVERY_PORT_NUM = 30303
ENGINE_RPC_PORT_NUM = 8551
METRICS_PORT_NUM = 9001

# Port IDs
RPC_PORT_ID = "eth-json-rpc"
WS_PORT_ID = "eth-json-rpc-ws"
TCP_DISCOVERY_PORT_ID = "tcp-discovery"
UDP_DISCOVERY_PORT_ID = "udp-discovery"
ENGINE_RPC_PORT_ID = "engine-rpc"
ENGINE_WS_PORT_ID = "engineWs"
METRICS_PORT_ID = "metrics"

# Because structs are immutable, we pass around a map to allow full modification up until we create the final ServiceConfig
def get_default_service_config(service_name, node_module):
    sc = service_config_lib.get_service_config_template(
        name = service_name,
        image = node_module.IMAGE,
        ports = node_module.USED_PORTS_TEMPLATE,
        entrypoint = node_module.ENTRYPOINT,
        cmd = node_module.CMD,
        files = node_module.FILES,
    )

    return sc

def upload_global_files(plan, node_modules):
    genesis_file = plan.upload_files(
        src = "../../networks/kurtosis-devnet/network-configs/genesis.json",
        name = "genesis_file",
    )

    jwt_file = plan.upload_files(
        src = constants.JWT_FILEPATH,
        name = "jwt_file",
    )

    kzg_trusted_setup_file = plan.upload_files(
        src = constants.KZG_TRUSTED_SETUP_FILEPATH,
        name = "kzg_trusted_setup",
    )

    for node_module in node_modules.values():
        for global_file in node_module.GLOBAL_FILES:
            plan.upload_files(
                src = global_file[0],
                name = global_file[1],
            )

    return jwt_file, kzg_trusted_setup_file

def get_enode_addr(plan, el_service, el_service_name, el_type):
    extract_statement = {"enode": """.result.enode | split("?") | .[0]"""}

    request_recipe = PostHttpRequestRecipe(
        endpoint = "",
        body = '{"method":"admin_nodeInfo","params":[],"id":1,"jsonrpc":"2.0"}',
        content_type = "application/json",
        port_id = RPC_PORT_ID,
        extract = extract_statement,
    )

    response = plan.request(
        service_name = el_service_name,
        recipe = request_recipe,
    )

    enode = response["extract.enode"]
    return enode

def add_bootnodes(node_module, config, bootnodes):
    if type(bootnodes) == builtins.types.list:
        if len(bootnodes) > 0:
            cmdList = config["cmd"][:]
            cmdList.append(node_module.BOOTNODE_CMD)
            config["cmd"] = cmdList

            bootnodes_str = ",".join(bootnodes)
            config["cmd"].append(bootnodes_str)
    elif type(bootnodes) == builtins.types.str:
        if len(bootnodes) > 0:
            config["cmd"].append(node_module.BOOTNODE_CMD)
            config["cmd"].append(bootnodes)
    else:
        fail("Bootnodes was not a list or string, but instead a {}", type(bootnodes))

    return config

def deploy_node(plan, config):
    service_config = service_config_lib.create_from_config(config)

    return plan.add_service(
        name = config["name"],
        config = service_config,
    )

def create_node(plan, node_modules, node, node_type = "validator", index = 0, bootnode_enode_addrs = []):
    el_type = node.el_type
    node_module = node_modules[el_type]
    el_service_name = "el-{}-{}-{}".format(node_type, el_type, index)

    # 4a. Launch EL
    el_service_config_dict = get_default_service_config(el_service_name, node_module)
    el_service_config_dict = add_bootnodes(node_module, el_service_config_dict, bootnode_enode_addrs)
    el_client_service = deploy_node(plan, el_service_config_dict)

    enode_addr = get_enode_addr(plan, el_client_service, el_service_name, el_type)
    return {
        "name": el_service_name,
        "service": el_client_service,
        "enode_addr": enode_addr,
    }
