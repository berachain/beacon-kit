def parse_service_from_dict(service):
    if "replicas" not in service:
        service["replicas"] = 1
    if "client" not in service:
        service["client"] = None
    if "verifier_image" not in service:
        service["verifier_image"] = None
    return struct(
        name = service["name"],
        replicas = service["replicas"],
        client = service["client"],
        verifier_image = service["verifier_image"],
    )
