# get returns a string that can be used to install a package with apk
# get :: [String] -> String
def get(deps):
    return "apk add {} > /dev/null 2>&1".format(" ".join(deps))
