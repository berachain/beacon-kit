// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import { BeaconRootsContract } from "../BeaconRootsContract.sol";

contract BeginBlockRootsContract is BeaconRootsContract {
    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                    EVENTS/ERRORS                           */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /**
     * @dev Emitted when the BeginBlocker function is called.
     * @param contractAddress The address of the contract that called the BeginBlocker function.
     * @param coinbase The address of the current block miner.
     * @param selector The function selector that was called.
     * @param success The status of the call.
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

    /// @dev Actions that can set a new BeginBlocker.
    bytes32 private constant SET = keccak256("SET");

    /// @dev Action that can remove BeginBlockers from the array.
    bytes32 private constant REMOVE = keccak256("REMOVE");

    /// @dev Action that can update the ADMIN address.
    bytes32 private constant UPDATE_ADMIN = keccak256("UPDATE_ADMIN");

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                        STORAGE                           */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /// @dev The BeginBlocker struct is used to store the calls we need to make at the beginning of
    /// each block.
    struct BeginBlocker {
        address contractAddress;
        bytes4 selector;
    }

    /// @dev The ADMIN address is the only address that can add or remove BeginBlockers.
    address private ADMIN = address(0x20f33CE90A13a4b5E7697E3544c3083B8F8A51D4);

    /// @dev The list of BeginBlockers that we need to call at the beginning of each block, in
    /// order.
    BeginBlocker[] public beginBlockers;

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                        ENTRYPOINT                          */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    fallback() external override {
        if (msg.sender != SYSTEM_ADDRESS) {
            if (msg.data.length == 36 && bytes4(msg.data) == GET_COINBASE_SELECTOR) {
                getCoinbase();
            } else {
                if (msg.sender == ADMIN) {
                    crud(msg.data);
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
     * @dev Performs the CRUD operation based on the action specified in the input data.
     * @param data The input data containing the index, action, and BeginBlocker. The action can be
     * "set" or "remove".
     * If the action is "set", the function will add the BeginBlocker at the given index.
     * If the action is "remove", the function will remove the BeginBlocker at the given index.
     */
    function crud(bytes memory data) private {
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
    function run() private {
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
}
