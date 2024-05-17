// SPDX-License-Identifier: MIT
pragma solidity ^0.8.21;

import { SSZ } from "./SSZ.sol";
import { Verifier } from "./Verifier.sol";

/// @author [madlabman](https://github.com/madlabman/eip-4788-proof)
contract ValidatorVerifier is Verifier {
    uint64 internal constant VALIDATOR_REGISTRY_LIMIT = 1 << 40;

    /// @dev Generalized index of the first validator struct root in the
    /// registry.
    uint256 public immutable gIndex;

    event Accepted(uint64 indexed validatorIndex);

    constructor(uint256 _gIndex) {
        gIndex = _gIndex;
    }

    function proveValidator(
        bytes32[] calldata validatorProof,
        SSZ.Validator calldata validator,
        uint64 validatorIndex,
        uint64 ts
    )
        public
    {
        if (validatorIndex >= VALIDATOR_REGISTRY_LIMIT) {
            revert IndexOutOfRange();
        }

        uint256 gI = gIndex + validatorIndex;
        bytes32 validatorRoot = SSZ.validatorHashTreeRoot(validator);
        bytes32 blockRoot = getParentBlockRoot(ts);

        if (!SSZ.verifyProof(validatorProof, blockRoot, validatorRoot, gI)) {
            revert InvalidProof();
        }

        emit Accepted(validatorIndex);
    }
}
