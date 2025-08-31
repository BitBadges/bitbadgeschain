package gamm

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	"github.com/bitbadges/bitbadgeschain/x/gamm/keeper"
)

// NewGammProposalHandler is a handler for governance proposals for the GAMM module.
func NewGammProposalHandler(k keeper.Keeper) govtypesv1.Handler {
	return func(ctx sdk.Context, content govtypesv1.Content) error {
		switch c := content.(type) {

		default:
			return errorsmod.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized migration record proposal content type: %T", c)
		}
	}
}
