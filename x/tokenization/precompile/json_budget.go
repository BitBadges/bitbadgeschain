package tokenization

import (
	"github.com/bitbadges/bitbadgeschain/pkg/evmcompat"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// meterJSONInput enforces structural size, depth, token, string, and numeric literal bounds
// before JSON decoding and semantic validation.
func meterJSONInput(ctx sdk.Context, input string) error {
	opts := evmcompat.DefaultJSONBudgetOptions()
	opts.MaxArrayLength = MaxApprovalRanges
	opts.ArrayLimits = map[string]int{
		"addresses": MaxAddressListEntries,
	}
	opts.MaxStringLength = MaxMetadataLength
	opts.GasPerValue = GasPerApprovalField
	return evmcompat.ValidateJSONInput(ctx, input, opts)
}
