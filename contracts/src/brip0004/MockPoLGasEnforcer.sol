// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

/// @title GasEnforcedPoLDistributor
/// @notice PoL distributor requiring exactly 29,999,646 gas for testing gas limits
/// @dev Reverts if called with incorrect gas amount
contract GasEnforcedPoLDistributor {
    function distributeFor(bytes calldata /*pubkey*/ ) public view {
        uint256 start_gas = gasleft();
        require(start_gas == 29_999_646, "Insufficient gas");
    }
}
