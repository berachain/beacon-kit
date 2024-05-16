geth = import_module("geth.star")
contracts = import_module("../packages/contracts.star")
optimism = import_module("../packages/optimism.star")

NAME = "op-node"
PATH="/optimism/op-node"

def init(plan, env, files):
    plan.run_sh(
        image = "golang:latest",
        run = "cd {} && go mod tidy && cd {} && go run cmd/main.go genesis l2 \
  --deploy-config {}/deploy-config/getting-started.json \
  --l1-deployments {}/deployments/getting-started/l1.json \
  --outfile.l2 genesis.json \
  --outfile.rollup rollup.json \
  --l1-rpc {} && cp genesis.json {}".format(
            optimism.PATH,
            PATH,
            contracts.PATH,
            contracts.PATH,
            env.l1_rpc_url,
            geth.GENESIS_PATH,
        ),
        files = {
            optimism.PATH: files.optimism,
            geth.PATH: files.op_geth,
        },
        store = [
            StoreSpec(src=optimism.PATH, name=files.optimism),
            StoreSpec(src=geth.PATH, name=files.op_geth),
        ],
    )

def launch(plan, image, files, env, l1, l2, node_rpc_port):
    node_rpc_url = "http://{}:{}".format(NAME, node_rpc_port)
    service = plan.add_service(
        name = NAME,
        config = ServiceConfig(
            image = image,
            cmd = [
                "op-node",
                "--l1={}".format(l1.rpc_url),
                "--l1.rpckind={}".format(l1.rpc_kind),
                "--l1.trustrpc=true",
                "--l2={}".format(l2.rpc_url),
                "--l2.jwt-secret={}".format(geth.JWT_PATH),
                "--sequencer.enabled",
                "--sequencer.l1-confs=5",
                "--verifier.l1-confs=4",
                "--rollup.config={}/rollup.json".format(PATH),
                "--rpc.addr={}".format(NAME),
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
            files = {
                optimism.PATH: files.optimism,
                geth.PATH: files.op_geth,
            },
        ),
    )

    return service.ports["rpc"].url
