// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import { SoladyTest } from "@solady/test/utils/SoladyTest.sol";
import { ERC20 } from "@solady/src/tokens/ERC20.sol";
import { IBeaconDepositContract } from
    "../src/staking/IBeaconDepositContract.sol";
import { SoladyTest } from "@solady/test/utils/SoladyTest.sol";
import { BeaconDepositContract } from "@src/staking/BeaconDepositContract.sol";

// Mock ERC20 token that we will use as the stake token.
contract ERC20Test is ERC20 {
    function mint(address to, uint256 amount) public {
        _mint(to, amount);
    }

    function burn(address from, uint256 amount) public {
        _burn(from, amount);
    }

    function name() public pure override returns (string memory) {
        return "STAKE";
    }

    function symbol() public pure override returns (string memory) {
        return "STAKE";
    }
}

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
    address internal DEPOSIT_CONTRACT_ADDRESS =
        0x00000000219ab540356cBB839Cbe05303d7705Fa;

    /// @dev the deposit contract.
    BeaconDepositContract internal depositContract;

    /// @dev the snapshot we will revert to after each test.
    uint256 internal snapshot;

    /// @dev the native snapshot.
    uint256 internal nativeSnapshot;

    /// @dev the STAKE token that we will use.
    ERC20Test internal stakeToken;

    function setUp() public virtual {
        DEPOSIT_CONTRACT_ADDRESS = address(new BeaconDepositContract());
        depositContract = BeaconDepositContract(DEPOSIT_CONTRACT_ADDRESS);
        // Deploy the STAKE token.
        stakeToken = new ERC20Test();
        // Set the STAKE_ASSET to the STAKE token.
        bytes32 stakeAssetSlot = bytes32(uint256(0));
        bytes32 stakeAssetValue = bytes32(uint256(uint160(address(stakeToken))));
        vm.store(DEPOSIT_CONTRACT_ADDRESS, stakeAssetSlot, stakeAssetValue);
        snapshot = vm.snapshot();

        // Set the STAKE_ASSET to the NATIVE token.
        address NATIVE_ASSET = 0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE;
        bytes32 nativeAssetSlot = bytes32(uint256(0));
        bytes32 nativeAssetValue = bytes32(uint256(uint160(NATIVE_ASSET)));
        vm.store(DEPOSIT_CONTRACT_ADDRESS, nativeAssetSlot, nativeAssetValue);
        nativeSnapshot = vm.snapshot();
    }

    function testFuzz_DepositsWrongPubKey(bytes calldata pubKey) public {
        vm.revertTo(snapshot);
        vm.assume(pubKey.length != 96);
        vm.expectRevert(IBeaconDepositContract.InvalidPubKeyLength.selector);
        depositContract.deposit(
            bytes("wrong_public_key"),
            STAKING_CREDENTIALS,
            32e9,
            _create96Byte()
        );
    }

    function test_DepositWrongPubKey() public {
        vm.revertTo(snapshot);
        vm.expectRevert(IBeaconDepositContract.InvalidPubKeyLength.selector);
        depositContract.deposit(
            bytes("wrong_public_key"),
            STAKING_CREDENTIALS,
            32e9,
            _create96Byte()
        );
    }

    function testFuzz_DepositWrongcredentials(bytes calldata credentials)
        public
    {
        vm.revertTo(snapshot);
        vm.assume(credentials.length != 32);

        vm.expectRevert(
            IBeaconDepositContract.InvalidCredentialsLength.selector
        );
        depositContract.deposit(
            _create48Byte(), credentials, 32e9, _create96Byte()
        );
    }

    function test_DepositWrongcredentials() public {
        vm.revertTo(snapshot);
        vm.expectRevert(
            IBeaconDepositContract.InvalidCredentialsLength.selector
        );
        depositContract.deposit(
            VALIDATOR_PUBKEY, bytes("wrong_credentials"), 32e9, _create96Byte()
        );
    }

    function testFuzz_DepositWrongAmount(uint256 amount) public {
        vm.revertTo(snapshot);
        amount = _bound(amount, 1, 32e9 - 1);
        stakeToken.mint(depositor, amount * 1e9);

        vm.prank(depositor);
        vm.expectRevert(IBeaconDepositContract.InsufficientDeposit.selector);
        depositContract.deposit(
            VALIDATOR_PUBKEY,
            STAKING_CREDENTIALS,
            uint64(amount),
            _create96Byte()
        );
    }

    function test_DepositWrongAmount() public {
        vm.revertTo(snapshot);
        stakeToken.mint(depositor, (32e9 - 1) * 1e9);
        vm.expectRevert(IBeaconDepositContract.InsufficientDeposit.selector);
        vm.prank(depositor);
        depositContract.deposit(
            VALIDATOR_PUBKEY, STAKING_CREDENTIALS, 32e9 - 1, _create96Byte()
        );
    }

    function test_Deposit() public {
        vm.revertTo(snapshot);
        stakeToken.mint(depositor, 32e9 * 1e9);

        vm.prank(depositor);
        vm.expectEmit(true, true, true, true);
        emit IBeaconDepositContract.Deposit(
            VALIDATOR_PUBKEY, STAKING_CREDENTIALS, 32e9, _create96Byte(), 0
        );
        depositContract.deposit(
            VALIDATOR_PUBKEY, STAKING_CREDENTIALS, 32e9, _create96Byte()
        );
    }

    function testFuzz_RedirectWrongFromPubKey(bytes calldata fromPubKey)
        public
    {
        vm.revertTo(snapshot);
        vm.assume(fromPubKey.length != 48);
        vm.expectRevert(IBeaconDepositContract.InvalidPubKeyLength.selector);
        depositContract.redirect(fromPubKey, VALIDATOR_PUBKEY, 32e9);
    }

    function test_RedirectWrongFromPubKey() public {
        vm.revertTo(snapshot);
        vm.expectRevert(IBeaconDepositContract.InvalidPubKeyLength.selector);
        depositContract.redirect(bytes("wrong_pub_key"), VALIDATOR_PUBKEY, 32e9);
    }

    function testFuzz_WrongToPubKey(bytes calldata toPubKey) public {
        vm.revertTo(snapshot);
        vm.assume(toPubKey.length != 48);
        vm.expectRevert(IBeaconDepositContract.InvalidPubKeyLength.selector);
        depositContract.redirect(VALIDATOR_PUBKEY, toPubKey, 32e9);
    }

    function test_RedirectWrongToPubKey() public {
        vm.revertTo(snapshot);
        vm.expectRevert(IBeaconDepositContract.InvalidPubKeyLength.selector);
        depositContract.redirect(VALIDATOR_PUBKEY, bytes("wrong_pub_key"), 32e9);
    }

    function testFuzz_RedirectWrongAmount(uint256 amount) public {
        vm.revertTo(snapshot);
        amount = _bound(amount, 1, 32e9 / 10 - 1);

        vm.expectRevert(
            IBeaconDepositContract.InsufficientRedirectAmount.selector
        );
        depositContract.redirect(
            VALIDATOR_PUBKEY, VALIDATOR_PUBKEY, uint64(amount)
        );
    }

    function test_RedirectWrongAmount() public {
        vm.revertTo(snapshot);
        vm.expectRevert(
            IBeaconDepositContract.InsufficientRedirectAmount.selector
        );
        depositContract.redirect(
            VALIDATOR_PUBKEY, VALIDATOR_PUBKEY, 32e9 / 10 - 1
        );
    }

    function testRedirect() public {
        vm.revertTo(snapshot);
        vm.expectEmit(true, true, true, true);

        vm.prank(depositor);
        emit IBeaconDepositContract.Redirect(
            VALIDATOR_PUBKEY, VALIDATOR_PUBKEY, _credential(depositor), 32e9, 0
        );
        depositContract.redirect(VALIDATOR_PUBKEY, VALIDATOR_PUBKEY, 32e9);
    }

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
            IBeaconDepositContract.InsufficientWithdrawalAmount.selector
        );
        depositContract.withdraw(
            VALIDATOR_PUBKEY, WITHDRAWAL_CREDENTIALS, uint64(amount)
        );
    }

    function test_WithdrawWrongAmount() public {
        vm.revertTo(snapshot);
        vm.expectRevert(
            IBeaconDepositContract.InsufficientWithdrawalAmount.selector
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
            VALIDATOR_PUBKEY,
            _credential(depositor),
            WITHDRAWAL_CREDENTIALS,
            32e9,
            0
        );
        depositContract.withdraw(VALIDATOR_PUBKEY, WITHDRAWAL_CREDENTIALS, 32e9);
    }

    function testFuzz_DepositNativeWrongMinAmount(uint256 amountInEther)
        public
    {
        vm.revertTo(nativeSnapshot);
        amountInEther = _bound(amountInEther, 1, 31);
        uint256 amountInGwei = amountInEther * 1 gwei;
        vm.deal(depositor, amountInGwei);
        vm.expectRevert(IBeaconDepositContract.InsufficientDeposit.selector);
        depositContract.deposit{ value: amountInGwei }(
            VALIDATOR_PUBKEY, STAKING_CREDENTIALS, 0, _create96Byte()
        );
    }

    function test_DepositNativeWrongMinAmount() public {
        vm.revertTo(nativeSnapshot);
        uint256 amount = 31 gwei;
        vm.deal(depositor, amount);
        vm.expectRevert(IBeaconDepositContract.InsufficientDeposit.selector);
        depositContract.deposit{ value: amount }(
            VALIDATOR_PUBKEY, STAKING_CREDENTIALS, 0, _create96Byte()
        );
    }

    function testFuzz_DepositNativeNotDivisibleByGwei(uint256 amount) public {
        vm.revertTo(nativeSnapshot);
        amount = _bound(amount, 31e9 + 1, uint256(type(uint64).max));
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

        amount = 32e9 - 1;
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
        vm.deal(depositor, 32 ether);
        vm.expectEmit(true, true, true, true);
        emit IBeaconDepositContract.Deposit(
            VALIDATOR_PUBKEY, STAKING_CREDENTIALS, 32 gwei, _create96Byte(), 0
        );
        depositContract.deposit{ value: 32 ether }(
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
    }
}
