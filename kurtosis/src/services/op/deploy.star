# build_deploy_config sets the rollup deploy configuration defined
# in the environment variables to rollup.json
def build_deploy_config(plan, files, env, chain_id):
    add_dependencies = "apk add bash jq"
    replace_cmd = "jq --argjson chainId {} \
  --argjson blockTime {} \
  '.l1ChainID = $chainId | .l1BlockTime = $blockTime | .finalizationPeriodSeconds = $blockTime' \
  deploy-config/getting-started.json > tmp.json && mv tmp.json deploy-config/getting-started.json".format(chain_id, env["L1_BLOCK_TIME"])
    
    plan.run_sh(
        image="ghcr.io/foundry-rs/foundry:latest",
        run="{} && cd contracts-bedrock && scripts/getting-started/config.sh && {}".format(
            add_dependencies,
            replace_cmd,
        ),
        env_vars={
            "GS_ADMIN_ADDRESS": env["GS_ADMIN_ADDRESS"],
            "GS_BATCHER_ADDRESS": env["GS_BATCHER_ADDRESS"],
            "GS_PROPOSER_ADDRESS": env["GS_PROPOSER_ADDRESS"],
            "GS_SEQUENCER_ADDRESS": env["GS_SEQUENCER_ADDRESS"],
            "L1_RPC_URL": env["L1_RPC_URL"],
        },
        files = {"/contracts-bedrock": files.contracts},
        store=[
            StoreSpec(
                src="/contracts-bedrock",
                name=files.contracts,
            )
        ],
    )
    
def deploy_create2(plan, env):
    codesize_output = plan.run_sh(
        image="ghcr.io/foundry-rs/foundry:latest",
        run="cast codesize 0x4e59b44847b379578588920cA78FbF26c0B4956C --rpc-url {}".format(
            env["L1_RPC_URL"],
        ),
    ).output

    # if codesize_output == "0":
    #     plan.run_sh(
    #         image="ghcr.io/foundry-rs/foundry:latest",
    #         run="cast send --private-key {} 0x3fAB184622Dc19b6109349B94811493BF2a45362 --value 1ether --rpc-url {} --legacy && \
    #             cast publish --rpc-url {} 0xf8a58085174876e800830186a08080b853604580600e600039806000f350fe7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe03601600081602082378035828234f58015156039578182fd5b8082525050506014600cf31ba02222222222222222222222222222222222222222222222222222222222222222a02222222222222222222222222222222222222222222222222222222222222222 \
    #             ".format(
    #             env["PRIVATE_KEY"],
    #             env["L1_RPC_URL"],
    #             env["L1_RPC_URL"],
    #         ),
    #     )
    #     codesize_output = plan.run_sh(
    #         image="ghcr.io/foundry-rs/foundry:latest",
    #         run="cast codesize 0x4e59b44847b379578588920cA78FbF26c0B4956C --rpc-url {}".format(
    #             env["L1_RPC_URL"],
    #         ),
    #     ).output.strip()
    #     plan.verify(value=codesize_output, assertion="!=", target_value="0")
    
    plan.verify(value=codesize_output, assertion="==", target_value="69\n")


def build_getting_started_dir(plan, env, files, l1):
    chain_id_cmd = "echo -n {} > getting-started/.chainId".format(l1.chain_id)
    deploy_cmd = "echo -n '{}' > getting-started/.deploy"
    plan.run_sh(
        image="alpine:latest",
        run="cd contracts-bedrock/deployments && mkdir getting-started && {} && {}".format(
            chain_id_cmd,
            deploy_cmd,
        ),
        files = {"/contracts-bedrock": files.contracts},
        store = [
            StoreSpec(
                src="contracts-bedrock",
                name=files.contracts,
            )
        ]
    )

def generate_jwt_secret(plan, files):
    plan.run_sh(
        image="alpine/openssl:latest",
        run="openssl rand -hex 32 > /config/jwt.txt",
        files = {"/config": files.config},
        store=[
            StoreSpec(
                src="/config",
                name=files.config,
            )
        ]
    )
