eth_static_files = import_module("github.com/kurtosis-tech/ethereum-package/src/static_files/static_files.star")
input_parser = import_module("github.com/kurtosis-tech/ethereum-package/src/package_io/input_parser.star")

reth = import_module("./reth/launcher.star")
geth = import_module("./geth/launcher.star")
execution_types = import_module("./types.star")
constants = import_module("../../constants.star")

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

# Returns the el client context
def get_client(plan, client_type, evm_genesis_data, jwt_file, el_service_name, network_params, existing_el_clients = []):
    if client_type == execution_types.CLIENTS.reth:
        return reth.get(plan, evm_genesis_data, jwt_file, el_service_name, network_params, existing_el_clients)

def get_default_service_config(service_name, client_type):
    if client_type == execution_types.CLIENTS.geth:
        return geth.get_default_service_config(service_name)

def upload_global_files(plan):
    genesis_file = plan.upload_files(
        src = "../../networks/kurtosis-devnet/network-configs/genesis.json",
        name = "genesis_file",
    )
    jwt_file = plan.upload_files(
        src = constants.KURTOSIS_ETH_PACKAGE_URL + eth_static_files.JWT_PATH_FILEPATH,
        name = "jwt_file",
    )
    geth.upload_global_files(plan)

    return jwt_file

# Expects a list of enode strings in the format "enode://<enode_id>@<old_ip>:<old_port>#<new_ip>:<new_port>"
def parse_proper_enode_ids(plan, enodes):
    result = plan.run_python(
        run = """import sys
enodes = []
for enode in sys.argv[1:]:
    parsed = enode.split('#')
    en = parsed[0]
    ip = parsed[1]
    enodes.append(en.split('@')[0] + "@" + ip + ":30303")
enode_str = ",".join(enodes)
print(enode_str)
""",
        args = enodes,
    )

    peer_nodes = result.output
    return peer_nodes

def get_enode_addr(plan, el_service, el_service_name, el_client_type):
    request_recipe = None
    if el_client_type == execution_types.CLIENTS.reth:
        request_recipe = PostHttpRequestRecipe(
            endpoint = "",
            body = '{"method":"admin_nodeInfo","params":[],"id":1,"jsonrpc":"2.0"}',
            content_type = "application/json",
            port_id = RPC_PORT_ID,
            extract = {
                "enode": """.result.id | split("?") | .[0][2:] | ("enode://" + .)""",
            },
        )
    elif el_client_type == execution_types.CLIENTS.geth:
        request_recipe = PostHttpRequestRecipe(
            endpoint = "",
            body = '{"method":"admin_nodeInfo","params":[],"id":1,"jsonrpc":"2.0"}',
            content_type = "application/json",
            port_id = RPC_PORT_ID,
            extract = {
                "enode": """.result.enode | split("?") | .[0]""",
            },
        )

    response = plan.request(
        service_name = el_service_name,
        recipe = request_recipe,
    )

    enode = response["extract.enode"]
    return enode + "@" + el_service.ip_address + ":" + str(DISCOVERY_PORT_NUM) if el_client_type == execution_types.CLIENTS.reth else enode
