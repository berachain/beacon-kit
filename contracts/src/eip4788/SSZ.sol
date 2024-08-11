// SPDX-License-Identifier: MIT
pragma solidity ^0.8.21;

/// @author [madlabman](https://github.com/madlabman/eip-4788-proof)
library SSZ {
    /// @dev SHA256 precompile address.
    uint8 internal constant SHA256 = 0x02;
    /// @dev Length of the validator pubkey in bytes.
    uint8 internal constant VALIDATOR_PUBKEY_LENGTH = 48;

    error BranchHasMissingItem();
    error BranchHasExtraItem();
    error InvalidValidatorPubkeyLength();

    function validatorPubkeyHashTreeRoot(
        bytes memory pubkey
    )
        internal
        view
        returns (bytes32 root)
    {
        if (pubkey.length != VALIDATOR_PUBKEY_LENGTH) {
            revert InvalidValidatorPubkeyLength();
        }

        assembly {
            // Call sha256 precompile with the pubkey pointer
            let result :=
                staticcall(gas(), SHA256, add(pubkey, 32), 0x40, 0x00, 0x20)
            // Precompile returns no data on OutOfGas error.
            if eq(result, 0) { revert(0, 0) }

            root := mload(0x00)
        }
    }

    function addressHashTreeRoot(
        address v
    )
        internal
        pure
        returns (bytes32 root)
    {
        return bytes32(bytes20(v));
    }

    function uint64HashTreeRoot(uint64 v) internal pure returns (bytes32) {
        v = ((v & 0xFF00FF00FF00FF00) >> 8) | ((v & 0x00FF00FF00FF00FF) << 8);
        v = ((v & 0xFFFF0000FFFF0000) >> 16) | ((v & 0x0000FFFF0000FFFF) << 16);
        v = (v >> 32) | (v << 32);
        return bytes32(uint256(v) << 192);
    }

    /// @notice Modified version of `verify` from `MerkleProofLib` to support
    /// generalized indices and sha256 precompile.
    /// @dev Returns whether `leaf` exists in the Merkle tree with `root`,
    /// given `proof`.
    function verifyProof(
        bytes32[] calldata proof,
        bytes32 root,
        bytes32 leaf,
        uint256 index
    )
        internal
        view
        returns (bool isValid)
    {
        /// @solidity memory-safe-assembly
        assembly {
            if proof.length {
                // Left shift by 5 is equivalent to multiplying by 0x20.
                let end := add(proof.offset, shl(5, proof.length))
                // Initialize `offset` to the offset of `proof` in the calldata.
                let offset := proof.offset
                // Iterate over proof elements to compute root hash.
                for { } 1 { } {
                    // Slot of `leaf` in scratch space.
                    // If the condition is true: 0x20, otherwise: 0x00.
                    let scratch := shl(5, and(index, 1))
                    index := shr(1, index)
                    if iszero(index) {
                        // revert BranchHasExtraItem()
                        mstore(0x00, 0x5849603f)
                        revert(0x1c, 0x04)
                    }
                    // Store elements to hash contiguously in scratch space.
                    // Scratch space is 64 bytes (0x00 - 0x3f) and both elements
                    // are 32 bytes.
                    mstore(scratch, leaf)
                    mstore(xor(scratch, 0x20), calldataload(offset))
                    // Call sha256 precompile
                    let result :=
                        staticcall(gas(), SHA256, 0x00, 0x40, 0x00, 0x20)

                    if eq(result, 0) { revert(0, 0) }

                    // Reuse `leaf` to store the hash to reduce stack operations.
                    leaf := mload(0x00)
                    offset := add(offset, 0x20)
                    if iszero(lt(offset, end)) { break }
                }
            }
            // index != 1
            if gt(sub(index, 1), 0) {
                // revert BranchHasMissingItem()
                mstore(0x00, 0x1b6661c3)
                revert(0x1c, 0x04)
            }
            isValid := eq(leaf, root)
        }
    }
}
