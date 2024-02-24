// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

/// @notice the interface for the distribution contract from the bera chain.
interface IDistributor {
    /**
     * @notice Distribute the rewards to the cutting board receivers.
     * @dev This is only callable by the caller.
     * @param valCoinBase The address of the coinbase that we are distributing the rewards to.
     */
    function distribute(address valCoinBase) external;
}
