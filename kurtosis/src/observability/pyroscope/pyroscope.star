"""
A package to install pyroscope:
"""

def run(plan, name = "pyroscope"):
    plan.add_service(name = name, config = ServiceConfig(
        image = "grafana/pyroscope:latest",
        ports = {
            "pyroscope": PortSpec(
                number = 4040,
                transport_protocol = "TCP",
                application_protocol = "http",
            ),
        },
    ))
