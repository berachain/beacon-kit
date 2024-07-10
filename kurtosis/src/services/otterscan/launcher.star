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
        client_from_user):
    # Currently fetching the public port from service config is not allowed,
    # .ports["eth-json-rpc"].number returns the private port. Hence hardcoding the public port.
    if client_from_user.split("-")[2] != "erigon":
        fail("Currently only erigon client is supported for otterscan")
    config = get_config()
    plan.add_service(SERVICE_NAME, config)

def get_config():
    public_rpc_port_num = 8547
    el_client_rpc_url = "http://localhost:{}/".format(
        public_rpc_port_num,
    )
    return ServiceConfig(
        image = IMAGE_NAME,
        ports = USED_PORTS,
        env_vars = {
            "ERIGON_URL": el_client_rpc_url,
        },
    )
