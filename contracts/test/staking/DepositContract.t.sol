// SPDX-License-Identifier: MIT
pragma solidity ^0.8.25;

import { SoladyTest } from "@solady/test/utils/SoladyTest.sol";
import { IDepositContract } from "@src/staking/IDepositContract.sol";
import { PermissionedDepositContract } from "./PermissionedDepositContract.sol";

contract DepositContractTest is SoladyTest {
    /// @dev The depositor address.
    address internal constant DEPOSITOR = 0x20f33CE90A13a4b5E7697E3544c3083B8F8A51D4;

    /// @dev The owner address.
    address internal constant OWNER = 0x6969696969696969696969696969696969696969;

    /// @dev The validator public key.
    bytes internal VALIDATOR_PUBKEY = _create48Byte();

    /// @dev The withdrawal credentials that we will use.
    bytes internal WITHDRAWAL_CREDENTIALS = _credential(address(this));

    /// @dev The staking credentials that are right.
    bytes internal STAKING_CREDENTIALS = _credential(DEPOSITOR);

    bytes32 internal constant STAKING_ASSET_SLOT = bytes32(0); // No usage yet

    /// @dev the deposit contract.
    PermissionedDepositContract internal depositContract;

    function setUp() public virtual {
        depositContract = new PermissionedDepositContract(OWNER);
        vm.prank(OWNER);
        depositContract.allowDeposit(DEPOSITOR, 100);
    }

    function testFuzz_DepositsWrongPubKey(bytes memory pubKey) public {
        vm.assume(pubKey.length != 96);
        vm.expectRevert(IDepositContract.InvalidPubKeyLength.selector);
        vm.prank(DEPOSITOR);
        depositContract.deposit(
            bytes("wrong_public_key"),
            STAKING_CREDENTIALS,
            32e9,
            _create96Byte()
        );
    }

    function test_DepositWrongPubKey() public {
        vm.expectRevert(IDepositContract.InvalidPubKeyLength.selector);
        vm.prank(DEPOSITOR);
        depositContract.deposit(
            bytes("wrong_public_key"),
            STAKING_CREDENTIALS,
            32e9,
            _create96Byte()
        );
    }

    function testFuzz_DepositWrongCredentials(bytes memory credentials) public {
        vm.assume(credentials.length != 32);

        vm.expectRevert(IDepositContract.InvalidCredentialsLength.selector);
        vm.prank(DEPOSITOR);
        depositContract.deposit(
            _create48Byte(), credentials, 32e9, _create96Byte()
        );
    }

    function test_DepositWrongCredentials() public {
        vm.expectRevert(IDepositContract.InvalidCredentialsLength.selector);
        vm.prank(DEPOSITOR);
        depositContract.deposit(
            VALIDATOR_PUBKEY, bytes("wrong_credentials"), 32e9, _create96Byte()
        );
    }

    function testFuzz_DepositWrongAmount(uint256 amount) public {
        amount = _bound(amount, 1, 32e9 - 1);
        vm.deal(DEPOSITOR, amount);
        vm.prank(DEPOSITOR);
        vm.expectRevert(IDepositContract.InsufficientDeposit.selector);
        depositContract.deposit(
            VALIDATOR_PUBKEY,
            STAKING_CREDENTIALS,
            uint64(amount),
            _create96Byte()
        );
    }

    function test_DepositWrongAmount() public {
        vm.expectRevert(IDepositContract.InsufficientDeposit.selector);
        vm.prank(DEPOSITOR);
        depositContract.deposit(
            VALIDATOR_PUBKEY, STAKING_CREDENTIALS, 32e9 - 1, _create96Byte()
        );
    }

    function test_Deposit() public {
        vm.deal(DEPOSITOR, 32 ether);
        vm.prank(DEPOSITOR);
        vm.expectEmit(true, true, true, true);
        emit IDepositContract.Deposit(
            VALIDATOR_PUBKEY, STAKING_CREDENTIALS, 32e9, _create96Byte(), 0
        );
        depositContract.deposit{ value: 32 ether }(
            VALIDATOR_PUBKEY, STAKING_CREDENTIALS, 32e9, _create96Byte()
        );
    }

    function testFuzz_DepositNativeWrongMinAmount(uint256 amountInGwei) public {
        amountInGwei = _bound(amountInGwei, 1, 31);
        vm.deal(DEPOSITOR, amountInGwei);
        vm.prank(DEPOSITOR);
        vm.expectRevert(IDepositContract.InsufficientDeposit.selector);
        depositContract.deposit{ value: amountInGwei }(
            VALIDATOR_PUBKEY, STAKING_CREDENTIALS, 0, _create96Byte()
        );
    }

    function test_DepositNativeWrongMinAmount() public {
        uint256 amount = 31 gwei;
        vm.deal(DEPOSITOR, amount);
        vm.prank(DEPOSITOR);
        vm.expectRevert(IDepositContract.InsufficientDeposit.selector);
        depositContract.deposit{ value: amount }(
            VALIDATOR_PUBKEY, STAKING_CREDENTIALS, 0, _create96Byte()
        );
    }

    function testFuzz_DepositNativeNotDivisibleByGwei(uint256 amount) public {
        amount = _bound(amount, 31e9 + 1, uint256(type(uint64).max));
        vm.assume(amount % 1e9 != 0);
        vm.deal(DEPOSITOR, amount);

        vm.prank(DEPOSITOR);
        vm.expectRevert(IDepositContract.DepositNotMultipleOfGwei.selector);
        depositContract.deposit{ value: amount }(
            VALIDATOR_PUBKEY, STAKING_CREDENTIALS, 0, _create96Byte()
        );
    }

    function test_DepositNativeNotDivisibleByGwei() public {
        uint256 amount = 32e9 + 1;
        vm.deal(DEPOSITOR, amount);
        vm.expectRevert(IDepositContract.DepositNotMultipleOfGwei.selector);
        vm.prank(DEPOSITOR);
        depositContract.deposit{ value: amount }(
            VALIDATOR_PUBKEY, STAKING_CREDENTIALS, 0, _create96Byte()
        );

        amount = 32e9 - 1;
        vm.deal(DEPOSITOR, amount);
        vm.expectRevert(IDepositContract.DepositNotMultipleOfGwei.selector);
        vm.prank(DEPOSITOR);
        depositContract.deposit{ value: amount }(
            VALIDATOR_PUBKEY, STAKING_CREDENTIALS, 0, _create96Byte()
        );
    }

    function test_DepositNative() public {
        vm.deal(DEPOSITOR, 32 ether);
        vm.prank(DEPOSITOR);
        vm.expectEmit(true, true, true, true);
        emit IDepositContract.Deposit(
            VALIDATOR_PUBKEY, STAKING_CREDENTIALS, 32e9, _create96Byte(), 0
        );
        depositContract.deposit{ value: 32 ether }(
            VALIDATOR_PUBKEY, STAKING_CREDENTIALS, 32e9, _create96Byte()
        );
    }

    function testFuzz_DepositCount(uint256 count) public {
        count = _bound(count, 1, 100);
        vm.deal(DEPOSITOR, 32 ether * count);
        vm.startPrank(DEPOSITOR);
        uint64 depositCount;
        for (uint256 i; i < count; ++i) {
            vm.expectEmit(true, true, true, true);
            emit IDepositContract.Deposit(
                VALIDATOR_PUBKEY,
                STAKING_CREDENTIALS,
                32e9,
                _create96Byte(),
                depositCount
            );
            depositContract.deposit{ value: 32 ether }(
                VALIDATOR_PUBKEY, STAKING_CREDENTIALS, 32e9, _create96Byte()
            );
            ++depositCount;
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
}
