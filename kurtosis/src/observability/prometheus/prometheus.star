prometheus = import_module("github.com/kurtosis-tech/prometheus-package/main.star")

def start(plan, services):
    metrics_jobs = []
    for service in services:
        metrics_job = {
            "Name": "{0} {1}".format(service["type"], service["index"]),
            "Endpoint": "{0}:{1}".format(service["service"].ip_address, service["service"].ports["metrics"].number),
            "MetricsPath": service["metrics_path"],
            "ScrapeInterval": "1s",
        }
        metrics_jobs.append(metrics_job)

    prometheus_url = prometheus.run(plan, metrics_jobs)
    plan.print(prometheus_url)
    return prometheus_url
