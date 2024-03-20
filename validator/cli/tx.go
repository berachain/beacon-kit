package cli

import (
	"fmt"

	beacontypes "github.com/berachain/beacon-kit/runtime/modules/beacon/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/genutil/types"
)

func ProvideMessageValidator() types.MessageValidator {
	return DefaultMessageValidator
}

func DefaultMessageValidator(msgs []sdk.Msg) error {
	if len(msgs) != 1 {
		return fmt.Errorf("unexpected number of GenTx messages; got: %d, expected: 1", len(msgs))
	}
	if _, ok := msgs[0].(*beacontypes.MsgCreateValidatorX); !ok {
		return fmt.Errorf("unexpected GenTx message type; expected: MsgCreateValidator, got: %T", msgs[0])
	}

	if m, ok := msgs[0].(sdk.HasValidateBasic); ok {
		if err := m.ValidateBasic(); err != nil {
			return fmt.Errorf("invalid GenTx '%s': %w", msgs[0], err)
		}
	}

	return nil
}
