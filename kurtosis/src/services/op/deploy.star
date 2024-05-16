deps = import_module("packages/deps.star")
contracts = import_module("packages/contracts.star")

# build_deploy_config sets the rollup deploy configuration defined
# in the environment variables to rollup.json
def build_deploy_config(plan, files, env, chain_id):
    replace_cmd = "jq --argjson chainId {} \
  --argjson blockTime {} \
  '.l1ChainID = $chainId | .l1BlockTime = $blockTime | .finalizationPeriodSeconds = $blockTime' \
  deploy-config/getting-started.json > tmp.json && mv tmp.json deploy-config/getting-started.json".format(chain_id, env.l1_block_time)

    plan.run_sh(
        image = "ghcr.io/foundry-rs/foundry:latest",
        run = "{} && cd {} && scripts/getting-started/config.sh && {}".format(
            deps.get(["bash", "jq"]),
            contracts.PATH,
            replace_cmd,
        ),
        env_vars = {
            "GS_ADMIN_ADDRESS": env.admin_address,
            "GS_BATCHER_ADDRESS": env.batcher_address,
            "GS_PROPOSER_ADDRESS": env.proposer_address,
            "GS_SEQUENCER_ADDRESS": env.sequencer_address,
            "L1_RPC_URL": env.l1_rpc_url,
        },
        files = {files.optimism: files.optimism},
        store = [StoreSpec(src = files.optimism, name = files.optimism)],
    )

def build_getting_started_dir(plan, env, files, l1):
    chain_id_cmd = "echo -n {} > getting-started/.chainId".format(l1.chain_id)
    deploy_cmd = "echo -n '{}' > getting-started/.deploy"
    plan.run_sh(
        image = "alpine:latest",
        run = "cd {}/deployments && mkdir getting-started && {} && {}".format(
            contracts.PATH,
            chain_id_cmd,
            deploy_cmd,
        ),
        files = {files.optimism: files.optimism},
        store = [StoreSpec(src = files.optimism, name = files.optimism)],
    )

def generate_jwt_secret(plan, files):
    plan.run_sh(
        image = "alpine/openssl:latest",
        run = "openssl rand -hex 32 > /config/jwt.txt",
        files = {"/config": files.config},
        store = [StoreSpec(src = "/config", name = files.config)],
    )
