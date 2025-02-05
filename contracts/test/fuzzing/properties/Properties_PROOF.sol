// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.0;

import "./PropertiesBase.sol";

/**
 * @title Properties_PROOF
 * @author 0xScourgedev, Rappie
 * @notice Contains all PROOF invariants
 */
abstract contract Properties_PROOF is PropertiesBase {
    ///////////////////////////////////////////////////////////////////////////////////////////////
    //                                       INVARIANTS                                          //
    ///////////////////////////////////////////////////////////////////////////////////////////////

    /**
     * @custom:invariant PROOF-01: If the zeroValidatorPubkeyGIndex is different, the proof should never be valid
     */
    function invariant_PROOF_01() internal {
        fl.eq(states[1].zeroValidatorPubkeyGIndex, ZERO_VALIDATOR_PUBKEY_G_INDEX, PROOF_01);
    }

    /**
     * @custom:invariant PROOF-02: If the proposerIndexGIndex is different, the proof should never be valid
     */
    function invariant_PROOF_02() internal {
        fl.eq(states[1].proposerIndexGIndex, PROPOSER_INDEX_G_INDEX, PROOF_02);
    }

    /**
     * @custom:invariant PROOF-03: If the proof for verifyValidatorPubkey was not modified post-generation,
     * then the proof should always be valid
     */
    function invariant_PROOF_03(bool proofMod) internal {
        fl.t(!proofMod, PROOF_03);
    }

    /**
     * @custom:invariant PROOF-04: If the proof for verifyProposerIndex was not modified post-generation,
     * then the proof should always be valid
     */
    function invariant_PROOF_04(bool proofMod) internal {
        fl.t(!proofMod, PROOF_04);
    }

    /**
     * @custom:invariant PROOF-05: If the zeroValidatorPubkeyGIndex is the same and the proof for
     * verifyValidatorPubkey was modified post-generation, then the proof should never be valid
     */
    function invariant_PROOF_05(bool proofMod) internal {
        if (states[0].zeroValidatorPubkeyGIndex == ZERO_VALIDATOR_PUBKEY_G_INDEX) {
            fl.t(proofMod, PROOF_05);
        }
    }

    /**
     * @custom:invariant PROOF-06: If the proposerIndexGIndex is the same and the proof for
     * verifyProposerIndex was modified post-generation, then the proof should never be valid
     */
    function invariant_PROOF_06(bool proofMod) internal {
        if (states[0].proposerIndexGIndex == PROPOSER_INDEX_G_INDEX) {
            fl.t(proofMod, PROOF_06);
        }
    }
}
