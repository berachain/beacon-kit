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

    /// @dev HISTORY_BUFFER_LENGTH is the length of the circular buffer for storing beacon roots
    /// and coinbases.
    uint256 private constant HISTORY_BUFFER_LENGTH = 256;

    /// @dev SYSTEM_ADDRESS is the address that is allowed to call the set function as defined in
    /// EIP-4788.
    address private constant SYSTEM_ADDRESS = 0xffffFFFfFFffffffffffffffFfFFFfffFFFfFFfE;

    /// @dev The selector for "getCoinbase(uint256)".
    bytes4 private constant GET_COINBASE_SELECTOR = 0xe8e284b9;

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                        STORAGE                           */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /// @dev The circular buffer for storing timestamps.
    uint256[HISTORY_BUFFER_LENGTH] private _timestamps;

    /// @dev The circular buffer for storing beacon roots.
    bytes32[HISTORY_BUFFER_LENGTH] private _beaconRoots;

    /// @dev The circular buffer for storing coinbases.
    address[HISTORY_BUFFER_LENGTH] private _coinbases;

    /// @dev The mapping of timestamps to block numbers.
    mapping(uint256 timestamp => uint256 blockNumber) private _blockNumbers;

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                        ENTRYPOINT                          */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /// @notice Conforming to EIP-4788, this contract follows two execution paths:
    /// 1. If it is called by the SYSTEM_ADDRESS, the calldata is the 32-byte encoded beacon block
    /// root.
    /// 2. If it is called by any other address, there are two possible scenarios:
    ///    a. If the calldata is the 32-byte encoded timestamp, the function will return the beacon
    /// block root.
    ///    b. If the calldata is the 4-bytes selector for "getCoinbase(uint256)" appended with the
    /// 32-byte encoded
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

    /// @dev Retrieves the beacon root for a given timestamp.
    /// This function is called internally and utilizes assembly for direct storage access.
    /// Reverts if the calldata is not a 32-byte timestamp or if the timestamp is 0.
    /// Reverts if the timestamp is not within the circular buffer.
    /// @return The beacon root associated with the given timestamp.
    function get() internal view returns (bytes32) {
        assembly ("memory-safe") {
            if iszero(and(eq(calldatasize(), 0x20), gt(calldataload(0), 0))) { revert(0, 0) }
            // index block number from timestamp
            mstore(0, calldataload(0))
            mstore(0x20, _blockNumbers.slot)
            let block_number := sload(keccak256(0, 0x40))
            let block_idx := mod(block_number, HISTORY_BUFFER_LENGTH)
            let _timestamp := sload(block_idx)
            if iszero(eq(_timestamp, calldataload(0))) { revert(0, 0) }
            let root_idx := add(block_idx, _beaconRoots.slot)
            mstore(0, sload(root_idx))
            return(0, 0x20)
        }
    }

    /// @dev Sets the beacon root and coinbase for the current block.
    /// This function is called internally and utilizes assembly for direct storage access.
    function set() internal {
        assembly ("memory-safe") {
            let block_idx := mod(number(), HISTORY_BUFFER_LENGTH)
            // clean the key in the mapping for the stale timestamp in the block index to be
            // overridden
            let stale_timestamp := sload(block_idx)
            mstore(0, stale_timestamp)
            mstore(0x20, _blockNumbers.slot)
            sstore(keccak256(0, 0x40), 0)
            // override the timestamp
            sstore(block_idx, timestamp())
            // set the current block number in the mapping
            mstore(0, timestamp())
            // 0x20 is already set
            sstore(keccak256(0, 0x40), number())
            // set the beacon root
            let root_idx := add(block_idx, _beaconRoots.slot)
            sstore(root_idx, calldataload(0))
            // set the coinbase
            let coinbase_idx := add(block_idx, _coinbases.slot)
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
        assembly ("memory-safe") {
            let block_idx := mod(blockNumber, HISTORY_BUFFER_LENGTH)
            let coinbase_idx := add(block_idx, _coinbases.slot)
            mstore(0, sload(coinbase_idx))
            return(0, 0x20)
        }
    }
}
