// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.0;

import "./FunctionCalls.sol";
import "../helper/BeforeAfter.sol";

/**
 * @title PreconditionsBase
 * @author 0xScourgedev, Rappie
 * @notice Contains the base for all preconditions
 */
abstract contract PreconditionsBase is FunctionCalls, BeforeAfter {
    error ClampFail(string);
    error MultiCallFail(string);

    /**
     * @notice modifier to set the current actor to the sender
     */
    modifier setCurrentActor() {
        if (_setActor) {
            currentActor = address(uint160(msg.sender) + 1);
        }
        _;
    }
}
