// SPDX-License-Identifier: MIT

pragma solidity ^0.8.0;


import { Errors } from "./Errors.sol";
import { RootFollower } from "./RootFollower.sol";
import { BeaconRootsContract } from "../eip4788/BeaconRootsContract.sol";
import { BeaconRootsContractBaseTest } from "../eip4788/BeaconRootsContract.t.sol";
import "forge-std/Test.sol";

contract RootFollowerUser is RootFollower { }

contract RootFollowerTest is BeaconRootsContractBaseTest {
    event AdvancedBlock(uint256 blockNum);

    RootFollowerUser c;
    address internal constant BEACON_ROOT_ADDRESS = 0x000F3df6D732807Ef1319fB7B8bB8522d0Beac02;

    function setUp() override public virtual {
        c = new RootFollowerUser();
        // vm.etch(BEACON_ROOT_ADDRESS, vm.getDeployedCode("BeaconRootsContract"));
        super.setUp();
    }

    function test_GetAndIncrementBlock() public {
        // Setup the block number
        uint256 blockNum = 1;
        vm.roll(blockNum);
        address expected = address(0);

        // Set the coinbase of block 1 to be address(0)
        vm.mockCall(
            BEACON_ROOT_ADDRESS,
            abi.encodeWithSelector(GET_COINBASE_SELECTOR, abi.encode(blockNum)),
            abi.encode(expected)
        );

        // Get the next actionable block and assert it
        assertEq(1, c.getNextActionableBlock());

        // Get the coinbase of block 1 and assert it
        address received = c.getCoinbase(blockNum);
        assertEq(expected, received);

        // Increment the block
        vm.expectEmit(address(c));
        emit AdvancedBlock(0);
        c.incrementBlock();

        // Check the last actioned block
        assertEq(1, c.getLastActionedBlock());
        assertEq(2, c.getNextActionableBlock());
    }

    function test_OutOfBuffer() public {
        // Advance the block num to a number out of the buffer
        // for (uint256 i = 0; i < 500; i++) {
        //     setStorage(, startTimestamp, length);
        // }
        
        c.getCoinbase(1);
        vm.roll(500);
        
        assertEq(500 - 256, c.getNextActionableBlock());

        // Incrementing the block should fail now because out of buffer
        vm.expectRevert(Errors.AttemptedToIncrementOutOfBuffer.selector);
        c.incrementBlock();

        // Getting an out of buffer coinbase should result in a revert
        vm.expectRevert();
        c.getCoinbase(1);
    }
}
