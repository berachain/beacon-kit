// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.0;

import "@perimetersec/fuzzlib/src/FuzzBase.sol";

import "./helper/FuzzStorageVariables.sol";

/**
 * @title FuzzSetup
 * @author 0xScourgedev, Rappie
 * @notice Setup for the fuzzing suite
 */
contract FuzzSetup is FuzzBase, FuzzStorageVariables {
    ///////////////////////////////////////////////////////////////////////////////////////////////
    //                                    SETUP CONTRACTS                                        //
    ///////////////////////////////////////////////////////////////////////////////////////////////

    function setup() internal {
        beaconVerifier = address(
            new BeaconVerifier(ZERO_VALIDATOR_PUBKEY_G_INDEX, PROPOSER_INDEX_G_INDEX)
        );
    }
}
