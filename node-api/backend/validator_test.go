package backend

import (
	"fmt"
	"testing"

	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/node-api/handlers/beacon/types"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

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
			name:     "success - filter by single id",
			slot:     math.Slot(1),
			ids:      []string{"0"},
			statuses: []string{},
			mockSetup: func(ms *mockState) {
				// Setup mock validator with index 0
				validator := createMockValidator(0, "active")
				ms.On("GetValidators").Return([]*types.Validator{validator}, nil)
				ms.On("ValidatorIndexByPubkey", validator.PublicKey).Return(math.U64(0), nil)
				ms.On("GetBalance", math.U64(0)).Return(math.U64(32e9), nil)
			},
			expectedCount: 1,
			expectedErr:   nil,
		},
		{
			name:     "success - filter by status",
			slot:     math.Slot(1),
			ids:      []string{},
			statuses: []string{"active"},
			mockSetup: func(ms *mockState) {
				// Setup multiple validators with different statuses
				validators := []*types.Validator{
					createMockValidator(0, "active"),
					createMockValidator(1, "pending"),
					createMockValidator(2, "active"),
				}
				ms.On("GetValidators").Return(validators, nil)
				// Setup mock responses for each validator
				for i, v := range validators {
					ms.On("ValidatorIndexByPubkey", v.PublicKey).Return(math.U64(i), nil)
					ms.On("GetBalance", math.U64(i)).Return(math.U64(32e9), nil)
				}
			},
			expectedCount: 2, // Only active validators
			expectedErr:   nil,
		},
		{
			name:     "success - filter by both id and status",
			slot:     math.Slot(1),
			ids:      []string{"0", "1"},
			statuses: []string{"active"},
			mockSetup: func(ms *mockState) {
				validators := []*types.Validator{
					createMockValidator(0, "active"),
					createMockValidator(1, "pending"),
				}
				ms.On("GetValidators").Return(validators, nil)
				for i, v := range validators {
					ms.On("ValidatorIndexByPubkey", v.PublicKey).Return(math.U64(i), nil)
					ms.On("GetBalance", math.U64(i)).Return(math.U64(32e9), nil)
				}
			},
			expectedCount: 1, // Only validator 0 is active
			expectedErr:   nil,
		},
		{
			name:     "error - invalid slot",
			slot:     math.Slot(999999),
			ids:      []string{},
			statuses: []string{},
			mockSetup: func(ms *mockState) {
				ms.On("GetValidators").Return(nil, errors.New("invalid slot"))
			},
			expectedCount: 0,
			expectedErr:   errors.New("invalid slot"),
		},
		{
			name:     "error - invalid validator id",
			slot:     math.Slot(1),
			ids:      []string{"invalid"},
			statuses: []string{},
			mockSetup: func(ms *mockState) {
				ms.On("GetValidators").Return(nil, errors.New("invalid validator id"))
			},
			expectedCount: 0,
			expectedErr:   errors.New("invalid validator id"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock state
			mockState := newMockState(t)
			tt.mockSetup(mockState)

			// Create backend with mock state
			backend := Backend{
				sb: mockState,
			}

			// Call the method
			validators, err := backend.FilteredValidators(tt.slot, tt.ids, tt.statuses)

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
					if len(tt.statuses) > 0 {
						require.Contains(t, tt.statuses, v.Status)
					}
				}
			}

			// Verify all mock expectations were met
			mockState.AssertExpectations(t)
		})
	}
}

// Helper function to create mock validators
func createMockValidator(index uint64, status string) *types.Validator {
	return &types.Validator{
		PublicKey: []byte(fmt.Sprintf("pubkey-%d", index)),
		Status:    status,
		// Add other necessary fields
	}
}

// Mock state interface
type mockState struct {
	mock.Mock
}

func newMockState(t *testing.T) *mockState {
	return &mockState{}
}

// Implement necessary mock methods
func (m *mockState) GetValidators() ([]*types.Validator, error) {
	args := m.Called()
	if v := args.Get(0); v != nil {
		return v.([]*types.Validator), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockState) ValidatorIndexByPubkey(pubkey []byte) (math.U64, error) {
	args := m.Called(pubkey)
	return args.Get(0).(math.U64), args.Error(1)
}

func (m *mockState) GetBalance(index math.U64) (math.U64, error) {
	args := m.Called(index)
	return args.Get(0).(math.U64), args.Error(1)
}
