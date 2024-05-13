# fund_wallets funds component wallets with 10 ether
def fund_wallets(plan, env):
    def build_fund_wallet_cmd(address):
        return "cast send --private-key {} {} --value 10ether --rpc-url {} --legacy".format(
            env["PRIVATE_KEY"],
            address,
            env["L1_RPC_URL"],
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
    replace_cmd = 'jq --argjson chainId {} \
  --argjson blockTime {} \
  ".l1ChainID = $chainId | .l1BlockTime = $blockTime | .finalizationPeriodSeconds = $blockTime" \
  deploy-config/getting-started.json > tmp.json && mv tmp.json deploy.json'.format(chain_id, env["L1_BLOCK_TIME"])
    
    plan.run_sh(
        run="./scripts/getting-started/config.sh && {}".format(replace_cmd),
        env_vars=env,
        store=[
            StoreSpec(
                src="deploy.json",
                name="/deploy/deploy.json",
            )
        ],
    )
    
def deploy_create2(plan, env):
    codesize_output = plan.run_sh(
        image="ghcr.io/foundry-rs/foundry:latest",
        run="cast codesize 0x4e59b44847b379578588920cA78FbF26c0B4956C --rpc-url {}".format(
            env["L1_RPC_URL"],
        ),
    ).output.strip()

    if codesize_output == "0":
        plan.run_sh(
            image="ghcr.io/foundry-rs/foundry:latest",
            run="cast send --private-key {} 0x3fAB184622Dc19b6109349B94811493BF2a45362 --value 1ether --rpc-url {} --legacy && \
                cast publish --rpc-url {} 0xf8a58085174876e800830186a08080b853604580600e600039806000f350fe7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe03601600081602082378035828234f58015156039578182fd5b8082525050506014600cf31ba02222222222222222222222222222222222222222222222222222222222222222a02222222222222222222222222222222222222222222222222222222222222222 \
                ".format(
                env["PRIVATE_KEY"],
                env["L1_RPC_URL"],
                env["L1_RPC_URL"],
            ),
            env_vars=env,
        )
        codesize_output = plan.run_sh(
            image="ghcr.io/foundry-rs/foundry:latest",
            run="cast codesize 0x4e59b44847b379578588920cA78FbF26c0B4956C --rpc-url {}".format(
                env["L1_RPC_URL"],
            ),
        ).output.strip()
        plan.verify(value=codesize_output, assertion="!=", target_value="0")
    
    plan.verify(value=codesize_output, assertion="==", target_value="69")

def deploy_l1_contracts(plan, env, files):
    plan.run_sh(
        image="ghcr.io/foundry-rs/foundry:latest",
        run="forge script scripts/Deploy.s.sol:Deploy --private-key {} --broadcast --rpc-url {} --legacy".format(
            env["GS_ADMIN_PRIVATE_KEY"],
            env["L1_RPC_URL"],
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
        image="alpine:latest",
        run="apk add openssl && openssl rand -hex 32 > jwt.txt",
        store=[
            StoreSpec(
                src="jwt.txt",
                name="/config/jwt.txt"
            )
        ]
    )