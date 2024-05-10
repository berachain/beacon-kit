# fund_wallets funds component wallets with 10 ether
def fund_wallets(plan, env):
    def build_fund_wallet_cmd(address):
        return "cast send --private-key {0} {1} --value 10ether --rpc-url {2} --legacy".format(
            env["PRIVATE_KEY"],
            address,
            env["L1_NODE_RPC_URL"],
        )
    plan.run_sh(
        image="ghcr.io/foundry-rs/foundry:latest",
        run="{0} && {1} && {2}".format(
            build_fund_wallet_cmd(env["GS_ADMIN_ADDRESS"]),
            build_fund_wallet_cmd(env["GS_BATCHER_ADDRESS"]),
            build_fund_wallet_cmd(env["GS_PROPOSER_ADDRESS"]),
        ),
        env_vars=env,
    )


# build_deploy_config sets the rollup deploy configuration defined
# in the environment variables to rollup.json
def build_deploy_config(plan, env, chain_id):
    replace_cmd = 'jq --argjson chainId {0} \
  --argjson blockTime {1} \
  ".l1ChainID = $chainId | .l1BlockTime = $blockTime | .finalizationPeriodSeconds = $blockTime" \
  deploy-config/getting-started.json > tmp.json && mv tmp.json deploy.json'.format(chain_id, env["L1_BLOCK_TIME"])
    
    plan.run_sh(
        run="./scripts/getting-started/config.sh && {0}".format(replace_cmd),
        env_vars=env,
        store=[
            StoreSpec(
                src="deploy.json",
                name="/deploy/deploy.json",
            )
        ],
    )

def deploy_l1_contracts(plan, env, files):
    plan.run_sh(
        image="ghcr.io/foundry-rs/foundry:latest",
        run="forge script scripts/Deploy.s.sol:Deploy --private-key {0} --broadcast --rpc-url {1} --legacy".format(
            env["GS_ADMIN_PRIVATE_KEY"],
            env["L1_NODE_RPC_URL"],
        ) + "cp /contracts-bedrock/deployments/getting-started/.deploy l1.json",
        files={"/contracts-bedrock/": files.contracts},
        store=[
            StoreSpec(
                src="deployed.json",
                name="/deployed.json",
            )
        ],
    )

def generate_jwt_secret(plan):
    plan.run_sh(
        run="openssl rand -hex 32 > jwt.txt",
        store=[
            StoreSpec(
                src="jwt.txt",
                name="/config/jwt.txt"
            )
        ]
    )