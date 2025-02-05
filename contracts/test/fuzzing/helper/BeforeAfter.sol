// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.0;

import "../FuzzSetup.sol";

/**
 * @title BeforeAfter
 * @author 0xScourgedev, Rappie
 * @notice Contains the states of the system before and after calls
 */
abstract contract BeforeAfter is FuzzSetup {
    ///////////////////////////////////////////////////////////////////////////////////////////////
    //                                         STRUCTS                                           //
    ///////////////////////////////////////////////////////////////////////////////////////////////

    struct State {
        uint256 zeroValidatorPubkeyGIndex;
        uint256 proposerIndexGIndex;
    }

    ///////////////////////////////////////////////////////////////////////////////////////////////
    //                                         VARIABLES                                         //
    ///////////////////////////////////////////////////////////////////////////////////////////////

    // callNum => State
    mapping(uint8 => State) states;

    ///////////////////////////////////////////////////////////////////////////////////////////////
    //                                         FUNCTIONS                                         //
    ///////////////////////////////////////////////////////////////////////////////////////////////

    function _before() internal {
        _setStates(0);

        if (DEBUG) debugBefore();
    }

    function _after() internal {
        _setStates(1);

        if (DEBUG) debugAfter();
    }

    function _setStates(uint8 callNum) internal {
        states[callNum].zeroValidatorPubkeyGIndex = BeaconVerifier(beaconVerifier).zeroValidatorPubkeyGIndex();
        states[callNum].proposerIndexGIndex = BeaconVerifier(beaconVerifier).proposerIndexGIndex();
    }

    function debugBefore() internal {
        debugState(0);
    }

    function debugAfter() internal {
        debugState(1);
    }

    function debugState(uint8 callNum) internal {
        fl.log("Call Number", callNum);
        fl.log("zeroValidatorPubkeyGIndex", states[callNum].zeroValidatorPubkeyGIndex);
        fl.log("proposerIndexGIndex", states[callNum].proposerIndexGIndex);
    }
}
