grafana = import_module('github.com/kurtosis-tech/grafana-package/main.star')

def start(plan, prometheus_url):
    grafana.run(plan, prometheus_url, "github.com/berachain/beacon-kit/kurtosis/src/observability/grafana/dashboards")