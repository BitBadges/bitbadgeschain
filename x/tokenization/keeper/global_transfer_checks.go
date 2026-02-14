package keeper

import (
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GlobalTransferChecker defines the interface for checking transfers at a global level.
// These checks run before HandleTransfer() and can validate transfers based on
// from, to, initiatedBy, collection, transfer balances, and memo.
// All implementations should return (deterministicErrorMsg, error) where:
// - deterministicErrorMsg is a user-friendly error message if the check fails
// - error is nil if the check passes, or an error if the check fails
type GlobalTransferChecker interface {
	Name() string
	Check(
		ctx sdk.Context,
		from string,
		to string,
		initiatedBy string,
		collection *types.TokenCollection,
		transferBalances []*types.Balance,
		memo string,
	) (detErrMsg string, err error)
}
