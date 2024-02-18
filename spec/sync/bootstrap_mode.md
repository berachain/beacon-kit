# Bootstrap Mode

- We cannot make `engine_` JSON-RPC requests during `FinalizeBlock()` as this will
cause the sync'ing process to go extremely slowly. The only operations that happen during this
period of ABCI should be updating the Beacon chain's view of what the finalized and safe
blocks are of the execution chain.
- However we need to come up with some sort of bootstrapping mode, to handle this. This is because if we sync the entire beacon chain from genesis and then only start syncing the 
execution chain once the beacon chain is fully sync'd we will run into problems.




# High Level Plan

- Define an `sync` queue.

- We should have some sort of global sync status that we can read off of the rpc.
- Always start in some sort of bootstrapping mode.

Case 1: The beacon chain is sync'd, execution chain is not.
- After `app.Load()` we kickoff an asynchronous job onto the queue that will
call `forkchoiceUpdate` over and over to get the execution chain caughtup.
- Once this is complete, we enter prepare / process proposal and life is good.

Case 2: The beacon chain is not sync'd, execution chain is sync'd.
- This is easy, (though in practice doesn't really work, due to timing), in
theory we just replay and then click in, however during replay, the execution chain
will ultimately fall behind. So once we fully sync the beacon, the execution chain will need to sync some blocks as well, and we potentially see a STATUS_SYNCING. We will however be able to.

Case 3: We see a STATUS_SYNCING or STATUS_ACCEPTED in Prepare or Process Proposal.
- We need to fire into some boostrap mode and pause preparing or proposing until we are caught up. We should just insant vote nil.