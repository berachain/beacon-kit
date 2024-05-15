OP_STACK_GITHUB_URL = "https://github.com/ethereum-optimism/optimism.git"
BRANCH = "tutorials/chain"
CONTRACTS_PATH = "packages/contracts_bedrock"

def clone_dir(plan):
    output = plan.run_sh(
        image = "alpine/git:latest",
        run = 'mkdir contracts && cd contracts && git init && git remote add origin https://github.com/ethereum-optimism/optimism.git && git config core.sparseCheckout true && echo "packages/contracts-bedrock/*" > .git/info/sparse-checkout && git pull --depth=1 origin tutorials/chain'.format(
            OP_STACK_GITHUB_URL,
            CONTRACTS_PATH,
            BRANCH,
        ),
        store = [StoreSpec(src="./git/contracts", name="contracts")],
    )
    return output.files_artifacts[0]

def install(plan, files):
    plan.run_sh(
        image="ghcr.io/foundry-rs/foundry:latest",
        run="cd /contracts/packages/contracts-bedrock && forge install",
        wait="3600s", # this takes hella long
        files={"/contracts": files.contracts},
        store=[
            StoreSpec(
                src="/contracts/packages/contracts-bedrock",
                name=files.contracts,
            )
        ],
    )

def deploy_l1(plan, env, files):
    get_deps_cmd = "apk add bash jq "
    copy_cmd = "cp deployments/getting-started/.deploy l1.json"
    plan.run_sh(
        image="ghcr.io/foundry-rs/foundry:latest",
        run="cd contracts-bedrock && forge script scripts/Deploy.s.sol:Deploy --private-key {} --broadcast --rpc-url {} --legacy && {}".format(
            env["GS_ADMIN_PRIVATE_KEY"],
            env["L1_RPC_URL"],
            copy_cmd,
        ),
        wait="3600s", # this takes hella long
        files={
            "/contracts-bedrock": files.contracts,    
        },
        store=[
            StoreSpec(
                src="contracts-bedrock/l1.json",
                name=files.config,
            )
        ],
        env_vars={
            "GS_ADMIN_ADDRESS": env["GS_ADMIN_ADDRESS"],
            "GS_ADMIN_PRIVATE_KEY": env["GS_ADMIN_PRIVATE_KEY"],
            "GS_BATCHER_ADDRESS": env["GS_BATCHER_ADDRESS"],
            "GS_BATCHER_PRIVATE_KEY": env["GS_BATCHER_PRIVATE_KEY"],
            "GS_PROPOSER_ADDRESS": env["GS_PROPOSER_ADDRESS"],
            "GS_PROPOSER_PRIVATE_KEY": env["GS_PROPOSER_PRIVATE_KEY"],
            "GS_SEQUENCER_ADDRESS": env["GS_SEQUENCER_ADDRESS"],
            "L1_RPC_KIND": env["L1_RPC_KIND"],
            "L1_RPC_URL": env["L1_RPC_URL"],
            "L1_BLOCK_TIME": env["L1_BLOCK_TIME"],
            "IMPL_SALT": env["IMPL_SALT"],
            "DEPLOYMENT_CONTEXT": env["DEPLOYMENT_CONTEXT"],
            "PRIVATE_KEY": env["PRIVATE_KEY"],
        },
    )