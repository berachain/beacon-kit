// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.0;

import "./Properties_PROOF.sol";
import "./Properties_REV.sol";

/**
 * @title Properties
 * @author 0xScourgedev, Rappie
 * @notice Composite contract for all of the properties, and contains general invariants
 */
abstract contract Properties is Properties_PROOF, Properties_REV { }
