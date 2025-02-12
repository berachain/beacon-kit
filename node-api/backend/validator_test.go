package backend_test

import (
	"testing"

	"github.com/berachain/beacon-kit/config/spec"
	consensustypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/node-api/backend"
	"github.com/berachain/beacon-kit/primitives/bytes"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/transition"
	"github.com/berachain/beacon-kit/state-transition/core/state"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockState struct {
	mock.Mock
}

func (m *mockState) GetValidators() ([]*consensustypes.Validator, error) {
	args := m.Called()
	if v := args.Get(0); v != nil {
		return v.([]*consensustypes.Validator), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockState) ValidatorIndexByPubkey(pubkey bytes.B48) (math.U64, error) {
	args := m.Called(pubkey)
	return args.Get(0).(math.U64), args.Error(1)
}

func (m *mockState) GetBalance(index math.U64) (math.U64, error) {
	args := m.Called(index)
	return args.Get(0).(math.U64), args.Error(1)
}

func (m *mockState) SlotToEpoch(slot math.Slot) (math.Epoch, error) {
	args := m.Called(slot)
	return args.Get(0).(math.Epoch), args.Error(1)
}

func (m *mockState) GetSlot() (math.U64, error) {
	args := m.Called()
	return args.Get(0).(math.U64), args.Error(1)
}

func (m *mockState) ProcessSlots(st *state.StateDB, slot math.U64) (transition.ValidatorUpdates, error) {
	args := m.Called(st, slot)
	return args.Get(0).(transition.ValidatorUpdates), args.Error(1)
}

func TestFilteredValidators(t *testing.T) {
	tests := []struct {
		name          string
		slot          math.Slot
		ids           []string
		statuses      []string
		mockSetup     func(*mockState)
		expectedCount int
		expectedErr   error
	}{
		{
			name:     "success - filter by pubkey",
			slot:     math.Slot(1),
			ids:      []string{"0x93247f2209abcacf57b75a51dafae777f9dd38bc7053d1af526f220a7489a6d3a2753e5f3e8b1cfe39b56f43611df74a"},
			statuses: []string{},
			mockSetup: func(ms *mockState) {
				// Mock state
				ms.On("GetSlot").Return(math.U64(1), nil)
				ms.On("ProcessSlots", mock.Anything, mock.Anything).Return(transition.ValidatorUpdates{}, nil)

				pubkey := bytes.B48{0x93, 0x24, 0x7f, 0x22, 0x09, 0xab, 0xca, 0xcf}
				validator := &consensustypes.Validator{Pubkey: pubkey}
				ms.On("GetValidators").Return([]*consensustypes.Validator{validator}, nil)
				ms.On("ValidatorIndexByPubkey", pubkey).Return(math.U64(0), nil)
				ms.On("GetBalance", math.U64(0)).Return(math.U64(32e9), nil)
				ms.On("SlotToEpoch", math.Slot(1)).Return(math.Epoch(0), nil)
			},
			expectedCount: 1,
			expectedErr:   nil,
		},
		{
			name:     "success - filter by index",
			slot:     math.Slot(1),
			ids:      []string{"0"},
			statuses: []string{},
			mockSetup: func(ms *mockState) {
				ms.On("GetSlot").Return(math.U64(1), nil)
				ms.On("ProcessSlots", mock.Anything, mock.Anything).Return(transition.ValidatorUpdates{}, nil)

				validator := &consensustypes.Validator{Pubkey: bytes.B48{1}}
				ms.On("GetValidators").Return([]*consensustypes.Validator{validator}, nil)
				ms.On("ValidatorIndexByPubkey", validator.Pubkey).Return(math.U64(0), nil)
				ms.On("GetBalance", math.U64(0)).Return(math.U64(32e9), nil)
				ms.On("SlotToEpoch", math.Slot(1)).Return(math.Epoch(0), nil)
			},
			expectedCount: 1,
			expectedErr:   nil,
		},
		{
			name:     "error - invalid slot",
			slot:     math.Slot(999999),
			ids:      []string{},
			statuses: []string{},
			mockSetup: func(ms *mockState) {
				ms.On("GetSlot").Return(math.U64(0), errors.New("invalid slot"))
			},
			expectedCount: 0,
			expectedErr:   errors.New("failed to get state from slot"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock state
			mockState := &mockState{}
			tt.mockSetup(mockState)

			cs, err := spec.DevnetChainSpec()
			require.NoError(t, err)

			// Create backend using New
			b := backend.New(nil, cs, mockState)
			b.AttachQueryBackend(mockState)

			// Call the method
			validators, err := b.FilteredValidators(tt.slot, tt.ids, tt.statuses)

			// Check results
			if tt.expectedErr != nil {
				require.Error(t, err)
				require.ErrorContains(t, err, tt.expectedErr.Error())
				require.Nil(t, validators)
			} else {
				require.NoError(t, err)
				require.NotNil(t, validators)
				require.Len(t, validators, tt.expectedCount)

				// Additional validation for successful cases
				for _, v := range validators {
					require.NotNil(t, v)
					require.NotZero(t, v.Balance)
				}
			}

			// Verify all mock expectations were met
			mockState.AssertExpectations(t)
		})
	}
} 