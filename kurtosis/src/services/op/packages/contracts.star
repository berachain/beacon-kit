OP_STACK_GITHUB_URL = "https://github.com/ethereum-optimism/optimism.git"
BRANCH = "tutorials/chain"
CONTRACTS_PATH = "packages/contracts_bedrock"

def clone_dir(plan):
    output = plan.run_sh(
        image = "alpine/git:latest",
        run = 'mkdir contracts && cd contracts && git init && git remote add origin https://github.com/ethereum-optimism/optimism.git && git config core.sparseCheckout true && echo "packages/contracts-bedrock/*" > .git/info/sparse-checkout && git pull --depth=1 origin tutorials/chain && cd .. && pwd && cd contracts && ls'.format(
            OP_STACK_GITHUB_URL,
            CONTRACTS_PATH,
            BRANCH,
        ),
        # run = "mkdir contracts && ls && pwd",
        store = [StoreSpec(src="./git/contracts/packages/contracts-bedrock", name="contracts")],
    )
    return output.files_artifacts[0]