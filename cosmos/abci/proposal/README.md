# Block Proposals

## Overview

There are two types of block proposals:

**Vote extension enabled block proposal**: When vote extensions are enabled, the proposer must include the vote extensions they received from comet within their block proposal - specifically in the first slot of the proposal. We choose to use the proposer to determine the cannonical vote extensions for a height because each validator may be maintaining different
vote extensions which may lead to in-deterministics results.

**Vote extension disabled block proposal**: When vote extensions are disabled, the proposer does not need to inject any vote extensions
into their block proposal.



## Prepare Proposal

When vote extensions are enabled, the proposer will then inject the vote extensions into their block proposal. The reason why this is done is to ensure that the vote extensions can be aggregated and written to state in `PreBlock`. `PreBlock` does not have access to the vote extensions used to create the block proposal, so the vote extensions must be injected into the block proposal. Additionally, this allows the network to have a cannonical set of vote extensions for a given height.

In the case where the validator does not have valid vote extensions, a new round of voting will be triggered. The validator will then wait for the next round of voting to complete before creating a new block proposal.

The process of constructing the rest of the block is left to the `PrepareProposalHandler` which is passed into the constructor. This means that process of 'oracle' block building can be compatible with the Block-SDK, which is used to build highly custom blocks.

## Process Proposal

When vote extensions are enabled, the validator will first verify that the block contains the block proposer's vote extensions. If the block does not contain the block proposer's vote extensions, the block will be rejected. If the block contains the block proposer's vote extensions, the validator will do a basic check to ensure the vote extensions are valid before verifying the rest of the proposal in accordance with the preferences of the `ProcessProposalHandler` which is passed into the constructor.
