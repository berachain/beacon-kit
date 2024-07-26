SOURCE_DIR_PATH = "/app/contracts"
IMAGE_FOUNDRY = "ghcr.io/foundry-rs/foundry:latest"

def run(plan, deployment = {}):
    deployed_contracts = deploy_contracts(plan, deployment)
    plan.print("deployed_contracts", str(deployed_contracts))

# Define the function to run the Forge script for deployment
def deploy_contracts(plan, deployment):
    contract_name = deployment["contract_name"]
    script_path = deployment["script_path"]
    repository = deployment["repository"]
    rpc_url = deployment["rpc_url"]
    private_key = deployment["private_key"]

    folder = plan.upload_files(src = repository, name = "contracts")

    ENTRYPOINT = ["/bin/sh"]
    service = plan.add_service(
        name = "foundry",
        config = ServiceConfig(
            image = IMAGE_FOUNDRY,
            entrypoint = ENTRYPOINT,
            files = {
                SOURCE_DIR_PATH: "contracts",
            },
        ),
    )

    exec_output = plan.exec(
        service_name = "foundry",
        recipe = ExecRecipe(
            command = ["/bin/sh", "-c", "cd /app/contracts && forge build"],
        ),
    )

    exec_output = plan.exec(
        service_name = "foundry",
        recipe = ExecRecipe(
            command = ["/bin/sh", "-c", "cd /app/contracts && forge script {}:{} --broadcast --rpc-url {} --private-key {} --json  --skip test > output.json ".format(script_path, contract_name, rpc_url, private_key)],
        ),
    )

    plan.print("exec_output", str(exec_output)) 

    # plan.store_service_files(service_name = "foundry", src = "/app/contracts/broadcast/", name = "broadcast_artifacts")

    # One way is to read the output and grep from it
    # command = "grep -l {} {} ".format("Transactions saved to:",exec_output["output"])
    # plan.print("command", command)
    # exec_output = plan.exec(
    #     service_name = "foundry",
    #     recipe = ExecRecipe(
    #         command = ["/bin/sh", "-c", command],
    #     ),
    # )
    # plan.print("exec_output", str(exec_output))

    # Another way is get the forge script output in a json file and grep from it

    # transaction_file="grep {} output.json | awk -F": " '{print $2}'".format("Transactions saved to:")

    # command = "grep "Transactions saved to:" /app/contracts/output.json"
    # plan.print("command", command)
    # exec_output = plan.exec(
    #     service_name = "foundry",
    #     recipe = ExecRecipe(
    #         # command = ["/bin/sh", "-c", "cd /app/contracts && echo {}".format(transaction_file)],
    #         command    = ["/bin/sh", "-c", command]
    #     ),
    # )
    # plan.print("exec_output", str(exec_output))

    # exec_output = plan.exec(
    #     service_name = "foundry",
    #     recipe = ExecRecipe(
    #         command = ["/bin/sh", "-c", "cat {}".format(exec_output["output"])],
    #     ),
    # )
    # plan.print("exec_output", str(exec_output))

    # output = plan.exec(
    #     service_name = "foundry",
    #     recipe = ExecRecipe(
    #     command = ["/bin/sh","-c","cd /app/contracts && cat /app/contracts/broadcast/NFT.s.sol/80084/run-latest.json" ],
    #     ),
    # )
    # plan.print("output", str(output["output"]))
    return service
