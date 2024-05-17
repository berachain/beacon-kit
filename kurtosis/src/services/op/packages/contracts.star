deps = import_module("deps.star")
optimism = import_module("optimism.star")

PATH = "/optimism/packages/contracts-bedrock"

# TODO: Run this concurrently with the rest of setup
def install(plan, files):
    plan.run_sh(
        image = "ghcr.io/foundry-rs/foundry:latest",
        run = "cd {} && forge install".format(PATH),
        files = {optimism.PATH: files.optimism},
        store = [StoreSpec(src = optimism.PATH, name = files.optimism)],
    )

def deploy_create2(plan, env):
    codesize_output = plan.run_sh(
        image = "ghcr.io/foundry-rs/foundry:latest",
        run = "cast codesize 0x4e59b44847b379578588920cA78FbF26c0B4956C --rpc-url {}".format(
            env.l1_rpc_url,
        ),
    ).output

    # TODO: Fix this logic: since output is a future value, this condition
    # will always be false even if the result is 0
    # if codesize_output == "0\n":
    #     plan.run_sh(
    #         image="ghcr.io/foundry-rs/foundry:latest",
    #         run="cast send --private-key {} 0x3fAB184622Dc19b6109349B94811493BF2a45362 --value 1ether --rpc-url {} --legacy && \
    #             cast publish --rpc-url {} 0xf8a58085174876e800830186a08080b853604580600e600039806000f350fe7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe03601600081602082378035828234f58015156039578182fd5b8082525050506014600cf31ba02222222222222222222222222222222222222222222222222222222222222222a02222222222222222222222222222222222222222222222222222222222222222 \
    #             ".format(
    #             env.pk,
    #             env.l1_rpc_url,
    #             env.l1_rpc_url,
    #         ),
    #     )
    #     codesize_output = plan.run_sh(
    #         image="ghcr.io/foundry-rs/foundry:latest",
    #         run="cast codesize 0x4e59b44847b379578588920cA78FbF26c0B4956C --rpc-url {}".format(
    #             env.l1_rpc_url,
    #         ),
    #     ).output.strip()
    #     plan.verify(value=codesize_output, assertion="!=", target_value="0\n")

    plan.verify(value = codesize_output, assertion = "==", target_value = "69\n")

# deploy_l1 deploys the L1 contracts and builds required files for the L2 deployment
def deploy_l1(plan, env, files):
    copy_cmd = "cp deployments/getting-started/.deploy deployments/getting-started/l1.json"
    plan.run_sh(
        image = "ghcr.io/foundry-rs/foundry:latest",
        run = "{} && cd {} && forge script scripts/Deploy.s.sol:Deploy --private-key {} --broadcast --rpc-url {} --legacy && {}".format(
            deps.get(["bash", "jq"]),
            PATH,
            env.admin_pk,
            env.l1_rpc_url,
            copy_cmd,
        ),
        files = {optimism.PATH: files.optimism},
        store = [StoreSpec(src = optimism.PATH, name = files.optimism)],
        env_vars = {
            "GS_ADMIN_ADDRESS": env.admin_address,
            "GS_ADMIN_PRIVATE_KEY": env.admin_pk,
            "GS_BATCHER_ADDRESS": env.batcher_address,
            "GS_BATCHER_PRIVATE_KEY": env.batcher_pk,
            "GS_PROPOSER_ADDRESS": env.proposer_address,
            "GS_PROPOSER_PRIVATE_KEY": env.proposer_pk,
            "GS_SEQUENCER_ADDRESS": env.sequencer_address,
            "GS_SEQUENCER_PRIVATE_KEY": env.sequencer_pk,
            "L1_RPC_KIND": env.l1_rpc_kind,
            "L1_RPC_URL": env.l1_rpc_url,
            "L1_BLOCK_TIME": env.l1_block_time,
            "IMPL_SALT": env.impl_salt,
            "DEPLOYMENT_CONTEXT": env.deployment_context,
            "PRIVATE_KEY": env.pk,
        },
    )
