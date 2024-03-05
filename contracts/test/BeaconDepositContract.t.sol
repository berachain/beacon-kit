// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import { SoladyTest } from "@solady/test/utils/SoladyTest.sol";
import { BeaconDepositContract } from "@src/staking/BeaconDepositContract.sol";

/// @title BeaconDepositContractTest
contract BeaconDepositContractTest is SoladyTest {
    address internal constant BEACON_DEPOSIT_ADDRESS =
        0x00000000219ab540356cBB839Cbe05303d7705Fa;
    bytes internal WITHDRAWAL_CREDENTIALS = abi.encodePacked(address(this));
    bytes internal VALIDATOR_PUBKEY = "validator_pubkey";

    BeaconDepositContract internal depositContract =
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

<<<<<<< HEAD
        // Setup and call the withdrawal function
        depositContract.withdraw(VALIDATOR_PUBKEY, amount);
=======
    function testFuzz_WithdrawWrongPubKey(bytes calldata pubKey) public {
        vm.revertTo(snapshot);
        vm.assume(pubKey.length != 48);
        vm.expectRevert(IBeaconDepositContract.InvalidPubKeyLength.selector);
        depositContract.withdraw(pubKey, WITHDRAWAL_CREDENTIALS, 32e9);
    }

    function test_WithdrawWrongPubKey() public {
        vm.revertTo(snapshot);
        vm.expectRevert(IBeaconDepositContract.InvalidPubKeyLength.selector);
        depositContract.withdraw(
            bytes("wrong_pub_key"), WITHDRAWAL_CREDENTIALS, 32e9
        );
    }

    function testFuzz_WithdrawWrongWithdrawalCredentials(
        bytes calldata withdrawalCredentials
    )
        public
    {
        vm.revertTo(snapshot);
        vm.assume(withdrawalCredentials.length != 32);
        vm.expectRevert(
            IBeaconDepositContract.InvalidCredentialsLength.selector
        );
        depositContract.withdraw(VALIDATOR_PUBKEY, withdrawalCredentials, 32e9);
    }

    function test_WithdrawWrongWithdrawCredentials() public {
        vm.revertTo(snapshot);
        vm.expectRevert(
            IBeaconDepositContract.InvalidCredentialsLength.selector
        );
        depositContract.withdraw(
            VALIDATOR_PUBKEY, bytes("wrong_credentials"), 32e9
        );
    }

    function testFuzz_WithdrawWrongAmount(uint256 amount) public {
        vm.revertTo(snapshot);
        amount = _bound(amount, 1, 32e9 / 10 - 1);

        vm.expectRevert(
            IBeaconDepositContract.InsufficientWithdrawAmount.selector
        );
        depositContract.withdraw(
            VALIDATOR_PUBKEY, WITHDRAWAL_CREDENTIALS, uint64(amount)
        );
    }

    function test_WithdrawWrongAmount() public {
        vm.revertTo(snapshot);
        vm.expectRevert(
            IBeaconDepositContract.InsufficientWithdrawAmount.selector
        );
        depositContract.withdraw(
            VALIDATOR_PUBKEY, WITHDRAWAL_CREDENTIALS, 32e9 / 10 - 1
        );
    }

    function testWithdraw() public {
        vm.revertTo(snapshot);
        vm.expectEmit(true, true, true, true);

        vm.prank(depositor);
        emit IBeaconDepositContract.Withdrawal(
            VALIDATOR_PUBKEY, WITHDRAWAL_CREDENTIALS, 32e9
        );
        depositContract.withdraw(VALIDATOR_PUBKEY, WITHDRAWAL_CREDENTIALS, 32e9);
    }

    function testFuzz_DepositNativeWrongMinAmount(uint256 amount) public {
        vm.revertTo(nativeSnapshot);
        amount = _bound(amount, 1, 32 gwei - 1);
        vm.deal(depositor, amount);
        vm.expectRevert(IBeaconDepositContract.InsufficientDeposit.selector);
        depositContract.deposit{ value: amount }(
            VALIDATOR_PUBKEY, STAKING_CREDENTIALS, 0, _create96Byte()
        );
    }

    function test_DepositNativeWrongMinAmount() public {
        vm.revertTo(nativeSnapshot);
        uint256 amount = 32 gwei - 1;
        vm.deal(depositor, amount);
        vm.expectRevert(IBeaconDepositContract.InsufficientDeposit.selector);
        depositContract.deposit{ value: amount }(
            VALIDATOR_PUBKEY, STAKING_CREDENTIALS, 0, _create96Byte()
        );
    }

    function testFuzz_DepositNativeNotDivisibleByGwei(uint256 amount) public {
        vm.revertTo(nativeSnapshot);
        amount = _bound(amount, 32e9 + 1, uint256(type(uint64).max));
        vm.assume(amount % 1e9 != 0);
        vm.deal(depositor, amount);

        vm.prank(depositor);
        vm.expectRevert(
            IBeaconDepositContract.DepositNotMultipleOfGwei.selector
        );
        depositContract.deposit{ value: amount }(
            VALIDATOR_PUBKEY, STAKING_CREDENTIALS, 0, _create96Byte()
        );
    }

    function test_DepositNativeNotDivisibleByGwei() public {
        vm.revertTo(nativeSnapshot);
        uint256 amount = 32e9 + 1;
        vm.deal(depositor, amount);
        vm.expectRevert(
            IBeaconDepositContract.DepositNotMultipleOfGwei.selector
        );
        vm.prank(depositor);
        depositContract.deposit{ value: amount }(
            VALIDATOR_PUBKEY, STAKING_CREDENTIALS, 0, _create96Byte()
        );
    }

    function test_DepositNative() public {
        vm.revertTo(nativeSnapshot);
        vm.deal(depositor, 1 ether);
        vm.expectEmit(true, true, true, true);
        emit IBeaconDepositContract.Deposit(
            VALIDATOR_PUBKEY, STAKING_CREDENTIALS, 32 gwei, _create96Byte()
        );
        depositContract.deposit{ value: 32 gwei }(
            VALIDATOR_PUBKEY, STAKING_CREDENTIALS, 0, _create96Byte()
        );
    }

    function _credential(address addr) internal pure returns (bytes memory) {
        return abi.encodePacked(bytes1(0x01), bytes11(0x0), addr);
    }

    function _create96Byte() internal pure returns (bytes memory) {
        return abi.encodePacked(bytes32("32"), bytes32("32"), bytes32("32"));
    }

    function _create48Byte() internal pure returns (bytes memory) {
        return abi.encodePacked(bytes32("32"), bytes16("16"));
>>>>>>> 9d0aa54e (feat(types): fix types for kzg n stuff)
    }
}
