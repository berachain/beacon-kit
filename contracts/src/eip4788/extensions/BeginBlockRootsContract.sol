// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

contract BeginBlockRootsContract {
    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                    EVENTS/ERRORS                           */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /**
     * @dev Emitted when the BeginBlocker function is called.
     * @param contractAddress The address of the contract that called the BeginBlocker function.
     * @param coinbase The address of the current block miner.
     * @param selector The function selector that was called.
     */
    event BeginBlockerCalled(
        address indexed contractAddress, address indexed coinbase, bytes4 selector, bool success
    );

    /**
     * @dev Emitted when a BeginBlocker with the specified index does not exist.
     * @param i The index of the BeginBlocker.
     */
    error BeginBlockerDoesNotExist(uint256 i);

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

    /// @dev The BeginBlocker struct is used to store the calls we need to make at the beginning of
    /// each block.
    struct BeginBlocker {
        address contractAddress;
        bytes4 selector;
    }

    /// @dev Actions that can set a new BeginBlocker.
    bytes32 private constant SET = keccak256("SET");

    /// @dev Action that can remove BeginBlockers from the array.
    bytes32 private constant REMOVE = keccak256("REMOVE");

    /// @dev Action that can update the ADMIN address.
    bytes32 private constant UPDATE_ADMIN = keccak256("UPDATE_ADMIN");

    /// @dev The ADMIN address is the only address that can add or remove BeginBlockers.
    address private ADMIN = address(0x20f33CE90A13a4b5E7697E3544c3083B8F8A51D4);

    /// @dev The list of BeginBlockers that we need to call at the beginning of each block, in
    /// order.
    BeginBlocker[] public beginBlockers;

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                        ENTRYPOINT                          */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    fallback() external {
        if (msg.sender != SYSTEM_ADDRESS) {
            if (msg.data.length == 36 && bytes4(msg.data) == GET_COINBASE_SELECTOR) {
                getCoinbase();
            } else if (msg.sender == ADMIN) {
                // Only the ADMIN can crud BeginBlockers and ADMIN.
                _crud(msg.data);
            } else {
                get();
            }
        } else {
            set();

            // Run all the BeginBlockers.
            _run();
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
    /// @return The coinbase for the given block number.
    function getCoinbase() internal view returns (address) {
        assembly ("memory-safe") {
            let block_idx := mod(calldataload(4), HISTORY_BUFFER_LENGTH)
            let coinbase_idx := add(block_idx, _coinbases.slot)
            mstore(0, sload(coinbase_idx))
            return(0, 0x20)
        }
    }

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                        BeginBlocker                        */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /**
     * @dev Parses the BeginBlocker message data.
     * @notice The input data must be encoded as: abi.encode(i, contractAddress, selector)
     * @param data The input data containing the BeginBlocker message.
     * @return i The index of the BeginBlocker.
     * @return BeginBlocker The BeginBlocker struct containing the contract address and the
     * selector.
     */
    function _parse(bytes memory data)
        private
        pure
        returns (uint256, bytes32, BeginBlocker memory, address)
    {
        // Decode the data to get the BeginBlocker struct, user must send a message that is
        // encoded:
        // abi.encode(i, action, contractAddress, selector, admin)
        (uint256 i, bytes32 action, address contractAddress, bytes4 selector, address admin) =
            abi.decode(data, (uint256, bytes32, address, bytes4, address));
        return (i, action, BeginBlocker(contractAddress, selector), admin);
    }

    /**
     * @dev Sets the BeginBlocker at the given index.
     * @param i The index of the BeginBlocker.
     * @param beginBlocker The BeginBlocker struct containing the contract address and the
     * selector.
     */
    function _add(uint256 i, BeginBlocker memory beginBlocker) private {
        // If the index is greater than the length of the array, we need to append the BeginBlocker
        // to the end of the array.
        if (i >= beginBlockers.length) {
            beginBlockers.push(beginBlocker);
            return;
        }

        // Shift all the elements after the index to the right by one.
        for (uint256 j = beginBlockers.length; j > i; j--) {
            beginBlockers[j] = beginBlockers[j - 1];
        }

        // Insert the BeginBlocker at the index.
        beginBlockers[i] = beginBlocker;
    }

    /**
     * @dev Removes the BeginBlocker at the given index.
     * @param i The index of the BeginBlocker.
     */
    function _remove(uint256 i) private {
        // Check if we are trying to remove a BeginBlocker that does not exist.
        if (i >= beginBlockers.length) {
            revert BeginBlockerDoesNotExist(i);
        }

        // Shift all the elements after the index to the left by one.
        for (uint256 j = i; j < beginBlockers.length - 1; j++) {
            beginBlockers[j] = beginBlockers[j + 1];
        }

        // Remove the last element from the array.
        beginBlockers.pop();
    }

    /**
     * @dev Performs the CRUD operation based on the action specified in the input data.
     * @param data The input data containing the index, action, and BeginBlocker. The action can be
     * "set" or "remove".
     * If the action is "set", the function will add the BeginBlocker at the given index.
     * If the action is "remove", the function will remove the BeginBlocker at the given index.
     */
    function _crud(bytes memory data) private {
        // Decode the data we get from the message.
        (uint256 i, bytes32 action, BeginBlocker memory beginBlocker, address admin) = _parse(data);

        // Prefrom the CRUD operation.
        if (action == SET) {
            // If the action is "SET", we need to add the BeginBlocker at the given index.
            _add(i, beginBlocker);
        } else if (action == REMOVE) {
            // If the action is "REMOVE", we need to remove the BeginBlocker at the given index.
            _remove(i);
        } else if (action == UPDATE_ADMIN) {
            // If the action is "UPDATE_ADMIN", we need to update the ADMIN address.
            // This can only be done by the current ADMIN, since we check that in the fallback
            // method.
            if (admin != address(0)) {
                ADMIN = admin;
            }
        }
    }

    /**
     * @dev Runs all the BeginBlocker functions stored in the contract.
     * It iterates over the array of BeginBlockers and calls each one using its contract address
     * and selector.
     * If the call is successful, it emits a BeginBlockerCalled event with the contract address,
     * the current block miner's address, the selector, and the success status.
     */
    function _run() private {
        for (uint256 i = 0; i < beginBlockers.length; i++) {
            (bool success,) = beginBlockers[i].contractAddress.call(
                abi.encodeWithSelector(beginBlockers[i].selector)
            );
            emit BeginBlockerCalled(
                beginBlockers[i].contractAddress,
                block.coinbase,
                beginBlockers[i].selector,
                success
            );
        }
    }
}
