# engine

Each folder within `types/engine` package matches a hard fork":

- **v1:** Pre-Shapella
- **v2:** Shapella (Shanghai + Capella)
- **v3:** Dencun (Deneb + Cancun)
- **v4:** Pralectra (Electra + Prague)

New fork versions may require updates to these types for new or altered functionalities. Updates are made only when necessary, preferring to reuse types from previous versions if their structure or functionality remains unchanged.

For instance, if a new fork version introduces a change that affects the execution payload structure but leaves the fork choice state unchanged, the new version (`v2`) of the `types/engine` package would include an updated `ExecutionPayloadCapella` message to reflect the changes specific to the new fork. Meanwhile, it would continue to reference the `ForkchoiceState` message from `v1` if no modifications to the fork choice logic are required.

This methodology ensures that each version within the `types/engine` package is a clear and concise representation of the data structures needed for that particular fork version, making it easier to manage and understand over time.
