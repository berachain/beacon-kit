// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

contract BeaconKitRootsContract {

    uint256 constant HISTORY_BUFFER_LENGTH = 256;
    fallback() external {
        address systemAddress = 0xffffFFFfFFffffffffffffffFfFFFfffFFFfFFfE;
        if (msg.sender == systemAddress) {
            // Call set function
            set();
        } else {
            // Call get function
            get();
        }
    }


    // From: https://eips.ethereum.org/EIPS/eip-4788
    //
    // def set():
    //     timestamp_idx = to_uint256_be(evm.timestamp) % HISTORY_BUFFER_LENGTH
    //     root_idx = timestamp_idx + HISTORY_BUFFER_LENGTH

    //     storage.set(timestamp_idx, evm.timestamp)
    //     storage.set(root_idx, evm.calldata)
    function set() internal {
        assembly {
            let timestamp_idx := mod(timestamp(), HISTORY_BUFFER_LENGTH)
            let root_idx := add(timestamp_idx, HISTORY_BUFFER_LENGTH)
            sstore(timestamp_idx, timestamp())
            sstore(root_idx, calldataload(0))
        }
    }
    

    // From: https://eips.ethereum.org/EIPS/eip-4788
    //
    // def get():
    //     if len(evm.calldata) != 32:
    //         evm.revert()

    //     if to_uint256_be(evm.calldata) == 0:
    //         evm.revert()

    //     timestamp_idx = to_uint256_be(evm.calldata) % HISTORY_BUFFER_LENGTH
    //     timestamp = storage.get(timestamp_idx)

    //     if timestamp != evm.calldata:
    //         evm.revert()

    //     root_idx = timestamp_idx + HISTORY_BUFFER_LENGTH
    //     root = storage.get(root_idx)

    //     evm.return(root)
    function get() internal view returns (bytes32) {
        assembly {
            if iszero(eq(calldatasize(), 32)) {
                revert(0, 0)
            }
            if eq(calldataload(0), 0) {
                revert(0, 0)
            }
            let timestamp_idx := mod(calldataload(0), HISTORY_BUFFER_LENGTH)
            let _timestamp := sload(timestamp_idx)
            if iszero(eq(_timestamp, calldataload(0))) {
                revert(0, 0)
            }
            let root_idx := add(timestamp_idx, HISTORY_BUFFER_LENGTH)
            let root := sload(root_idx)
            return(root, 32)
        }
    }
    
}





