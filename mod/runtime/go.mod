module github.com/berachain/beacon-kit/mod/runtime

go 1.23.0

replace (
	// The following are required to build with the latest version of the cosmos-sdk main branch:
	cosmossdk.io/api => cosmossdk.io/api v0.7.3-0.20240806152830-8fb47b368cd4
	cosmossdk.io/core => cosmossdk.io/core v0.12.1-0.20240806152830-8fb47b368cd4
	cosmossdk.io/core/testing => cosmossdk.io/core/testing v0.0.0-20240806152830-8fb47b368cd4
	cosmossdk.io/x/auth => cosmossdk.io/x/auth v0.0.0-20240806152830-8fb47b368cd4

	cosmossdk.io/x/consensus => cosmossdk.io/x/consensus v0.0.0-20240806152830-8fb47b368cd4
	cosmossdk.io/x/staking => cosmossdk.io/x/staking v0.0.0-20240806152830-8fb47b368cd4
	github.com/cosmos/cosmos-sdk => github.com/berachain/cosmos-sdk v0.46.0-beta2.0.20240808182639-7bdbf06a94f2
)
