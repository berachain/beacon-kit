// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "@solady/test/utils/SoladyTest.sol";
import "@solady/src/utils/FixedPointMathLib.sol";
import "@src/staking/DepositContract.sol";
import "@solady/src/tokens/ERC20.sol";
import { IDepositContract } from "../src/staking/IDepositContract.sol";

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
    DepositContract internal depositContract =
        DepositContract(DEPOSIT_CONTRACT_ADDRESS);

    /// @dev the snapshot we will revert to after each test.
    uint256 internal snapshot;

    /// @dev the native snapshot.
    uint256 internal nativeSnapshot;

    /// @dev the STAKE token that we will use.
    ERC20Test internal stakeToken;

    function setUp() public virtual {
        // etch the DepositContract.
        vm.etch(
            DEPOSIT_CONTRACT_ADDRESS, vm.getDeployedCode("DepositContract.sol")
        );
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
        vm.expectRevert(IDepositContract.InvalidPubKeyLength.selector);
        depositContract.deposit(
            bytes("wrong_public_key"),
            STAKING_CREDENTIALS,
            32e9,
            _create96Byte()
        );
    }

    function testFuzz_DepositWrongStakingCredentials(
        bytes calldata stakingCredentials
    )
        public
    {
        vm.revertTo(snapshot);
        vm.assume(stakingCredentials.length != 32);

        vm.expectRevert(IDepositContract.InvalidCredentialsLength.selector);
        depositContract.deposit(
            _create48Byte(), stakingCredentials, 32e9, _create96Byte()
        );
    }

    function testFuzz_DepositWrongAmount(uint256 amount) public {
        vm.revertTo(snapshot);
        amount = _bound(amount, 1, 32e9 - 1);

        vm.startPrank(depositor);
        stakeToken.mint(depositor, amount);

        vm.expectRevert(IDepositContract.InsufficientDeposit.selector);
        depositContract.deposit(
            VALIDATOR_PUBKEY,
            STAKING_CREDENTIALS,
            uint64(amount),
            _create96Byte()
        );
    }

    function testFuzz_DepositWrongMultiple(uint256 amount) public {
        vm.revertTo(snapshot);
        amount = _bound(amount, 32e9 + 1, type(uint64).max);
        vm.assume(amount % 1e9 != 0);

        vm.startPrank(depositor);
        stakeToken.mint(depositor, amount);

        vm.expectRevert(IDepositContract.DepositNotMultipleOfGwei.selector);
        depositContract.deposit(
            VALIDATOR_PUBKEY,
            STAKING_CREDENTIALS,
            uint64(amount),
            _create96Byte()
        );
    }

    function test_Deposit() public {
        vm.revertTo(snapshot);
        vm.startPrank(depositor);
        stakeToken.mint(depositor, 32e9);

        vm.expectEmit(true, true, true, true);
        emit IDepositContract.Deposit(
            VALIDATOR_PUBKEY, STAKING_CREDENTIALS, 32e9, _create96Byte()
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
        vm.expectRevert(IDepositContract.InvalidPubKeyLength.selector);
        depositContract.redirect(fromPubKey, VALIDATOR_PUBKEY, 32e9);
    }

    function test_fuzzWrongToPubKey(bytes calldata toPubKey) public {
        vm.revertTo(snapshot);
        vm.assume(toPubKey.length != 48);
        vm.expectRevert(IDepositContract.InvalidPubKeyLength.selector);
        depositContract.redirect(VALIDATOR_PUBKEY, toPubKey, 32e9);
    }

    function test_fuzzRedirectWrongAmount(uint256 amount) public {
        vm.revertTo(snapshot);
        amount = _bound(amount, 1, 32e9 / 10 - 1);

        vm.expectRevert(IDepositContract.InsufficientRedirectAmount.selector);
        depositContract.redirect(
            VALIDATOR_PUBKEY, VALIDATOR_PUBKEY, uint64(amount)
        );
    }

    function testRedirect() public {
        vm.revertTo(snapshot);
        vm.expectEmit(true, true, true, true);

        vm.startBroadcast(depositor);
        emit IDepositContract.Redirect(
            VALIDATOR_PUBKEY, VALIDATOR_PUBKEY, _credential(depositor), 32e9
        );
        depositContract.redirect(VALIDATOR_PUBKEY, VALIDATOR_PUBKEY, 32e9);
    }

    function test_fuzzWithdrawWrongPubKey(bytes calldata pubKey) public {
        vm.revertTo(snapshot);
        vm.assume(pubKey.length != 48);
        vm.expectRevert(IDepositContract.InvalidPubKeyLength.selector);
        depositContract.withdraw(pubKey, WITHDRAWAL_CREDENTIALS, 32e9);
    }

    function test_fuzzWithdrawWrongWithdrawalCredentials(
        bytes calldata withdrawalCredentials
    )
        public
    {
        vm.revertTo(snapshot);
        vm.assume(withdrawalCredentials.length != 32);
        vm.expectRevert(IDepositContract.InvalidCredentialsLength.selector);
        depositContract.withdraw(VALIDATOR_PUBKEY, withdrawalCredentials, 32e9);
    }

    function test_WithdrawWrongAmount(uint256 amount) public {
        vm.revertTo(snapshot);
        amount = _bound(amount, 1, 32e9 / 10 - 1);

        vm.expectRevert(IDepositContract.InsufficientWithdrawAmount.selector);
        depositContract.withdraw(
            VALIDATOR_PUBKEY, WITHDRAWAL_CREDENTIALS, uint64(amount)
        );
    }

    function testWithdraw() public {
        vm.revertTo(snapshot);
        vm.expectEmit(true, true, true, true);

        vm.startBroadcast(depositor);
        emit IDepositContract.Withdraw(
            VALIDATOR_PUBKEY, WITHDRAWAL_CREDENTIALS, 32e9
        );
        depositContract.withdraw(VALIDATOR_PUBKEY, WITHDRAWAL_CREDENTIALS, 32e9);
    }

    function test_fuzzDepositNativeWrongMinAmount(uint256 amount) public {
        vm.revertTo(nativeSnapshot);
        amount = _bound(amount, 1, 32 gwei - 1);
        vm.deal(depositor, amount);
        vm.expectRevert(IDepositContract.InsufficientDeposit.selector);
        depositContract.deposit{ value: amount }(
            VALIDATOR_PUBKEY, STAKING_CREDENTIALS, 0, _create96Byte()
        );
    }

    function test_DepositNativeNotDivisibleByGwei(uint256 amount) public {
        vm.revertTo(nativeSnapshot);
        amount = _bound(amount, 32e9 + 1, uint256(type(uint64).max));
        vm.assume(amount % 1e9 != 0);
        vm.deal(depositor, amount);
        vm.startPrank(depositor);

        vm.expectRevert(IDepositContract.DepositNotMultipleOfGwei.selector);
        depositContract.deposit{ value: amount }(
            VALIDATOR_PUBKEY, STAKING_CREDENTIALS, 0, _create96Byte()
        );
    }

    function test_DepositNative() public {
        vm.revertTo(nativeSnapshot);
        vm.deal(depositor, 1 ether);

        vm.expectEmit(true, true, true, true);
        emit IDepositContract.Deposit(
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
    }
}
