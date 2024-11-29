// SPDX-License-Identifier: MIT
pragma solidity ^0.8.25;

import "forge-std/Test.sol";

import { SoladyTest } from "@solady/test/utils/SoladyTest.sol";
import { IDepositContract } from "@src/staking/IDepositContract.sol";
import { PermissionedDepositContract } from "./PermissionedDepositContract.sol";

contract DepositContractTest is SoladyTest, StdCheats {
    /// @dev The depositor address.
    address internal depositor = 0x20f33CE90A13a4b5E7697E3544c3083B8F8A51D4;

    /// @dev The owner address.
    address owner = 0x6969696969696969696969696969696969696969;

    /// @dev The validator public key.
    bytes internal VALIDATOR_PUBKEY = _create48Byte();

    /// @dev The withdrawal credentials that we will use.
    bytes internal WITHDRAWAL_CREDENTIALS = _credential(address(this));

    /// @dev The staking credentials that are right.
    bytes internal STAKING_CREDENTIALS = _credential(depositor);

    bytes32 internal constant STAKING_ASSET_SLOT = bytes32(0);

    /// @dev the deposit contract.
    PermissionedDepositContract internal depositContract;

    function setUp() public virtual {
        depositContract = new PermissionedDepositContract(owner);
        vm.prank(owner);
        depositContract.allowDeposit(depositor, 100);
    }

    function testFuzz_DepositsWrongPubKey(bytes calldata pubKey) public {
        vm.assume(pubKey.length != 96);
        vm.expectRevert(IDepositContract.InvalidPubKeyLength.selector);
        vm.deal(depositor, 32 ether);
        vm.prank(depositor);
        depositContract.deposit{ value: 32 ether }(
            bytes("wrong_public_key"),
            STAKING_CREDENTIALS,
            _create96Byte(),
            depositor
        );
    }

    function test_DepositWrongPubKey() public {
        vm.expectRevert(IDepositContract.InvalidPubKeyLength.selector);
        vm.deal(depositor, 32 ether);
        vm.prank(depositor);
        depositContract.deposit{ value: 32 ether }(
            bytes("wrong_public_key"),
            STAKING_CREDENTIALS,
            _create96Byte(),
            depositor
        );
    }

    function testFuzz_DepositWrongCredentials(bytes calldata credentials)
        public
    {
        vm.assume(credentials.length != 32);

        vm.deal(depositor, 32 ether);
        vm.expectRevert(IDepositContract.InvalidCredentialsLength.selector);
        vm.prank(depositor);
        depositContract.deposit{ value: 32 ether }(
            _create48Byte(), credentials, _create96Byte(), depositor
        );
    }

    function test_DepositWrongCredentials() public {
        vm.expectRevert(IDepositContract.InvalidCredentialsLength.selector);
        vm.deal(depositor, 32 ether);
        vm.prank(depositor);
        depositContract.deposit{ value: 32 ether }(
            VALIDATOR_PUBKEY,
            bytes("wrong_credentials"),
            _create96Byte(),
            depositor
        );
    }

    function test_DepositWrongAmount() public {
        vm.expectRevert(IDepositContract.InsufficientDeposit.selector);
        vm.deal(depositor, 31 ether);
        vm.prank(depositor);
        depositContract.deposit{ value: 31 ether }(
            VALIDATOR_PUBKEY, STAKING_CREDENTIALS, _create96Byte(), depositor
        );
    }

    function testFuzz_DepositNativeWrongMinAmount(uint256 amountInEther)
        public
    {
        amountInEther = _bound(amountInEther, 1, 31);
        uint256 amountInGwei = amountInEther * 1 gwei;
        vm.deal(depositor, amountInGwei);
        vm.prank(depositor);
        vm.expectRevert(IDepositContract.InsufficientDeposit.selector);
        depositContract.deposit{ value: amountInGwei }(
            VALIDATOR_PUBKEY, STAKING_CREDENTIALS, _create96Byte(), depositor
        );
    }

    function test_DepositNativeWrongMinAmount() public {
        uint256 amount = 31 gwei;
        vm.deal(depositor, amount);
        vm.prank(depositor);
        vm.expectRevert(IDepositContract.InsufficientDeposit.selector);
        depositContract.deposit{ value: amount }(
            VALIDATOR_PUBKEY, STAKING_CREDENTIALS, _create96Byte(), depositor
        );
    }

    function testFuzz_DepositNativeNotDivisibleByGwei(uint256 amount) public {
        amount = _bound(amount, 31e9 + 1, uint256(type(uint64).max));
        vm.assume(amount % 1e9 != 0);
        vm.deal(depositor, amount);

        vm.prank(depositor);
        vm.expectRevert(IDepositContract.DepositNotMultipleOfGwei.selector);
        depositContract.deposit{ value: amount }(
            VALIDATOR_PUBKEY, STAKING_CREDENTIALS, _create96Byte(), depositor
        );
    }

    function test_DepositNativeNotDivisibleByGwei() public {
        uint256 amount = 32e9 + 1;
        vm.deal(depositor, amount);
        vm.expectRevert(IDepositContract.DepositNotMultipleOfGwei.selector);
        vm.prank(depositor);
        depositContract.deposit{ value: amount }(
            VALIDATOR_PUBKEY, STAKING_CREDENTIALS, _create96Byte(), depositor
        );

        amount = 32e9 - 1;
        vm.deal(depositor, amount);
        vm.expectRevert(IDepositContract.DepositNotMultipleOfGwei.selector);
        vm.prank(depositor);
        depositContract.deposit{ value: amount }(
            VALIDATOR_PUBKEY, STAKING_CREDENTIALS, _create96Byte(), depositor
        );
    }

    function test_DepositNative() public {
        vm.deal(depositor, 32 ether);
        vm.prank(depositor);
        vm.expectEmit(true, true, true, true);
        emit IDepositContract.Deposit(
            VALIDATOR_PUBKEY, STAKING_CREDENTIALS, 32 gwei, _create96Byte(), 0
        );
        depositContract.deposit{ value: 32 ether }(
            VALIDATOR_PUBKEY, STAKING_CREDENTIALS, _create96Byte(), depositor
        );
    }

    function testFuzz_DepositCount(uint256 count) public {
        count = _bound(count, 1, 100);
        uint64 depositCount;
        for (uint256 i; i < count; ++i) {
            depositor = makeAddr(vm.toString(i));
            vm.deal(depositor, 32 ether);

            vm.startPrank(owner);
            depositContract.allowDeposit(depositor, 1);
            vm.stopPrank();

            vm.startPrank(depositor);
            vm.expectEmit(true, true, true, true);
            emit IDepositContract.Deposit(
                _newPubkey(i),
                STAKING_CREDENTIALS,
                32 gwei,
                _create96Byte(),
                depositCount++
            );
            depositContract.deposit{ value: 32 ether }(
                _newPubkey(i), STAKING_CREDENTIALS, _create96Byte(), depositor
            );
            vm.stopPrank();
        }
        assertEq(depositContract.depositCount(), depositCount);
    }

    function _credential(address addr) internal pure returns (bytes memory) {
        return abi.encodePacked(bytes1(0x01), bytes11(0x0), addr);
    }

    function _create96Byte() internal pure returns (bytes memory) {
        return abi.encodePacked(bytes32("32"), bytes32("32"), bytes32("32"));
    }

    function _create48Byte() internal pure returns (bytes memory) {
        return abi.encodePacked(bytes32("32"), bytes16("16"));
    }

    function _newPubkey(uint256 i) internal pure returns (bytes memory) {
        return abi.encodePacked(bytes32(i), bytes16("16"));
    }
}
