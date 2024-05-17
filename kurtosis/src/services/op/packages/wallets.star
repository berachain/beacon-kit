images = import_module("../constants/images.star")

bash = import_module("bash.star")
deps = import_module("deps.star")
contracts = import_module("contracts.star")
optimism = import_module("optimism.star")

ARTIFACT_NAME = "wallets"
PATH = "/wallets/wallets.txt"

# get returns wallets (<address, pk> pair) in the order of:
#   GS_ADMIN, GS_BATCHER, GS_PROPOSER, GS_SEQUENCER
def get(plan, op):
    wallets = plan.run_sh(
        description = "Generating op-component wallets",
        image = images.FOUNDRY,
        run = bash.run([
            deps.get(["bash"]),
            "cd {}".format(contracts.PATH),
            "scripts/getting-started/wallets.sh | grep '_ADDRESS\\|_PRIVATE_KEY' | cut -d '=' -f 2 > {}".format(
                "wallets.txt",
            ),
        ]),
        files = {optimism.PATH: op},
        store = [StoreSpec(src = "{}/wallets.txt".format(contracts.PATH), name = ARTIFACT_NAME)],
    )
    return wallets.files_artifacts[0]

# get_by_index returns the line at the given (1-indexed) index
# a wallet <address, pk> pair is: <wallets.txt[index], wallets.txt[index+1]>
# requires: wallets/wallets.txt to be a valid file artifact
def get_by_index(plan, index):
    wallet = plan.run_sh(
        image = images.ALPINE,
        run = "sed -n '{}p' {} | tr -d '\n'".format(index, PATH),
        files = {"/{}".format(ARTIFACT_NAME): ARTIFACT_NAME},
    )
    return wallet.output

# fund funds each component wallet with 10 ether
def fund(plan, env):
    wallet_cmd = "cast send --private-key {} --value 10ether --rpc-url {} --legacy {} "
    plan.run_sh(
        description = "Funding op-component wallets with 10 ether",
        image = images.FOUNDRY,
        run = bash.run([
            wallet_cmd.format(env.pk, env.l1_rpc_url, env.admin_address),
            wallet_cmd.format(env.pk, env.l1_rpc_url, env.batcher_address),
            wallet_cmd.format(env.pk, env.l1_rpc_url, env.proposer_address),
        ]),
    )
