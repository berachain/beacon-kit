reth = import_module("./reth/launcher.star")
input_parser = import_module("github.com/kurtosis-tech/ethereum-package/src/package_io/input_parser.star")
execution_types = import_module("./types.star")

# Returns the el client context
def get_client(plan, client_type, evm_genesis_data, jwt_file, el_service_name, network_params, existing_el_clients = []):
    if client_type == execution_types.CLIENTS.reth:
        return reth.get(plan, evm_genesis_data, jwt_file, el_service_name, network_params, existing_el_clients)
