// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.0;

import "@perimetersec/fuzzlib/src/IHevm.sol";

/**
 * @title FuzzConstants
 * @author 0xScourgedev, Rappie
 * @notice Constants and assumptions for the fuzzing suite
 */
abstract contract FuzzConstants {
    ///////////////////////////////////////////////////////////////////////////////////////////////
    //                                         FUZZ CONFIGS                                      //
    ///////////////////////////////////////////////////////////////////////////////////////////////

    address internal constant USER1 = address(0x10001);
    address internal constant USER2 = address(0x20001);
    address internal constant USER3 = address(0x30001);
    address[] internal USERS = [USER1, USER2, USER3];

    ///////////////////////////////////////////////////////////////////////////////////////////////
    //                                         CONSTANTS                                         //
    ///////////////////////////////////////////////////////////////////////////////////////////////

    uint64 internal constant ZERO_VALIDATOR_PUBKEY_G_INDEX = 3_254_554_418_216_960;
    uint64 internal constant PROPOSER_INDEX_G_INDEX = 9;

    uint64 constant timestamp = 31_337;

    uint256 constant MAX_VALIDATORS_LENGTH = 100;
}
