# State Processor

## Validators handling

Currently:

- Any validator whose effective balance is above `EjectionBalance` will stay a validator forever, as we have not (yet) implemented withdrawals facilities, nor we do slash.
- Withdrawals are automatically generated only if a validator effective balance goes beyond `MaxEffectiveBalance`. In this case some of the balance is scheduled for withdrawal, just enough to make validator's effective balance equal to `MaxEffectiveBalance`. Since `MaxEffectiveBalance` > `EjectionBalance`, the validator will keep being a validator.
- If a deposit is made for a validator with a balance smaller or equal to `EjectionBalance`, no validator will be created[^1] because of the insufficient balance. However currently the whole deposited balance is **not** scheduled for withdrawal at the next epoch.
- `EffectiveBalance`s are updated one per epoch. Following Eth2.0 specs, the whole validators list is scanned and `EffectiveBalance` is updated only if the difference among `Balance` and `EffectiveBalance` is larger than a (upward or downward) threshold, set considering `EffectiveBalanceIncrement` and hysteresis.
- Validators returned to consensus engine are guaranteed to have their effective balance ranging between `EjectionBalance` excluded (by filtering out state validators with smaller balance) and `MaxEffectiveBalance` included (by validators construction). Moreover only diffs with respect to previous epoch validator set are returned as an optimization measure.

[^1]: Technically a validator is made in the BeaconKit state to track the deposit, but such a validator is never returned to the consensus engine.
