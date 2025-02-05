// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.0;

import "@perimetersec/fuzzlib/src/FuzzBase.sol";

import "../helper/FuzzStorageVariables.sol";

/**
 * @title FunctionCalls
 * @author 0xScourgedev, Rappie
 * @notice Contains the function calls for all of the handlers
 */
abstract contract FunctionCalls is FuzzBase, FuzzStorageVariables {
    ///////////////////////////////////////////////////////////////////////////////////////////////
    //                                           EVENTS                                          //
    ///////////////////////////////////////////////////////////////////////////////////////////////

    event SetExecutionFeeRecipientGIndexCall(uint256 _executionFeeRecipientGIndex);
    event SetProposerIndexGIndexCall(uint256 _proposerIndexGIndex);
    event SetZeroValidatorPubkeyGIndexCall(uint256 _zeroValidatorPubkeyGIndex);
    event VerifyValidatorPubkeyCall(
        bytes32 beaconBlockRoot, uint64 proposerIndex, bytes proposerPubkey, bytes32[] proposerPubkeyProof
    );
    event VerifyProposerIndexCall(bytes32 beaconBlockRoot, uint64 proposerIndex, bytes32[] proposerPubkeyProof);

    ///////////////////////////////////////////////////////////////////////////////////////////////
    //                                       FUNCTIONS                                           //
    ///////////////////////////////////////////////////////////////////////////////////////////////

    function _setProposerIndexGIndexCall(uint256 _proposerIndexGIndex)
        internal
        returns (bool success, bytes memory returnData)
    {
        emit SetProposerIndexGIndexCall(_proposerIndexGIndex);

        (success, returnData) = beaconVerifier.call{ gas: 1_000_000 }(
            abi.encodeWithSelector(
                BeaconVerifier.setProposerIndexGIndex.selector,
                _proposerIndexGIndex
            )
        );
    }

    function _setZeroValidatorPubkeyGIndexCall(
        uint256 _zeroValidatorPubkeyGIndex
    )
        internal
        returns (bool success, bytes memory returnData)
    {
        emit SetZeroValidatorPubkeyGIndexCall(_zeroValidatorPubkeyGIndex);

        (success, returnData) = beaconVerifier.call{ gas: 1_000_000 }(
            abi.encodeWithSelector(
                BeaconVerifier.setZeroValidatorPubkeyGIndex.selector,
                _zeroValidatorPubkeyGIndex
            )
        );
    }

    function _verifyValidatorPubkeyCall(
        bytes32 beaconBlockRoot,
        bytes32[] memory proposerPubkeyProof,
        bytes memory proposerPubkey,
        uint64 proposerIndex
    ) internal returns (bool success, bytes memory returnData) {
        emit VerifyValidatorPubkeyCall(beaconBlockRoot, proposerIndex, proposerPubkey, proposerPubkeyProof);

        (success, returnData) = beaconVerifier.staticcall(
            abi.encodeWithSelector(
                BeaconVerifier.verifyValidatorPubkeyInBeaconBlock.selector,
                beaconBlockRoot,
                proposerPubkeyProof,
                proposerPubkey,
                proposerIndex
            )
        );
    }

    function _verifyProposerIndexCall(
        bytes32 beaconBlockRoot,
        uint64 proposerIndex,
        bytes32[] memory proposerPubkeyProof
    ) internal returns (bool success, bytes memory returnData) {
        emit VerifyProposerIndexCall(beaconBlockRoot, proposerIndex, proposerPubkeyProof);

        (success, returnData) = beaconVerifier.staticcall(
            abi.encodeWithSelector(
                BeaconVerifier.verifyProposerIndexInBeaconBlock.selector,
                beaconBlockRoot,
                proposerPubkeyProof,
                proposerIndex
            )
        );
    }
}
