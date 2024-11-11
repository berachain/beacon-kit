# State Processor

## Validators handling

Currently:

- Any validator whose effective balance is above `EjectionBalance` will stay a validator forever, as we have not (yet) implemented withdrawals facilities.
- Withdrawals are automatically generated only if a validator effective balance goes beyond `MaxEffectiveBalance`. In this case enough balance is scheduled for withdrawal, just enough to make validator's effective balance equal to `MaxEffectiveBalance`. Since `MaxEffectiveBalance` > `EjectionBalance`, the validator will keep being a validator.
- If a deposit is made for a validator with a balance smaller or equal to `EjectionBalance`, no validator will be created[^1] because of the insufficient balance. However currently the whole deposited balance is **not** scheduled for withdrawal at the next epoch.
- Validators returned to consensus engine are guaranteed to have their effective balance ranging between `EjectionBalance` excluded (by filtering out state validators with smaller balance) and `MaxEffectiveBalance` included (by valdiators construction). Moreover only diffs with respect to previous epoch validator set are returned as an optimization measure.

[^1]: Technically a validator is made in the BeaconKit state to track the deposit, but such a validator is never returned to the consensus engine. Moreover the deposit should be evicted at the next epoch.
