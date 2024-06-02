CONSENSUS_DEFAULT_SETTINGS = {
    "min_cpu": 0,
    "max_cpu": 2000,
    "min_memory": 0,
    "max_memory": 2048,
}

EXECUTION_DEFAULT_SETTINGS = {
    "min_cpu": 0,
    "max_cpu": 2000,
    "min_memory": 0,
    "max_memory": 2048,
}

CL_TYPE = "beaconkit"

def parse_nodes_from_dict(vals, settings):
    node_type = vals["type"]
    node_list = []

    count = 0
    for val_configuration in vals["nodes"]:
        replicas = 1
        if "replicas" in val_configuration:
            replicas = val_configuration["replicas"]

        if replicas != 0:
            for i in range(replicas):
                val_struct = parse_node_from_dict(node_type, val_configuration, settings, count)
                node_list.append(val_struct)
                count += 1

    return node_list

def parse_node_from_dict(node_type, val, settings, index):
    return struct(
        node_type = node_type,
        el_type = val["el_type"],
        el_image = settings["execution_settings"]["images"][val["el_type"]],
        cl_type = CL_TYPE,
        cl_image = settings["consensus_settings"]["images"][CL_TYPE],
        index = index,
        cl_service_name = "cl-{}-{}-{}".format(node_type, CL_TYPE, index),
        el_service_name = "el-{}-{}-{}".format(node_type, val["el_type"], index),
    )

def parse_node_settings(settings):
    default_settings = {}
    if not "default" in settings:
        default_settings = {
            "consensus": {
                "min_cpu": CONSENSUS_DEFAULT_SETTINGS["min_cpu"],
                "max_cpu": CONSENSUS_DEFAULT_SETTINGS["max_cpu"],
                "min_memory": CONSENSUS_DEFAULT_SETTINGS["min_memory"],
                "max_memory": CONSENSUS_DEFAULT_SETTINGS["max_memory"],
            },
            "execution": {
                "min_cpu": EXECUTION_DEFAULT_SETTINGS["min_cpu"],
                "max_cpu": EXECUTION_DEFAULT_SETTINGS["max_cpu"],
                "min_memory": EXECUTION_DEFAULT_SETTINGS["min_memory"],
                "max_memory": EXECUTION_DEFAULT_SETTINGS["max_memory"],
            },
        }
    else:
        default_settings = dict(settings["default"])

def get_settings_struct(settings):
    return struct(
        consensus_settings = struct(
            min_cpu = settings["consensus"]["min_cpu"],
            max_cpu = settings["consensus"]["max_cpu"],
            min_memory = settings["consensus"]["min_memory"],
            max_memory = settings["consensus"]["max_memory"],
        ),
        execution_settings = struct(
            min_cpu = settings["execution"]["min_cpu"],
            max_cpu = settings["execution"]["max_cpu"],
            min_memory = settings["execution"]["min_memory"],
            max_memory = settings["execution"]["max_memory"],
        ),
    )
