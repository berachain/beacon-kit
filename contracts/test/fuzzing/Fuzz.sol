// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.0;

import "./FuzzIntegrityBeaconVerifier.sol";

/**
 * @title Fuzz
 * @author 0xScourgedev, Rappie
 * @notice Composite contract for all of the handlers
 */
contract Fuzz is FuzzBeaconVerifierIntegrity {
    constructor() payable {
        setup();
    }
}
