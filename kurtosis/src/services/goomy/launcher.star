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
    goomy_blob_params,
):
    config = get_config(
        funding_account,
        rpc_endpoint,
        goomy_blob_params,
    )
    plan.add_service(SERVICE_NAME, config)


def get_config(
    funding_account,
    rpc_endpoint,
    goomy_blob_args,
):
    """_summary_

    Args:
        funding_account (_type_): _description_
        rpc_endpoint (_type_): _description_
        goomy_blob_args (_type_): _description_

    Returns:
        _type_: _description_
    """

    address = "http://" + rpc_endpoint.ip_address + ":" + str(rpc_endpoint.ports["http"])
    goomy_cli_args = []
    goomy_cli_args.append(
        "-h " + address,
    )

    goomy_args = " ".join(goomy_blob_args)
    if goomy_args == "":
        goomy_args = "combined -b 2 -t 2 --max-pending 3"
    goomy_cli_args.append(goomy_args)

    return ServiceConfig(
        image=IMAGE_NAME,
        entrypoint=ENTRYPOINT_ARGS,
        cmd=[
            " && ".join(
                [
                    "apt-get update",
                    "apt-get install -y curl jq",
                    "./blob-spammer -p {0} {1}".format(
                        funding_account.private_key,
                        " ".join(goomy_cli_args),
                    ),
                ]
            )
        ],
        min_cpu=MIN_CPU,
        max_cpu=MAX_CPU,
        min_memory=MIN_MEMORY,
        max_memory=MAX_MEMORY,
        # node_selectors=node_selectors,
    )