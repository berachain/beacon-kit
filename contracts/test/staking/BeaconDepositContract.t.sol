// SPDX-License-Identifier: MIT
pragma solidity ^0.8.25;

import { SoladyTest } from "@solady/test/utils/SoladyTest.sol";
import { IBeaconDepositContract } from "@src/staking/IBeaconDepositContract.sol";
import { SoladyTest } from "@solady/test/utils/SoladyTest.sol";
import { BeaconDepositContract } from "@src/staking/BeaconDepositContract.sol";

contract DepositContractTest is SoladyTest {
    /// @dev The depositor address.
    address internal depositor = 0x20f33CE90A13a4b5E7697E3544c3083B8F8A51D4;

    /// @dev The validator public key.
    bytes internal VALIDATOR_PUBKEY = _create48Byte();

    /// @dev The withdrawal credentials that we will use.
    bytes internal WITHDRAWAL_CREDENTIALS = _credential(address(this));

    /// @dev The staking credentials that are right.
    bytes internal STAKING_CREDENTIALS = _credential(depositor);

    /// @dev the deposit contract address.
    address internal constant DEPOSIT_CONTRACT_ADDRESS =
        0x4242424242424242424242424242424242424242;

    bytes32 internal constant STAKING_ASSET_SLOT = bytes32(0);

    address internal gov = 0x8a73D1380345942F1cb32541F1b19C40D8e6C94B;

    /// @dev the deposit contract.
    BeaconDepositContract internal depositContract;

    function setUp() public virtual {
        // Set the STAKE_ASSET to the NATIVE token.
        depositContract = new BeaconDepositContract(gov);
        vm.prank(gov);
        depositContract.allowDeposit(depositor, 100);
    }

    function testFuzz_DepositsWrongPubKey(bytes calldata pubKey) public {
        vm.assume(pubKey.length != 96);
        vm.expectRevert(IBeaconDepositContract.InvalidPubKeyLength.selector);
        vm.prank(depositor);
        depositContract.deposit(
            bytes("wrong_public_key"), STAKING_CREDENTIALS, _create96Byte()
        );
    }

    function test_DepositWrongPubKey() public {
        vm.expectRevert(IBeaconDepositContract.InvalidPubKeyLength.selector);
        vm.prank(depositor);
        depositContract.deposit(
            bytes("wrong_public_key"), STAKING_CREDENTIALS, _create96Byte()
        );
    }

    function testFuzz_DepositWrongCredentials(bytes calldata credentials)
        public
    {
        vm.assume(credentials.length != 32);

        vm.expectRevert(
            IBeaconDepositContract.InvalidCredentialsLength.selector
        );
        vm.prank(depositor);
        depositContract.deposit(_create48Byte(), credentials, _create96Byte());
    }

    function test_DepositWrongCredentials() public {
        vm.expectRevert(
            IBeaconDepositContract.InvalidCredentialsLength.selector
        );
        vm.prank(depositor);
        depositContract.deposit(
            VALIDATOR_PUBKEY, bytes("wrong_credentials"), _create96Byte()
        );
    }

    function testFuzz_DepositWrongMinAmount(uint256 amountInEther) public {
        amountInEther = _bound(amountInEther, 1, 31);
        uint256 amountInGwei = amountInEther * 1 gwei;
        vm.deal(depositor, amountInGwei);
        vm.prank(depositor);
        vm.expectRevert(IBeaconDepositContract.InsufficientDeposit.selector);
        depositContract.deposit{ value: amountInGwei }(
            VALIDATOR_PUBKEY, STAKING_CREDENTIALS, _create96Byte()
        );
    }

    function test_DepositWrongMinAmount() public {
        uint256 amount = 31 gwei;
        vm.deal(depositor, amount);
        vm.prank(depositor);
        vm.expectRevert(IBeaconDepositContract.InsufficientDeposit.selector);
        depositContract.deposit{ value: amount }(
            VALIDATOR_PUBKEY, STAKING_CREDENTIALS, _create96Byte()
        );
    }

    function testFuzz_DepositNotDivisibleByGwei(uint256 amount) public {
        amount = _bound(amount, 31e9 + 1, uint256(type(uint64).max));
        vm.assume(amount % 1e9 != 0);
        vm.deal(depositor, amount);

        vm.prank(depositor);
        vm.expectRevert(
            IBeaconDepositContract.DepositNotMultipleOfGwei.selector
        );
        depositContract.deposit{ value: amount }(
            VALIDATOR_PUBKEY, STAKING_CREDENTIALS, _create96Byte()
        );
    }

    function test_DepositNotDivisibleByGwei() public {
        uint256 amount = 32e9 + 1;
        vm.deal(depositor, amount);
        vm.expectRevert(
            IBeaconDepositContract.DepositNotMultipleOfGwei.selector
        );
        vm.prank(depositor);
        depositContract.deposit{ value: amount }(
            VALIDATOR_PUBKEY, STAKING_CREDENTIALS, _create96Byte()
        );

        amount = 32e9 - 1;
        vm.deal(depositor, amount);
        vm.expectRevert(
            IBeaconDepositContract.DepositNotMultipleOfGwei.selector
        );
        vm.prank(depositor);
        depositContract.deposit{ value: amount }(
            VALIDATOR_PUBKEY, STAKING_CREDENTIALS, _create96Byte()
        );
    }

    function test_Deposit() public {
        vm.deal(depositor, 32 ether);
        vm.prank(depositor);
        vm.expectEmit(true, true, true, true);
        emit IBeaconDepositContract.Deposit(
            VALIDATOR_PUBKEY, STAKING_CREDENTIALS, 32 gwei, _create96Byte(), 0
        );
        depositContract.deposit{ value: 32 ether }(
            VALIDATOR_PUBKEY, STAKING_CREDENTIALS, _create96Byte()
        );
    }

    function testFuzz_DepositCount(uint256 count) public {
        count = _bound(count, 1, 100);
        vm.deal(depositor, 32 ether * count);
        vm.startPrank(depositor);
        uint64 depositCount;
        for (uint256 i; i < count; ++i) {
            vm.expectEmit(true, true, true, true);
            emit IBeaconDepositContract.Deposit(
                VALIDATOR_PUBKEY,
                STAKING_CREDENTIALS,
                32 gwei,
                _create96Byte(),
                depositCount
            );
            depositContract.deposit{ value: 32 ether }(
                VALIDATOR_PUBKEY, STAKING_CREDENTIALS, _create96Byte()
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
