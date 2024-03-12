shared_utils = import_module("github.com/kurtosis-tech/ethereum-package/src/shared_utils/shared_utils.star")

service_config_lib = import_module("../../lib/service_config.star")
port_spec_lib = import_module("../../lib/port_spec.star")

LOKI_IMAGE = "grafana/loki:main"
PROMTAIL_IMAGE = "grafana/promtail:main"

LOCAL_CONFIG_PATH = "./config"
CONTAINER_CONFIG_PATH = "/mnt/config"
CONFIG_FOLDER = "config"
LOKI_CONFIG_FILENAME = "loki-config.yaml"
PROMTAIL_CONFIG_FILENAME = "promtail-config.yaml"
LOKI_CONFIG_ARTIFACT = "loki-config"
PROMTAIL_CONFIG_ARTIFACT = "promtail-config"

LOKI_PORT_ID = "http"
LOKI_PORT_NUM = 3100
LOKI_PORT_SPEC_TEMPLATE = port_spec_lib.get_port_spec_template(LOKI_PORT_NUM, "TCP", shared_utils.HTTP_APPLICATION_PROTOCOL)

# LOKI_ENTRYPOINT = ["sh", "-c"]
# LOKI_CMD =

LOKI_FILES = {
    "/mnt/config": LOKI_CONFIG_ARTIFACT,
}

PROMTAIL_PORT_ID = "http"
PROMTAIL_PORT_NUM = 9080
PROMTAIL_PORT_SPEC_TEMPLATE = port_spec_lib.get_port_spec_template(PROMTAIL_PORT_NUM, "TCP", shared_utils.HTTP_APPLICATION_PROTOCOL)
PROMTAIL_CMD = ["-config.file=/mnt/config/promtail-config.yaml"]


PROMTAIL_FILES = {
    "/mnt/config": PROMTAIL_CONFIG_ARTIFACT
}

def upload_global_files(plan):
    plan.upload_files(
        src = LOCAL_CONFIG_PATH + "/" + LOKI_CONFIG_FILENAME,
        name = LOKI_CONFIG_ARTIFACT,
    )

    plan.upload_files(
        src = LOCAL_CONFIG_PATH + "/" + PROMTAIL_CONFIG_FILENAME,
        name = PROMTAIL_CONFIG_ARTIFACT,
    )

def start(plan):
    loki_config = service_config_lib.get_service_config_template(
        "loki",
        LOKI_IMAGE,
        ports = {LOKI_PORT_ID: LOKI_PORT_SPEC_TEMPLATE},
        files = LOKI_FILES,
        # entrypoint = node_module.ENTRYPOINT,
        # cmd = node_module.CMD
    )

    loki_service_config = service_config_lib.create_from_config(loki_config)
    plan.add_service(
        name = "loki",
        config = loki_service_config,
    )

    promtail_config = service_config_lib.get_service_config_template(
        "promtail",
        PROMTAIL_IMAGE,
        ports = {PROMTAIL_PORT_ID: PROMTAIL_PORT_SPEC_TEMPLATE},
        files = PROMTAIL_FILES,
        # entrypoint = node_module.ENTRYPOINT,
        cmd = PROMTAIL_CMD
    )

    promtail_service_config = service_config_lib.create_from_config(promtail_config)
    plan.add_service(
        name = "promtail",
        config = promtail_service_config,
    )
