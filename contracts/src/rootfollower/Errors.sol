// SPDX-License-Identifier: MIT

pragma solidity >=0.8.4;

/**
 * @title Errors Library
 * @dev Provides custom error definitions for the RootFollower contract operations.
 */
library Errors {
    /// @dev Unauthorized caller
    error UNAUTHORIZED(address);
    /// @dev The queried block is not in the buffer range
    error BLOCK_NOT_IN_BUFFER(uint256);
    /// @dev Increment was called with a block number no longer in the buffer range
    error ATTEMPTED_TO_INCREMENT_OUT_OF_BUFFER();
    /// @dev The block number does not exist yet
    error BLOCK_DOES_NOT_EXIST();
}
