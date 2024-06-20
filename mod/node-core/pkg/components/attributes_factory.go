package components

import (
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/payload/pkg/attributes"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
)

// ProvideAttributesFactory provides an AttributesFactory for the client.
func ProvideAttributesFactory(
	chainSpec primitives.ChainSpec,
	logger log.Logger[any],
	suggestedFeeRecipient common.ExecutionAddress,
) (*AttributesFactory, error) {
	return attributes.NewAttributesFactory[BeaconState, *Withdrawal](
		chainSpec,
		logger,
		suggestedFeeRecipient,
	), nil
}
