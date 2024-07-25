SOURCE_DIR_PATH = "/app/source"
IMAGE_FOUNDRY = "ghcr.io/foundry-rs/foundry:latest"

def run(plan, deployment = {}):
    deployed_contracts = deploy_contracts(plan, deployment)
    # for name, address in deployed_contracts.items():
    #     print("Contract " + name + " deployed at: " + address)
    # interact_with_contracts(deployed_contracts, deployment["interactions"])

# Define the function to run the Forge script for deployment
def deploy_contracts(plan, deployment):
    contract_name = deployment["contract_name"]
    script_path = deployment["script_path"]
    repository = deployment["repository"]
    rpc_url = deployment["rpc_url"]
    private_key = deployment["private_key"]

    folder = plan.upload_files(src = repository, name = "folder")

    run_command = "cd {} && forge build && forge script {}:{} --broadcast --rpc-url {} --private-key {} --json --skip test > output.json && sleep 20".format(SOURCE_DIR_PATH, script_path, contract_name, rpc_url, private_key)

    plan.print("run_command: " + str(run_command))

    # add a service
    service = plan.add_service(
        name = "foundry",
        config = ServiceConfig(
            image = IMAGE_FOUNDRY,
            cmd = [
                "-c",
                "cd {} && forge build && sleep infinity".format(SOURCE_DIR_PATH),
            ],
            entrypoint = ["/bin/sh", "-c"],
            files = {
                SOURCE_DIR_PATH: folder,
            },
        ),
    )

    result = plan.run_sh(
        run = run_command,
        image = IMAGE_FOUNDRY,
        files = {
            SOURCE_DIR_PATH: folder,
        },
        description = "Deploying smart contract",
        wait = "4m",
        # store = [
        # StoreSpec(src = "/app/source/", name = "output.json"),
        # ]
    )
    plan.print("result: " + str(result))

#     artifact_name = plan.store_service_files(
#     service_name = service.name,
#     src = "/app/source/output.json",
#     name = "output.json",
#     description = "storing some files"
# )

#     plan.print("artifact_name: " + str(artifact_name))
# How to access files from a service

# output_result = plan.run_sh(
#     run = "cd {} && cat output.json".format(SOURCE_DIR_PATH),
#     image = IMAGE_FOUNDRY,
#     files = {
#         SOURCE_DIR_PATH: folder,
#         "output.json": "output.json",

#     },
#     description = "Output of the deployment",
#     store = [
#     StoreSpec(src = "/app/source/", name = "output.json"),
#     ]
# )

# plan.print("output_result: " + str(output_result))
# plan.print("output: "+ str(output_result.output))

# shell_command = ["/bin/sh","-c",run_command]

# exec_output = plan.exec(
#     service_name = service.name,
#     recipe = ExecRecipe(
#         command = ["cat", "/app/source/output.json"],
#     ),
# )

# plan.print("exec_output: " + str(exec_output))

# exec_output = plan.exec(
#     service_name = service.name,
#     recipe = ExecRecipe(
#         command = shell_command,
#     ),
# )

# plan.print("exec_output: " + str(exec_output))

# recipe_result = plan.wait(service_name=service.name,
#     recipe= ExecRecipe(
#         command = shell_command,
#     ),
#     timeout = "2m",
#     field = "code",
#     assertion = "==",
#     target_value = 200,
# )

# plan.print(recipe_result["code"])

# plan.print("receipe_result" +str(recipe_result))

# result = run_command(deploy_command)

# # Extract contract addresses from the output
# deployed_contracts = extract_addresses(result.stdout)
# return deployed_contracts

# Define the contract interaction function
# def interact_with_contracts(deployed_contracts, interactions):
#     for contract_name, contract_address in deployed_contracts.items():
#         for interaction in interactions.get(contract_name, []):
#             function_name = interaction["function"]
#             args = interaction.get("args", [])
#             interaction_command = ["forge", "call", contract_address, function_name] + [str(arg) for arg in args]

#             result = run_command(interaction_command)
#             if result.return_code != 0:
#                 fail("Interaction with " + contract_name + " failed: " + result.stderr)

# Utility function to run a command
# def run_command(command):
#     return ExecRecipe({
#         "args": command,
#         "env_vars": {},
#         "workdir": "",
#         "timeout_seconds": 60,
#     })

# Utility function to extract contract addresses from the deployment output
# def extract_addresses(output):
#     # Implement logic to extract contract addresses from the output
#     # Example dummy implementation (update this to parse your actual output):
#     deployed_contracts = {
#         "ContractA": "0xContractAAddress",
#         "ContractB": "0xContractBAddress",
#     }
#     return deployed_contracts
