// SPDX-License-Identifier: MIT
pragma solidity ^0.8.25;

/// @title RandaoTester
/// @dev This contract is used during the integration testing of
///      EIP-4399 in BeaconKit.
///      DO NOT USE THIS FOR GENERATING RANDOMNESS IN PRODUCTION
//       FOR ANYTHING IMPORTANT.
/// @author https://eips.ethereum.org/EIPS/eip-4399
/// @author itsdevbear@berachain.com
contract RandaoTester {
    /// @notice Stores the last retrieved RANDAO mix.
    uint256 public lastValue;

    /// @notice Retrieves and stores the previous RANDAO mix.
    /// @dev Accesses the `prevrandao` property from the block global variable.
    /// @return The last retrieved RANDAO mix.
    function storePrevRandao() external returns (uint256) {
        lastValue = block.prevrandao;
        return block.prevrandao;
    }
}
