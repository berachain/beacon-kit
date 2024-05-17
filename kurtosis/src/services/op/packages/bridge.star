deps = import_module("deps.star")
contracts = import_module("contracts.star")
optimism = import_module("optimism.star")

# send bridges amount (in ether) from the L1 to the L2 for pk's wallet
def send(plan, files, env, amount, pk):
    bridge_contract_addr = "BRIDGE_CONTRACT=$(jq -r '.L1StandardBridgeProxy' {}/deployments/getting-started/l1.json)".format(
        contracts.PATH,
    )
    plan.run_sh(
        image = "ghcr.io/foundry-rs/foundry:latest",
        run = "{} && {} && cast send $BRIDGE_CONTRACT --value {}ether --private-key {} --legacy --rpc-url {}".format(
            deps.get(["bash", "jq"]),
            bridge_contract_addr,
            amount,
            pk,
            env.l1_rpc_url,
        ),
        files = {optimism.PATH: files.optimism},
    )
