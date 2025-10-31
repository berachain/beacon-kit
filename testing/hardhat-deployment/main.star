SOURCE_DIR_PATH = "/app/contracts"
IMAGE_NODE = "node:current-alpine3.20"
DEPENDENCY_DIR_PATH = "/app/dependency"
LOCAL_DEPENDENCY = "local"
GIT_DEPENDENCY = "git"

def run(plan, deployment = {}):
    repository = deployment["repository"]
    contracts_path = deployment["contracts_path"]
    script_path = deployment["script_path"]
    wallet = deployment["wallet"]
    network = deployment["network"]
    dependency = deployment["dependency"]
    dependency_type = dependency["type"]

    if wallet["type"] != "private_key":
        fail("Wallet type {} not supported.".format(wallet["type"]))

    wallet_value = wallet["value"]

    plan.upload_files(src = repository, name = "contracts")

    dependency_artifact_name = ""
    if dependency_type == LOCAL_DEPENDENCY or dependency_type == GIT_DEPENDENCY:
        dependency_path = dependency["path"]
        plan.upload_files(src = "dependency", name = "dependency")
        dependency_artifact_name = "dependency"

    node_service = plan.add_service(
        name = "hardhat",
        config = get_service_config(wallet_value, dependency_artifact_name),
    )

    if contracts_path:
        contract_path = "{}/{}".format(SOURCE_DIR_PATH, contracts_path)
    else:
        contract_path = SOURCE_DIR_PATH

    if dependency_type == LOCAL_DEPENDENCY:
        exec_on_service(plan, node_service.name, "sh {}/{}".format(DEPENDENCY_DIR_PATH, dependency_path))
    elif dependency_type == GIT_DEPENDENCY:
        exec_on_service(plan, node_service.name, "cd {} && sh {}".format(contract_path, dependency_path))

    # Compile the contracts
    result = exec_on_service(plan, node_service.name, "cd {} && npx hardhat compile --network {}".format(contract_path, network))
    plan.verify(result["code"], "==", 0)

    # Deploy the contracts
    result = exec_on_service(plan, node_service.name, "cd {} && npx hardhat run {}".format(contract_path, script_path))
    plan.verify(result["code"], "==", 0)

def get_service_config(wallet, dependency_artifact_name = None):
    files = {
        SOURCE_DIR_PATH: "contracts",
    }

    if dependency_artifact_name:
        files[DEPENDENCY_DIR_PATH] = dependency_artifact_name

    return ServiceConfig(
        image = IMAGE_NODE,
        files = files,
        env_vars = {
            "WALLET_KEY": wallet,
        },
    )

def exec_on_service(plan, service_name, command):
    return plan.exec(
        service_name = service_name,
        recipe = ExecRecipe(
            command = ["/bin/sh", "-c", command],
        ),
    )
