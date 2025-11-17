package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgUpdateRateLimit{}

// ValidateBasic does a sanity check on the provided data.
func (m *MsgUpdateRateLimit) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}

	if err := m.RateLimit.Validate(); err != nil {
		return errorsmod.Wrap(err, "invalid rate limit config")
	}

	return nil
}

