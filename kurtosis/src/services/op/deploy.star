images = import_module("constants/images.star")

bash = import_module("packages/bash.star")
deps = import_module("packages/deps.star")
contracts = import_module("packages/contracts.star")
optimism = import_module("packages/optimism.star")

# build_deploy_config sets the rollup deploy configuration defined
# in the environment variables to rollup.json
def build_deploy_config(plan, files, env, chain_id):
    plan.run_sh(
        description = "Building the L1 deploy config",
        image = images.FOUNDRY,
        run = bash.run([
            deps.get(["bash", "jq"]),
            "cd {}".format(contracts.PATH),
            "./scripts/getting-started/config.sh",
            "jq --argjson chainId {} \
  --argjson blockTime {} \
  '.l1ChainID = $chainId | .l1BlockTime = $blockTime | .finalizationPeriodSeconds = $blockTime' \
  deploy-config/getting-started.json > tmp.json".format(
                chain_id,
                env.l1_block_time,
            ),
            "mv tmp.json deploy-config/getting-started.json",
        ]),
        env_vars = {
            "GS_ADMIN_ADDRESS": env.admin_address,
            "GS_BATCHER_ADDRESS": env.batcher_address,
            "GS_PROPOSER_ADDRESS": env.proposer_address,
            "GS_SEQUENCER_ADDRESS": env.sequencer_address,
            "L1_RPC_URL": env.l1_rpc_url,
        },
        files = {optimism.PATH: files.optimism},
        store = [StoreSpec(src = optimism.PATH, name = files.optimism)],
    )

def build_getting_started_dir(plan, env, files, l1):
    plan.run_sh(
        image = images.ALPINE,
        run = bash.run([
            "cd {}/deployments".format(contracts.PATH),
            "mkdir getting-started",
            "echo -n {} > getting-started/.chainId".format(l1.chain_id),
            "echo -n '{}' > getting-started/.deploy",
        ]),
        files = {optimism.PATH: files.optimism},
        store = [StoreSpec(src = optimism.PATH, name = files.optimism)],
    )
