package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TransferKeeper is the minimal interface we need from the IBC transfer keeper.
type TransferKeeper interface {
	DenomPathFromHash(ctx sdk.Context, denom string) (string, error)
}

// ExpandIBCDenomToFullPath converts a local ibc/<hash> denom to its full
// denom trace path using the transfer keeper's DenomPathFromHash method.
// If the denom is not an ibc/<hash> denom, it returns the original denom unchanged.
func ExpandIBCDenomToFullPath(ctx sdk.Context, denom string, transferKeeper TransferKeeper) (string, error) {
	if transferKeeper == nil {
		return "", fmt.Errorf("transferKeeper is nil")
	}

	// DenomPrefix is "ibc" (no slash); ICS20 hashed denoms are "ibc/<hash>"
	const hashedPrefix = "ibc/"
	if !strings.HasPrefix(denom, hashedPrefix) {
		return denom, nil
	}

	// Use the transfer keeper's built-in method to resolve the hash to full path
	return transferKeeper.DenomPathFromHash(ctx, denom)
}
