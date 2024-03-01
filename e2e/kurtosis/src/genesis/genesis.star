el_cl_genesis_data = import_module(
    "github.com/kurtosis-tech/ethereum-package/src/prelaunch_data_generator/el_cl_genesis/el_cl_genesis_data.star"
)
static_files = import_module('../static_files/static_files.star')

def get_genesis_data(plan):
    genesis_file = plan.upload_files(
        static_files.GENESIS_FILEPATH,
        "el_cl_genesis_data",
    )


    return el_cl_genesis_data.new_el_cl_genesis_data(
        genesis_file,

        # The following fields are not relevant for our testing, but are required by the parent
        "", # genesis_validators_root
        0,  # cancun_time
        0,  # prague_time
    )