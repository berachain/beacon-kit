CONSENSUS_DEFAULT_SETTINGS = {
    "specs": {
        "min_cpu": 0,
        "max_cpu": 2000,
        "min_memory": 0,
        "max_memory": 2048,
    },
    "images": {
        "beaconkit": "beacond-dev:kurtosis-local",
    },
    "labels": {},
    "node_selectors": {},
    "config": {
        "timeout_propose": "3s",
        "timeout_prevote": "1s",
        "timeout_precommit": "1s",
        "timeout_commit": "1s",
        "max_num_inbound_peers": 40,
        "max_num_outbound_peers": 10,
    },
    "app": {
        "payload-timeout": "1.5s",
        "enable_optimistic_payload_builds": "false",
    },
}

EXECUTION_DEFAULT_SETTINGS = {
    "specs": {
        "min_cpu": 0,
        "max_cpu": 2000,
        "min_memory": 0,
        "max_memory": 2048,
    },
    "images": {
        "besu": "hyperledger/besu:latest",
        "erigon": "thorax/erigon:v2.60.1",
        "ethereumjs": "ethpandaops/ethereumjs:stable",
        "geth": "ethereum/client-go:latest",
        "nethermind": "nethermind/nethermind:latest",
        "reth": "ghcr.io/paradigmxyz/reth:latest",
    },
    "labels": {},
    "node_selectors": {},
}

CL_TYPE = "beaconkit"

def parse_nodes_from_dict(vals, settings):
    node_type = vals["type"]
    node_list = []

    consensus_settings = parse_consensus_settings(settings)
    execution_settings = parse_execution_settings(settings)

    count = 0
    for val_configuration in vals["nodes"]:
        replicas = 1
        if "replicas" in val_configuration:
            replicas = val_configuration["replicas"]
        if replicas != 0:
            for i in range(replicas):
                val_struct = parse_node_from_dict(node_type, val_configuration, consensus_settings, execution_settings, count)
                node_list.append(val_struct)
                count += 1

    return node_list

def parse_node_from_dict(node_type, val, consensus_settings, execution_settings, index):
    # if kzg implementation is not provided, give default
    kzg_impl = "crate-crypto/go-kzg-4844"
    if "kzg_impl" in val:
        kzg_impl = val["kzg_impl"]

    return struct(
        node_type = node_type,
        el_type = val["el_type"],
        el_image = execution_settings.images[val["el_type"]],
        cl_type = CL_TYPE,
        cl_image = consensus_settings.images[CL_TYPE],
        index = index,
        cl_service_name = "cl-{}-{}-{}".format(node_type, CL_TYPE, index),
        el_service_name = "el-{}-{}-{}".format(node_type, val["el_type"], index),
        consensus_settings = consensus_settings,
        execution_settings = execution_settings,
        kzg_impl = kzg_impl,
    )

def parse_consensus_settings(settings):
    consensus_settings = {}
    if "consensus_settings" in settings:
        consensus_settings = dict(settings["consensus_settings"])
    consensus_settings = parse_default_node_settings(consensus_settings, CONSENSUS_DEFAULT_SETTINGS)
    consensus_settings = parse_extra_consensus_settings(consensus_settings)

    return get_consensus_struct(consensus_settings)

def parse_extra_consensus_settings(settings):
    config_settings = parse_consensus_config_settings(settings["config"]) if "config" in settings else parse_consensus_config_settings(CONSENSUS_DEFAULT_SETTINGS["config"])
    app_settings = parse_consensus_app_settings(settings["app"]) if "app" in settings else parse_consensus_app_settings(CONSENSUS_DEFAULT_SETTINGS["app"])

    consensus_settings = dict(settings)
    consensus_settings["config"] = config_settings
    consensus_settings["app"] = app_settings

    return consensus_settings

