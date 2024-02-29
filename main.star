
eth_constants = import_module('github.com/kurtosis-tech/ethereum-package/src/package_io/constants.star')
reth = import_module('github.com/kurtosis-tech/ethereum-package/src/el/reth/reth_launcher.star')
input_parser = import_module("github.com/kurtosis-tech/ethereum-package/src/package_io/input_parser.star")
el_cl_genesis_data_generator = import_module(
    "github.com/kurtosis-tech/ethereum-package/src/prelaunch_data_generator/el_cl_genesis/el_cl_genesis_generator.star"
)
el_cl_genesis_data = import_module(
    "github.com/kurtosis-tech/ethereum-package/src/prelaunch_data_generator/el_cl_genesis/el_cl_genesis_data.star"
)
eth_static_files = import_module("github.com/kurtosis-tech/ethereum-package/src/static_files/static_files.star")
participant_network = import_module("github.com/kurtosis-tech/ethereum-package/src/participant_network.star")


beacond = import_module('./e2e/kurtosis/src/beacond/beacond_launcher.star')
static_files = import_module('./e2e/kurtosis/src/static_files/static_files.star')

KURTOSIS_ETH_PACKAGE_URL = "github.com/kurtosis-tech/ethereum-package"

EL_CLIENT_TYPE = struct(
    reth="reth"
)

CL_CLIENT_TYPE = struct(
    beacond="beacond"
)
def run(plan, args={}):
    jwt_file = plan.upload_files(
        src=KURTOSIS_ETH_PACKAGE_URL + eth_static_files.JWT_PATH_FILEPATH,
        name="jwt_file",
    )



    plan.print("Your args: {}".format(args))
    args_with_right_defaults = input_parser.input_parser(plan, args)
    num_participants = len(args_with_right_defaults.participants)
    network_params = args_with_right_defaults.network_params
    mev_params = args_with_right_defaults.mev_params
    parallel_keystore_generation = args_with_right_defaults.parallel_keystore_generation
    persistent = args_with_right_defaults.persistent
    xatu_sentry_params = args_with_right_defaults.xatu_sentry_params
    global_tolerations = args_with_right_defaults.global_tolerations
    global_node_selectors = args_with_right_defaults.global_node_selectors

    plan.print(network_params)



    genesis_file = plan.upload_files(
        static_files.GENESIS_FILEPATH,
        "el_cl_genesis_data",
    )

    el_cl_data = el_cl_genesis_data.new_el_cl_genesis_data(
        genesis_file,
        "",
        0,
        0,
    )

    plan.print(el_cl_data)

    el_launchers = get_el_launcher(el_cl_data, jwt_file, network_params)
    el_launcher = el_launchers[EL_CLIENT_TYPE.reth]["launcher"]
    launch_method = el_launchers[EL_CLIENT_TYPE.reth]["launch_method"]

    el_service_name = "el-0-reth-beaconkit"

    el_client_context = launch_method(
        plan,
        el_launcher,
        el_service_name,
        input_parser.DEFAULT_EL_IMAGES["reth"],
        "info",
        "",
        [],
        0,
        0,
        0,
        0,
        [],
        {},
        {},
        False,
        0,
        [],
        {},
    )
    

    # plan.print()
    plan.print(el_client_context)

    beacond_config = beacond.get_config(jwt_file, "http://el-0-reth-beaconkit:8551")
    plan.add_service(
        name = "beaconkit-node",
        config = beacond_config,
    )

    return el_client_context

    


def get_el_launcher(el_cl_data, jwt_file, network_params):    
    el_launchers = {
        EL_CLIENT_TYPE.reth: {
            "launcher": reth.new_reth_launcher(
                el_cl_data,
                jwt_file,
                network_params.network,
            ),
            "launch_method": reth.launch,
        },
    }

    # cl_launchers = {
    #     CL_CLIENT_TYPE.beacond: {
    #         "launcher": beacond.new_beacond_launcher(
    #         ),
    #         "launch_method": reth.launch,
    #     },
    # }
    
    
    return el_launchers
    


