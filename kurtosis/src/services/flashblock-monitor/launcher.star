# Flashblock monitor service launcher for viewing flashblocks from the sequencer.
# This service connects to the sequencer's WebSocket and displays received flashblocks.

SERVICE_NAME = "flashblock-monitor"

# Use alpine with required packages
DEFAULT_IMAGE = "alpine:3.19"

# Resource limits (minimal since this is just a monitoring service)
MIN_CPU = 100
MAX_CPU = 500
MIN_MEMORY = 128
MAX_MEMORY = 512

def launch_flashblock_monitor(
        plan,
        sequencer_el_service_name,
        flashblock_ws_port = 8548,
        image = DEFAULT_IMAGE):
    """
    Launch the flashblock monitor service.

    Args:
        plan: The Kurtosis plan object.
        sequencer_el_service_name: The service name of the sequencer's EL client.
        flashblock_ws_port: The WebSocket port for flashblocks (default: 8548).
        image: The Docker image to use.

    Returns:
        The service context for the launched monitor.
    """
    # Get the sequencer service to find its IP
    sequencer_service = plan.get_service(sequencer_el_service_name)
    sequencer_ip = sequencer_service.ip_address

    ws_url = "ws://{}:{}".format(sequencer_ip, flashblock_ws_port)

    # Build the monitor command directly - avoids template issues
    # Output raw JSON from websocat
    monitor_cmd = "echo '=============================================='" + \
        " && echo '  Flashblock Monitor'" + \
        " && echo '  Sequencer: " + ws_url + "'" + \
        " && echo '=============================================='" + \
        " && apk add --no-cache wget > /dev/null 2>&1" + \
        " && echo 'Installing websocat...'" + \
        " && wget -q -O /usr/local/bin/websocat https://github.com/vi/websocat/releases/download/v1.13.0/websocat.x86_64-unknown-linux-musl" + \
        " && chmod +x /usr/local/bin/websocat" + \
        " && echo 'Connecting to flashblock WebSocket...'" + \
        " && while true; do /usr/local/bin/websocat --text -E '" + ws_url + "' 2>/dev/null || echo 'Reconnecting...'; sleep 3; done"

    config = ServiceConfig(
        image = image,
        entrypoint = ["sh", "-c"],
        cmd = [monitor_cmd],
        min_cpu = MIN_CPU,
        max_cpu = MAX_CPU,
        min_memory = MIN_MEMORY,
        max_memory = MAX_MEMORY,
    )

    return plan.add_service(SERVICE_NAME, config)
