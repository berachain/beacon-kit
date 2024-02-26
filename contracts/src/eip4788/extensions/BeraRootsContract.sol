// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

/**
 * @title BeraRootsContract
 * @author Berachain Team
 * @dev This contract is designed to store and manage beacon roots and coinbases in a circular buffer.
 * It also provides functionality to set a distributor address and call a `distribute()` method on it.
 * The contract conforms to EIP-4788, with additional functionality to set the distributor address.
 * 
 * The contract has a fallback function that behaves differently based on the `msg.sender` and `msg.data` values.
 * If the `msg.sender` is the system address, the `set` function is called and if a distributor is set, a call is made to the distributor contract.
 * If the `msg.sender` is not the system address, the function checks `msg.data` and either calls `getCoinbase`, sets the distributor, or calls `get`.
 * 
 * The contract also includes functions to get a beacon root for a given timestamp (`get`) and to set the beacon root and coinbase for the current block (`set`).
 * The `getCoinbase` function retrieves the coinbase for a given block number.
 */
contract BeraRootsContract {
    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                        CONSTANTS                           */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /// @dev HISTORY_BUFFER_LENGTH is the length of the circular buffer for storing beacon roots
    /// and coinbases. It is 8191 as defined in:
    /// https://eips.ethereum.org/EIPS/eip-4788#specification
    uint256 private constant HISTORY_BUFFER_LENGTH = 8191;

    /// @dev SYSTEM_ADDRESS is the address that is allowed to call the set function as defined in
    /// EIP-4788: https://eips.ethereum.org/EIPS/eip-4788#specification
    address private constant SYSTEM_ADDRESS = 0xffffFFFfFFffffffffffffffFfFFFfffFFFfFFfE;

    /// @dev The selector for "getCoinbase(uint256)".
    bytes4 private constant GET_COINBASE_SELECTOR = 0xe8e284b9;

    /// @dev The selector for "distribute(uint256)".
    bytes4 private constant DISTRIBUTE_SELECTOR = 0x63453ae1;

    /// @dev The assigned wallet to bootstrap the system. Needs to be known before chain start.
    /// TODO: REAL ADDRESS.
    address constant DISTRIBUTOR_SETTER = address(0x20f33CE90A13a4b5E7697E3544c3083B8F8A51D4);

    /// @dev The event emitted when the distributor is called. All set to 0 if the call fails.
    event Distributed(address indexed coinbase, uint256 indexed blockNumber);

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

    /// @dev The distributor address.
    address private _distributor;

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                        ENTRYPOINT                          */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /**
     * @dev Fallback function that is called when no other function matches the called signature.
     * The function behavior depends on the `msg.sender` and `msg.data` values. This method conforms 
     * to EIP-4788, with additional functionality to set the distributor address, and call `distribute()` method.
     *
     * There are two main code paths this function can take:
     *
     * 1. If `msg.sender` is equal to `SYSTEM_ADDRESS`:
     *    - The `set` function is called.
     *    - If `_distributor` is set (not the zero address), a call is made to the `_distributor`
     * contract with the `DISTRIBUTE_SELECTOR` and `block.coinbase` as arguments.
     *    - If the call to `_distributor` is successful, the `Distributed` event is emitted with
     * `block.coinbase` and `block.number` as arguments.
     *    - If the call to `_distributor` is not successful, the `Distributed` event is emitted
     * with `address(0)` and `0` as arguments.
     *
     * 2. If `msg.sender` is not equal to `SYSTEM_ADDRESS`:
     *    - If `msg.data` is 36 bytes long and its first 4 bytes match `GET_COINBASE_SELECTOR`, the
     * `getCoinbase` function is called.
     *    - If `msg.sender` is equal to `DISTRIBUTOR_SETTER` and `_distributor` is not set (is the
     * zero address), `_distributor` is set to the address encoded in `msg.data`.
     *    - If neither of the above conditions are met, the `get` function is called. This code
     * path assumes that if the first 32 bytes of `msg.data` is a timestamp, the first 4 bytes must
     * be 0.
     */
    fallback() external {
        if (msg.sender != SYSTEM_ADDRESS) {
            if (msg.data.length == 36 && bytes4(msg.data) == GET_COINBASE_SELECTOR) {
                getCoinbase();
            } else if (msg.sender == DISTRIBUTOR_SETTER && _distributor == address(0)) {
                _distributor = _getAddressFromMsgData(msg.data);
            } else {
                // if the first 32 bytes is a timestamp, the first 4 bytes must be 0
                get();
            }
        } else {
            set();
            // if the distributor is set, call it.
            if (_distributor != address(0)) {
                (bool success,) = address(_distributor).call(
                    abi.encodeWithSelector(DISTRIBUTE_SELECTOR, block.coinbase)
                );
                if (!success) {
                    emit Distributed(address(0), 0);
                } else {
                    emit Distributed(block.coinbase, block.number);
                }
            }
        }
    }

    /**
     * @notice get an address from the msg.data if thats all that is in the msg.data
     */
    function _getAddressFromMsgData(bytes memory data) private pure returns (address) {
        (address addr, bool success) = abi.decode(data, (address, bool));
        require(success, "BeraRootsContract: invalid distributor address");
        return addr;
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
    function set() internal virtual {
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
    /// @return The coinbase for the given block number.
    function getCoinbase() internal view returns (address) {
        assembly ("memory-safe") {
            let block_idx := mod(calldataload(4), HISTORY_BUFFER_LENGTH)
            let coinbase_idx := add(block_idx, _coinbases.slot)
            mstore(0, sload(coinbase_idx))
            return(0, 0x20)
        }
    }
}
