// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.0;

/**
 * @title PropertiesDescriptions
 * @author 0xScourgedev, Rappie
 * @notice Descriptions strings for the invariants
 */
abstract contract PropertiesDescriptions {
    ///////////////////////////////////////////////////////////////////////////////////////////////
    //                                       DESCRIPTIONS                                        //
    ///////////////////////////////////////////////////////////////////////////////////////////////

    string internal constant PROOF_01 =
        "PROOF-01: If the zeroValidatorPubkeyGIndex is different, the proof should never be valid";
    string internal constant PROOF_02 =
        "PROOF-02: If the proposerIndexGIndex is different, the proof should never be valid";
    string internal constant PROOF_03 =
        "PROOF-03: If the proof for verifyValidatorPubkey was not modified post-generation, then the proof should always be valid";
    string internal constant PROOF_04 =
        "PROOF-04: If the proof for verifyProposerIndex was not modified post-generation, then the proof should always be valid";
    string internal constant PROOF_05 =
        "PROOF-05: If the zeroValidatorPubkeyGIndex is the same and the proof for verifyValidatorPubkey was modified post-generation, then the proof should never be valid";
    string internal constant PROOF_06 =
        "PROOF-06: If the proposerIndexGIndex is the same and the proof for verifyProposerIndex was modified post-generation, then the proof should never be valid";

    string internal constant REV_01 =
        "REV-01: setZeroValidatorPubkeyGIndex never reverts";
    string internal constant REV_02 =
        "REV-02: setProposerIndexGIndex never reverts";
}
