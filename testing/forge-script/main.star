import shlex

SOURCE_DIR_PATH = "/app/contracts"
IMAGE_FOUNDRY = "ghcr.io/foundry-rs/foundry:latest"
ENTRYPOINT = ["/bin/sh"]
DEPENDENCY_DIR_PATH = "/app/dependency"
OUTPUT_FILE = "output.json"
CONTRACTS_DIR = "contracts"
DEPENDENCY_DIR = "dependency"
TRANSACTION_FILE_CMD = "grep 'Transactions saved to' output.json | awk -F': ' '{print $2}'"
testBaseFee = 123

def run(plan, deployment={}):
    deploy_contracts(plan, deployment)

# Function to run commands and check for errors
def run_command(plan, service_name, command, error_message):
    result = plan.exec(
        service_name=service_name,
        recipe=ExecRecipe(
            command=["/bin/sh", "-c", command],
        ),
    )
    if result["code"] != 0:
        fail(error_message)
    return result

# Get the wallet command
def get_wallet_command(wallet):
    wallet_type = wallet["type"]
    if wallet_type == "private_key":
        return "--private-key {}".format(shlex.quote(wallet["value"]))
    else:
        fail("Wallet type {} not supported.".format(wallet_type))

# Main function for deploying contracts
def deploy_contracts(plan, deployment):
    # Fetching data from deployment with checks
    repository = deployment.get("repository")
    if not repository:
        fail("Repository not specified in deployment")
    contracts_path = deployment.get("contracts_path", "")
    script_path = deployment.get("script_path", "")
    contract_name = deployment.get("contract_name", "")
    rpc_url = shlex.quote(deployment.get("rpc_url", ""))
    wallet = deployment.get("wallet")
    if not wallet:
        fail("Wallet not specified")
    dependency = deployment.get("dependency", {})
    dependency_type = dependency.get("type", "")

    wallet_command = get_wallet_command(wallet)

    folder = plan.upload_files(src=repository, name=CONTRACTS_DIR)

    dependency_artifact_name = ""
    if dependency_type == "local" or dependency_type == "git":
        dependency_path = dependency.get("path", "")
        plan.upload_files(src=DEPENDENCY_DIR, name=DEPENDENCY_DIR)
        dependency_artifact_name = DEPENDENCY_DIR

    foundry_service = plan.add_service(
        name="foundry",
        config=get_service_config(dependency_artifact_name),
    )

    contract_path = "{}/{}".format(SOURCE_DIR_PATH, contracts_path) if contracts_path else SOURCE_DIR_PATH

    if dependency_type == "local":
        run_command(plan, foundry_service.name, "sh {}/{}".format(DEPENDENCY_DIR_PATH, dependency_path), "Local dependency execution failed")
    elif dependency_type == "git":
        run_command(plan, foundry_service.name, "cd {} && sh {}".format(contract_path, dependency_path), "Git dependency execution failed")

    if script_path:
        # Logging the build and deployment process
        plan.print("Running deployment script: {}".format(script_path))
        run_command(plan, foundry_service.name, "cd {} && forge build".format(contract_path), "Forge build failed")

        # Executing forge script
        script_output = exec_on_service(
            plan,
            foundry_service.name,
            "cd {} && forge script {}:{} --broadcast --rpc-url {} {} --json --skip test > {}".format(
                contract_path,
                script_path,
                contract_name,
                rpc_url,
                wallet_command,
                OUTPUT_FILE
            ),
        )

    exec_on_service(
        plan,
        foundry_service.name,
        "cat {}/{}".format(contract_path, OUTPUT_FILE),
    )

    if script_path:
        # Fetch and check the transaction file
        plan.print("Retrieving transaction file...")
        transaction_file_details = exec_on_service(
            plan,
            foundry_service.name,
            "cd {} && {}".format(contract_path, TRANSACTION_FILE_CMD)
        )
        
        if not transaction_file_details["output"]:
            fail("Transaction file not found.")

        exec_output = exec_on_service(
            plan,
            foundry_service.name,
            "chmod -R 777 /app/contracts && cat {}".format(transaction_file_details["output"])
        )
        plan.verify(exec_output["code"], "==", 0)

# Function to execute commands inside the service
def exec_on_service(plan, service_name, command):
    return plan.exec(
        service_name=service_name,
        recipe=ExecRecipe(
            command=["/bin/sh", "-c", command],
        ),
    )

# Service configuration
def get_service_config(dependency_artifact_name=None):
    files = {
        SOURCE_DIR_PATH: CONTRACTS_DIR,
    }

    if dependency_artifact_name:
        files[DEPENDENCY_DIR_PATH] = dependency_artifact_name

    return ServiceConfig(
        image=IMAGE_FOUNDRY,
        entrypoint=ENTRYPOINT,
        files=files,
    )
