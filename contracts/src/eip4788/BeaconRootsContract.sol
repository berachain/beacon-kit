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
    uint256 private constant HISTORY_BUFFER_LENGTH = 256;
    uint256 private constant BEACON_ROOT_OFFSET = HISTORY_BUFFER_LENGTH;
    uint256 private constant COINBASE_OFFSET = BEACON_ROOT_OFFSET + HISTORY_BUFFER_LENGTH;

    // SYSTEM_ADDRESS is the address that is allowed to call the set function
    // as defined in EIP-4788: https://eips.ethereum.org/EIPS/eip-4788
    address private constant SYSTEM_ADDRESS = 0xffffFFFfFFffffffffffffffFfFFFfffFFFfFFfE;

    /// @dev The selector for "getCoinbase(uint256)".
    bytes4 private constant GET_COINBASE_SELECTOR = 0xe8e284b9;

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                        ENTRYPOINT                          */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /// @notice Conforming to EIP-4788, this contract follows two execution paths:
    /// 1. If it is called by the SYSTEM_ADDRESS, the calldata is the 32-byte encoded beacon block root.
    /// 2. If it is called by any other address, there are two possible scenarios:
    ///    a. If the calldata is the 32-byte encoded timestamp, the function will return the beacon block root.
    ///    b. If the calldata is the 4-bytes selector for "getCoinbase(uint256)" appended with the 32-byte encoded
    ///       block number, the function will return the coinbase for the given block number.
    fallback() external {
        if (msg.sender != SYSTEM_ADDRESS) {
            if (msg.data.length == 36 && bytes4(msg.data) == GET_COINBASE_SELECTOR) {
                getCoinbase(uint256(bytes32(msg.data[4:36])));
            } else {
                // if the first 32 bytes is a timestamp, the first 4 bytes must be 0
                get();
            }
        } else {
            set();
        }
    }

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                       BEACON ROOT                          */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

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
    /// Reverts if the calldata is not a 32-byte timestamp or if the timestamp is 0.
    /// Reverts if the timestamp is not within the circular buffer.
    /// @return The beacon root associated with the given timestamp.
    function get() internal view returns (bytes32) {
        assembly ("memory-safe") {
            if iszero(and(eq(calldatasize(), 0x20), gt(calldataload(0), 0))) { revert(0, 0) }
        }
        uint256 block_idx = binarySearch();
        assembly ("memory-safe") {
            let _timestamp := sload(block_idx)
            if iszero(eq(_timestamp, calldataload(0))) { revert(0, 0) }
            let root_idx := add(block_idx, BEACON_ROOT_OFFSET)
            mstore(0, sload(root_idx))
            return(0, 0x20)
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
        uint256 _COINBASE_OFFSET = COINBASE_OFFSET;
        assembly ("memory-safe") {
            let block_idx := mod(number(), HISTORY_BUFFER_LENGTH)
            // set the timestamp
            sstore(block_idx, timestamp())
            let root_idx := add(block_idx, BEACON_ROOT_OFFSET)
            // set the beacon root
            sstore(root_idx, calldataload(0))
            let coinbase_idx := add(block_idx, _COINBASE_OFFSET)
            // set the coinbase
            sstore(coinbase_idx, coinbase())
        }
    }

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                         COINBASE                           */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /// @notice Retrieves the coinbase for a given block number.
    /// @dev if called with a block number that is before the history buffer
    /// it will return the coinbase for blockNumber + HISTORY_BUFFER_LENGTH * A
    /// Where A is the number of times the buffer has cycled since the blockNumber
    /// @param blockNumber The block number for which to retrieve the coinbase.
    /// @return The coinbase for the given block number.
    function getCoinbase(uint256 blockNumber) internal view returns (address) {
        uint256 _COINBASE_OFFSET = COINBASE_OFFSET;
        assembly ("memory-safe") {
            let block_idx := mod(blockNumber, HISTORY_BUFFER_LENGTH)
            let coinbase_idx := add(block_idx, _COINBASE_OFFSET)
            mstore(0, sload(coinbase_idx))
            return(0, 0x20)
        }
    }

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                 TIMESTAMP TO BLOCK NUMBER                  */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /// @dev Retrieves the block index for a given timestamp using binary search on the circular buffer.
    /// @dev Precondition: Any two consecutive timestamps in the circular buffer are strictly increasing.
    function binarySearch() internal view returns (uint256 block_idx) {
        assembly ("memory-safe") {
            let high := mod(number(), HISTORY_BUFFER_LENGTH)
            let low := mod(add(number(), 1), HISTORY_BUFFER_LENGTH)
            // revert if the timestamp is not within the circular buffer
            if or(lt(calldataload(0), sload(low)), gt(calldataload(0), sload(high))) { revert(0, 0) }
            for {} 1 {} {
                let high_adjusted := add(high, mul(lt(high, low), HISTORY_BUFFER_LENGTH))
                block_idx := mod(shr(1, add(low, high_adjusted)), HISTORY_BUFFER_LENGTH)
                if eq(low, high) { break }
                if lt(sload(block_idx), calldataload(0)) {
                    low := mod(add(block_idx, 1), HISTORY_BUFFER_LENGTH)
                    continue
                }
                high := block_idx
            }
        }
    }
}
