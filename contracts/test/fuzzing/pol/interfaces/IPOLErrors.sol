// SPDX-License-Identifier: MIT
pragma solidity >=0.8.4;

import { IStakingRewardsErrors } from "../../base/IStakingRewardsErrors.sol";

/// @notice Interface of POL errors
interface IPOLErrors is IStakingRewardsErrors {
    // Signature: 0xf2d81d95
    error NotApprovedSender();
    // Signature: 0x02e6c295
    error NotRootFollower();
    // Signature: 0x1db3b859
    error NotDelegate();
    // Signature: 0x53f0a596
    error NotBGT();
    // Signature: 0x1b0eb4ec
    error NotBlockRewardController();
    // Signature: 0x385296d5
    error NotDistributor();
    // Signature: 0x73fcd3fe
    error NotFeeCollector();
    // Signature: 0x36407850
    error NotWhitelistedVault();
    // Signature: 0xb56f932c
    error NotGovernance();
    // Signature:0x7c214f04
    error NotOperator();
    // Signature: 0xad3a8b9e
    error NotEnoughBalance();
    // Signature: 0xadd377f6
    error InvalidActivateBoostDelay();
    // Signature: 0x2f14f4f9
    error InvalidDropBoostDelay();
    // Signature: 0x14969061
    error NotEnoughBoostedBalance();
    // Signature: 0xe8966d7a
    error NotEnoughTime();
    // Signature: 0xec2caa0d
    error InvalidStartBlock();
    // Signature: 0x3be31f8c
    error RewardAllocationAlreadyQueued();
    // Signature: 0x13134d24
    error InvalidRewardAllocationWeights();
    // Signature: 0xf6fae721
    error TooManyWeights();
    // Signature: 0x3e38573f
    error InvalidateDefaultRewardAllocation();
    // Signature: 0x0dc149f0
    error AlreadyInitialized();
    // Signature: 0x04aabf33
    error VaultAlreadyExists();
    // Signature: 0xd92e233d
    error ZeroAddress();
    // Signature: 0x462a2bb2
    error RewardAllocationBlockDelayTooLarge();
    // Signature: 0x08519afa
    error NotFactoryVault();
    // Signature: 0x978dc040
    error ZeroPercentageWeight();

    /*                   BLOCK REWARD CONTROLLLER                  */

    // Signature: 0x2e2dab43
    error InvalidBaseRate();
    // Signature: 0x22be2284
    error InvalidRewardRate();
    // Signature: 0x15482337
    error InvalidMinBoostedRewardRate();
    // Signature: 0xb7b2319a
    error InvalidBoostMultiplier();
    // Signature: 0x347f95b2
    error InvalidRewardConvexity();

    /*                           STAKING                           */

    // Signature: 0x09ee12d5
    error NotAContract();
    // Signature: 0x9e1ccc6f
    error InvalidRewardVaultAddress();
    // Signature: 0xe4ea100b
    error CannotRecoverRewardToken();
    // Signature: 0x1b813803
    error CannotRecoverStakingToken();
    // Signature: 0x2899103f
    error CannotRecoverIncentiveToken();
    // Signature: 0x38432c89
    error IncentiveRateTooHigh();
    // Signature: 0x5ee4de0e
    error NotIncentiveManager();

    // Signature: 0xf84835a0
    error TokenNotWhitelisted();
    // Signature: 0x8d1473a6
    error InsufficientDelegateStake();
    // Signature: 0x08e88f46
    error InsufficientSelfStake();
    // Signature: 0xfbf97e07
    error TokenAlreadyWhitelistedOrLimitReached();
    // Signature: 0xad57d95d
    error AmountLessThanMinIncentiveRate();
    // Signature: 0xfbf1123c
    error InvalidMaxIncentiveTokensCount();

    // Signature: 0x546c7600
    error PayoutAmountIsZero();
    // Signature: 0x89c622a2
    error DonateAmountLessThanPayoutAmount();
    // Signature: 0xa4cc22ed
    error MaxNumWeightsPerRewardAllocationIsZero();
    // Signature: 0x0b5c3aff
    error MinIncentiveRateIsZero();

    // Signature: 0x8e7572da
    error InvariantCheckFailed();

    /*                         BEACON ROOTS                        */

    // Signature: 0x1390f2a1
    error IndexOutOfRange();
    // Signature: 0x09bde339
    error InvalidProof();
    // Signature: 0x0a431b2a
    error TimestampAlreadyProcessed();

    /*                        BEACON DEPOSIT                       */

    /// @dev Error thrown when the deposit amount is too small, to prevent dust deposits.
    // Signature: 0x0e1eddda
    error InsufficientDeposit();

    /// @dev Error thrown when the deposit amount is not a multiple of Gwei.
    // Signature: 0x40567b38
    error DepositNotMultipleOfGwei();

    /// @dev Error thrown when the deposit amount is too high, since it is a uint64.
    // Signature: 0x2aa66734
    error DepositValueTooHigh();

    /// @dev Error thrown when the public key length is not 48 bytes.
    // Signature: 0x9f106472
    error InvalidPubKeyLength();

    /// @dev Error thrown when the withdrawal credentials length is not 32 bytes.
    // Signature: 0xb39bca16
    error InvalidCredentialsLength();

    /// @dev Error thrown when the signature length is not 96 bytes.
    // Signature: 0x4be6321b
    error InvalidSignatureLength();

    /// @dev Error thrown when the input operator is zero address on the first deposit.
    // Signature: 0x51969a7a
    error ZeroOperatorOnFirstDeposit();

    /// @dev Error thrown when the operator is already set and caller passed non-zero operator.
    // Signature: 0xc4142b41
    error OperatorAlreadySet();

    /// @dev Error thrown when the caller is not the current operator.
    // Signature: 0x819a0d0b
    error NotNewOperator();
}
