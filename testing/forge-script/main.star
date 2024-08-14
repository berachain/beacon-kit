SOURCE_DIR_PATH = "/app/contracts"
IMAGE_FOUNDRY = "ghcr.io/foundry-rs/foundry:latest"
ENTRYPOINT = ["/bin/sh"]
DEPENDENCY_DIR_PATH = "/app/dependency"

def run(plan, deployment = {}):
    deploy_contracts(plan, deployment)

# Define the function to run the Forge script for deployment
def deploy_contracts(plan, deployment):
    repository = deployment["repository"]
    contracts_path = deployment["contracts_path"]
    script_path = deployment["script_path"]
    contract_name = deployment["contract_name"]
    rpc_url = deployment["rpc_url"]
    wallet = deployment["wallet"]
    dependency = deployment["dependency"]
    dependency_status = dependency["status"]

    # TODO: Support other wallet options such as mnemonics, keystore, hardware wallets.
    if wallet["type"] == "private_key":
        wallet_command = "--private-key {}".format(wallet["value"])
    else:
        fail("Wallet type {} not supported.".format(wallet["type"]))

    folder = plan.upload_files(src = repository, name = "contracts")

    foundry_service = plan.add_service(
        name = "foundry",
        config = ServiceConfig(
            image = IMAGE_FOUNDRY,
            entrypoint = ENTRYPOINT,
            files = {
                SOURCE_DIR_PATH: "contracts",
            },
        ),
    )

    if contracts_path:
        contract_path = "{}/{}".format(SOURCE_DIR_PATH, contracts_path)
    else:
        contract_path = SOURCE_DIR_PATH

    # Check if dependency status is set to true
    if dependency_status:
        plan.print("Dependency status is true", dependency_status)
        dependency_path = dependency["path"]
        plan.upload_files(src = dependency_path, name = "dependency")

        artifact_name = plan.store_service_files(
            service_name = foundry_service.name,
            # The path on the service's container that will be copied into a files artifact.
            src = "{}/{}".format(DEPENDENCY_DIR_PATH, dependency_path),
            name = "dependency",
        )

        # Run shell script
        plan.exec(
            service_name = foundry_service.name,
            recipe = ExecRecipe(
                command = ["/bin/sh", "-c", "sh {}".format(DEPENDENCY_DIR_PATH + "/" + dependency_path)],
            ),
        )

    result = plan.exec(
        service_name = foundry_service.name,
        recipe = ExecRecipe(
            command = ["/bin/sh", "-c", "cd {} && forge build".format(contract_path)],
        ),
    )
    plan.verify(result["code"], "==", 0)

    script_output = exec_on_service(
        plan,
        foundry_service.name,
        "cd {} && forge script {}:{} --broadcast --rpc-url {} {} --json  --skip test > output.json ".format(
            contract_path,
            script_path,
            contract_name,
            rpc_url,
            wallet_command,
        ),
    )

    exec_on_service(
        plan,
        foundry_service.name,
        "cat {}/output.json ".format(contract_path),
    )

    # Get the forge script output in a output.json file and grep from it
    transaction_file = "grep 'Transactions saved to' output.json | awk -F': ' '{print $2}'"
    plan.print("transaction_file", transaction_file)

    transaction_file_details = exec_on_service(plan, foundry_service.name, "cd {} && {}".format(contract_path, transaction_file))

    if not transaction_file_details["output"]:
        fail("Transaction file not found.")
    exec_output = exec_on_service(plan, foundry_service.name, "cat {}".format(transaction_file_details["output"]))
    plan.verify(exec_output["code"], "==", 0)

def exec_on_service(plan, service_name, command):
    return plan.exec(
        service_name = service_name,
        recipe = ExecRecipe(
            command = ["/bin/sh", "-c", command],
        ),
    )
