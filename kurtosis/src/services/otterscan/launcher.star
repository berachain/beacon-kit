shared_utils = import_module("github.com/ethpandaops/ethereum-package/src/shared_utils/shared_utils.star")

SERVICE_NAME = "otterscan"
IMAGE_NAME = "otterscan/otterscan:latest"
HTTP_PORT_ID = "http"
HTTP_PORT_NUMBER = 80
USED_PORTS = {
    HTTP_PORT_ID: shared_utils.new_port_spec(
        HTTP_PORT_NUMBER,
        shared_utils.TCP_PROTOCOL,
        shared_utils.HTTP_APPLICATION_PROTOCOL,
    ),
}

def launch_otterscan(
        plan,
        full_node_el_clients,
        client_from_user):
    el_client_info = {}

    # TODO: If client_from_user is other than erigon node, give error or something to user.
    # Get the full_node_el_clients that match the client_from_user
    for full_node_el_client_name, full_node_el_client_service in full_node_el_clients.items():
        if full_node_el_client_name in client_from_user:
            rpc_port = full_node_el_client_service.ports["eth-json-rpc"].number
            name = full_node_el_client_name
            ip_address = full_node_el_client_service.ip_address

            el_client_info = get_el_client_info(
                ip_address,
                rpc_port,
                name,
            )
            break
    config = get_config(el_client_info.get("RPC_Url"))
    plan.add_service(SERVICE_NAME, config)

def get_config(RPC_Url):
    return ServiceConfig(
        image = IMAGE_NAME,
        ports = USED_PORTS,
        env_vars = {
            "ERIGON_URL": RPC_Url,
        },
    )

def get_el_client_info(ip_addr, rpc_port_num, full_name):
    el_client_rpc_url = "http://{}:{}/".format(
        "localhost",
        rpc_port_num,
    )
    el_client_type = full_name.split("-")[2]
    return {
        "RPC_Url": el_client_rpc_url,
        "Eth_Type": el_client_type,
    }
