# How to run the models with the model checker.

- Clone the [tla-bin](https://github.com/pmer/tla-bin) repository.
- Download the nightly build: `$ ./download_or_update_tla.sh --nightly`
- Install: `$ ./install.sh ~/.local`
- The specifications have been already transpiled from PlusCal to TLA+,
  New modifications require running the transpiler again with: `$ pcal -nocfg Model.tla`
- Run each model inside its own directory, for example:
```
$ cd deposits
$ tlc -deadlock -cleanup -workers auto BeaconKitDeposits.tla
```

Speficially regarding the deposit system model, when run without the
`-deadlock`, the model checker will check for deadlocks, which are indeed
possible as confirmed by the eventual output trace: in the circumstance
where the Beacon chain is ahead of the the execution layer's chain,
the deposit store could be non empty but the beacon blocks could not
add any deposit since no new incoming Eth1 block was produced. Or,
otherwise said, stale deposits could only be popped as long as the Eth1
chain advances. This is the reason why Beacon kit goes to extreme length
in keeping the two chains in lockstep.
