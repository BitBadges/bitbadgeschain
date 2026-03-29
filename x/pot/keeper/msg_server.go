package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/pot/types"
)

type msgServer struct {
	Keeper
}

// Ensure msgServer implements the MsgServer interface.
var _ types.MsgServer = msgServer{}

// NewMsgServerImpl returns an implementation of the MsgServer interface.
func NewMsgServerImpl(keeper Keeper) msgServer {
	return msgServer{Keeper: keeper}
}

// UpdateParams handles MsgUpdateParams.
func (k msgServer) UpdateParams(goCtx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	if k.GetAuthority() != msg.Authority {
		return nil, types.ErrInvalidSigner.Wrapf("invalid authority; expected %s, got %s", k.GetAuthority(), msg.Authority)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Read old params before updating so we can detect credential collection changes.
	oldParams := k.GetParams(ctx)

	if err := k.SetParams(ctx, msg.Params); err != nil {
		return nil, err
	}

	// Fix 4: If CredentialCollectionId changed, clear all compliance-jailed entries.
	// All validators will be re-evaluated on the next EndBlocker with the new collection.
	if oldParams.CredentialCollectionId != msg.Params.CredentialCollectionId {
		for _, addrBytes := range k.GetAllComplianceJailed(ctx) {
			consAddr := sdk.ConsAddress(addrBytes)
			k.RemoveComplianceJailed(ctx, consAddr)

			// Enable if safe (respects tombstone/slashing for staking, always true for PoA).
			if k.validatorSet.CanSafelyEnable(ctx, consAddr) {
				k.safeEnable(ctx, consAddr, consAddr.String())
			}
		}
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"pot_params_updated",
			sdk.NewAttribute("authority", msg.Authority),
			sdk.NewAttribute("mode", msg.Params.Mode),
		),
	)

	return &types.MsgUpdateParamsResponse{}, nil
}
