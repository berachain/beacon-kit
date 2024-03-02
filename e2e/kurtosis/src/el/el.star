reth = import_module("./reth_launcher.star")
input_parser = import_module("github.com/kurtosis-tech/ethereum-package/src/package_io/input_parser.star")
constants = import_module("../constants.star")

# Returns the el client context
def get_el(plan, client_type, evm_genesis_data, jwt_file, el_service_name, network_params):
    if client_type == constants.EL_CLIENT_TYPE.reth:
        return reth.get_reth(plan, evm_genesis_data, jwt_file, el_service_name, network_params)