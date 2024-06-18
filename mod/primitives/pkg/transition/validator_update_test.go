package transition_test

import (
	"testing"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
	"github.com/stretchr/testify/assert"
)

func TestValidatorUpdates_RemoveDuplicates(t *testing.T) {
	pubkey1 := crypto.BLSPubkey{1}
	pubkey2 := crypto.BLSPubkey{2}

	updates := transition.ValidatorUpdates{
		&transition.ValidatorUpdate{Pubkey: pubkey1, EffectiveBalance: math.Gwei(1000)},
		&transition.ValidatorUpdate{Pubkey: pubkey1, EffectiveBalance: math.Gwei(1000)},
		&transition.ValidatorUpdate{Pubkey: pubkey2, EffectiveBalance: math.Gwei(2000)},
	}

	expected := transition.ValidatorUpdates{
		&transition.ValidatorUpdate{Pubkey: pubkey1, EffectiveBalance: math.Gwei(1000)},
		&transition.ValidatorUpdate{Pubkey: pubkey2, EffectiveBalance: math.Gwei(2000)},
	}

	result := updates.RemoveDuplicates()
	assert.Equal(t, expected, result)
}

func TestValidatorUpdates_Sort(t *testing.T) {
	pubkey1 := crypto.BLSPubkey{1}
	pubkey2 := crypto.BLSPubkey{2}
	pubkey3 := crypto.BLSPubkey{3}

	updates := transition.ValidatorUpdates{
		&transition.ValidatorUpdate{Pubkey: pubkey3, EffectiveBalance: math.Gwei(3000)},
		&transition.ValidatorUpdate{Pubkey: pubkey1, EffectiveBalance: math.Gwei(1000)},
		&transition.ValidatorUpdate{Pubkey: pubkey2, EffectiveBalance: math.Gwei(2000)},
	}

	expected := transition.ValidatorUpdates{
		&transition.ValidatorUpdate{Pubkey: pubkey1, EffectiveBalance: math.Gwei(1000)},
		&transition.ValidatorUpdate{Pubkey: pubkey2, EffectiveBalance: math.Gwei(2000)},
		&transition.ValidatorUpdate{Pubkey: pubkey3, EffectiveBalance: math.Gwei(3000)},
	}

	result := updates.Sort()
	assert.Equal(t, expected, result)
}
