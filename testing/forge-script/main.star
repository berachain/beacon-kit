# Main function to load config, deploy contracts, and interact with them
CONFIG_DIR_PATH = "/config"

def run(plan,deployment = {}):
    deployed_contracts = deploy_contracts(plan,deployment)
    # for name, address in deployed_contracts.items():
    #     print("Contract " + name + " deployed at: " + address)
    # interact_with_contracts(deployed_contracts, deployment["interactions"])

# Define the function to run the Forge script for deployment
def deploy_contracts(plan,deployment):

    script_path = deployment["script_path"]
    plan.print("script_path: " + script_path)
    # Pull the file from github using curl
    # Fetch the content of the file from script_path
    script_content = plan.upload_files(src = script_path, name = "script",description = "Uploading deployment script")

    plan.print("script_content: " + script_content)
    gas = str(deployment["params"]["gas"])
    plan.print("gas: " + gas)


    # add a service
    service = plan.add_service(
        name = "foundry",
        config = ServiceConfig(
        image = "ghcr.io/foundry-rs/foundry:latest",
        cmd = [
        "-c",
        "sleep 99",
        ],
        files = {
            CONFIG_DIR_PATH: script_content,
        },
        )
        )

    plan.print("service: " + str(service))  
    # deploy_command = ["forge", "script", str(script_content), "--broadcast", "--gas-limit", gas]
    # deploy_command = "tail -f /dev/null"
    plan.print("Service name", service.name)
    # plan.exec(
    #     service_name = service.name,
    #     recipe = ExecRecipe(
    #         command = [deploy_command],
    #     ),
    # )

    # plan.print("deploy_command: " + str(deploy_command))
    # result = run_command(deploy_command)
    # if result.return_code != 0:
    #     fail("Deployment script failed: " + result.stderr)

    # # Extract contract addresses from the output
    # deployed_contracts = extract_addresses(result.stdout)
    # return deployed_contracts

# Define the contract interaction function
def interact_with_contracts(deployed_contracts, interactions):
    for contract_name, contract_address in deployed_contracts.items():
        for interaction in interactions.get(contract_name, []):
            function_name = interaction["function"]
            args = interaction.get("args", [])
            interaction_command = ["forge", "call", contract_address, function_name] + [str(arg) for arg in args]

            result = run_command(interaction_command)
            if result.return_code != 0:
                fail("Interaction with " + contract_name + " failed: " + result.stderr)


# Utility function to run a command
def run_command(command):
    return ExecRecipe({
        "args": command,
        "env_vars": {},
        "workdir": "",
        "timeout_seconds": 60
    })

# Utility function to extract contract addresses from the deployment output
def extract_addresses(output):
    # Implement logic to extract contract addresses from the output
    # Example dummy implementation (update this to parse your actual output):
    deployed_contracts = {
        "ContractA": "0xContractAAddress",
        "ContractB": "0xContractBAddress"
    }
    return deployed_contracts