CONFIG_DIR_PATH = "/app/scripts"
IMAGE_FOUNDRY = "ghcr.io/foundry-rs/foundry:latest"
def run(plan, deployment = {}):
    deployed_contracts = deploy_contracts(plan, deployment)
    # for name, address in deployed_contracts.items():
    #     print("Contract " + name + " deployed at: " + address)
    # interact_with_contracts(deployed_contracts, deployment["interactions"])

# Define the function to run the Forge script for deployment
def deploy_contracts(plan, deployment):
    script_path = deployment["script_path"]
    plan.print("script_path: " + script_path)

    contract_name = deployment["contract_name"]
    plan.print("contract_name: " + contract_name)

    rpc_url = deployment["rpc_url"]
    plan.print("rpc_url: " + rpc_url)
    # TODO: Pull the file from github using curl
    # Fetch the content of the file from script_path
    script_content = plan.upload_files(src = script_path, name = "deployscript", description = "Uploading deployment script")

    folder = plan.upload_files(src = "github.com/nidhi-singh02/solidity-scripting/", name = "folder")


    plan.print("script_content: " + script_content)
    gas = str(deployment["params"]["gas"])
    plan.print("gas: " + gas)

    # add a service
    service = plan.add_service(
        name = "foundry",
        config = ServiceConfig(
            image = IMAGE_FOUNDRY,
            cmd = [
                "-c",
                "sleep infinity",
            ],
            files = {
                CONFIG_DIR_PATH: script_content,
                "/app/folder": folder,
            },
        ),
    )

    # plan.print("service: " + str(service))
    # plan.print("Service name", service.name)
    deploy_command = "forge script "+ str(CONFIG_DIR_PATH) + " --broadcast --gas-limit "+ gas
    plan.print("deploy_command: " + str(deploy_command))

    # deploy_command = "tail -f /dev/null"


    plan.print("Script content name: " + str(script_content))
    # run_command = "forge script "+str(CONFIG_DIR_PATH)+"/NFT.s.sol"+":"+contract_name+" --broadcast --gas-limit "+ gas
    # plan.print("run_command: " + str(run_command))

# "forge script "+str(CONFIG_DIR_PATH)+"/"+script_content+":"+contract_name+" --broadcast --gas-limit "+ gas
    run_command = "cd /app/folder && forge build && forge script /app/folder/script/NFT.s.sol"+":"+contract_name+" --broadcast --gas-limit "+ gas+ " --rpc-url " + rpc_url + " -vvvv"
    plan.print("run_command: " + str(run_command))

    plan.run_sh(
        run = run_command,
        image = IMAGE_FOUNDRY,
        files = {
            CONFIG_DIR_PATH: script_content,
            "/app/folder": folder,
        },
        description = "Deploying smart contract",
    )

    # plan.exec(
    #     service_name = service.name,
    #     recipe = ExecRecipe(
    #         command = [run_command],
    #     ),
    # )

    
    # result = run_command(deploy_command)
    # if result.return_code != 0:
    #     fail("Deployment script failed: " + result.stderr)

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
