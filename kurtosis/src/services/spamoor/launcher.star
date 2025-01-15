constants = import_module("../../constants.star")
execution = import_module("../../nodes/execution/execution.star")

SERVICE_NAME = "spamoor"
IMAGE_NAME = "ethpandaops/spamoor:latest"

# The min/max CPU/memory that spamoor can use
MIN_CPU = 100
MAX_CPU = 500
MIN_MEMORY = 20
MAX_MEMORY = 300

def launch_spamoors(plan, replicas, next_free_prefunded_account, full_node_el_client_configs, full_node_el_clients):
    configs = {}
    for i in range(replicas):
        rpc_node = full_node_el_clients[full_node_el_client_configs[i % len(full_node_el_client_configs)]["name"]]
        configs[SERVICE_NAME + "-" + str(i)] = ServiceConfig(
            image = IMAGE_NAME,
            entrypoint = ["/bin/sh", "-c"],
            cmd = ["./spamoor blob-combined -p {0} -b 3 -t 10 --max-pending 100 -h {1}".format(
                constants.PRE_FUNDED_ACCOUNTS[next_free_prefunded_account].private_key,
                "http://{}:{}".format(rpc_node.ip_address, execution.RPC_PORT_NUM),
            )],
            min_cpu = MIN_CPU,
            max_cpu = MAX_CPU,
            min_memory = MIN_MEMORY,
            max_memory = MAX_MEMORY,
        )
        next_free_prefunded_account += 1

    plan.add_services(configs)
    plan.print(configs)

    return next_free_prefunded_account
