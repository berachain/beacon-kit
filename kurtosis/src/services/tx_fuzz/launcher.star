shared_utils = import_module("github.com/ethpandaops/ethereum-package/src/shared_utils/shared_utils.star")
constants = import_module("../../constants.star")
execution = import_module("../../nodes/execution/execution.star")
SERVICE_NAME = "tx-fuzz"

# The min/max CPU/memory that tx-spammer can use
MIN_CPU = 100
MAX_CPU = 2000
MIN_MEMORY = 128
MAX_MEMORY = 1024

def launch_tx_fuzzes(plan, amount, next_free_prefunded_account, full_node_el_client_configs, full_node_el_clients, tx_spammer_extra_args):
    tx_fuzz_service_configs = {}
    for i in range(amount):
        full_node_service_name = full_node_el_client_configs[i % len(full_node_el_client_configs)]["name"]
        fuzzing_node = full_node_el_clients[full_node_service_name]
        tx_fuzz_config = get_config(
            constants.PRE_FUNDED_ACCOUNTS[next_free_prefunded_account].private_key,
            "http://{}:{}".format(fuzzing_node.ip_address, execution.RPC_PORT_NUM),
            tx_spammer_extra_args,
        )
        tx_fuzz_service_configs[SERVICE_NAME + "-" + str(i)] = tx_fuzz_config
        next_free_prefunded_account += 1

    plan.add_services(tx_fuzz_service_configs)
    return next_free_prefunded_account

def launch_tx_fuzzes_gang(plan, amount, next_free_prefunded_account, tx_spammer_extra_args):
    blutgang_ip = plan.get_service("blutgang").ip_address
    blutgang_port = plan.get_service("blutgang").ports["http"].number
    tx_fuzz_service_configs = {}
    for i in range(amount):
        tx_fuzz_config = get_config(
            constants.PRE_FUNDED_ACCOUNTS[next_free_prefunded_account].private_key,
            "http://{}:{}".format(blutgang_ip, blutgang_port),
            tx_spammer_extra_args,
        )
        tx_fuzz_service_configs[SERVICE_NAME + "-" + str(i)] = tx_fuzz_config
        next_free_prefunded_account += 1

    plan.add_services(tx_fuzz_service_configs)
    return next_free_prefunded_account

def launch_tx_fuzz(
        plan,
        id,
        prefunded_private_key,
        el_uri,
        tx_spammer_extra_args):
    config = get_config(
        prefunded_private_key,
        el_uri,
        tx_spammer_extra_args,
    )
    plan.add_service(
        SERVICE_NAME + "-" + str(id),
        config,
    )

def get_config(
        prefunded_private_key,
        el_uri,
        tx_spammer_extra_args):
    tx_spammer_image = "ethpandaops/tx-fuzz:master"

    entrypoint = [
        "/bin/sh",
        "-c",
    ]

    # A sleep is added to ensure the full node is up in single-node deployments
    cmd = " ".join([
        "sleep",
        "5",
        "&&",
        "/tx-fuzz.bin",
        "spam",
        "--rpc={}".format(el_uri),
        "--sk={0}".format(prefunded_private_key),
        "--accounts=100",
        "--txcount=100",
        "--slot-time=2",
    ])

    if len(tx_spammer_extra_args) > 0:
        cmd.extend([param for param in tx_spammer_extra_args])

    return ServiceConfig(
        image = tx_spammer_image,
        cmd = [cmd],
        entrypoint = entrypoint,
        min_cpu = MIN_CPU,
        max_cpu = MAX_CPU,
        min_memory = MIN_MEMORY,
        max_memory = MAX_MEMORY,
    )