def parse_consensus_config_settings(config_settings):
    config_settings = dict(config_settings)

    if "timeout_propose" not in config_settings:
        config_settings["timeout_propose"] = CONSENSUS_DEFAULT_SETTINGS["config"]["timeout_propose"]
    if "timeout_prevote" not in config_settings:
        config_settings["timeout_prevote"] = CONSENSUS_DEFAULT_SETTINGS["config"]["timeout_prevote"]
    if "timeout_precommit" not in config_settings:
        config_settings["timeout_precommit"] = CONSENSUS_DEFAULT_SETTINGS["config"]["timeout_precommit"]
    if "timeout_commit" not in config_settings:
        config_settings["timeout_commit"] = CONSENSUS_DEFAULT_SETTINGS["config"]["timeout_commit"]
    if "max_num_inbound_peers" not in config_settings:
        config_settings["max_num_inbound_peers"] = CONSENSUS_DEFAULT_SETTINGS["config"]["max_num_inbound_peers"]
    if "max_num_outbound_peers" not in config_settings:
        config_settings["max_num_outbound_peers"] = CONSENSUS_DEFAULT_SETTINGS["config"]["max_num_outbound_peers"]

    return struct(
        timeout_propose = config_settings["timeout_propose"],
        timeout_prevote = config_settings["timeout_prevote"],
        timeout_precommit = config_settings["timeout_precommit"],
        timeout_commit = config_settings["timeout_commit"],
        max_num_inbound_peers = config_settings["max_num_inbound_peers"],
        max_num_outbound_peers = config_settings["max_num_outbound_peers"],
    )

def parse_consensus_app_settings(app_settings):
    app_settings = dict(app_settings)

    if "payload_timeout" not in app_settings:
        app_settings["payload_timeout"] = CONSENSUS_DEFAULT_SETTINGS["app"]["payload_timeout"]
    if "enable_optimistic_payload_builds" not in app_settings:
        app_settings["enable_optimistic_payload_builds"] = CONSENSUS_DEFAULT_SETTINGS["app"]["enable_optimistic_payload_builds"]

    return struct(
        payload_timeout = app_settings["payload_timeout"],
        enable_optimistic_payload_builds = app_settings["enable_optimistic_payload_builds"],
    )

def parse_execution_settings(settings):
    execution_settings = {}
    if "execution_settings" in settings:
        execution_settings = dict(settings["execution_settings"])
    execution_settings = parse_default_node_settings(execution_settings, EXECUTION_DEFAULT_SETTINGS)
    return get_execution_struct(execution_settings)

def parse_default_node_settings(settings, default_settings):
    node_settings = dict(settings)
    if "specs" not in node_settings:
        node_settings["specs"] = default_settings["specs"]

    node_specs = dict(node_settings["specs"])
    if "min_cpu" not in node_specs:
        node_specs["min_cpu"] = default_settings["specs"]["min_cpu"]
    if "max_cpu" not in node_specs:
        node_specs["max_cpu"] = default_settings["specs"]["max_cpu"]
    if "min_memory" not in node_specs:
        node_specs["min_memory"] = default_settings["specs"]["min_memory"]
    if "max_memory" not in node_specs:
        node_specs["max_memory"] = default_settings["specs"]["max_memory"]
    if "images" not in node_settings:
        node_settings["images"] = default_settings["images"]
    if "labels" not in node_settings:
        node_settings["labels"] = default_settings["labels"]
    if "node_selectors" not in node_settings:
        node_settings["node_selectors"] = default_settings["node_selectors"]

    node_specs = struct(
        min_cpu = node_specs["min_cpu"],
        max_cpu = node_specs["max_cpu"],
        min_memory = node_specs["min_memory"],
        max_memory = node_specs["max_memory"],
    )
    node_settings["specs"] = node_specs

    return node_settings

def get_consensus_struct(consensus_settings):
    return struct(
        specs = consensus_settings["specs"],
        images = consensus_settings["images"],
        labels = consensus_settings["labels"],
        node_selectors = consensus_settings["node_selectors"],
        config = consensus_settings["config"],
        app = consensus_settings["app"],
    )

def get_execution_struct(execution_settings):
    return struct(
        specs = execution_settings["specs"],
        images = execution_settings["images"],
        labels = execution_settings["labels"],
        node_selectors = execution_settings["node_selectors"],
    )
