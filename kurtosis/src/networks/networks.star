el_cl_genesis_data = import_module(
    "github.com/ethpandaops/ethereum-package/src/prelaunch_data_generator/el_cl_genesis/el_cl_genesis_data.star",
)

NETWORKS_DIR_PATH = "/kurtosis/src/networks/"

NETWORKS = struct(
    kurtosis_devnet = "kurtosis-devnet/",
)

def get_genesis_data(plan):
    genesis_file = plan.upload_files(
        NETWORKS_DIR_PATH + NETWORKS.kurtosis_devnet,
        "el_cl_genesis_data",
    )

    return el_cl_genesis_data.new_el_cl_genesis_data(
        genesis_file,

        # The following fields are not relevant for our testing, but are required by the parent
        "",  # genesis_validators_root
        0,  # prague_time
    )

def get_genesis_path(network = "kurtosis-devnet"):
    return NETWORKS_DIR_PATH + network
