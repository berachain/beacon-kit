// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import { BeaconRootsContract } from "../BeaconRootsContract.sol";

/**
 * @title BeginBlockRootsContract
 * @author Berachain Team.
 *
 * @dev This contract extends the BeaconRootsContract and adds the BeginBlocker
 * functionality, where the specified BeginBlockers are called at the beginning
 * of each block.
 * @dev ADD THE STORAGE FOR ADMIN at slot 24_574 at genesis.
 * @dev Set slot 24_574 to an ADMIN address in your genesis file.
 *
 *  dP""b8  dP"Yb  .dP"Y8 8b    d8  dP"Yb  .dP"Y8     .dP"Y8 8888b.  88  dP
 * dP   `" dP   Yb `Ybo." 88b  d88 dP   Yb `Ybo."     `Ybo."  8I  Yb 88odP
 * Yb      Yb   dP o.`Y8b 88YbdP88 Yb   dP o.`Y8b     o.`Y8b  8I  dY 88"Yb
 *  YboodP  YbodP  8bodP' 88 YY 88  YbodP  8bodP'     8bodP' 8888Y"  88  Yb
 *
 *  Beacon-Kit BeginBlock:
 * |
 * |--- "Borrow" the logic of BeginBlock from
 *  https://github.com/cosmos/cosmos-sdk/blob/main/types/module/module.go.
 *
 * |    |
 * |    `--- "Appreciate" the key components and behaviors
 * |
 * `--- Implement in Solidity (with a cheeky grin)
 *      |
 *      |--- Create a struct to represent BeginBlocker
 *      |    |
 *      |    `--- Stuff it with necessary fields (e.g., contract address, selector)
 *      |
 *      |--- Implement CRUD operations for BeginBlocker
 *      |    |
 *      |    |--- Add: Sneak a BeginBlocker into a specific index, make others scooch over
 *      |    |
 *      |    |--- Remove: Yank a BeginBlocker from a specific index, others scooch back
 *      |    |
 *      |    `--- Update: Give a BeginBlocker a makeover at a specific index
 *      |
 *      `--- Implement the logic to call BeginBlocker
 *           |
 *           `--- For each BeginBlocker in the array, dial the contract and selector specified
 */
