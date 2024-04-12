def init_beacond(plan, chain_id, moniker, home, is_first_validator, cl_service_name):
    genesis_file = "{}/config/genesis.json".format(home)

    # Check if genesis file exists, if not then initialize the beacond
    command_check = "if [ ! -f {} ]; then /usr/bin/beacond init --chain-id {} {} --home {} --beacon-kit.accept-tos; fi".format(genesis_file, chain_id, moniker, home)
    plan.exec(
        service_name = cl_service_name,
        recipe = ExecRecipe(
            command = ["bash", "-c", command_check],
        ),
    )
    if is_first_validator == True:
        create_beacond_config_directory(plan, home, cl_service_name)
    add_validator(plan, home, cl_service_name)
    collect_validator(plan, home, cl_service_name)

def create_beacond_config_directory(plan, home, cl_service_name):
    GENESIS = "{}/config/genesis.json".format(home)
    TMP_GENESIS = "{}/config/tmp_genesis.json".format(home)
    command = "jq '.consensus.params.validator.pub_key_types += [\"bls12_381\"] | .consensus.params.validator.pub_key_types -= [\"ed25519\"]' {} > {}".format(GENESIS, TMP_GENESIS)
    plan.exec(
        service_name = cl_service_name,
        recipe = ExecRecipe(
            command = ["bash", "-c", command],
        ),
    )
    command_mv = "mv {} {}".format(TMP_GENESIS, GENESIS)
    plan.exec(
        service_name = cl_service_name,
        recipe = ExecRecipe(
            command = ["bash", "-c", command_mv],
        ),
    )

def add_validator(plan, home, cl_service_name):
    command = "/usr/bin/beacond genesis add-validator --home {} --beacon-kit.accept-tos".format(home)
    plan.exec(
        service_name = cl_service_name,
        recipe = ExecRecipe(
            command = ["bash", "-c", command],
        ),
    )

def collect_validator(plan, home, cl_service_name):
    command = "/usr/bin/beacond genesis collect-validators --home {}".format(home)
    plan.exec(
        service_name = cl_service_name,
        recipe = ExecRecipe(
            command = ["bash", "-c", command],
        ),
    )
