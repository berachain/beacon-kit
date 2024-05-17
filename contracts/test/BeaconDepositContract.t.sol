// SPDX-License-Identifier: MIT
pragma solidity ^0.8.25;

import { SoladyTest } from "@solady/test/utils/SoladyTest.sol";
import { IBeaconDepositContract } from
    "../src/staking/IBeaconDepositContract.sol";
import { SoladyTest } from "@solady/test/utils/SoladyTest.sol";
import { SSZ } from "@src/eip4788/SSZ.sol";
import { BeaconDepositContract } from "@src/staking/BeaconDepositContract.sol";

contract DepositContractTest is SoladyTest {
    /// @dev The depositor address.
    address internal depositor = 0x20f33CE90A13a4b5E7697E3544c3083B8F8A51D4;

    /// @dev The validator public key.
    bytes internal VALIDATOR_PUBKEY = _create48Byte();

    /// @dev The withdrawal credentials that we will use.
    bytes internal WITHDRAWAL_CREDENTIALS = _credential(address(this));

    /// @dev The staking credentials that are right.
    bytes internal SIGNATURE = _create96Byte();

    bytes32 internal DEPOSIT_DATA_ROOT;

    /// @dev the deposit contract.
    BeaconDepositContract internal depositContract;

    function setUp() public virtual {
        // Create the deposit contract.
        depositContract = new BeaconDepositContract();
        DEPOSIT_DATA_ROOT = this.computeDepositDataRoot(32 gwei, SIGNATURE);
    }

    function testFuzz_DepositsWrongPubKey(bytes calldata pubKey) public {
        vm.assume(pubKey.length != 96);
        vm.expectRevert(IBeaconDepositContract.InvalidPubKeyLength.selector);
        depositContract.deposit{ value: 32 ether }(
            bytes("wrong_public_key"),
            WITHDRAWAL_CREDENTIALS,
            SIGNATURE,
            DEPOSIT_DATA_ROOT
        );
    }

    function test_DepositWrongPubKey() public {
        vm.expectRevert(IBeaconDepositContract.InvalidPubKeyLength.selector);
        depositContract.deposit{ value: 32 ether }(
            bytes("wrong_public_key"),
            WITHDRAWAL_CREDENTIALS,
            SIGNATURE,
            DEPOSIT_DATA_ROOT
        );
    }

    function testFuzz_DepositWrongCredentials(bytes calldata credentials)
        public
    {
        vm.assume(credentials.length != 32);
        vm.expectRevert(
            IBeaconDepositContract.InvalidCredentialsLength.selector
        );
        depositContract.deposit{ value: 32 ether }(
            VALIDATOR_PUBKEY, credentials, SIGNATURE, DEPOSIT_DATA_ROOT
        );
    }

    function test_DepositWrongCredentials() public {
        vm.expectRevert(
            IBeaconDepositContract.InvalidCredentialsLength.selector
        );
        depositContract.deposit{ value: 32 ether }(
            VALIDATOR_PUBKEY,
            bytes("wrong_credentials"),
            SIGNATURE,
            DEPOSIT_DATA_ROOT
        );
    }

    function testFuzz_DepositWrongSignature(bytes calldata signature) public {
        vm.assume(signature.length != 96);
        vm.expectRevert(IBeaconDepositContract.InvalidSignatureLength.selector);
        depositContract.deposit{ value: 32 ether }(
            VALIDATOR_PUBKEY,
            WITHDRAWAL_CREDENTIALS,
            signature,
            DEPOSIT_DATA_ROOT
        );
    }

    function test_DepositWrongSignature() public {
        vm.expectRevert(IBeaconDepositContract.InvalidSignatureLength.selector);
        depositContract.deposit{ value: 32 ether }(
            VALIDATOR_PUBKEY,
            WITHDRAWAL_CREDENTIALS,
            bytes("wrong_signature"),
            DEPOSIT_DATA_ROOT
        );
    }

    function testFuzz_DepositWrongAmount(uint256 amount) public {
        amount = _bound(amount, 1, 32e9 - 1);
        vm.deal(depositor, amount * 1e9);
        vm.prank(depositor);
        vm.expectRevert(IBeaconDepositContract.DepositValueTooLow.selector);
        depositContract.deposit{ value: (amount * 1 gwei) }(
            VALIDATOR_PUBKEY,
            WITHDRAWAL_CREDENTIALS,
            SIGNATURE,
            DEPOSIT_DATA_ROOT
        );
    }

    function test_DepositWrongAmount() public {
        vm.deal(depositor, (32e9 - 1) * 1e9);
        vm.prank(depositor);
        vm.expectRevert(IBeaconDepositContract.DepositValueTooLow.selector);
        depositContract.deposit{ value: (32e9 - 1) * 1e9 }(
            VALIDATOR_PUBKEY,
            WITHDRAWAL_CREDENTIALS,
            SIGNATURE,
            DEPOSIT_DATA_ROOT
        );
    }

    function test_DepositFailsWithMaxAmount() public {
        vm.deal(depositor, uint256(type(uint64).max) * 2 gwei);
        vm.prank(depositor);
        vm.expectRevert(IBeaconDepositContract.DepositValueTooHigh.selector);
        depositContract.deposit{ value: uint256(type(uint64).max) * 2 gwei }(
            VALIDATOR_PUBKEY, WITHDRAWAL_CREDENTIALS, SIGNATURE, bytes32("")
        );
    }

    function test_DepositFailsWithInvalidDepositRoot() public {
        vm.deal(depositor, 32 ether);
        vm.prank(depositor);
        vm.expectRevert(IBeaconDepositContract.InvalidDepositDataRoot.selector);
        depositContract.deposit{ value: 32 ether }(
            VALIDATOR_PUBKEY, WITHDRAWAL_CREDENTIALS, SIGNATURE, bytes32("")
        );
    }

    function test_Deposit() public {
        vm.deal(depositor, 32 ether);
        vm.prank(depositor);
        vm.expectEmit(true, true, true, true);
        emit IBeaconDepositContract.Deposit(
            VALIDATOR_PUBKEY, WITHDRAWAL_CREDENTIALS, 32e9, SIGNATURE, 0
        );
        depositContract.deposit{ value: 32 ether }(
            VALIDATOR_PUBKEY,
            WITHDRAWAL_CREDENTIALS,
            SIGNATURE,
            DEPOSIT_DATA_ROOT
        );
    }

    function testFuzz_DepositAmountNotDivisibleByGwei(uint256 amount) public {
        amount = _bound(amount, 31e9 + 1, uint256(type(uint64).max));
        vm.assume(amount % 1e9 != 0);
        vm.deal(depositor, amount);
        vm.prank(depositor);
        vm.expectRevert(
            IBeaconDepositContract.DepositNotMultipleOfGwei.selector
        );
        depositContract.deposit{ value: amount }(
            VALIDATOR_PUBKEY,
            WITHDRAWAL_CREDENTIALS,
            SIGNATURE,
            DEPOSIT_DATA_ROOT
        );
    }

    function test_DepositAmountNotDivisibleByGwei() public {
        uint256 amount = 32e9 + 1;
        vm.deal(depositor, amount);
        vm.prank(depositor);
        vm.expectRevert(
            IBeaconDepositContract.DepositNotMultipleOfGwei.selector
        );
        depositContract.deposit{ value: amount }(
            VALIDATOR_PUBKEY,
            WITHDRAWAL_CREDENTIALS,
            SIGNATURE,
            DEPOSIT_DATA_ROOT
        );

        amount = 32e9 - 1;
        vm.deal(depositor, amount);
        vm.prank(depositor);
        vm.expectRevert(
            IBeaconDepositContract.DepositNotMultipleOfGwei.selector
        );
        depositContract.deposit{ value: amount }(
            VALIDATOR_PUBKEY,
            WITHDRAWAL_CREDENTIALS,
            SIGNATURE,
            DEPOSIT_DATA_ROOT
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
                WITHDRAWAL_CREDENTIALS,
                32e9,
                SIGNATURE,
                depositCount
            );
            depositContract.deposit{ value: 32 ether }(
                VALIDATOR_PUBKEY,
                WITHDRAWAL_CREDENTIALS,
                SIGNATURE,
                DEPOSIT_DATA_ROOT
            );
            ++depositCount;
        }
        assertEq(depositContract.depositCount(), depositCount);
        assertEq(address(depositContract).balance, 0);
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

    function computeDepositDataRoot(
        uint64 amountInGwei,
        bytes calldata signature
    )
        public
        view
        returns (bytes32)
    {
        bytes32 amount = SSZ.toLittleEndian(amountInGwei);
        // Compute deposit data root (`DepositData` hash tree root)
        bytes32 pubkey_root =
            sha256(abi.encodePacked(VALIDATOR_PUBKEY, bytes16(0)));
        bytes32 signature_root = sha256(
            abi.encodePacked(
                sha256(signature[:64]),
                sha256(abi.encodePacked(signature[64:], bytes32(0)))
            )
        );
        bytes32 node = sha256(
            abi.encodePacked(
                sha256(abi.encodePacked(pubkey_root, WITHDRAWAL_CREDENTIALS)),
                sha256(abi.encodePacked(amount, signature_root))
            )
        );
        return node;
    }
}
