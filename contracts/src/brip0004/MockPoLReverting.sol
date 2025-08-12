// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;


import "./MockValidatorRegistry.sol";

/// @title RevertingPoLDistributor
/// @notice Mock PoL distributor that reverts after 10 distributions for testing
/// @dev Interacts with ValidatorRegistry to test multi-contract state changes
contract RevertingPoLDistributor {
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
    
    /// @dev Restricts access to system address (execution layer client only)
    modifier onlySystemCall() {
        if (msg.sender != SYSTEM_ADDRESS) {
            revert NotSystemAddress();
        }
        _;
    }
    
    /// @notice Main function called by execution client
    /// @param pubkey The validator public key
    /// @dev Calls ValidatorRegistry to test multi-contract state changes
    // slither-disable-next-line reentrancy-events
    function distributeFor(bytes calldata pubkey) external onlySystemCall {
        require(totalDistributions < 10, "Max distributions reached");
        totalDistributions++;
        VALIDATOR_REGISTRY.recordValidatorActivity(pubkey);
        
        emit PoLDistributed(pubkey);
    }
}

