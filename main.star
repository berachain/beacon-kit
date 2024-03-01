eth_constants = import_module('github.com/kurtosis-tech/ethereum-package/src/package_io/constants.star')
reth = import_module('github.com/kurtosis-tech/ethereum-package/src/el/reth/reth_launcher.star')
input_parser = import_module("github.com/kurtosis-tech/ethereum-package/src/package_io/input_parser.star")
el_cl_genesis_data_generator = import_module(
    "github.com/kurtosis-tech/ethereum-package/src/prelaunch_data_generator/el_cl_genesis/el_cl_genesis_generator.star"
)

eth_static_files = import_module("github.com/kurtosis-tech/ethereum-package/src/static_files/static_files.star")
participant_network = import_module("github.com/kurtosis-tech/ethereum-package/src/participant_network.star")

el = import_module('./e2e/kurtosis/src/el/el.star')
beacond = import_module('./e2e/kurtosis/src/beacond/beacond_launcher.star')
static_files = import_module('./e2e/kurtosis/src/static_files/static_files.star')
constants = import_module('./e2e/kurtosis/src/constants.star')
genesis = import_module('./e2e/kurtosis/src/genesis/genesis.star')



def run(plan, args={}):
    plan.print("Your args: {}".format(args))
    args_with_right_defaults = input_parser.input_parser(plan, args)
    num_participants = len(args_with_right_defaults.participants)
    network_params = args_with_right_defaults.network_params

    # 1. Initialize genesis data
    el_cl_data = genesis.get_genesis_data(plan)
    plan.print(el_cl_data)

    # 2. Upload jwt
    jwt_file = plan.upload_files(
        src=constants.KURTOSIS_ETH_PACKAGE_URL + eth_static_files.JWT_PATH_FILEPATH,
        name="jwt_file",
    )

    # 3. Launch EL
    el_service_name = "el-0-reth-beaconkit"
    el_client_context = el.get_el(plan, constants.EL_CLIENT_TYPE.reth, el_cl_data, jwt_file, el_service_name, network_params)


    # 4. Launch CL
    beacond_config = beacond.get_config(jwt_file, "http://el-0-reth-beaconkit:8551")
    plan.add_service(
        name = "beaconkit-node",
        config = beacond_config,
    )

    return el_client_context