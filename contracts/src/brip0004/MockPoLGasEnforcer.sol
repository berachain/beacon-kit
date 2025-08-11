// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

/// @notice A contract that requires exactly 29_999_646 gas to execute the `distributeFor` function.
/// @dev This contract is used to test gas limits and ensure precise gas consumption.
/// @dev The function will revert if the gas used is not exactly 29_999_646,
contract GasEnforcedPoLDistributor {
        function distributeFor(bytes calldata /*pubkey*/) public view {
          uint256 start_gas = gasleft();
          require(start_gas == 29_999_646, "Insufficient gas");
    }
}
