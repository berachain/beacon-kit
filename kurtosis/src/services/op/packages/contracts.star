OP_STACK_GITHUB_URL = "https://github.com/ethereum-optimism/optimism.git"
BRANCH = "tutorials/chain"
CONTRACTS_PATH = "packages/contracts_bedrock"

def clone_dir():
    output = run_sh(
        image = "alpine/git:latest",
        run = """mkdir temp && cd temp && git init && git remote add origin {0} &&
                git sparse-checkout init && git sparse-checkout set {1} &&
                git fetch --depth=1 origin {2} && git checkout {3}""".format(
            OP_STACK_GITHUB_URL,
            CONTRACTS_PATH,
            BRANCH,
            BRANCH
        ),
        store = [src="temp", name="contracts"]
    )
    return output.files_artifacts[0]