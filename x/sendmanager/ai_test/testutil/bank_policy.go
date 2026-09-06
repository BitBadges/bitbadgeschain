package testutil

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (*MockBankKeeper) IsSendEnabledCoins(context.Context, ...sdk.Coin) error { return nil }
func (*MockBankKeeper) BlockedAddr(sdk.AccAddress) bool                       { return false }
