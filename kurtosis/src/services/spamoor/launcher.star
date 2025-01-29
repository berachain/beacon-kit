SERVICE_NAME = "spamoor"
IMAGE_NAME = "ethpandaops/spamoor:latest"

ENTRYPOINT_ARGS = ["/bin/sh", "-c"]

# The min/max CPU/memory that spamoor can use
MIN_CPU = 100
MAX_CPU = 500
MIN_MEMORY = 20
MAX_MEMORY = 300

def launch_spamoor(plan, funding_account, rpc_endpoint):
    config = get_config(funding_account, rpc_endpoint)
    plan.add_service(SERVICE_NAME, config)
    plan.print(config)

def get_config(funding_account, rpc_endpoint):
    blob_cmd = "./spamoor blob-combined -p {0} -b 6 -t 3 --max-pending 9 -h {1}".format(
        funding_account.private_key,
        rpc_endpoint,
    )

    return ServiceConfig(
        image = IMAGE_NAME,
        entrypoint = ENTRYPOINT_ARGS,
        cmd = [blob_cmd],
        min_cpu = MIN_CPU,
        max_cpu = MAX_CPU,
        min_memory = MIN_MEMORY,
        max_memory = MAX_MEMORY,
    )
