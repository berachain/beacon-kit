// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "../../lib/solady/test/utils/SoladyTest.sol";
import "../../lib/solady/src/utils/FixedPointMathLib.sol";
import { BeaconDepositContract } from "./BeaconDepositContract.sol";

/// @title BeaconDepositContractTest
contract BeaconDepositContractTest is SoladyTest {
    address internal constant BEACON_DEPOSIT_ADDRESS =
        0x00000000219ab540356cBB839Cbe05303d7705Fa;
    bytes WITHDRAWAL_CREDENTIALS = abi.encodePacked(address(this));
    bytes VALIDATOR_PUBKEY = "validator_pubkey";

    BeaconDepositContract depositContract =
        BeaconDepositContract(BEACON_DEPOSIT_ADDRESS);
    uint256 internal snapshot;

    /// @dev Set up the test environment by deploying a new BeaconRootsContract.
    function setUp() public virtual {
        // etch the BeaconDepositContract to the BEACON_DEPOSIT_ADDRESS
        vm.etch(
            BEACON_DEPOSIT_ADDRESS,
            vm.getDeployedCode("BeaconDepositContract.sol")
        );
        // take a snapshot of the clean state
        snapshot = vm.snapshot();
    }

    /// @dev Tests the deposit functionality with a valid amount.
    function testDepositValidAmount() public {
        // revert to the snapshot to get a fresh storage
        vm.revertTo(snapshot);
        // Assuming Gwei to Ether conversion for simplicity
        uint64 amount = 32 gwei;

        // Expect the Deposit event to be emitted with correct parameters
        vm.expectEmit(true, true, true, true);
        emit BeaconDepositContract.Deposit(
            VALIDATOR_PUBKEY, WITHDRAWAL_CREDENTIALS, amount
        );

        // Call the deposit function
        depositContract.deposit{ value: amount }(VALIDATOR_PUBKEY, amount);
    }

    /// @dev Tests the deposit functionality with an amount below the minimum
    /// required.
    function testDepositBelowMinimum() public {
        // revert to the snapshot to get a fresh storage
        vm.revertTo(snapshot);
        uint64 amount = 0.9 gwei; // Below the 1 ether minimum

        // Expect the transaction to revert with the InsufficientDeposit error
        vm.expectRevert(BeaconDepositContract.InsufficientDeposit.selector);

        // Call the deposit function
        depositContract.deposit{ value: amount }(VALIDATOR_PUBKEY, amount);
    }

    /// @dev Tests the withdrawal functionality with a valid request.
    function testWithdrawalValidRequest() public {
        // revert to the snapshot to get a fresh storage
        vm.revertTo(snapshot);
        // Assuming Gwei to Ether conversion for simplicity
        uint64 amount = 32 gwei;

        // Expect the Withdrawal event to be emitted with correct parameters
        vm.expectEmit(true, true, true, true);
        emit BeaconDepositContract.Withdrawal(
            VALIDATOR_PUBKEY, WITHDRAWAL_CREDENTIALS, amount
        );

        // Setup and call the withdrawal function
        depositContract.withdraw(VALIDATOR_PUBKEY, amount);
    }
}
