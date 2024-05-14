// SPDX-License-Identifier: MIT
pragma solidity ^0.8.25;

import { SoladyTest } from "@solady/test/utils/SoladyTest.sol";
import { ERC20 } from "@solady/src/tokens/ERC20.sol";
import { IBeaconDepositContract } from
    "../src/staking/IBeaconDepositContract.sol";
import { SoladyTest } from "@solady/test/utils/SoladyTest.sol";
import { BeaconDepositContract } from "@src/staking/BeaconDepositContract.sol";

contract DepositContractTest is SoladyTest {
    /// @dev The depositor address.
    address internal depositor = 0x20f33CE90A13a4b5E7697E3544c3083B8F8A51D4;

    /// @dev The validator public key.
    bytes internal VALIDATOR_PUBKEY = _create48Byte();

    /// @dev The withdrawal credentials that we will use.
    bytes internal WITHDRAWAL_CREDENTIALS = _credential(address(this));

    bytes internal SIGNATURE = _create96Byte();

    bytes32 internal DEPOSIT_DATA_ROOT;

    /// @dev the deposit contract.
    BeaconDepositContract internal depositContract;

    function setUp() public virtual {
        // Set the STAKE_ASSET to the NATIVE token.
        depositContract = new BeaconDepositContract();
        DEPOSIT_DATA_ROOT = this.createDepositDataRoot(32 gwei, SIGNATURE);
    }

    function test_DepositFailsIfInvalidPubKeySize() public {
        vm.deal(depositor, 32 ether);
        vm.expectRevert(IBeaconDepositContract.InvalidPubKeyLength.selector);
        depositContract.deposit{ value: 32 ether }(
            abi.encode(bytes32("32")),
            WITHDRAWAL_CREDENTIALS,
            SIGNATURE,
            DEPOSIT_DATA_ROOT
        );
    }

    function test_DepositFailsIfInvalidSignatureSize() public {
        vm.deal(depositor, 32 ether);
        vm.expectRevert(IBeaconDepositContract.InvalidSignatureLength.selector);
        depositContract.deposit{ value: 32 ether }(
            VALIDATOR_PUBKEY,
            WITHDRAWAL_CREDENTIALS,
            abi.encode(bytes32("32")),
            DEPOSIT_DATA_ROOT
        );
    }

    function test_DepositFailsIfInvalidWithdrawalCredentialsSize() public {
        vm.deal(depositor, 32 ether);
        vm.expectRevert(
            IBeaconDepositContract.InvalidCredentialsLength.selector
        );
        depositContract.deposit{ value: 32 ether }(
            VALIDATOR_PUBKEY, bytes(""), SIGNATURE, DEPOSIT_DATA_ROOT
        );
    }

    function testFuzz_DepositWrongMinAmount(uint256 amountInEther) public {
        amountInEther = _bound(amountInEther, 0, 31);
        uint256 amountInETH = amountInEther * 1 ether;
        vm.deal(depositor, amountInETH);
        vm.expectRevert(IBeaconDepositContract.DepositValueTooLow.selector);
        depositContract.deposit{ value: amountInETH }(
            VALIDATOR_PUBKEY, WITHDRAWAL_CREDENTIALS, SIGNATURE, bytes32("")
        );
    }

    function test_DepositWrongMinAmount() public {
        uint256 amount = 31 ether;
        vm.deal(depositor, amount);
        vm.expectRevert(IBeaconDepositContract.DepositValueTooLow.selector);
        depositContract.deposit{ value: amount }(
            VALIDATOR_PUBKEY, WITHDRAWAL_CREDENTIALS, SIGNATURE, bytes32("")
        );
    }

    function test_DepositFailsWithMaxAmount() public {
        vm.deal(depositor, uint256(type(uint64).max) * 2 gwei);
        vm.expectRevert(IBeaconDepositContract.DepositValueTooHigh.selector);
        depositContract.deposit{ value: uint256(type(uint64).max) * 2 gwei }(
            VALIDATOR_PUBKEY, WITHDRAWAL_CREDENTIALS, SIGNATURE, bytes32("")
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
            VALIDATOR_PUBKEY, WITHDRAWAL_CREDENTIALS, SIGNATURE, bytes32("")
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
            VALIDATOR_PUBKEY, WITHDRAWAL_CREDENTIALS, SIGNATURE, bytes32("")
        );

        amount = 32e9 - 1;
        vm.deal(depositor, amount);
        vm.expectRevert(
            IBeaconDepositContract.DepositNotMultipleOfGwei.selector
        );
        vm.prank(depositor);
        depositContract.deposit{ value: amount }(
            VALIDATOR_PUBKEY, WITHDRAWAL_CREDENTIALS, SIGNATURE, bytes32("")
        );
    }

    function test_DepositFailsWithInvalidDepositRoot() public {
        vm.deal(depositor, 32 ether);
        vm.expectRevert(IBeaconDepositContract.InvalidDepositDataRoot.selector);
        depositContract.deposit{ value: 32 ether }(
            VALIDATOR_PUBKEY, WITHDRAWAL_CREDENTIALS, SIGNATURE, bytes32("")
        );
    }

    function test_Deposit() public {
        vm.deal(depositor, 32 ether);
        vm.expectEmit(true, true, true, true);
        emit IBeaconDepositContract.Deposit(
            VALIDATOR_PUBKEY, WITHDRAWAL_CREDENTIALS, 32 gwei, SIGNATURE, 0
        );
        depositContract.deposit{ value: 32 ether }(
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
                32 gwei,
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

    function createDepositDataRoot(
        uint64 amountInGwei,
        bytes calldata signature
    )
        public
        view
        returns (bytes32)
    {
        bytes memory amount = _toLittleEndian64(amountInGwei);
        // Compute deposit data root (`DepositData` hash tree root)
        bytes32 pubkey_root =
            sha256(abi.encodePacked(VALIDATOR_PUBKEY, bytes16(0)));
        bytes32 signature_root = sha256(
            abi.encodePacked(
                sha256(abi.encodePacked(signature[:64])),
                sha256(abi.encodePacked(signature[64:], bytes32(0)))
            )
        );
        bytes32 node = sha256(
            abi.encodePacked(
                sha256(abi.encodePacked(pubkey_root, WITHDRAWAL_CREDENTIALS)),
                sha256(abi.encodePacked(amount, bytes24(0), signature_root))
            )
        );
        return node;
    }

    function _toLittleEndian64(uint64 value)
        internal
        pure
        returns (bytes memory ret)
    {
        ret = new bytes(8);
        bytes8 bytesValue = bytes8(value);
        // Byteswapping during copying to bytes.
        ret[0] = bytesValue[7];
        ret[1] = bytesValue[6];
        ret[2] = bytesValue[5];
        ret[3] = bytesValue[4];
        ret[4] = bytesValue[3];
        ret[5] = bytesValue[2];
        ret[6] = bytesValue[1];
        ret[7] = bytesValue[0];
    }
}
