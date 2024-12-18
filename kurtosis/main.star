# Importing modules
modules = {
    "el_cl_genesis_data_generator": "github.com/ethpandaops/ethereum-package/src/prelaunch_data_generator/el_cl_genesis/el_cl_genesis_generator.star",
    "execution": "./src/nodes/execution/execution.star",
    "service_module": "./src/services/service.star",
    "beacond": "./src/nodes/consensus/beacond/launcher.star",
    "networks": "./src/networks/networks.star",
    "port_spec_lib": "./src/lib/port_spec.star",
    "nodes": "./src/nodes/nodes.star",
    "nginx": "./src/services/nginx/nginx.star",
    "constants": "./src/constants.star",
    "goomy_blob": "./src/services/goomy/launcher.star",
    "prometheus": "./src/observability/prometheus/prometheus.star",
    "grafana": "./src/observability/grafana/grafana.star",
    "pyroscope": "./src/observability/pyroscope/pyroscope.star",
    "tx_fuzz": "./src/services/tx_fuzz/launcher.star",
    "blutgang": "./src/services/blutgang/launcher.star",
    "blockscout": "./src/services/blockscout/launcher.star",
    "otterscan": "./src/services/otterscan/launcher.star"
}
for name, path in modules.items():
    globals()[name] = import_module(path)

def run(plan, network_configuration={}, node_settings={}, eth_json_rpc_endpoints=[], additional_services=[], metrics_enabled_services=[]):
    """Initiates the execution plan with the specified configuration."""

    validators = nodes.parse_nodes_from_dict(network_configuration["validators"], node_settings)
    full_nodes = nodes.parse_nodes_from_dict(network_configuration["full_nodes"], node_settings)
    seed_nodes = nodes.parse_nodes_from_dict(network_configuration["seed_nodes"], node_settings)
    all_nodes = validators + seed_nodes + full_nodes
    num_validators = len(validators)

    # Initialize modules, files, and genesis data
    node_modules = {node.el_type: import_module(f"./src/nodes/execution/{node.el_type}/config.star") for node in all_nodes if node.el_type not in node_modules}
    jwt_file, kzg_trusted_setup = execution.upload_global_files(plan, node_modules)
    evm_genesis_data = networks.get_genesis_data(plan)
    metrics_enabled_services = metrics_enabled_services[:]
    
    # Process nodes
    def configure_nodes(node_list, is_seed=False):
        client_configs, configs = [], {}
        for node in node_list:
            client_config = execution.generate_node_config(plan, node_modules, node, [el_enode_addrs] if not is_seed else [])
            client_configs.append(client_config)
            configs[node.cl_service_name] = beacond.create_node_config(plan, node, consensus_node_peering_info, node.el_service_name, jwt_file, kzg_trusted_setup)
        return client_configs, configs

    el_enode_addrs, consensus_node_peering_info, all_consensus_peering_info = [], [], {}

    # Deploy seed nodes
    seed_client_configs, seed_configs = configure_nodes(seed_nodes, is_seed=True)
    seed_clients = execution.deploy_nodes(plan, seed_client_configs)
    update_metrics_and_peering(seed_nodes, seed_clients)

    # Deploy full nodes
    full_client_configs, full_configs = configure_nodes(full_nodes)
    full_clients = execution.deploy_nodes(plan, full_client_configs, True)
    update_metrics_and_peering(full_nodes, full_clients)

    # Deploy validators
    validator_client_configs, validator_configs = configure_nodes(validators)
    validator_clients = execution.deploy_nodes(plan, validator_client_configs)
    update_metrics_and_peering(validators, validator_clients)

    # Establish RPC endpoints
    setup_rpc_endpoints(plan, eth_json_rpc_endpoints, full_clients)

    # Start additional services
    start_additional_services(plan, additional_services, metrics_enabled_services)

    plan.print("Successfully launched development network")


def update_metrics_and_peering(nodes, clients):
    """Updates metrics and consensus node peering information for each node."""
    for node in nodes:
        peer_info = beacond.get_peer_info(plan, node.cl_service_name)
        consensus_node_peering_info.append(peer_info)
        metrics_enabled_services.append({
            "name": node.cl_service_name,
            "service": clients[node.cl_service_name],
            "metrics_path": beacond.METRICS_PATH,
        })
        all_consensus_peering_info[node.cl_service_name] = peer_info


def setup_rpc_endpoints(plan, eth_json_rpc_endpoints, full_clients):
    """Configures RPC endpoints based on the provided endpoint type."""
    endpoint = eth_json_rpc_endpoints[0]
    endpoint_type = endpoint["type"]
    if endpoint_type == "nginx":
        plan.print("Launching RPCs for", endpoint_type)
        nginx.get_config(plan, endpoint["clients"])
    elif endpoint_type == "blutgang":
        plan.print("Launching blutgang")
        blutgang_config_template = read_file(constants.BLUTGANG_CONFIG_TEMPLATE_FILEPATH)
        blutgang.launch_blutgang(plan, blutgang_config_template, full_clients, endpoint["clients"], "kurtosis")
    else:
        plan.print("Invalid type for eth_json_rpc_endpoint")


def start_additional_services(plan, additional_services, metrics_enabled_services):
    """Starts additional services specified in the configuration."""
    prometheus_url = ""
    for s_dict in additional_services:
        s = service_module.parse_service_from_dict(s_dict)
        if s.name == "goomy-blob":
            start_goomy_blob(plan, s)
        elif s.name == "tx-fuzz":
            start_tx_fuzz(plan, s)
        elif s.name == "prometheus":
            prometheus_url = prometheus.start(plan, metrics_enabled_services)
        elif s.name == "grafana":
            grafana.start(plan, prometheus_url)
        elif s.name == "pyroscope":
            pyroscope.run(plan)
        elif s.name == "blockscout":
            blockscout.launch_blockscout(plan, full_clients, s.client, False)
        elif s.name == "otterscan":
            otterscan.launch_otterscan(plan, s.client)


def start_goomy_blob(plan, service):
    """Launches Goomy Blob Spammer with provided configuration."""
    ip = plan.get_service("endpoint_type").ip_address
    port = plan.get_service("endpoint_type").ports["http"].number
    goomy_blob.launch_goomy_blob(
        plan,
        constants.PRE_FUNDED_ACCOUNTS[next_free_prefunded_account],
        f"http://{ip}:{port}",
        {"goomy_blob_args": []}
    )


def start_tx_fuzz(plan, service):
    """Launches transaction fuzzing services."""
    service.replicas = service.replicas if "replicas" in service_dict else 1
    tx_fuzz.launch_tx_fuzzes(
        plan,
        service.replicas,
        next_free_prefunded_account,
        full_node_el_client_configs,
        full_clients,
        []
    )
