eth_static_files = import_module("github.com/kurtosis-tech/ethereum-package/src/static_files/static_files.star")
input_parser = import_module("github.com/kurtosis-tech/ethereum-package/src/package_io/input_parser.star")

execution_types = import_module("./types.star")
constants = import_module("../../constants.star")
service_config_lib = import_module("../../lib/service_config.star")
builtins = import_module("../../lib/builtins.star")

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
        src = constants.KURTOSIS_ETH_PACKAGE_URL + eth_static_files.JWT_PATH_FILEPATH,
        name = "jwt_file",
    )
    for node_module in node_modules.values():
        for global_file in node_module.GLOBAL_FILES:
            plan.upload_files(
                src = global_file[0],
                name = global_file[1],
            )

    return jwt_file

def get_enode_addr(plan, el_service, el_service_name, el_type):
    extract_statement = {"enode": """.result.enode | split("?") | .[0]"""}
    if el_type == execution_types.CLIENTS.reth:
        extract_statement = {"enode": """.result.id | split("?") | .[0][2:] | ("enode://" + .)"""}

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
    return enode + "@" + el_service.ip_address + ":" + str(DISCOVERY_PORT_NUM) if el_type == execution_types.CLIENTS.reth else enode

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
