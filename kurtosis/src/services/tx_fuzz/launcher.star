shared_utils = import_module("github.com/kurtosis-tech/ethereum-package/src/shared_utils/shared_utils.star")
SERVICE_NAME = "tx-fuzz"

# The min/max CPU/memory that tx-spammer can use
MIN_CPU = 100
MAX_CPU = 1000
MIN_MEMORY = 20
MAX_MEMORY = 300

def launch_tx_fuzz(
        plan,
        prefunded_private_key,
        el_uri,
        tx_spammer_extra_args):
    config = get_config(
        prefunded_private_key,
        el_uri,
        tx_spammer_extra_args,
    )
    plan.add_service(SERVICE_NAME, config)

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
        "3",
        "&&",
        "/tx-fuzz.bin",
        "spam",
        "--rpc={}".format(el_uri),
        "--sk={0}".format(prefunded_private_key),
        "--accounts=100",
        "--txcount=100",
        "--slot-time=3",
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
