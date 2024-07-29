SOURCE_DIR_PATH = "/app/contracts"
IMAGE_FOUNDRY = "ghcr.io/foundry-rs/foundry:latest"
ENTRYPOINT = ["/bin/sh"]

def run(plan, deployment = {}):
    deploy_contracts(plan, deployment)

# Define the function to run the Forge script for deployment
def deploy_contracts(plan, deployment):
    contract_name = deployment["contract_name"]
    script_path = deployment["script_path"]
    repository = deployment["repository"]
    rpc_url = deployment["rpc_url"]
    private_key = deployment["private_key"]

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

    result = plan.exec(
        service_name = foundry_service.name,
        recipe = ExecRecipe(
            command = ["/bin/sh", "-c", "cd /app/contracts && forge build"],
        ),
    )
    plan.verify(result["code"], "==", 0)

    script_output = exec_on_service(
        plan,
        foundry_service.name,
        "cd /app/contracts && forge script {}:{} --broadcast --rpc-url {} --private-key {} --json  --skip test > output.json ".format(
            script_path,
            contract_name,
            rpc_url,
            private_key,
        ),
    )

    # Get the forge script output in a output.json file and grep from it
    transaction_file = "grep {} output.json | awk -F{} {}".format("'Transactions saved to'", "': '", "'{print $2}'")
    plan.print("transaction_file", transaction_file)

    transaction_file_details = exec_on_service(plan, foundry_service.name, "cd /app/contracts && {}".format(transaction_file))

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
