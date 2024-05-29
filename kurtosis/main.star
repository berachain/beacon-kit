el_cl_genesis_data_generator = import_module(
    "github.com/kurtosis-tech/ethereum-package/src/prelaunch_data_generator/el_cl_genesis/el_cl_genesis_generator.star",
)

execution = import_module("./src/nodes/execution/execution.star")
service_module = import_module("./src/services/service.star")
beacond = import_module("./src/nodes/consensus/beacond/launcher.star")
networks = import_module("./src/networks/networks.star")
port_spec_lib = import_module("./src/lib/port_spec.star")
nodes = import_module("./src/nodes/nodes.star")
nginx = import_module("./src/services/nginx/nginx.star")
constants = import_module("./src/constants.star")
goomy_blob = import_module("./src/services/goomy/launcher.star")
prometheus = import_module("./src/observability/prometheus/prometheus.star")
grafana = import_module("./src/observability/grafana/grafana.star")
pyroscope = import_module("./src/observability/pyroscope/pyroscope.star")
tx_fuzz = import_module("./src/services/tx_fuzz/launcher.star")
blutgang = import_module("./src/services/blutgang/launcher.star")

def run(plan, validators, full_nodes = [], eth_json_rpc_endpoints = [], boot_sequence = {"type": "sequential"}, additional_services = [], metrics_enabled_services = []):
    """
    Initiates the execution plan with the specified number of validators and arguments.

    Args:
    plan: The execution plan to be run.
    args: Additional arguments to configure the plan. Defaults to an empty dictionary.
    """

    next_free_prefunded_account = 0
    validators = nodes.parse_nodes_from_dict(validators)
    full_nodes = nodes.parse_nodes_from_dict(full_nodes)
    num_validators = len(validators)

    # 1. Initialize EVM genesis data
    evm_genesis_data = networks.get_genesis_data(plan)

    node_modules = {}
    for node in validators:
        if node.el_type not in node_modules.keys():
            node_path = "./src/nodes/execution/{}/config.star".format(node.el_type)
            node_module = import_module(node_path)
            node_modules[node.el_type] = node_module

    for node in full_nodes:
        if node.el_type not in node_modules.keys():
            node_path = "./src/nodes/execution/{}/config.star".format(node.el_type)
            node_module = import_module(node_path)
            node_modules[node.el_type] = node_module

    # 2. Upload files
    jwt_file, kzg_trusted_setup = execution.upload_global_files(plan, node_modules)

    # 3. Perform genesis ceremony
    if boot_sequence["type"] == "sequential":
        beacond.perform_genesis_ceremony(plan, validators, jwt_file)
    else:
        beacond.perform_genesis_ceremony_parallel(plan, validators, jwt_file)

    el_enode_addrs = []
    metrics_enabled_services = metrics_enabled_services[:]

    consensus_node_peering_info = []

    # 4. Start network validators
    validator_node_el_clients = []
    for n, validator in enumerate(validators):
        el_client = execution.create_node(plan, node_modules, validator, "validator", n, el_enode_addrs)
        validator_node_el_clients.append(el_client)
        el_enode_addrs.append(el_client["enode_addr"])

        # As ethereumjs currently does not support metrics, we only add the metrics path for other clients
        if validator.el_type != "ethereumjs":
            metrics_enabled_services.append({
                "name": el_client["name"],
                "service": el_client["service"],
                "metrics_path": node_modules[validator.el_type].METRICS_PATH,
            })

        # 4b. Launch CL
        beacond_service = beacond.create_node(plan, validator.cl_image, consensus_node_peering_info[:n], el_client["name"], jwt_file, kzg_trusted_setup, n)
        peer_info = beacond.get_peer_info(plan, beacond_service.name)
        consensus_node_peering_info.append(peer_info)
        if validator.el_type != "ethereumjs":
            metrics_enabled_services.append({
                "name": beacond_service.name,
                "service": beacond_service,
                "metrics_path": beacond.METRICS_PATH,
            })

    # 5. Start full nodes (rpcs)
    full_node_configs = {}
    full_node_el_clients = []
    for n, full in enumerate(full_nodes):
        el_client = execution.create_node(plan, node_modules, full, "full", n, el_enode_addrs)
        full_node_el_clients.append(el_client)
        el_enode_addrs.append(el_client["enode_addr"])

        if full.el_type != "ethereumjs":
            metrics_enabled_services.append({
                "name": el_client["name"],
                "service": el_client["service"],
                "metrics_path": node_modules[full.el_type].METRICS_PATH,
            })

        # 4b. Launch CL
        cl_service_name = "cl-full-beaconkit-{}".format(n)
        full_node_config = beacond.create_full_node_config(plan, full.cl_image, consensus_node_peering_info, el_client["name"], jwt_file, kzg_trusted_setup, n)
        full_node_configs[cl_service_name] = full_node_config

    if full_node_configs != {}:
        services = plan.add_services(
            configs = full_node_configs,
        )

        for name, service in services.items():
            # excluding ethereumjs from metrics as it is the last full node in the args file beaconkit-all.yaml, TO-DO: to improve this later
            if name != cl_service_name:
                metrics_enabled_services.append({
                    "name": name,
                    "service": service,
                    "metrics_path": beacond.METRICS_PATH,
                })

    # 6. Start RPCs
    #  check the "type" value inside of rpc_endpoints to determine which rpc endpoint to launch

    # Get only the first rpc endpoint
    eth_json_rpc_endpoint = eth_json_rpc_endpoints[0]
    endpoint_type = eth_json_rpc_endpoint["type"]
    plan.print("RPC Endpoint Type:", endpoint_type)
    if endpoint_type == "nginx":
        plan.print("Launching RPCs for ", endpoint_type)
        nginx.get_config(plan, eth_json_rpc_endpoint["clients"])

    elif endpoint_type == "blutgang":
        plan.print("Launching blutgang")
        blutgang_config_template = read_file(
            constants.BLUTGANG_CONFIG_TEMPLATE_FILEPATH,
        )
        blutgang.launch_blutgang(
            plan,
            blutgang_config_template,
            full_node_el_clients,
            eth_json_rpc_endpoint["clients"],
            "kurtosis",
        )

    else:
        plan.print("Invalid type for eth_json_rpc_endpoint")

    # 7. Start additional services
    for s_dict in additional_services:
        s = service_module.parse_service_from_dict(s_dict)
        if s.name == "goomy_blob":
            plan.print("Launching Goomy the Blob Spammer")
            rpc_endpoint_goomy_blob = plan.get_service(endpoint_type).ports["http"].url
            plan.print("Launching goomy blob for rpc endpoint: ", rpc_endpoint_goomy_blob)
            goomy_blob_args = {"goomy_blob_args": []}
            goomy_blob.launch_goomy_blob(
                plan,
                constants.PRE_FUNDED_ACCOUNTS[next_free_prefunded_account],
                rpc_endpoint_goomy_blob,
                goomy_blob_args,
            )
            next_free_prefunded_account += 1
            plan.print("Successfully launched goomy the blob spammer")
        elif s.name == "tx-fuzz":
            plan.print("Launching tx-fuzz")
            fuzzing_node = validator_node_el_clients[0]["service"]
            if "replicas" not in s_dict:
                s.replicas = 1

            for i in range(s.replicas):
                if i > 0:
                    fuzzing_node = full_node_el_clients[i % len(full_node_el_clients)]["service"]
                tx_fuzz.launch_tx_fuzz(
                    plan,
                    i,
                    constants.PRE_FUNDED_ACCOUNTS[next_free_prefunded_account].private_key,
                    "http://{}:{}".format(fuzzing_node.ip_address, execution.RPC_PORT_NUM),
                    [],
                )
                next_free_prefunded_account += 1
        elif s.name == "prometheus":
            prometheus_url = prometheus.start(plan, metrics_enabled_services)
            if "grafana" in additional_services:
                grafana.start(plan, prometheus_url)
            if "pyroscope" in additional_services:
                pyroscope.run(plan)

    plan.print("Successfully launched development network")
