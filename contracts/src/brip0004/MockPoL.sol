// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import "./MockValidatorRegistry.sol";

/// @title SimplePoLDistributor
/// @notice Simple mock PoL distributor for testing bera-reth
/// @dev Updated to interact with ValidatorRegistry to test multi-contract state changes
contract SimplePoLDistributor {
    /// @notice System address that can call distributeFor (execution layer client)
    address private constant SYSTEM_ADDRESS = 0xffffFFFfFFffffffffffffffFfFFFfffFFFfFFfE;
    
    /// @notice The validator registry contract address (hardcoded for genesis deployment)
    ValidatorRegistry private constant VALIDATOR_REGISTRY = ValidatorRegistry(0x4200000000000000000000000000000000000043);
    
    /// @notice Event emitted when distributeFor is called
    event PoLDistributed(bytes pubkey);
    
    /// @notice Counter for total distributions
    uint256 public totalDistributions;
    
    /// @notice Error thrown when caller is not the system address
    error NotSystemAddress();
    
    /// @dev Modifier to restrict function access to system address.
    /// @dev This ensures only the execution layer client can call `distributeFor` function.
    modifier onlySystemCall() {
        if (msg.sender != SYSTEM_ADDRESS) {
            revert NotSystemAddress();
        }
        _;
    }
    
    /// @notice Main function that bera-reth will call
    /// @param pubkey The validator public key
    /// @dev Now calls the ValidatorRegistry to test multi-contract state changes
    function distributeFor(bytes calldata pubkey) external onlySystemCall {
        require(totalDistributions < 10, "Max distributions reached");
        // Update state in this contract
        totalDistributions++;
        
        // Call another contract to update its state
        // This tests whether system calls properly capture state changes from multiple contracts
        VALIDATOR_REGISTRY.recordValidatorActivity(pubkey);
        
        emit PoLDistributed(pubkey);
    }
}
