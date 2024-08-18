package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	badgetypes "bitbadgeschain/x/badges/types"
)

var _ sdk.Msg = &MsgCreateProposal{}

func NewMsgCreateProposal(creator string, parties []*Parties, validTimes []*UintRange, creatorMustFinalize bool, anyoneCanFinalize bool) *MsgCreateProposal {
	return &MsgCreateProposal{
		Creator:             creator,
		Parties:             parties,
		ValidTimes:          validTimes,
		CreatorMustFinalize: creatorMustFinalize,
		AnyoneCanFinalize:   anyoneCanFinalize,
	}
}

func (msg *MsgCreateProposal) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if err := ValidateParties(msg.Parties); err != nil {
		return err
	}

	if len(msg.Parties) == 0 {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid parties (%s)", msg.Parties)
	}

	//ensure no duplicate creators
	creators := make(map[string]bool)
	for _, party := range msg.Parties {
		if creators[party.Creator] {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "duplicate creator (%s)", party.Creator)
		}
		creators[party.Creator] = true
	}

	//Ensure all accepts are false
	for _, party := range msg.Parties {
		if party.Accepted {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid party accepted status")
		}
	}

	castedRanges := CastUintRanges(msg.ValidTimes)
	if err := badgetypes.ValidateRangesAreValid(castedRanges, false, true); err != nil {
		return err
	}

	return nil
}

func ValidateParties(parties []*Parties) error {
	for _, party := range parties {
		if party.Creator == "" {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid party creator address (%s)", party.Creator)
		}

		if party.MsgsToExecute == nil || len(party.MsgsToExecute) == 0 {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid party msgs to execute (%s)", party.MsgsToExecute)
		}
	}

	return nil
}
