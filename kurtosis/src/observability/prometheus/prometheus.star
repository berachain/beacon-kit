prometheus = import_module("github.com/kurtosis-tech/prometheus-package/main.star")

"""
A service should follow the format of below example:

{
    ## required
    "name": "api_service",
    "service": Service(
        ip_address="0.0.0.0",
        ports={
            "metrics": PortSpec(
                number=8080
            )
        }
    ),
    "metrics_path": "/metrics",

    ## optional
    "labels": {
        "service_type": "api"
    },
    "scrape_interval": "60s"
}
"""

def start(plan, services):
    metrics_jobs = []
    for service in services:
        constant_labels = {}  # use no constant labels if none provided
        if "labels" in service:
            constant_labels = service["labels"]

        scrape_interval = prometheus.DEFAULT_SCRAPE_INTERVAL  # use 5s as default scrape interval
        if "scrape_interval" in service:
            scrape_interval = service["scrape_interval"]

        metrics_job = {
            "Name": "{0}".format(service["name"]),
            "Endpoint": "{0}:{1}".format(service["service"].ip_address, service["service"].ports["metrics"].number),
            "Labels": constant_labels,
            "MetricsPath": service["metrics_path"],
            "ScrapeInterval": scrape_interval,
        }
        metrics_jobs.append(metrics_job)

    prometheus_url = prometheus.run(plan, metrics_jobs)
    plan.print(prometheus_url)
    return prometheus_url
