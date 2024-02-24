// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

/// @title BeaconRootCaller
/// @dev Designed for testing to call and retrieve the beacon block root.
/// Interacts with a predefined address to fetch the beacon block root using
/// the current block timestamp.
contract BeaconRootCaller {
    /// @notice The beacon root address, set as a constant for this test.
    address internal constant BEACON_ROOT_ADDRESS = 0x000F3df6D732807Ef1319fB7B8bB8522d0Beac02;

    /// @notice Calls the beacon root address to retrieve the beacon block root.
    /// @dev Performs a static call to `BEACON_ROOT_ADDRESS` with the current
    /// block timestamp encoded as the call data. Requires that the call
    /// succeeds, otherwise reverts with an error message.
    /// @return The beacon block root as a bytes32 value.
    function getBeaconBlockRoot() external view returns (bytes32) {
        (bool ok, bytes memory result) =
            BEACON_ROOT_ADDRESS.staticcall(abi.encode(block.timestamp));
        require(ok, "BeaconRootCaller: call failed");
        return abi.decode(result, (bytes32));
    }
}
