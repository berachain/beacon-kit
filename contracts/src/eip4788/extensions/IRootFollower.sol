// SPDX-License-Identifier: MIT
pragma solidity >=0.8.4;

/// @title IRootFollower
/// @author Berachain
/// @dev The interface for an abstract follower of the beacon root contract.
interface IRootFollower {
    /// @notice Emitted when the block is advanced.
    /// @param blockNum The block number of the block just actioned upon.
    event AdvancedBlock(uint256 blockNum);

    /// @notice Emitted when the block count is skipped.
    /// @param start The start block number of the block just actioned upon.
    /// @param end The end block number of the block just actioned upon.
    event BlockCountReset(uint256 start, uint256 end);

    /// @notice Gets the address of the coinbase for the given block number.
    /// @dev The BeaconRootsContract stores the coinbase for the last 8191
    /// blocks.
    /// @dev Querying a block number greater than the last 8191 blocks will
    /// return an error.
    /// @dev This also implies that actions should be invoked within 8191 blocks
    /// of being proposed. Otherwise any intended actions that were supposed to
    /// occur will be missed as the coinbase for the given block will no longer
    /// be available from the beacon root contract.
    /// @param blockNum The address performing the mint.
    /// @return coinbase The address of the coinbase for the given block number.
    function getCoinbase(uint256 blockNum)
        external
        view
        returns (address coinbase);

    /// @notice Gets the next block to be rewarded.
    /// @dev This returns the greater of current previously invoked block + 1,
    /// or current block number - 8191 as that is the limitation on number of
    /// blocks that can be queried, and actioned upon.
    /// @return blockNum The block number of the next block to be invoked.
    function getNextActionableBlock()
        external
        view
        returns (uint256 blockNum);

    /// @notice Gets the last block that was actioned upon.
    /// @return blockNum The block number of the last block that was actioned
    /// upon.
    function getLastActionedBlock() external view returns (uint256 blockNum);

    /// @notice Increments the block number to the next block.
    /// @dev This action should be permissioned to prevent unauthorized actors
    /// from modifying the block number inappropriately.
    function incrementBlock() external;

    /// @notice Resets the next actionable block number to _block, used when out
    /// of the beacon root buffer.
    /// @dev This action should be permissioned to prevent unauthorized actors
    /// from modifying the block number inappropriately.
    /// @param _block The block number to reset to.
    function resetCount(uint256 _block) external;
}
