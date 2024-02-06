// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

/// @title BeaconRootsContract
/// @dev This contract is an implementation of the BeaconRootsContract as defined in EIP-4788.
/// It has been extended to include a coinbase storage slot for each block for use with
/// the Berachain Proof-of-Liquidity protocol.
/// @author https://eips.ethereum.org/EIPS/eip-4788
/// @author itsdevbear@berachain.com
/// @author rusty@berachain.com
/// @author po@berachain.com
contract BeaconRootsContract {
    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                        CONSTANTS                           */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    // HISTORY_BUFFER_LENGTH is the length of the circular buffer for
    // storing beacon roots and coinbases.
    uint256 constant HISTORY_BUFFER_LENGTH = 256;

    // SYSTEM_ADDRESS is the address that is allowed to call the set function
    // as defined in EIP-4788: https://eips.ethereum.org/EIPS/eip-4788
    address constant SYSTEM_ADDRESS = 0xffffFFFfFFffffffffffffffFfFFFfffFFFfFFfE;

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                        ENTRYPOINT                          */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    fallback() external {
        if (msg.sender == SYSTEM_ADDRESS) {
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
    //
    /// @dev Sets the beacon root and coinbase for the current block.
    /// This function is called internally and utilizes assembly for direct storage access.
    function set() internal {
        setBeaconRoot();
        setCoinbase();
    }

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                       BEACON ROOT                          */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /// @dev Sets the beacon root and coinbase for the current block.
    /// This function is called internally and utilizes assembly for direct storage access.
    function setBeaconRoot() internal {
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
    //
    //     evm.return(root)
    //
    /// @dev Retrieves the beacon root for a given timestamp.
    /// This function is called internally and utilizes assembly for direct storage access.
    /// @return The beacon root associated with the given timestamp.
    function get() internal view returns (bytes32) {
        assembly {
            if iszero(eq(calldatasize(), 32)) { revert(0, 0) }
            if iszero(calldataload(0)) { revert(0, 0) }
            let timestamp_idx := mod(calldataload(0), HISTORY_BUFFER_LENGTH)
            let _timestamp := sload(timestamp_idx)
            if iszero(eq(_timestamp, calldataload(0))) { revert(0, 0) }
            let root_idx := add(timestamp_idx, HISTORY_BUFFER_LENGTH)
            let root := sload(root_idx)
            mstore(0, root)
            return(0, 32)
        }
    }

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                         COINBASE                           */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /// @dev Sets the coinbase for the current block in storage.
    /// This function is called internally and utilizes assembly for direct storage access.
    function setCoinbase() internal {
        assembly {
            sstore(mod(number(), HISTORY_BUFFER_LENGTH), coinbase())
        }
    }

    /// @dev Retrieves the coinbase for a given block number.
    /// @dev if called with a block number that is before the history buffer
    /// it will return the coinbase for blockNumber + HISTORY_BUFFER_LENGTH * A
    /// Where A is the number of times the buffer has cycled since the blockNumber
    /// @param blockNumber The block number for which to retrieve the coinbase.
    /// @return The coinbase for the given block number.
    function getCoinbase(uint256 blockNumber) external view returns (address) {
        assembly {
            mstore(0, sload(mod(blockNumber, HISTORY_BUFFER_LENGTH))) // fix collision
            return(0, 0x20) // or 32
        }
    }
}
