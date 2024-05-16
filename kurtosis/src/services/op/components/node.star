contracts = import_module("../packages/contracts.star")

NODE_RPC_URL_BASE = "op-node"

def init(plan, image, env, files):
    plan.run_sh(
        image = image,
        run = "genesis l2            --deploy-config {}/deploy-config/getting-started.json             --l1-deployments {}/deployments/getting-started/l1.json             --outfile.l2 genesis.json             --outfile.rollup rollup.json             --l1-rpc {}".format(
            contracts.PATH,
            contracts.PATH,
            env.l1_rpc_url,
        ),
        files = {files.optimism: files.optimism},
        store = [
            StoreSpec(src = "genesis.json", name = files.l2),
            StoreSpec(src = "rollup.json", name = files.rollup),
        ],
    )

def launch(plan, image, files, env, l1, l2, node_rpc_port):
    node_rpc_url = "http://{}:{}".format(NODE_RPC_URL_BASE, node_rpc_port)
    service = plan.add_service(
        name = "op-node",
        config = ServiceConfig(
            image = image,
            cmd = [
                "op-node",
                "--l1={}".format(l1.rpc_url),
                "--l1.rpckind={}".format(l1.rpc_kind),
                "--l1.trustrpc=true",
                "--l2={}".format(l2.rpc_url),
                "--l2.jwt-secret=./config/jwt.txt",
                "--sequencer.enabled",
                "--sequencer.l1-confs=5",
                "--verifier.l1-confs=4",
                "--rollup.config=./rollup/rollup.json",
                "--rpc.addr={}".format(NODE_RPC_URL_BASE),
                "--rpc.port={}".format(node_rpc_port),
                "--p2p.disable",
                "--rpc.enable-admin",
                "--p2p.sequencer.key={}".format(env.sequencer_pk),
            ],
            ports = {
                "rpc": PortSpec(
                    number = int(node_rpc_port),
                    url = node_rpc_url,
                ),
            },
            files = {"/config": files.config, "/rollup": files.rollup},
        ),
    )

    return service.ports["rpc"].url
