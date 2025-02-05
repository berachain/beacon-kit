// SPDX-License-Identifier: MIT
pragma solidity ^0.8.21;

/// @notice A library for reading beacon roots using EIP-4788.
/// @author Berachain Team
/// @author Inspired by [madlabman](https://github.com/madlabman/eip-4788-proof)
library BeaconRoots {
    /// @notice The address of the EIP-4788 Beacon Roots contract.
    address public constant ADDRESS = 0x000F3df6D732807Ef1319fB7B8bB8522d0Beac02;

    // Signature: 0x3033b0ff
    error RootNotFound();

    /// @notice Checks if the parent block root exists at a given timestamp.
    function isParentBlockRootAt(uint64 ts) internal view returns (bool success) {
        assembly ("memory-safe") {
            mstore(0, ts)
            success := staticcall(gas(), ADDRESS, 0, 0x20, 0, 0x20)
        }
    }

    /// @notice Get the parent block root at a given timestamp.
    /// @dev Reverts with `RootNotFound()` if the root is not found.
    function getParentBlockRootAt(uint64 ts) internal view returns (bytes32 root) {
        assembly ("memory-safe") {
            mstore(0, ts)
            let success := staticcall(gas(), ADDRESS, 0, 0x20, 0, 0x20)
            if iszero(success) {
                mstore(0, 0x3033b0ff) // RootNotFound()
                revert(0x1c, 0x04)
            }
            root := mload(0)
        }
    }
}
