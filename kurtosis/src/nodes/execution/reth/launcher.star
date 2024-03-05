reth = import_module("github.com/kurtosis-tech/ethereum-package/src/el/reth/reth_launcher.star")
input_parser = import_module("github.com/kurtosis-tech/ethereum-package/src/package_io/input_parser.star")
#execution = import_module("../execution.star")

def get(plan, evm_genesis_data, jwt_file, el_service_name, network_params, existing_el_clients = []):
    reth_launcher = get_launcher(evm_genesis_data, jwt_file, network_params)
    return get_context(plan, reth_launcher["launcher"], el_service_name, reth_launcher["launch_method"], existing_el_clients)

def get_launcher(evm_genesis_data, jwt_file, network_params):
    reth_launcher = {
        "launcher": reth.new_reth_launcher(
            evm_genesis_data,
            jwt_file,
            network_params.network,
        ),
        "launch_method": reth.launch,
    }

    return reth_launcher

def get_context(plan, el_launcher, el_service_name, launch_method, existing_el_clients = [], participant_log_level = "info"):
    extra_params = []

    if len(existing_el_clients) > 0:
        # enode_args = [
        #     ctx.enode + "#" + ctx.ip_addr
        #     for ctx in existing_el_clients
        # ]
        # peer_nodes = execution.parse_proper_enode_ids(plan,existing_el_clients):

        peer_nodes = ",".join(existing_el_clients)

        trusted_peers_cmd = "--trusted-peers"
        bootnodes_cmd = "--bootnodes"
        extra_params.append(trusted_peers_cmd)
        extra_params.append(peer_nodes)
        extra_params.append(bootnodes_cmd)
        extra_params.append(peer_nodes)

    return launch_method(
        plan,
        el_launcher,
        el_service_name,
        input_parser.DEFAULT_EL_IMAGES["reth"],
        participant_log_level,
        "",  # global_log_level: unused because we pass a default participant_log_level here
        [],  # existing_el_clients

        # min,max cpu and min,max mem
        # currently undefined
        # TODO(resources): Add support for specific resource management
        0,
        0,
        0,
        0,
        extra_params,  # extra_params
        {},  # extra_env_vars
        {},  # extra_labels
        False,  # persistent: Not using persistent storage for now
        0,  # el_volume_size: Using default
        [],  # tolerations: none
        {},  # node_selectors: none
    )
