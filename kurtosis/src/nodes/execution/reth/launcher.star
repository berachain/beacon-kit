reth = import_module("github.com/kurtosis-tech/ethereum-package/src/el/reth/reth_launcher.star")
input_parser = import_module("github.com/kurtosis-tech/ethereum-package/src/package_io/input_parser.star")

def get(plan, evm_genesis_data, jwt_file, el_service_name, network_params):
    reth_launcher = get_launcher(evm_genesis_data, jwt_file, network_params)
    return get_context(plan, reth_launcher["launcher"], el_service_name, reth_launcher["launch_method"])

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

def get_context(plan, el_launcher, el_service_name, launch_method, participant_log_level = "info"):
    return launch_method(
        plan,
        el_launcher,
        el_service_name,
        input_parser.DEFAULT_EL_IMAGES["reth"],
        participant_log_level,
        "",  # global_log_level: unused because we pass a default participant_log_level here
        [],  # existing_el_clients TODO(p2p): insert multiclient p2p support here

        # min,max cpu and min,max mem
        # currently undefined
        # TODO(resources): Add support for specific resource management
        0,
        0,
        0,
        0,
        [],  # extra_params
        {},  # extra_env_vars
        {},  # extra_labels
        False,  # persistent: Not using persistent storage for now
        0,  # el_volume_size: Using default
        [],  # tolerations: none
        {},  # node_selectors: none
    )
