// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.0;

/**
 * @title FuzzIntegrityBase
 * @author 0xScourgedev, Rappie
 * @notice Contains the base for all fuzz integrity contracts
 */
abstract contract FuzzIntegrityBase {
    /**
     * @notice Executes a delegatecall to the this contract with the given callData
     * @dev This function is used to call the handlers in order to test the integrity of the handlers
     * @param callData The data to be used in the delegatecall
     * @return the success of the delegatecall and the return data
     */
    function _testSelf(bytes memory callData) internal returns (bool, bytes4) {
        (bool success, bytes memory returnData) =
            address(this).delegatecall(callData);

        return (success, bytes4(returnData));
    }
}
