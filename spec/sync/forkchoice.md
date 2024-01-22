# Forkchoice


There are two ways we can handle forkchoices in Beaconkit.

1. We utilize a simlar format to Polaris, where the proposer includes 
an Execution Payload in their proposal. This proposal is then gossiped around and all validators will vote yes or no on inclusion.
2. We utilize Vote Extensions, all validators include their latest head in
the VE and then we perform the LMD ghost algorithmn on them. The LMD Ghost alogrithmn will output us a new Finalized block and then we can set this outputted block from the algo to the store in PreBlocker().

Additionally, we need to decide if we are going to gossip full payloads at the consensus layer (this is what we do now). For the VE model is this definitely bad.

- We also need to decide if our forkchoice state (HEAD, SAFE etc.) are being done on
the `State Root` of the Execution chain or the `Block Hash` (tbh just gotta read the prysm repo).