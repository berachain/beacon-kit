SERVICE_NAME = "goomy-blob-spammer"
IMAGE_NAME = "ethpandaops/goomy-blob:master-52526aa"

ENTRYPOINT_ARGS = ["/bin/sh", "-c"]

# The min/max CPU/memory that goomy can use
MIN_CPU = 100
MAX_CPU = 500
MIN_MEMORY = 20
MAX_MEMORY = 300

def launch_goomy_blob(
        plan,
        funding_account,
        rpc_endpoint,
        goomy_blob_params):
    config = get_config(
        funding_account,
        rpc_endpoint,
        goomy_blob_params,
    )
    plan.add_service(SERVICE_NAME, config)
    plan.print(config)

def get_config(
        funding_account,
        rpc_endpoint,
        goomy_blob_args):
    goomy_cli_args = []
    goomy_args = " ".join(goomy_blob_args)
    if goomy_args == "":
        goomy_args = "combined -b 2 -t 1 --max-pending 3"
    goomy_cli_args.append(goomy_args)

    blob_cmd = "./blob-spammer -p {0} combined -b 50 -t 50 --max-pending 50 -h {1}".format(
        funding_account.private_key,
        rpc_endpoint,
    )

    return ServiceConfig(
        image = IMAGE_NAME,
        entrypoint = ENTRYPOINT_ARGS,
        cmd = [
            " && ".join(
                [
                    "apt-get update",
                    "apt-get install -y curl",
                    blob_cmd,
                ],
            ),
        ],
        min_cpu = MIN_CPU,
        max_cpu = MAX_CPU,
        min_memory = MIN_MEMORY,
        max_memory = MAX_MEMORY,
    )
