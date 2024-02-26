

reth = import_module('github.com/kurtosis-tech/ethereum-package/src/el/reth/reth_launcher.star')
input_parser = import_module("github.com/kurtosis-tech/ethereum-package/src/package_io/input_parser.star")
beacond = import_module('./src/beacond/beacond_launcher.star')
EL_CLIENT_TYPE = struct(
    reth="reth"
)

CL_CLIENT_TYPE = struct(
    beacond="beacond"
)
def run(plan, args={}):
    args_with_right_defaults = input_parser.input_parser(plan, args)
    plan.print("printing this")


def get_el_launcher(jwt_file, network_params, network_id):
    el_cl_data = network_params.el_client_data
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

    cl_launchers = {
        CL_CLIENT_TYPE.beacond: {
            "launcher": beacond.new_beacond_launcher(
            ),
            "launch_method": reth.launch,
        },
    }
    
    
    return el_launchers[el_cl_data.client_type]
    