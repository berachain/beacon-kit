// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

/// @title ValidatorRegistry
/// @notice Simple registry contract for testing multi-contract state changes
/// @dev Called by PoL distributors to increment activity counter
contract ValidatorRegistry {
    /// @notice Activity counter incremented on each call
    uint256 public callCount;

    /// @notice Event emitted when activity is recorded
    event RegistryCalled(uint256 newCount);

    /// @notice Records validator activity by incrementing counter
    function recordValidatorActivity(bytes calldata /* pubkey */ ) external {
        callCount++;
        emit RegistryCalled(callCount);
    }
}
