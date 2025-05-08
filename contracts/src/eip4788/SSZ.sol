// SPDX-License-Identifier: MIT
pragma solidity ^0.8.21;

/// @notice Library for SSZ (Simple Serialize) proof verification.
/// @author [madlabman](https://github.com/madlabman/eip-4788-proof)
library SSZ {
    /// @dev SHA256 precompile address.
    uint8 internal constant SHA256 = 0x02;

    error BranchHasMissingItem();
    error BranchHasExtraItem();

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

/// @notice Contract for testing SSZ (Simple Serialize) proof verification with the `SSZ` library.
/// @author Inspired by [madlabman](https://github.com/madlabman/eip-4788-proof).
contract SSZTest {
    /// @notice The address of the EIP-4788 Beacon Roots contract.
    address public constant BEACON_ROOTS =
        0x000F3df6D732807Ef1319fB7B8bB8522d0Beac02;

    // Signature: 0x3033b0ff
    error RootNotFound();

    /// @notice Verifies a proof of inclusion for a given leaf in a Merkle tree.
    /// @param proof The proof of inclusion.
    /// @param root The root of the Merkle tree.
    /// @param leaf The leaf to verify.
    /// @param index The index of the leaf in the Merkle tree.
    /// @return isValid Whether the proof is valid.
    function verifyProof(
        bytes32[] calldata proof,
        bytes32 root,
        bytes32 leaf,
        uint256 index
    )
        external
        view
        returns (bool isValid)
    {
        isValid = SSZ.verifyProof(proof, root, leaf, index);
    }

    /// @notice Verifies a proof of inclusion for a given leaf in a Merkle tree. Reverts if the proof is invalid.
    /// @param proof The proof of inclusion.
    /// @param root The root of the Merkle tree.
    /// @param leaf The leaf to verify.
    /// @param index The index of the leaf in the Merkle tree.
    function mustVerifyProof(
        bytes32[] calldata proof,
        bytes32 root,
        bytes32 leaf,
        uint256 index
    )
        external
        view
    {
        if (!SSZ.verifyProof(proof, root, leaf, index)) {
            revert("Proof is invalid");
        }
    }

    /// @notice Get the parent block root at a given timestamp.
    /// @dev Reverts with `RootNotFound()` if the root is not found.
    function getParentBlockRootAt(uint64 ts)
        external
        view
        returns (bytes32 root)
    {
        assembly ("memory-safe") {
            mstore(0, ts)
            let success := staticcall(gas(), BEACON_ROOTS, 0, 0x20, 0, 0x20)
            if iszero(success) {
                mstore(0, 0x3033b0ff) // RootNotFound()
                revert(0x1c, 0x04)
            }
            root := mload(0)
        }
    }
}
