contracts = import_module("../packages/contracts.star")

PROPOSER_RPC_PORT_BASE = "op-proposer"

def launch(plan, image, files, env, l1, proposer_rpc_port, node_rpc_url):
    proposer_rpc_url = "http://{}:{}".format(PROPOSER_RPC_PORT_BASE, proposer_rpc_port)
    service = plan.add_service(
        name = "op-proposer",
        config = ServiceConfig(
            image = image,
            cmd = [
                "op-proposer",
                "--l1-eth-rpc={}".format(l1.rpc_url),
                "--poll-interval=12s",
                "--rpc.addr={}".format(PROPOSER_RPC_PORT_BASE),
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
            files = {files.optimism: files.optimism},
        ),
    )

    return service.ports["rpc"].url
