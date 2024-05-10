def get(
    plan,
    files,
    l1,
    private_key="fffdbb37105441e14b0ee6330d855d8504ff39e705c3afa8f859ac9865f99306",
):
    wallets = get_wallets(plan, files)
    return {
        "GS_ADMIN_ADDRESS": wallets[0],
        "GS_ADMIN_PRIVATE_KEY": wallets[1],
        "GS_BATCHER_ADDRESS": wallets[2],
        "GS_BATCHER_PRIVATE_KEY": wallets[3],
        "GS_PROPOSER_ADDRESS": wallets[4],
        "GS_PROPOSER_PRIVATE_KEY": wallets[5],
        "GS_SEQUENCER_ADDRESS": wallets[6],
        "GS_SEQUENCER_PRIVATE_KEY": wallets[7],
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


# Returns wallets (address, pk pair) in the order of:
#   GS_ADMIN, GS_BATCHER, GS_PROPOSER, GS_SEQUENCER
def get_wallets(plan, files):
    wallets=plan.run_sh(
        run="./contracts-bedrock/scripts/getting-started/wallets.sh | grep '_ADDRESS\|_PRIVATE_KEY' | cut -d '=' -f 2"
        files={"/contracts-bedrock/": files.contracts}
    )

    return wallets.output.split("\n")


def get_salt(plan):
    return plan.run_sh(
        run="openssl rand -hex 32"
    ).output
