// SPDX-License-Identifier: MIT
pragma solidity ^0.8.21;

abstract contract Verifier {
    address public constant BEACON_ROOTS =
        0x000F3df6D732807Ef1319fB7B8bB8522d0Beac02;

    error IndexOutOfRange();
    error InvalidProof();
    error RootNotFound();

    function getParentBlockRoot(uint64 ts)
        internal
        view
        returns (bytes32 root)
    {
        assembly ("memory-safe") {
            mstore(0, ts)
            let success := staticcall(gas(), BEACON_ROOTS, 0, 0x20, 0, 0x20)
            if iszero(success) {
                mstore(0, 0x3033b0ff) // RootNotFound()
                revert(0x1c, 0x04)
            }
            root := mload(0)
        }
    }
}
