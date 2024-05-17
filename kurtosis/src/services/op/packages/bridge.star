images = import_module("../constants/images.star")

bash = import_module("bash.star")
deps = import_module("deps.star")
contracts = import_module("contracts.star")
optimism = import_module("optimism.star")

# send bridges amount (in ether) from the L1 to the L2 for pk's wallet
def send(plan, files, env, amount, pk):
    plan.run_sh(
        description = "Bridging {} ether from L1 to L2".format(amount),
        image = images.FOUNDRY,
        run = bash.run([
            deps.get(["bash", "jq"]),
            "BRIDGE_CONTRACT=$(jq -r '.L1StandardBridgeProxy' {}/deployments/getting-started/l1.json)".format(
                contracts.PATH,
            ),
            "cast send $BRIDGE_CONTRACT --value {}ether --private-key {} --legacy --rpc-url {}".format(
                amount,
                pk,
                env.l1_rpc_url,
            ),
        ]),
        files = {optimism.PATH: files.optimism},
    )
