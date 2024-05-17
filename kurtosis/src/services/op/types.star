ports = import_module("constants/ports.star")
optimism = import_module("packages/optimism.star")
wallets = import_module("packages/wallets.star")
geth = import_module("components/geth.star")

def get_file_artifacts(plan):
    optimism_dir = optimism.clone(plan)
    return struct(
        optimism = optimism_dir,
        wallets = wallets.get(plan, optimism_dir),
        op_geth = geth.generate_jwt_secret(plan),
        l2 = "l2",  # gets set in deploy
    )

def get_l1(
        ip_address,
        rpc_kind = "any",
        block_time = "6",
        chain_id = "80087"):
    return struct(
        ip_address = ip_address,
        rpc_kind = rpc_kind,
        rpc_url = "http://{}:{}".format(ip_address, ports.L1_ETH_RPC),
        block_time = block_time,
        chain_id = chain_id,
    )

def get_l2(ip_address):
    return struct(
        rpc_url = "http://{}:{}".format(ip_address, ports.L1_ETH_RPC),
        auth_rpc_url = "http://{}:{}".format(ip_address, ports.L1_ENGINE_RPC),
    )
