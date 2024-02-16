// SPDX-License-Identifier: MIT

pragma solidity >=0.8.10;

import { Errors } from "./Errors.sol";
import { IRootFollower } from "./IRootFollower.sol";

import { Ownable } from "@solady/auth/Ownable.sol";
import { FixedPointMathLib } from "@solady/utils/FixedPointMathLib.sol";

abstract contract RootFollower is IRootFollower, Ownable {
    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                        CONSTANTS                           */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /// @dev The length of the history buffer.
    uint256 private constant HISTORY_BUFFER_LENGTH = 256;
    /// @dev The selector for "getCoinbase(uint256)"
    bytes4 private constant GET_COINBASE_SELECTOR = 0xe8e284b9;
    /// @dev The selector for "BytesNotInBuffer()"
    bytes4 private constant BYTES_NOT_IN_BUFFER_SELECTOR = 0x68c0ab1c;
    /// @dev The beacon roots contract address.
    address private constant BEACON_ROOT_ADDRESS = 0x000F3df6D732807Ef1319fB7B8bB8522d0Beac02;

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                          STORAGE                           */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /// @dev The last block number that was processed.
    uint256 private _LAST_PROCESSED_BLOCK;

    constructor() {
        _initializeOwner(msg.sender);
    }

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                   PUBLIC READ FUNCTIONS                    */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /// @inheritdoc IRootFollower
    function getCoinbase(uint256 _block) public view returns (address coinbase) {
        return _getCoinbase(_block);
    }

    /// @inheritdoc IRootFollower
    function getNextActionableBlock() public view returns (uint256 blockNum) {
        unchecked {
            return FixedPointMathLib.max(
                _LAST_PROCESSED_BLOCK + 1,
                FixedPointMathLib.zeroFloorSub(block.number, HISTORY_BUFFER_LENGTH)
            );
        }
    }

    /// @inheritdoc IRootFollower
    function getLastActionedBlock() public view returns (uint256 blockNum) {
        return _LAST_PROCESSED_BLOCK;
    }

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                     ADMIN FUNCTIONS                        */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /// @inheritdoc IRootFollower
    function incrementBlock() public onlyOwner {
        _incrementBlock();
    }

    /// @inheritdoc IRootFollower
    function resetCount(uint256 _block) public onlyOwner {
        _resetCount(_block);
    }

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                     INTERNAL FUNCTIONS                     */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /**
     * @dev Fetches the coinbase address for a block using inline assembly & `staticcall`.
     * @param _block The block number to query.
     * @return _coinbase The miner's address for the block.
     * Assumes BEACON_ROOT_ADDRESS contract returns the coinbase. Reverts on failure.
     */
    function _getCoinbase(uint256 _block) internal view returns (address _coinbase) {
        assembly ("memory-safe") {
            mstore(0, GET_COINBASE_SELECTOR)
            mstore(0x04, _block)
            if iszero(staticcall(gas(), BEACON_ROOT_ADDRESS, 0, 0x24, 0, 0x20)) {
                mstore(0, BYTES_NOT_IN_BUFFER_SELECTOR)
                revert(0, 0x04)
            }
            _coinbase := mload(0)
        }
    }

    /**
     * @dev Increments `_LAST_PROCESSED_BLOCK` if it's the next actionable block.
     * Reverts with `ATTEMPTED_TO_INCREMENT_OUT_OF_BUFFER` if the next block isn't actionable.
     * Emits `AdvancedBlock` event after incrementing.
     */
    function _incrementBlock() internal {
        uint256 processingBlock;
        unchecked {
            processingBlock = _LAST_PROCESSED_BLOCK + 1;
        }
        // Check if next block is actionable, revert if not.
        if (processingBlock != getNextActionableBlock()) {
            revert Errors.AttemptedToIncrementOutOfBuffer();
        }
        // Increment and emit event.
        _LAST_PROCESSED_BLOCK = processingBlock;
        unchecked {
            emit AdvancedBlock(processingBlock - 1);
        }
    }

    /**
     * @dev Resets the count to a specified block number.
     * @param _block The block number to reset the count to.
     */
    function _resetCount(uint256 _block) internal {
        // Reverts if the block number is in the future.
        if (_block > block.number) {
            revert Errors.BlockDoesNotExist();
        }
        // Reverts if the block number is before the next actionable block.
        if (_block < getNextActionableBlock()) {
            revert Errors.BlockNotInBuffer();
        }

        // Emit an event to capture a block count reset.
        emit BlockCountReset(_block, _LAST_PROCESSED_BLOCK);

        // Sets the last processed block to the specified block number.
        _LAST_PROCESSED_BLOCK = _block;
    }
}
