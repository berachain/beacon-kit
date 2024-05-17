images = import_module("constants/images.star")
wallets = import_module("packages/wallets.star")

def get(
        plan,
        l1,
        private_key = "fffdbb37105441e14b0ee6330d855d8504ff39e705c3afa8f859ac9865f99306"):
    return struct(
        admin_address = wallets.get_by_index(plan, 1),
        admin_pk = wallets.get_by_index(plan, 2),
        batcher_address = wallets.get_by_index(plan, 3),
        batcher_pk = wallets.get_by_index(plan, 4),
        proposer_address = wallets.get_by_index(plan, 5),
        proposer_pk = wallets.get_by_index(plan, 6),
        sequencer_address = wallets.get_by_index(plan, 7),
        sequencer_pk = wallets.get_by_index(plan, 8),
        l1_rpc_kind = l1.rpc_kind,
        l1_block_time = l1.block_time,
        l1_rpc_url = l1.rpc_url,
        impl_salt = generate_salt(plan),
        deployment_context = "getting-started",
        tenderly_project = "",
        tenderly_username = "",
        etherscan_api_key = "",
        pk = private_key,
    )

def generate_salt(plan):
    return plan.run_sh(
        image = images.ALPINE_OPENSSL,
        run = "openssl rand -hex 32",
    ).output
