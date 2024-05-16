contracts = import_module("../packages/contracts.star")
optimism = import_module("../packages/optimism.star")

NAME = "op-proposer"

def launch(plan, image, files, env, l1, proposer_rpc_port, node_rpc_url):
    proposer_rpc_url = "http://{}:{}".format(NAME, proposer_rpc_port)
    service = plan.add_service(
        name = NAME,
        config = ServiceConfig(
            image = image,
            cmd = [
                "op-proposer",
                "--l1-eth-rpc={}".format(l1.rpc_url),
                "--poll-interval=12s",
                "--rpc.addr={}".format(NAME),
                "--rpc.port={}".format(proposer_rpc_port),
                "--rollup-rpc={}".format(node_rpc_url),
                "--l2oo-address=$(jq -r '.L2OutputOracleProxy' {}/deployments/getting-started/l1.json)".format(contracts.PATH),
                "--private-key={}".format(env.proposer_pk),
            ],
            ports = {
                "rpc": PortSpec(
                    number = int(proposer_rpc_port),
                    url = proposer_rpc_url,
                ),
            },
            files = {optimism.PATH: files.optimism},
        ),
    )

    return service.ports["rpc"].url
