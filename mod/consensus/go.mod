module github.com/berachain/beacon-kit/mod/consensus

go 1.22.4

replace (
	// The following are required to build with the latest version of the cosmos-sdk main branch:
	cosmossdk.io/api => cosmossdk.io/api v0.7.3-0.20240623110059-dec2d5583e39
	cosmossdk.io/core => cosmossdk.io/core v0.12.1-0.20240623110059-dec2d5583e39
	cosmossdk.io/core/testing => cosmossdk.io/core/testing v0.0.0-20240623110059-dec2d5583e39
	cosmossdk.io/server/v2 => cosmossdk.io/server/v2 v2.0.0-20240620161400-a6407f411e80
	cosmossdk.io/server/v2/appmanager => cosmossdk.io/server/v2/appmanager v0.0.0-20240623110059-dec2d5583e39
	cosmossdk.io/server/v2/stf => cosmossdk.io/server/v2/stf v0.0.0-20240623110059-dec2d5583e39
	cosmossdk.io/x/auth => cosmossdk.io/x/auth v0.0.0-20240623110059-dec2d5583e39
	cosmossdk.io/x/consensus => cosmossdk.io/x/consensus v0.0.0-20240623110059-dec2d5583e39
	cosmossdk.io/x/staking => cosmossdk.io/x/staking v0.0.0-20240623110059-dec2d5583e39
	github.com/cosmos/cosmos-sdk => github.com/berachain/cosmos-sdk v0.46.0-beta2.0.20240620224206-8693be92d561
)
