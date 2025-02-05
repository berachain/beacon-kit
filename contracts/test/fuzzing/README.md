## Overview

Berachain engaged [Perimeter](https://www.perimetersec.io) for two in-depth invariant suite for the proof verification contracts of Beacon Kit. The primary goal of this engagement was to rigorously test proof validation and to establish a robust foundation that can be easily extended in the future. To maintain consistency with the Beacon Kit code, we utilized many of its components, written in Golang, for the proof generation script. These proofs were then verified using the proof validation Solidity contracts.

The invariants test for false positive cases, false negative cases, and potential denial-of-service in the proof validation contracts.

This engagement was conducted by Lead Fuzzing Specialist [0xScourgedev](https://twitter.com/0xScourgedev) and Lead Fuzzing Specialist [Rappie](https://twitter.com/rappie_eth). The fuzzing suite was successfully delivered upon the engagement's conclusion.

## Contents

This fuzzing suite was originally created for the commit hash `dd024c5b196afe43cf871b5f825a7e371fc4542e`, and expanded upon for the scope below.

The primary goal of this engagement was to rigorously test proof validation and to establish a robust foundation that can be easily extended in the future. To maintain consistency with the Beacon Kit code, we utilized many of its components, written in Golang, for the proof generation script. All aspects of the beacon state and beacon header were randomized.

The Solidity contracts used for proof validation were obtained from a different repository, which can be accessed [here](https://github.com/berachain/contracts-monorepo/tree/84af1993df0c43265709bdbaaa30c5b3f6f6572b). Additionally, a minimal `BeaconVerifier` contract was developed for the fuzzing suite, as the visibility of the validation functions is internal.

The contracts that were utilized in the invariant suite can be found under `test/fuzzing/base`, `test/fuzzing/libraries` and `test/fuzzing/pol`.

## Setup

1. Installing Medusa

   a. Install Medusa, follow the steps here: [Installation Guide](https://github.com/crytic/medusa/blob/master/docs/src/getting_started/installation.md) using the latest master branch

2. To fuzz all Beacon Kit invariants using Medusa, run the command:

   ```
   cd contracts
   medusa fuzz
   ```

3. To run only the proof generation script, run the command:

   ```
   go run ./test/fuzzing/script/proof_gen -sel=<0 | 1>
   ```

   Flags:

   `-seed=<SEED>`
   The seed is an int64 value to seed the randomization of the proof generation. The default value is 0.

   `-valLen=<VALUE_LENGTH>`
   The number of validators to use. The default value is 5.

   `-hardcoded`
   This flag will use hardcoded values in the proof generation script. It will override the other flags

   `-debug`
   This flag will add additional text to help visualize the outputs

## Scope

Repo: [https://github.com/berachain/beacon-kit](https://github.com/berachain/beacon-kit)

Branch: `main`

Commit: `86b41a4d292028019458f921c738527cd095`

```
src
├── base
│   └── IStakingRewardsErrors.sol
├── libraries
│   ├── BeaconRoots.sol
│   ├── SSZ.sol
│   └── Utils.sol
└── pol
    ├── BeaconRootsHelper.sol
    ├── BeaconVerifier.sol
    └── interfaces
         └── IPOLErrors.sol
```
