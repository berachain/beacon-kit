load("@yaml//:parse.bzl", "parse")

# Load the YAML configuration
def load_config(file_path):
    config_content = read_file(file_path)
    return parse(config_content)

# Define the function to run the Forge script for deployment
def deploy_contracts(config):
    script_path = config["deployment"]["script_path"]
    deploy_command = ["forge", "script", script_path, "--broadcast", "--gas-limit", str(config["deployment"]["params"]["gas"])]

    result = run_command(deploy_command)
    if result.return_code != 0:
        fail("Deployment script failed: " + result.stderr)

    # Extract contract addresses from the output
    deployed_contracts = extract_addresses(result.stdout)
    return deployed_contracts

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

# Main function to load config, deploy contracts, and interact with them
def main():
    config = load_config("config.yaml")
    deployed_contracts = deploy_contracts(config)
    for name, address in deployed_contracts.items():
        print("Contract " + name + " deployed at: " + address)
    interact_with_contracts(deployed_contracts, config["interactions"])

# Utility function to read file content
def read_file(file_path):
    with open(file_path, 'r') as file:
        return file.read()

# Utility function to run a command
def run_command(command):
    return execute({
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

# Entry point
main()
