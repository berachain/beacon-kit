def exec_on_service(plan, service_name, command):
    return plan.exec(
        service_name = service_name,
        recipe = ExecRecipe(
            command = ["bash", "-c", command],
        ),
    )
