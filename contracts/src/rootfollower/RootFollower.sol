// SPDX-License-Identifier: MIT

pragma solidity >=0.8.10;

import {Errors} from "./Errors.sol";
import {IRootFollower} from "./IRootFollower.sol";
import {FixedPointMathLib} from "solady/src/utils/FixedPointMathLib.sol";
import {OwnableRoles} from "solady/src/auth/OwnableRoles.sol";

abstract contract RootFollower is IRootFollower, OwnableRoles {
    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                        CONSTANTS                           */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /// @dev The length of the history buffer.
    uint256 private constant HISTORY_BUFFER_LENGTH = 256;
    /// @dev The selector for "getCoinbase(uint256)"
    bytes4 private constant GET_COINBASE_SELECTOR = 0xe8e284b9;
    /// @dev The beacon roots contract address.
    address private constant BEACON_ROOT_ADDRESS =
        0x000F3df6D732807Ef1319fB7B8bB8522d0Beac02;

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                          STORAGE                           */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    uint256 private _LAST_PROCESSED_BLOCK;

    constructor() {
        _initializeOwner(msg.sender);
    }

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                   PUBLIC READ FUNCTIONS                    */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /// @inheritdoc IRootFollower
    function getCoinbase(
        uint256 _block
    ) public view returns (address coinbase) {
        return _getCoinbase(_block);
    }

    /// @inheritdoc IRootFollower
    function getNextActionableBlock() public view returns (uint256 blockNum) {
        return
            FixedPointMathLib.max(
                _LAST_PROCESSED_BLOCK + 1,
                FixedPointMathLib.zeroFloorSub(
                    block.number,
                    HISTORY_BUFFER_LENGTH
                )
            );
    }

    /// @inheritdoc IRootFollower
    function getLastActionedBlock() public view returns (uint256 blockNum) {
        return _LAST_PROCESSED_BLOCK;
    }

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                  PUBLIC UPDATE FUNCTIONS                   */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    function incrementBlock() public onlyOwner {
        _incrementBlock();
    }

    function resetCount(uint256 _block) public onlyOwner {
        _resetCount(_block);
    }

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                     INTERNAL FUNCTIONS                     */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    function _getCoinbase(
        uint256 _block
    ) internal view returns (address _coinbase) {
        assembly ("memory-safe") {
            mstore(0, GET_COINBASE_SELECTOR)
            mstore(0x04, _block)
            if iszero(
                staticcall(gas(), BEACON_ROOT_ADDRESS, 0, 0x24, 0, 0x20)
            ) {
                revert(0, 0)
            }
            _coinbase := mload(0)
        }
    }

    function _incrementBlock() internal {
        if ((_LAST_PROCESSED_BLOCK + 1) != getNextActionableBlock()) {
            revert Errors.ATTEMPTED_TO_INCREMENT_OUT_OF_BUFFER();
        }
        emit AdvancedBlock(++_LAST_PROCESSED_BLOCK);
    }

    function _resetCount(uint256 _block) internal {
        if (_block > block.number) {
            revert Errors.BLOCK_DOES_NOT_EXIST();
        }
        if (_block < getNextActionableBlock()) {
            revert Errors.BLOCK_NOT_IN_BUFFER(_block);
        }
        _LAST_PROCESSED_BLOCK = _block;
    }
}
