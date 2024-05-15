ARTIFACT_NAME="wallets"
FILEPATH="wallets.txt"

# Returns wallets (address, pk pair) in the order of:
#   GS_ADMIN, GS_BATCHER, GS_PROPOSER, GS_SEQUENCER
def get(plan, contracts):
    wallets=plan.run_sh(
        image="ghcr.io/foundry-rs/foundry:latest",
        run='apk add bash > /dev/null 2>&1 && ./contracts/packages/contracts-bedrock/scripts/getting-started/wallets.sh | grep "_ADDRESS\\|_PRIVATE_KEY" | cut -d "=" -f 2 > {}'.format(FILEPATH),
        files={"/contracts": contracts},
        store=[StoreSpec(src=FILEPATH, name=ARTIFACT_NAME)],
    )
    return wallets.files_artifacts[0]

def get_by_index(plan, index):
    wallet = plan.run_sh(
        image="alpine:latest",
        run="sed -n '{}p' /{}/{} | tr -d '\n'".format(index, ARTIFACT_NAME, FILEPATH),
        files={"/{}".format(ARTIFACT_NAME) : ARTIFACT_NAME},
    )
    return wallet.output

# fund funds component wallets with 10 ether
def fund(plan, env):
    wallet_cmd = "cast send --private-key {} --value 10ether --rpc-url {} --legacy {} "
    plan.run_sh(
        image="ghcr.io/foundry-rs/foundry:latest",
        run="{} && {} && {}".format(
            wallet_cmd.format(env["PRIVATE_KEY"], env["L1_RPC_URL"], env["GS_ADMIN_ADDRESS"]),
            wallet_cmd.format(env["PRIVATE_KEY"], env["L1_RPC_URL"], env["GS_BATCHER_ADDRESS"]),
            wallet_cmd.format(env["PRIVATE_KEY"], env["L1_RPC_URL"], env["GS_PROPOSER_ADDRESS"]),
        ),
    )