# State Processor

## Validators handling

As a general principle, BeaconKit strives to keep validators handling aligned with Ethereum 2.0 specs. There currently two notable exceptions to this principle:

- BeaconKit **does not** enforce a cap on validators churn, neither in the activation nor in the exit queue.
- BeaconKit **does** enforce an explicit cap on the validators set.
- BeaconKit **does not** currently support voluntary withdrawals.

We list below a few relevant facts.

## `Balance` and `EffectiveBalance`

BeaconKit distingishes a validator `Balance` and a validator `EffectiveBalance`.

- `Balance` is updated slot by slot, when a deposit in made over the deposit contract and events are subsequently processed by BeaconKit.
- `Balance` can increase only in multiples of `MinDepositAmount`, which is specified in the deposit contract. There is no cap on the `Balance`.
- `EffectiveBalance` is updated at the turn of every epoch.
- `EffectiveBalance` is enforced to be a multiple of `EffectiveBalanceIncrement`
- `EffectiveBalance` is capped at `MaxEffectiveBalance`. Any `Balance` in excess of `MaxEffectiveBalance` is automatically withdrawn.
- `EffectiveBalance` is aligned to `Balance`, within the constrains listed above, and accounting for hysteresys, at the turn of every epoch, via [`processEffectiveBalanceUpdates`](./state_processor.go#L491)

## Validator lifecycle

Validators are created by a deposit made to the deposit contract, as soon as the corresponding event is emitted.

Say a deposit is made at slot `S`, epoch `N`, that creates a validator.

The funds deposited will be locked and the validator will stay inactive until the a minimum staking balance is hit. The minimum staking balance is set to `EjectionBalance + EffectiveBalanceIncrement` to account for hysteresys. However if the active validator set has already reached the `ValidatorSetCap`, a new prospective validator must deposit more funds than the effective balance of the active validator with the lowest stake.

So let's assumed that `S`, epoch `N`, is the slot where finally the validator `Balance` equals the minimum required for staking (one or multiple deposits may have been done to get there). Then:

- The validator is marked as `EligibleForActivationQueue` as soon as epoch `N+1` starts. This is guaranteed since there is no cap on the activation queue size.
- The validator is marked as active as soon as epoch `N+2` starts. However
  - if the size of validator set goes beyond the `ValidatorSetCap` enough validators with the lowest stake are marked for eviction, to make the cap be fullfilled. Validators are sorted by increasing `EffectiveBalance` and ties are broken ordering their pub keys alphabetically.
- BeaconKit does not currently support voluntary withdrawals, nor slashing or inactivity leaks. Therefore a validator keeps validating indefinitely.
  - The only case in which a validator may be evicted from the validator set (and its funds returned) is when `ValidatorSetCap` is hit and a validator with greater priority is added (i.e. with larger `EffectiveBalance` or equal `EffectiveBalance` and larger PubKey in alphabetical order).
- Once a validator is marked as active, `CometBFT` consensus will reach it out for block proposals, validations and voting. The higher a validator `EffectiveBalance`, the higher its voting power the frequency it is polled for block proposal.

Now say the validator is marked for exit (currently only as a result of the validator cap being hit), at epoch `M`. Then its funds will be fully withdrawned at epoch `M+1`, since again BeaconKit does not currently enforce a cap on validators churn.
