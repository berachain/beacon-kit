// SPDX-License-Identifier: MIT
pragma solidity ^0.8.25;

import { SoladyTest } from "@solady/test/utils/SoladyTest.sol";
import { IDepositContract } from "@src/staking/IDepositContract.sol";
import { PermissionedDepositContract } from "./PermissionedDepositContract.sol";

contract DepositContractTest is SoladyTest {
    /// @dev Address allowed to make deposits.
    address internal depositor = 0x20f33CE90A13a4b5E7697E3544c3083B8F8A51D4;

    /// @dev Validator public key for testing.
    bytes internal VALIDATOR_PUBKEY = _create48Byte();

    /// @dev Withdrawal credentials used in deposits.
    bytes internal WITHDRAWAL_CREDENTIALS = _generateCredentials(address(this));

    /// @dev Staking credentials for the depositor.
    bytes internal STAKING_CREDENTIALS = _generateCredentials(depositor);

    /// @dev Storage slot for staking assets.
    bytes32 internal constant STAKING_ASSET_SLOT = bytes32(0);

    /// @dev Permissioned deposit contract instance.
    PermissionedDepositContract internal depositContract;

    /// @notice Sets up the test environment.
    function setUp() public virtual {
        address owner = 0x6969696969696969696969696969696969696969;
        depositContract = new PermissionedDepositContract(owner);
        vm.prank(owner);
        depositContract.allowDeposit(depositor, 100);
    }

    /// @notice Test deposit with invalid public key length.
    function testFuzz_DepositsWrongPubKey(bytes calldata pubKey) public {
        _assumeInvalidLength(pubKey, 96);
        _expectRevert(IDepositContract.InvalidPubKeyLength.selector);
        _prankDeposit(pubKey, STAKING_CREDENTIALS, 32e9);
    }

    function test_DepositWrongPubKey() public {
        _expectRevert(IDepositContract.InvalidPubKeyLength.selector);
        _prankDeposit(bytes("wrong_public_key"), STAKING_CREDENTIALS, 32e9);
    }

    /// @notice Test deposit with invalid credentials length.
    function testFuzz_DepositWrongCredentials(bytes calldata credentials)
        public
    {
        _assumeInvalidLength(credentials, 32);
        _expectRevert(IDepositContract.InvalidCredentialsLength.selector);
        _prankDeposit(_create48Byte(), credentials, 32e9);
    }

    function test_DepositWrongCredentials() public {
        _expectRevert(IDepositContract.InvalidCredentialsLength.selector);
        _prankDeposit(VALIDATOR_PUBKEY, bytes("wrong_credentials"), 32e9);
    }

    /// @notice Test deposit with invalid amount.
    function testFuzz_DepositWrongAmount(uint256 amount) public {
        amount = _bound(amount, 1, 32e9 - 1);
        vm.deal(depositor, amount);
        _expectRevert(IDepositContract.InsufficientDeposit.selector);
        _prankDeposit(VALIDATOR_PUBKEY, STAKING_CREDENTIALS, uint64(amount));
    }

    function test_DepositWrongAmount() public {
        _expectRevert(IDepositContract.InsufficientDeposit.selector);
        _prankDeposit(VALIDATOR_PUBKEY, STAKING_CREDENTIALS, 32e9 - 1);
    }

    /// @notice Valid deposit test.
    function test_Deposit() public {
        vm.deal(depositor, 32 ether);
        vm.prank(depositor);
        vm.expectEmit(true, true, true, true);
        emit IDepositContract.Deposit(
            VALIDATOR_PUBKEY, STAKING_CREDENTIALS, 32e9, _create96Byte(), 0
        );
        depositContract.deposit{ value: 32 ether }(
            VALIDATOR_PUBKEY, STAKING_CREDENTIALS, 32e9, _create96Byte()
        );
    }

    /// @notice Helper to generate staking credentials.
    function _generateCredentials(address addr) internal pure returns (bytes memory) {
        return abi.encodePacked(bytes1(0x01), bytes11(0x0), addr);
    }

    /// @notice Helper to create 96-byte placeholder.
    function _create96Byte() internal pure returns (bytes memory) {
        return abi.encodePacked(bytes32("32"), bytes32("32"), bytes32("32"));
    }

    /// @notice Helper to create 48-byte placeholder.
    function _create48Byte() internal pure returns (bytes memory) {
        return abi.encodePacked(bytes32("32"), bytes16("16"));
    }

    /// @notice Helper to assume invalid length for fuzz tests.
    function _assumeInvalidLength(bytes calldata data, uint256 expectedLength) internal pure {
        vm.assume(data.length != expectedLength);
    }

    /// @notice Helper to expect revert with specific selector.
    function _expectRevert(bytes4 selector) internal {
        vm.expectRevert(selector);
    }

    /// @notice Helper to perform deposit with `prank`.
    function _prankDeposit(bytes memory pubKey, bytes memory credentials, uint64 amount) internal {
        vm.prank(depositor);
        depositContract.deposit(pubKey, credentials, amount, _create96Byte());
    }
}
