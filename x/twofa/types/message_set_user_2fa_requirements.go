package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	badgestypes "github.com/bitbadges/bitbadgeschain/x/badges/types"
)

const TypeMsgSetUser2FARequirements = "msg_set_user_2fa_requirements"

var _ sdk.Msg = &MsgSetUser2FARequirements{}

func NewMsgSetUser2FARequirements(creator string, mustOwnTokens []*badgestypes.MustOwnTokens, dynamicStoreChallenges []*badgestypes.DynamicStoreChallenge) *MsgSetUser2FARequirements {
	return &MsgSetUser2FARequirements{
		Creator:                creator,
		MustOwnTokens:          mustOwnTokens,
		DynamicStoreChallenges: dynamicStoreChallenges,
	}
}

func (msg *MsgSetUser2FARequirements) Route() string {
	return RouterKey
}

func (msg *MsgSetUser2FARequirements) Type() string {
	return TypeMsgSetUser2FARequirements
}

func (msg *MsgSetUser2FARequirements) GetSigners() []sdk.AccAddress {
	// MustAccAddressFromBech32 panics if address is invalid, which is expected
	// since ValidateBasic() should have already validated the address
	creator := sdk.MustAccAddressFromBech32(msg.Creator)
	return []sdk.AccAddress{creator}
}

func (msg *MsgSetUser2FARequirements) GetSignBytes() []byte {
	bz := AminoCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSetUser2FARequirements) ValidateBasic() error {
	if len(msg.Creator) == 0 {
		return sdkerrors.Wrapf(ErrInvalidRequest, "creator address cannot be empty")
	}
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidRequest, "invalid creator address (%s)", err)
	}

	// Validate each MustOwnTokens requirement using shared validation function
	if err := badgestypes.ValidateMustOwnTokensList(msg.MustOwnTokens); err != nil {
		return err
	}

	// Validate each DynamicStoreChallenge requirement using shared validation function
	if err := badgestypes.ValidateDynamicStoreChallengesList(msg.DynamicStoreChallenges); err != nil {
		return err
	}

	return nil
}
