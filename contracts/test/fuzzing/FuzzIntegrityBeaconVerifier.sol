// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.0;

import "./helper/handlers/HandlerBeaconVerifier.sol";
import "./FuzzIntegrityBase.sol";

/**
 * @title FuzzBeaconVerifierIntegrity
 * @author 0xScourgedev, Rappie
 * @notice Checks for errors in the handlers for BeaconVerifier
 */
contract FuzzBeaconVerifierIntegrity is
    HandlerBeaconVerifier,
    FuzzIntegrityBase
{
    ///////////////////////////////////////////////////////////////////////////////////////////////
    //                                         INTEGRITY                                         //
    ///////////////////////////////////////////////////////////////////////////////////////////////

    /**
     * @notice Checks the integrity of handler_setProposerIndexGIndex
     */
    function fuzz_setProposerIndexGIndex(uint64 _proposerIndexGIndex)
        public
    {
        bytes memory callData = abi.encodeWithSelector(
            HandlerBeaconVerifier.handler_setProposerIndexGIndex.selector,
            _proposerIndexGIndex
        );

        (bool success, bytes4 errorSelector) = _testSelf(callData);
        if (!success) {
            bytes4[] memory allowedErrors = new bytes4[](1);
            allowedErrors[0] = ClampFail.selector;
            fl.errAllow(errorSelector, allowedErrors, "SELF-02");
        }
    }

    /**
     * @notice Checks the integrity of handler_setZeroValidatorPubkeyGIndex
     */
    function fuzz_setZeroValidatorPubkeyGIndex(
        uint64 _zeroValidatorPubkeyGIndex
    )
        public
    {
        bytes memory callData = abi.encodeWithSelector(
            HandlerBeaconVerifier.handler_setZeroValidatorPubkeyGIndex.selector,
            _zeroValidatorPubkeyGIndex
        );

        (bool success, bytes4 errorSelector) = _testSelf(callData);
        if (!success) {
            bytes4[] memory allowedErrors = new bytes4[](1);
            allowedErrors[0] = ClampFail.selector;
            fl.errAllow(errorSelector, allowedErrors, "SELF-03");
        }
    }

    /**
     * @notice Checks the integrity of fuzz_verifyValidatorPubkey
     */
    function fuzz_verifyValidatorPubkey(
        int64 seed,
        uint256 valLen,
        bool proofMod,
        uint16 shiftIndex
    )
        public
    {
        bytes memory callData = abi.encodeWithSelector(
            HandlerBeaconVerifier.handler_verifyValidatorPubkey.selector,
            seed,
            valLen,
            proofMod,
            shiftIndex
        );

        (bool success, bytes4 errorSelector) = _testSelf(callData);
        if (!success) {
            bytes4[] memory allowedErrors = new bytes4[](1);
            allowedErrors[0] = ClampFail.selector;
            fl.errAllow(errorSelector, allowedErrors, "SELF-04");
        }
    }

    /**
     * @notice Checks the integrity of handler_verifyProposerIndex
     */
    function fuzz_verifyProposerIndex (
        int64 seed,
        uint256 valLen,
        bool proofMod,
        uint16 shiftIndex
    )
        public
    {
        bytes memory callData = abi.encodeWithSelector(
            HandlerBeaconVerifier.handler_verifyProposerIndex.selector,
            seed,
            valLen,
            proofMod,
            shiftIndex
        );

        (bool success, bytes4 errorSelector) = _testSelf(callData);
        if (!success) {
            bytes4[] memory allowedErrors = new bytes4[](1);
            allowedErrors[0] = ClampFail.selector;
            fl.errAllow(errorSelector, allowedErrors, "SELF-05");
        }
    }
}
