// SPDX-License-Identifier: MIT

pragma solidity >=0.8.4;

/**
 * @title Errors Library
 * @dev Provides custom error definitions for the RootFollower contract operations.
 */
library Errors {
    /// @dev Unauthorized caller
    error Unauthorized(address);
    /// @dev The queried block is not in the buffer range
    error BlockNotInBuffer();
    /// @dev Increment was called with a block number no longer in the buffer range
    error AttemptedToIncrementOutOfBuffer();
    /// @dev The block number does not exist yet
    error BlockDoesNotExist();
}
