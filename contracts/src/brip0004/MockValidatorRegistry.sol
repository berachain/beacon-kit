// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

/// @title ValidatorRegistry
/// @notice A simple registry contract that increments a counter when called
/// @dev This contract will be called by SimplePoLDistributor to test multi-contract state changes
contract ValidatorRegistry {
    /// @notice Simple counter that increments on each call
    uint256 public callCount;
    
    /// @notice Event emitted when the contract is called
    event RegistryCalled(uint256 newCount);
    
    /// @notice Records activity - simply increments counter and emits event
    function recordValidatorActivity(bytes calldata /* pubkey */) external {
        callCount++;
        emit RegistryCalled(callCount);
    }
}