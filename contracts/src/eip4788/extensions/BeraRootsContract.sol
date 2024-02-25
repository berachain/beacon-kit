// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import { IDistributor } from "./IDistributor.sol";
import { BeaconRootsContract } from "../BeaconRootsContract.sol";

/// @title BeraRootsContract is an extension of the BeaconRootsContract that allows for the block
/// rewards
/// logic to be implemented and conducted block by block
/// @dev TODOS: Authenticate, test, and deploy this contract.
contract BeraRootsContract is BeaconRootsContract {
    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                        STORAGE                           */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /// @dev The distribution contract that we will be calling into every block on the `set()`
    /// method.
    IDistributor private __distributor;

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                       DISTRIBUTE                           */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    function setDistributor(address distributor) external {
        __distributor = IDistributor(distributor);
    }

    /**
     * @notice Distributes the block rewards to the distributor contract.
     * @dev This function is called every time the beacon root is set.
     * @param coinBase The coinbase of the block.
     */
    function _distribute(address coinBase) private {
        // Only distribute if the distributor is set.
        if (address(__distributor) != address(0)) {
            __distributor.distribute(coinBase);
        }
    }

    /// @dev Overriding the set function to include the distribution of the rewards.
    function set() internal override {
        // if the distributor is set, distribute the rewards.
        if (address(__distributor) != address(0)) {
            _distribute(block.coinbase);
        }
        
        // Set the beacon root.
        super.set();
    }
}
