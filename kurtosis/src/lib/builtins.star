"""
Helper library for type comparisons

Typical usage examples:
    type(None) == builtins.types.none # True

    hello = "hello"
    type(hello) == builtins.types.bool # False
"""

types = struct(
    none = "NoneType",  # the type of None
    bool = "bool",  # True or False
    int = "int",  # a signed integer of arbitrary magnitude
    float = "float",  # an IEEE 754 double-precision floating-point number
    string = "string",  # a text string, with Unicode encoded as UTF-8 or UTF-16
    bytes = "bytes",  # a byte string
    list = "list",  # a fixed-length sequence of values
    tuple = "tuple",  # a fixed-length sequence of values, unmodifiable
    dict = "dict",  # a mapping from values to values
    function = "function",  # a function
    serviceConfig = "ServiceConfig",  # a service configuration object
    portSpec = "PortSpec",  # a port spec object
    directory = "Directory",  # a directory object
    execRecipe = "ExecRecipe",  # an exec recipe object
    getHttpRequestRecipe = "GetHttpRequestRecipe",  # a get http request recipe object
    imageBuildSpec = "ImageBuildSpec",  # an image build spec object
    module = "module",  # a module
    # TODO(types): add remaining kurtosis types
)

# Types that can be used as keys in a dictionary due to being hashable
# Note: tuples can also be hashable but only if they contain hashable elements
hashable = [types.none, types.bool, types.int, types.float, types.string, types.bytes]
