// SPDX-License-Identifier: MIT
pragma solidity ^0.8.25;

import { SoladyTest } from "@solady/test/utils/SoladyTest.sol";
import { ERC20 } from "@solady/src/tokens/ERC20.sol";
import { IBeaconDepositContract } from
    "../src/staking/IBeaconDepositContract.sol";
import { SoladyTest } from "@solady/test/utils/SoladyTest.sol";
import { BeaconDepositContract } from "@src/staking/BeaconDepositContract.sol";
import { ERC20BeaconDepositContract } from
    "@src/staking/extensions/ERC20BeaconDepositContract.sol";

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
    address internal constant DEPOSIT_CONTRACT_ADDRESS =
        0x00000000219ab540356cBB839Cbe05303d7705Fa;

    bytes32 internal constant STAKING_ASSET_SLOT = bytes32(uint256(1));

    /// @dev the deposit contract.
    BeaconDepositContract internal depositContract;

    /// @dev the deposit contract.
    ERC20BeaconDepositContract internal erc20DepositContract;

    /// @dev the STAKE token that we will use.
    ERC20Test internal stakeToken;

    function setUp() public virtual {
        erc20DepositContract = new ERC20BeaconDepositContract();
        // Deploy the STAKE token.
        stakeToken = new ERC20Test();
        // Set the STAKE_ASSET to the STAKE token.
        bytes32 stakeAssetValue = bytes32(uint256(uint160(address(stakeToken))));
        vm.store(
            address(erc20DepositContract), STAKING_ASSET_SLOT, stakeAssetValue
        );

        // Set the STAKE_ASSET to the NATIVE token.
        depositContract = new BeaconDepositContract();
    }

    function testFuzz_DepositsWrongPubKey(bytes calldata pubKey) public {
        vm.assume(pubKey.length != 96);
        vm.expectRevert(IBeaconDepositContract.InvalidPubKeyLength.selector);
        erc20DepositContract.deposit(
            bytes("wrong_public_key"),
            STAKING_CREDENTIALS,
            32e9,
            _create96Byte()
        );
    }

    function test_DepositWrongPubKey() public {
        vm.expectRevert(IBeaconDepositContract.InvalidPubKeyLength.selector);
        erc20DepositContract.deposit(
            bytes("wrong_public_key"),
            STAKING_CREDENTIALS,
            32e9,
            _create96Byte()
        );
    }

    function testFuzz_DepositWrongcredentials(bytes calldata credentials)
        public
    {
        vm.assume(credentials.length != 32);

        vm.expectRevert(
            IBeaconDepositContract.InvalidCredentialsLength.selector
        );
        erc20DepositContract.deposit(
            _create48Byte(), credentials, 32e9, _create96Byte()
        );
    }

    function test_DepositWrongcredentials() public {
        vm.expectRevert(
            IBeaconDepositContract.InvalidCredentialsLength.selector
        );
        erc20DepositContract.deposit(
            VALIDATOR_PUBKEY, bytes("wrong_credentials"), 32e9, _create96Byte()
        );
    }

    function testFuzz_DepositWrongAmount(uint256 amount) public {
        amount = _bound(amount, 1, 32e9 - 1);
        stakeToken.mint(depositor, amount * 1e9);

        vm.prank(depositor);
        vm.expectRevert(IBeaconDepositContract.InsufficientDeposit.selector);
        erc20DepositContract.deposit(
            VALIDATOR_PUBKEY,
            STAKING_CREDENTIALS,
            uint64(amount),
            _create96Byte()
        );
    }

    function test_DepositWrongAmount() public {
        stakeToken.mint(depositor, (32e9 - 1) * 1e9);
        vm.expectRevert(IBeaconDepositContract.InsufficientDeposit.selector);
        vm.prank(depositor);
        erc20DepositContract.deposit(
            VALIDATOR_PUBKEY, STAKING_CREDENTIALS, 32e9 - 1, _create96Byte()
        );
    }

    function test_Deposit() public {
        stakeToken.mint(depositor, 32e9 * 1e9);

        vm.prank(depositor);
        vm.expectEmit(true, true, true, true);
        emit IBeaconDepositContract.Deposit(
            VALIDATOR_PUBKEY, STAKING_CREDENTIALS, 32e9, _create96Byte(), 0
        );
        erc20DepositContract.deposit(
            VALIDATOR_PUBKEY, STAKING_CREDENTIALS, 32e9, _create96Byte()
        );
    }

    function testFuzz_WithdrawWrongPubKey(bytes calldata pubKey) public {
        vm.assume(pubKey.length != 48);
        vm.expectRevert(IBeaconDepositContract.InvalidPubKeyLength.selector);
        erc20DepositContract.withdraw(pubKey, WITHDRAWAL_CREDENTIALS, 32e9);
    }

    function test_WithdrawWrongPubKey() public {
        vm.expectRevert(IBeaconDepositContract.InvalidPubKeyLength.selector);
        erc20DepositContract.withdraw(
            bytes("wrong_pub_key"), WITHDRAWAL_CREDENTIALS, 32e9
        );
    }

    function testFuzz_WithdrawWrongWithdrawalCredentials(
        bytes calldata withdrawalCredentials
    )
        public
    {
        vm.assume(withdrawalCredentials.length != 32);
        vm.expectRevert(
            IBeaconDepositContract.InvalidCredentialsLength.selector
        );
        erc20DepositContract.withdraw(
            VALIDATOR_PUBKEY, withdrawalCredentials, 32e9
        );
    }

    function test_WithdrawWrongWithdrawCredentials() public {
        vm.expectRevert(
            IBeaconDepositContract.InvalidCredentialsLength.selector
        );
        erc20DepositContract.withdraw(
            VALIDATOR_PUBKEY, bytes("wrong_credentials"), 32e9
        );
    }

    function testFuzz_WithdrawWrongAmount(uint256 amount) public {
        amount = _bound(amount, 1, 32e9 / 10 - 1);

        vm.expectRevert(
            IBeaconDepositContract.InsufficientWithdrawalAmount.selector
        );
        erc20DepositContract.withdraw(
            VALIDATOR_PUBKEY, WITHDRAWAL_CREDENTIALS, uint64(amount)
        );
    }

    function test_WithdrawWrongAmount() public {
        vm.expectRevert(
            IBeaconDepositContract.InsufficientWithdrawalAmount.selector
        );
        erc20DepositContract.withdraw(
            VALIDATOR_PUBKEY, WITHDRAWAL_CREDENTIALS, 32e9 / 10 - 1
        );
    }

    function testWithdraw() public {
        vm.expectEmit(true, true, true, true);

        vm.prank(depositor);
        emit IBeaconDepositContract.Withdrawal(
            VALIDATOR_PUBKEY,
            _credential(depositor),
            WITHDRAWAL_CREDENTIALS,
            32e9,
            0
        );
        erc20DepositContract.withdraw(
            VALIDATOR_PUBKEY, WITHDRAWAL_CREDENTIALS, 32e9
        );
    }

    function testFuzz_DepositNativeWrongMinAmount(uint256 amountInEther)
        public
    {
        amountInEther = _bound(amountInEther, 1, 31);
        uint256 amountInGwei = amountInEther * 1 gwei;
        vm.deal(depositor, amountInGwei);
        vm.expectRevert(IBeaconDepositContract.InsufficientDeposit.selector);
        depositContract.deposit{ value: amountInGwei }(
            VALIDATOR_PUBKEY, STAKING_CREDENTIALS, 0, _create96Byte()
        );
    }

    function test_DepositNativeWrongMinAmount() public {
        uint256 amount = 31 gwei;
        vm.deal(depositor, amount);
        vm.expectRevert(IBeaconDepositContract.InsufficientDeposit.selector);
        depositContract.deposit{ value: amount }(
            VALIDATOR_PUBKEY, STAKING_CREDENTIALS, 0, _create96Byte()
        );
    }

    function testFuzz_DepositNativeNotDivisibleByGwei(uint256 amount) public {
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
        vm.deal(depositor, 32 ether);
        vm.expectEmit(true, true, true, true);
        emit IBeaconDepositContract.Deposit(
            VALIDATOR_PUBKEY, STAKING_CREDENTIALS, 32 gwei, _create96Byte(), 0
        );
        depositContract.deposit{ value: 32 ether }(
            VALIDATOR_PUBKEY, STAKING_CREDENTIALS, 0, _create96Byte()
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
                VALIDATOR_PUBKEY, STAKING_CREDENTIALS, 0, _create96Byte()
            );
            ++depositCount;
        }
        assertEq(depositContract.depositCount(), depositCount);
    }

    function testFuzz_WithdrawCount(uint256 count) public {
        count = _bound(count, 1, 100);
        testFuzz_DepositCount(count);
        uint64 withdrawalCount;
        for (uint256 i; i < count; ++i) {
            vm.expectEmit(true, true, true, true);
            emit IBeaconDepositContract.Withdrawal(
                VALIDATOR_PUBKEY,
                _credential(depositor),
                WITHDRAWAL_CREDENTIALS,
                32 gwei,
                withdrawalCount
            );
            depositContract.withdraw(
                VALIDATOR_PUBKEY, WITHDRAWAL_CREDENTIALS, 32e9
            );
            ++withdrawalCount;
        }
        assertEq(depositContract.withdrawalCount(), withdrawalCount);
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
