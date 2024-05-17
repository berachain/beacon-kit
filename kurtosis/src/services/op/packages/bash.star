# run returns a string which executes each command in
# cmds in order, joined by &&.
def run(cmds):
    return " && ".join(cmds)
