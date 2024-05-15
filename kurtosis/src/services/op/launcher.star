env = import_module("env.star")
deployer = import_module("deploy.star")
types = import_module("types.star")

wallets = import_module("packages/wallets.star")
contracts = import_module("packages/contracts.star")

geth = import_module("components/geth.star")
batcher = import_module("components/batcher.star")
node = import_module("components/node.star")
proposer = import_module("components/proposer.star")

# TODO: Make these configurable
GETH_RPC_PORT="8545"
NODE_RPC_PORT="8547"
BATCHER_RPC_PORT="8548"
PROPOSER_RPC_PORT="8560"

def launch(plan, images, l1_rpc_url):
    l1 = types.get_l1(l1_rpc_url)
    files = types.get_files(plan)
    contracts.install(plan, files)
    e = env.get(plan, l1)

    wallets.fund(plan, e)
    deployer.build_deploy_config(plan, files, e, l1.chain_id)
    deployer.deploy_create2(plan, e)

    deployer.build_getting_started_dir(plan, e, files, l1)
    contracts.deploy_l1(plan, e, files)

    node.init(plan, images["node"], e, files)
    deployer.generate_jwt_secret(plan, files)
    geth.init(plan, images["geth"], files)

    # Deploy L2 Components
    geth_rpc_url = geth.launch(plan, images["geth"], l1, files)
    l2 = types.get_l2(geth_rpc_url)

    node_rpc_url = node.launch(plan, images["node"], files, e, l1, l2, NODE_RPC_PORT)
    batcher_rpc_url = batcher.launch(plan, images["batcher"], e, l1, l2, BATCHER_RPC_PORT, node_rpc_url)
    proposer_rpc_url = proposer.launch(plan, images["proposer"], files, e, l1, PROPOSER_RPC_PORT, node_rpc_url)

    # TODO: Bridge Tokens to Address
    