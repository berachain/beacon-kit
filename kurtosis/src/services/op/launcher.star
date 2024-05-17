execution = import_module("../../nodes/execution/execution.star")
constants = import_module("../../constants.star")

env = import_module("env.star")
deployer = import_module("deploy.star")
types = import_module("types.star")

optimism = import_module("packages/optimism.star")
wallets = import_module("packages/wallets.star")
contracts = import_module("packages/contracts.star")
bridge = import_module("packages/bridge.star")

geth = import_module("components/geth.star")
batcher = import_module("components/batcher.star")
node = import_module("components/node.star")
proposer = import_module("components/proposer.star")

# TODO: Make these configurable
GETH_RPC_PORT = "8545"
NODE_RPC_PORT = "8547"
BATCHER_RPC_PORT = "8548"
PROPOSER_RPC_PORT = "8560"

def launch(plan, images, l1_ip_address):
    l1 = types.get_l1(l1_ip_address)
    files = types.get_file_artifacts(plan)
    e = env.get(plan, l1)

    wallets.fund(plan, e)
    deployer.build_deploy_config(plan, files, e, l1.chain_id)
    contracts.install(plan, files)
    contracts.deploy_create2(plan, e)

    deployer.build_getting_started_dir(plan, e, files, l1)
    contracts.deploy_l1(plan, e, files)

    node.init(plan, e, files)
    geth.init(plan, images["geth"], files)

    # Deploy L2 Components
    geth_ip_address = geth.launch(plan, images["geth"], l1, files)
    l2 = types.get_l2(geth_ip_address)

    node_ip_address = node.launch(plan, images["node"], files, e, l1, l2, NODE_RPC_PORT)
    node_rpc_url = "http://{}:{}".format(node_ip_address, NODE_RPC_PORT)

    batcher_ip_address = batcher.launch(plan, images["batcher"], e, l1, l2, BATCHER_RPC_PORT, node_rpc_url)
    proposer_ip_address = proposer.launch(plan, images["proposer"], files, e, l1, PROPOSER_RPC_PORT, node_rpc_url)

    # Bridge Tokens to Address
    for wallet in constants.PRE_FUNDED_ACCOUNTS:
        bridge.send(plan, files, e, 10, wallet.private_key)
