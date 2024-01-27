# beacon 

## Subfolders

- `blockchain`: This folder contains the code for managing the state transition on
the beacon chain as well as orchestrating calls to the execution client.
- `execution`: This folder contains the code for communicating with the execution 
engine, including the Ethereum 1 client. The client is responsible for connecting to 
the Ethereum 1 node, processing blocks, and interacting with the blockchain.
- `initial-sync`: This folder contains the code for handling the synchronizatino of the 
execution client with respect to the consensus client.
as well as triggering the execution client to sync.
- `logs`: This folder contains the code for processing logs from the Ethereum 1 client. 
It includes the implementation of the Processor, which is responsible for processing 
logs, and the Handler interface, which defines the ABIEvents method.
- `withdrawals`: This folder contains the code for handling validator withdrawal operations,
including withdrawal proofs and exit routines.
