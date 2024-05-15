wallets = import_module("./packages/wallets.star")

def get(
    plan,
    l1,
    private_key="fffdbb37105441e14b0ee6330d855d8504ff39e705c3afa8f859ac9865f99306",
):
    return {
        "GS_ADMIN_ADDRESS": wallets.get_by_index(plan, 1),
        "GS_ADMIN_PRIVATE_KEY": wallets.get_by_index(plan, 2),
        "GS_BATCHER_ADDRESS": wallets.get_by_index(plan, 3),
        "GS_BATCHER_PRIVATE_KEY": wallets.get_by_index(plan, 4),
        "GS_PROPOSER_ADDRESS": wallets.get_by_index(plan, 5),
        "GS_PROPOSER_PRIVATE_KEY": wallets.get_by_index(plan, 6),
        "GS_SEQUENCER_ADDRESS": wallets.get_by_index(plan, 7),
        "GS_SEQUENCER_PRIVATE_KEY": wallets.get_by_index(plan, 8),
        "L1_RPC_KIND": l1.rpc_kind,
        "L1_BLOCK_TIME": l1.block_time,
        "L1_RPC_URL": l1.rpc_url,
        "IMPL_SALT": get_salt(plan),
        "DEPLOYMENT_CONTEXT": "getting-started",
        "TENDERLY_PROJECT": "",
        "TENDERLY_USERNAME": "",
        "ETHERSCAN_API_KEY": "",
        "PRIVATE_KEY": private_key
    }

def get_salt(plan):
    return plan.run_sh(
        image="alpine/openssl:latest",
        run="openssl rand -hex 32"
    ).output
