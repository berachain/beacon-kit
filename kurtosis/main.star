input_parser = import_module("github.com/kurtosis-tech/ethereum-package/src/package_io/input_parser.star")
el_cl_genesis_data_generator = import_module(
    "github.com/kurtosis-tech/ethereum-package/src/prelaunch_data_generator/el_cl_genesis/el_cl_genesis_generator.star",
)

participant_network = import_module("github.com/kurtosis-tech/ethereum-package/src/participant_network.star")

execution = import_module("./src/nodes/execution/execution.star")
execution_types = import_module("./src/nodes/execution/types.star")
beacond = import_module("./src/nodes/consensus/beacond/launcher.star")
genesis = import_module("./src/networks/networks.star")
port_spec_lib = import_module("./src/lib/port_spec.star")

def run(plan, args = {}):
    """
    Initiates the execution plan with the specified number of participants and arguments.

    Args:
      plan: The execution plan to be run.
      args: Additional arguments to configure the plan. Defaults to an empty dictionary.
    """

    plan.print("Your args: {}".format(args))
    args_with_right_defaults = input_parser.input_parser(plan, args)
    num_participants = len(args_with_right_defaults.participants)

    network_params = args_with_right_defaults.network_params

    # 1. Initialize EVM genesis data
    evm_genesis_data = genesis.get_genesis_data(plan)

    node_modules = {}
    for node in args_with_right_defaults.participants:
        if node.el_type not in node_modules.keys():
            node_path = "./src/nodes/execution/{}/config.star".format(node.el_type)
            node_module = import_module(node_path)
            node_modules[node.el_type] = node_module

    # 2. Upload jwt
    jwt_file = execution.upload_global_files(plan, node_modules)

    node_peering_info = []

    # 3. Perform genesis ceremony
    for n in range(num_participants):
        cl_service_name = "cl-{}-{}-beaconkit".format(n, args_with_right_defaults.participants[n].el_type)
        engine_dial_url = ""  # not needed for this step
        beacond_config = beacond.get_config(
            args_with_right_defaults.participants[n].cl_image,
            jwt_file,
            engine_dial_url,
            cl_service_name,
            expose_ports = False,
        )

        if n > 0:
            beacond_config.files["/root/.beacond/config"] = Directory(
                artifact_names = ["cosmos-genesis-{}".format(n - 1)],
            )

        if n == num_participants - 1 and n != 0:
            collected_gentx = []
            for other_participant_id in range(num_participants - 1):
                collected_gentx.append("cosmos-gentx-{}".format(other_participant_id))

            beacond_config.files["/root/.beacond/config/gentx"] = Directory(
                artifact_names = collected_gentx,
            )

        cl_service = plan.add_service(
            name = cl_service_name,
            config = beacond_config,
        )

        exec_recipe = None
        if n == 0:
            exec_recipe = ExecRecipe(
                # Initialize the Cosmos genesis file
                command = ["/usr/bin/init_first.sh"],
            )
        else:
            exec_recipe = ExecRecipe(
                # Initialize the Cosmos genesis file
                command = ["/usr/bin/init_others.sh"],
            )

        result = plan.exec(
            service_name = cl_service_name,
            recipe = exec_recipe,
        )

        peer_result = plan.exec(
            service_name = cl_service_name,
            recipe = ExecRecipe(
                command = ["bash", "-c", "/usr/bin/beacond comet show-node-id --home $BEACOND_HOME | tr -d '\n'"],
            ),
        )

        node_peering_info.append(peer_result["output"])

        file_suffix = "{}".format(n)
        if n == num_participants - 1:
            finalize_recipe = ExecRecipe(
                # Initialize the Cosmos genesis file
                command = ["/usr/bin/finalize.sh"],
            )
            result = plan.exec(
                service_name = cl_service_name,
                recipe = finalize_recipe,
            )
            file_suffix = "final"

        node_beacond_config = plan.store_service_files(
            service_name = cl_service_name,
            src = "/root/.beacond",
            name = "node-beacond-config-{}".format(n),
        )

        genesis_artifact = plan.store_service_files(
            # The service name of a preexisting service from which the file will be copied.
            service_name = cl_service_name,

            # The path on the service's container that will be copied into a files artifact.
            # MANDATORY
            src = "/root/.beacond/config/genesis.json",

            # The name to give the files artifact that will be produced.
            # If not specified, it will be auto-generated.
            # OPTIONAL
            name = "cosmos-genesis-{}".format(file_suffix),
        )

        gentx_artifact = plan.store_service_files(
            service_name = cl_service_name,
            src = "/root/.beacond/config/gentx/*",
            name = "cosmos-gentx-{}".format(n),
        )

        # Node has completed its genesis step. We will add it back later once genesis is complete
        plan.remove_service(
            cl_service_name,
        )

    el_enode_addrs = []

    # 4. Start network participants
    for n in range(num_participants):
        el_type = args_with_right_defaults.participants[n].el_type
        node_module = node_modules[el_type]
        el_service_name = "el-{}-{}-beaconkit".format(n, el_type)

        # 4a. Launch EL
        el_service_config_dict = execution.get_default_service_config(el_service_name, node_module)
        el_service_config_dict = execution.add_bootnodes(node_module, el_service_config_dict, el_enode_addrs)
        el_client_service = execution.deploy_node(plan, el_service_config_dict)

        enode_addr = execution.get_enode_addr(plan, el_client_service, el_service_name, el_type)
        el_enode_addrs.append(enode_addr)

        # 4b. Launch CL
        cl_service_name = "cl-{}-{}-beaconkit".format(n, el_type)
        engine_dial_url = "http://{}:{}".format(el_service_name, execution.ENGINE_RPC_PORT_NUM)

        # Get peers for the cl node
        my_peers = node_peering_info[:n]
        for i in range(len(my_peers)):
            peer_el_service_name = "cl-{}-{}-beaconkit".format(i, args_with_right_defaults.participants[i].el_type)
            peer_service = plan.get_service(peer_el_service_name)
            my_peers[i] = my_peers[i] + "@" + peer_service.ip_address + ":26656"
        persistent_peers = ",".join(my_peers)

        beacond_config = beacond.get_config(
            args_with_right_defaults.participants[n].cl_image,
            jwt_file,
            engine_dial_url,
            cl_service_name,
            entrypoint = ["bash"],
            cmd = ["-c", "/usr/bin/start.sh"],
            persistent_peers = persistent_peers,
        )

        # Add back in the node's config data and overwrite genesis.json with final genesis file
        beacond_config.files["/root"] = Directory(
            artifact_names = ["node-beacond-config-{}".format(n)],
        )
        beacond_config.files["/root/.tmp_genesis"] = Directory(artifact_names = ["cosmos-genesis-final"])

        plan.print(beacond_config)

        plan.add_service(
            name = cl_service_name,
            config = beacond_config,
        )
