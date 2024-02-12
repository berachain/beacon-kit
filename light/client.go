package light

// light client runs an infinite loop that ticks to query the current state of the beacon fork choice.
// It does this by querying the full node for the current head, and verifying the merkle proof.
// It is responsible for keeping the beacon chain in sync with the execution chain.
// Once it has the correct value, it will send a request to update the fork choice using UpdateForkChoice.
