def init_beacond(plan,chain_id, moniker, home, is_first_validator):
    command = "/usr/bin/beacond init --chain-id {} {} --home {} --beacon-kit.accept-tos".format(chain_id, moniker, home)
    plan.run_sh(command)
    if is_first_validator:
        create_beacond_config_directory(plan,home)
    add_validator(plan,home)
    collect_validator(plan,home)

def create_beacond_config_directory(plan,home):
    GENESIS = "{}/config/genesis.json".format(home)
    TMP_GENESIS = "{}/config/tmp_genesis.json".format(home)
    plan.run_sh("jq", ".consensus.params.validator.pub_key_types += ['bls12_381'] | .consensus.params.validator.pub_key_types -= ['ed25519']' {} > {}".format(GENESIS, TMP_GENESIS))
    command = "mv {} {}".format(TMP_GENESIS, GENESIS)
    plan.run_sh(command)

def add_validator(plan, home):
    command = "/usr/bin/beacond genesis add-validator --home {}".format(home)
    plan.run_sh(command)


def collect_validator(plan, home):
    command = "/usr/bin/beacond genesis collect-validators --home {}".format(home)
    plan.run_sh(command)