contract BeginBlockRootsContract is BeaconRootsContract {
    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                    EVENTS/ERRORS                           */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /**
     * @dev Emitted when the BeginBlocker function is called.
     * @param contractAddress The address of the contract that called the
     * BeginBlocker function.
     * @param coinbase The address of the current block miner.
     * @param selector The function selector that was called.
     * @param success The status of the call.
     */
    event BeginBlockerCalled(
        address indexed contractAddress,
        address indexed coinbase,
        bytes4 selector,
        bool success
    );

    /**
     * @dev Emitted when a BeginBlocker with the specified index does not exist.
     * @param i The index of the BeginBlocker.
     */
    error BeginBlockerDoesNotExist(uint256 i);

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                        CONSTANTS                           */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /// @dev Actions that can set a new BeginBlocker.
    bytes32 private constant SET = keccak256("SET");

    /// @dev Action that can remove BeginBlockers from the array.
    bytes32 private constant REMOVE = keccak256("REMOVE");

    /// @dev Action that can update the ADMIN address.
    bytes32 private constant UPDATE_ADMIN = keccak256("UPDATE_ADMIN");

    /// @dev The selector for "getBeginBlockers(uint256)".
    bytes4 private constant GET_BEGIN_BLOCKERS_SELECTOR =
        bytes4(keccak256("getBeginBlockers(uint256)"));

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                        STORAGE                           */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /// @dev The BeginBlocker struct is used to store the calls we need to make
    /// at the beginning of each block.
    struct BeginBlocker {
        address contractAddress;
        bytes4 selector;
    }

    /// @dev The ADMIN address is the only address that can add or remove
    /// BeginBlockers.
    address private ADMIN = address(0x20f33CE90A13a4b5E7697E3544c3083B8F8A51D4);

    /// @dev The list of BeginBlockers that we need to call at the beginning of
    /// each block, in order.
    BeginBlocker[] private beginBlockers;

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                        ENTRYPOINT                          */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /**
     * @dev Fallback function that is called when a non-function payload is sent
     * to the contract.
     * @dev This fallback adheres to the EIP-4788 specification with added
     * BeginBlocker functionality.
     *
     * The function behavior is as follows:
     *
     * Sender is not the system address:
     * |
     * |--- Sender is the admin:
     * |    |
     * |    `--- Perform CRUD operation with the message data.
     * |
     * |--- Sender is not the admin:
     *      |
     *      |--- Message data length is 36 and the first 4 bytes match the `GET_COINBASE_SELECTOR`:
     *      |    |
     *      |    `--- Call `getCoinbase()`.
     *      |
     *      `--- Message data length is not 36 or the first 4 bytes do not match the `GET_COINBASE_SELECTOR`:
     *           |
     *           `--- Call `get()`.
     *
     * Sender is the system address:
     * |
     * `--- Call `set()` and `run()`.
     */
    fallback() external override {
        if (msg.sender != SYSTEM_ADDRESS) {
            if (msg.sender == ADMIN) {
                crud();
            } else {
                if (
                    msg.data.length == 36
                        && bytes4(msg.data) == GET_COINBASE_SELECTOR
                ) {
                    getCoinbase();
                } else if (
                    msg.data.length == 36
                        && bytes4(msg.data) == GET_BEGIN_BLOCKERS_SELECTOR
                ) {
                    getBeginBlockers();
                } else {
                    get();
                }
            }
        } else {
            set();
            run();
        }
    }

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                        BeginBlocker                        */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /**
     * @dev Performs the CRUD operation based on the action specified in the
     * call data.
     * The call data contains the index, action, and BeginBlocker. The action
     * can be "set", "remove", or "update_admin".
     *
     * The function behavior is as follows:
     *
     * - If the action is "set", the function will add the BeginBlocker at the
     * given index.
     * - If the action is "remove", the function will remove the BeginBlocker at
     * the given index.
     * - If the action is "update_admin", the function will update the ADMIN
     * address if the new admin address is not zero.
     *
     *
     * Action received
     * |
     * |--- Action is "set":
     * |    |
     * |    `--- Add the BeginBlocker at the given index.
     * |
     * |--- Action is "remove":
     * |    |
     * |    `--- Remove the BeginBlocker at the given index.
     * |
     * `--- Action is "update_admin":
     *      |
     *      `--- If the new admin address is not zero, update the ADMIN address.
     */
    function crud() private {
        // Decode the data we get from the message.
        (
            uint256 i,
            bytes32 action,
            BeginBlocker memory beginBlocker,
            address admin
        ) = _parse();

        // Prefrom the CRUD operation.
        if (action == SET) {
            // If the action is "SET", we need to add the BeginBlocker at the
            // given index.
            _add(i, beginBlocker);
        } else if (action == REMOVE) {
            // If the action is "REMOVE", we need to remove the BeginBlocker at
            // the given index.
            _remove(i);
        } else if (action == UPDATE_ADMIN) {
            // If the action is "UPDATE_ADMIN", we need to update the ADMIN
            // address.
            // This can only be done by the current ADMIN, since we check that
            // in the fallback method.
            if (admin != address(0)) {
                ADMIN = admin;
            }
        }
    }

    /**
     * @dev Runs all the BeginBlocker functions stored in the contract.
     * It iterates over the array of BeginBlockers and calls each one using its
     * contract address and selector.
     * If the call is successful, it emits a BeginBlockerCalled event with the
     * contract address, the current block miner's address, the selector, and
     * the success status.
     */
    function run() private {
        uint256 length = beginBlockers.length;
        for (uint256 i; i < length;) {
            BeginBlocker storage beginBlocker = beginBlockers[i];
            address contractAddress = beginBlocker.contractAddress;
            bytes4 selector = beginBlocker.selector;
            bool success;
            assembly ("memory-safe") {
                mstore(0, selector)
                success := call(gas(), contractAddress, 0, 0, 4, 0, 0)
                i := add(i, 1)
            }
            emit BeginBlockerCalled(
                contractAddress, block.coinbase, selector, success
            );
        }
    }

    /**
     * @dev Parses the BeginBlocker message data.
     * @dev The call data must be encoded as: abi.encode(i, action,
     * contractAddress, selector, admin)
     * @return i The index of the BeginBlocker.
     * @return action The action to perform.
     * @return BeginBlocker The BeginBlocker struct containing the contract
     * address and the selector.
     * @return admin The new admin address.
     */
    function _parse()
        private
        pure
        returns (uint256, bytes32, BeginBlocker memory, address)
    {
        // Decode the data to get the BeginBlocker struct, user must send a
        // message that is encoded:
        // abi.encode(i, action, contractAddress, selector, admin)
        uint256 i;
        bytes32 action;
        address contractAddress;
        bytes4 selector;
        address admin;
        assembly ("memory-safe") {
            if iszero(eq(calldatasize(), 0xa0)) { revert(0, 0) }
            i := calldataload(0)
            action := calldataload(0x20)
            contractAddress := calldataload(0x40)
            selector := calldataload(0x60)
            admin := calldataload(0x80)
        }
        return (i, action, BeginBlocker(contractAddress, selector), admin);
    }

    /**
     * @dev Sets the BeginBlocker at the given index.
     * @param i The index of the BeginBlocker.
     * @param beginBlocker The BeginBlocker struct containing the contract
     * address and the selector.
     *
     * The function behavior is as follows:
     *
     * - If the index `i` is greater than the length of the `beginBlockers`
     * array, it reverts the transaction with a `BeginBlockerDoesNotExist`
     * error.
     * - If the index `i` is equal to the length of the `beginBlockers` array,
     * it adds the `beginBlocker` to the end of the array.
     * - If the index `i` is less than the length of the `beginBlockers` array,
     * it adds an empty element to the end of the array, then shifts all
     * elements from index `i` onwards one place to the right, and finally
     * places the `beginBlocker` at index `i`.
     *
     *
     * Index `i` received
     * |
     * |--- `i` is greater than the length of the `beginBlockers` array:
     * |    |
     * |    `--- Revert the transaction with a `BeginBlockerDoesNotExist` error.
     * |
     * |--- `i` is equal to the length of the `beginBlockers` array:
     * |    |
     * |    `--- Add the `beginBlocker` to the end of the array.
     * |
     * `--- `i` is less than the length of the `beginBlockers` array:
     *      |
     *      |--- Add an empty element to the end of the array.
     *      |
     *      |--- Shift all elements from index `i` onwards one place to the right.
     *      |
     *      `--- Place the `beginBlocker` at index `i`.
     */
    function _add(uint256 i, BeginBlocker memory beginBlocker) private {
        // cache the length of the array.
        uint256 length = beginBlockers.length;

        // Check that we are not trying to add a BeginBlocker at an index that
        // is greater than the length.
        if (i > length) {
            revert BeginBlockerDoesNotExist(i);
        }

        // push a new empty element at the end of the array since we are going
        // to fill it.
        beginBlockers.push();
        unchecked {
            for (uint256 j = length; j > i;) {
                beginBlockers[j] = beginBlockers[j - 1];
                --j;
            }
        }
        beginBlockers[i] = beginBlocker;
    }

    /**
     * @dev Removes the BeginBlocker at the given index.
     * @param i The index of the BeginBlocker.
     *
     * The function behavior is as follows:
     *
     * - If the index `i` is greater than or equal to the length of the
     * `beginBlockers` array, it
     * reverts the transaction with a `BeginBlockerDoesNotExist` error.
     * - If the index `i` is equal to the length of the `beginBlockers` array
     * minus 1 (i.e., it's the last element), it removes the last element from
     * the array.
     * - If the index `i` is less than the length of the `beginBlockers` array
     * minus 1, it shifts all elements from index `i+1` onwards one place to the
     * left, overwriting the element at index `i`, and then removes the last
     * element from the array.
     *
     * Index `i` received
     * |
     * |--- `i` is greater than or equal to the length of the `beginBlockers` array:
     * |    |
     * |    `--- Revert the transaction with a `BeginBlockerDoesNotExist` error.
     * |
     * |--- `i` is equal to the length of the `beginBlockers` array minus 1:
     * |    |
     * |    `--- Remove the last element from the array.
     * |
     * `--- `i` is less than the length of the `beginBlockers` array minus 1:
     *      |
     *      |--- Shift all elements from index `i+1` onwards one place to the left, overwriting the
     * element at index `i`.
     *      |
     *      `--- Remove the last element from the array.
     */
    function _remove(uint256 i) private {
        // Cache the length of the array.
        uint256 length = beginBlockers.length;

        if (i >= length) {
            revert BeginBlockerDoesNotExist(i);
        }

        unchecked {
            uint256 lastIndex = length - 1;
            for (uint256 j = i; j < lastIndex;) {
                beginBlockers[j] = beginBlockers[j + 1];
                ++j;
            }
            beginBlockers.pop();
        }
    }

    /**
     * @notice Reads the BeginBlockers array.
     * @dev The call data is encoded as: abi.encodePacked(GET_BEGIN_BLOCKERS_SELECTOR, index)
     * @return A BeginBlocker struct.
     */
    function getBeginBlockers() private view returns (BeginBlocker memory) {
        assembly ("memory-safe") {
            if iszero(lt(calldataload(4), sload(beginBlockers.slot))) {
                revert(0, 0)
            }
            mstore(0, beginBlockers.slot)
            let arraySlot := keccak256(0, 0x20)
            let data := sload(add(calldataload(4), arraySlot))
            let addr := shr(96, shl(96, data))
            let selector := shl(224, shr(160, data))
            mstore(0, addr)
            mstore(0x20, selector)
            // slither-disable-next-line incorrect-return
            return(0, 0x40)
        }
    }
}
