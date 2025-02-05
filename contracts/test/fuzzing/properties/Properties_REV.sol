// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.0;

import "./PropertiesBase.sol";

/**
 * @title Properties_REV
 * @author 0xScourgedev, Rappie
 * @notice Contains all REV invariants
 */
abstract contract Properties_REV is PropertiesBase {
    ///////////////////////////////////////////////////////////////////////////////////////////////
    //                                       INVARIANTS                                          //
    ///////////////////////////////////////////////////////////////////////////////////////////////

    /**
     * @custom:invariant REV-01: setZeroValidatorPubkeyGIndex never reverts
     */
    function invariant_REV_01(bytes4 errorSelector) internal {
        bytes4[] memory allowedErrors = new bytes4[](0);
        fl.errAllow(errorSelector, allowedErrors, REV_01);
    }

    /**
     * @custom:invariant REV-02: setProposerIndexGIndex never reverts
     */
    function invariant_REV_02(bytes4 errorSelector) internal {
        bytes4[] memory allowedErrors = new bytes4[](0);
        fl.errAllow(errorSelector, allowedErrors, REV_02);
    }
}
