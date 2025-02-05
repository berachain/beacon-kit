// SPDX-License-Identifier: MIT
pragma solidity >=0.8.4;

/// @notice Interface of staking rewards errors
interface IStakingRewardsErrors {
    // Signature: 0xf4ba521f
    error InsolventReward();
    // Signature: 0xf1bc94d2
    error InsufficientStake();
    // Signature: 0x49835af0
    error RewardCycleNotEnded();
    // Signature: 0x5ce91fd0
    error StakeAmountIsZero();
    // Signature: 0xe5cfe957
    error TotalSupplyOverflow();
    // Signature: 0xa393d14b
    error WithdrawAmountIsZero();
    // Signature: 0x359f174d
    error RewardsDurationIsZero();
}
