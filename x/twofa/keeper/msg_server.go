package keeper

import (
	"context"
	"fmt"

	sdkerrors "cosmossdk.io/errors"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	twofatypes "github.com/bitbadges/bitbadgeschain/x/twofa/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgServer defines the interface for the twofa message server
type MsgServer interface {
	SetUser2FARequirements(context.Context, *twofatypes.MsgSetUser2FARequirements) (*twofatypes.MsgSetUser2FARequirementsResponse, error)
}

var _ twofatypes.MsgServer = msgServer{}

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
func NewMsgServerImpl(keeper Keeper) twofatypes.MsgServer {
	return &msgServer{Keeper: keeper}
}

// SetUser2FARequirements allows a user to set their 2FA token requirements for transaction authorization.
func (k msgServer) SetUser2FARequirements(goCtx context.Context, msg *twofatypes.MsgSetUser2FARequirements) (*twofatypes.MsgSetUser2FARequirementsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate the message
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	// Validate the creator address
	if err := types.ValidateAddress(msg.Creator, false); err != nil {
		return nil, sdkerrors.Wrap(err, "Invalid creator address")
	}

	// Validate each MustOwnTokens requirement
	for idx, mustOwnToken := range msg.MustOwnTokens {
		if mustOwnToken == nil {
			return nil, sdkerrors.Wrapf(twofatypes.ErrInvalidRequest, "MustOwnTokens requirement at index %d is nil", idx)
		}

		// Validate collection exists
		_, found := k.badgesKeeper.GetCollectionFromStore(ctx, mustOwnToken.CollectionId)
		if !found {
			return nil, sdkerrors.Wrapf(twofatypes.ErrInvalidRequest, "Collection %s not found for requirement at index %d", mustOwnToken.CollectionId.String(), idx)
		}

		if mustOwnToken.AmountRange == nil {
			return nil, sdkerrors.Wrapf(twofatypes.ErrInvalidRequest, "AmountRange is required for requirement at index %d", idx)
		}

		// Validate ownership times if not using override
		if !mustOwnToken.OverrideWithCurrentTime && len(mustOwnToken.OwnershipTimes) == 0 {
			return nil, sdkerrors.Wrapf(twofatypes.ErrInvalidRequest, "OwnershipTimes must be set or OverrideWithCurrentTime must be true for requirement at index %d", idx)
		}
	}

	// Validate each DynamicStoreChallenge requirement
	for idx, challenge := range msg.DynamicStoreChallenges {
		if challenge == nil {
			return nil, sdkerrors.Wrapf(twofatypes.ErrInvalidRequest, "DynamicStoreChallenge at index %d is nil", idx)
		}

		if challenge.StoreId.IsNil() || challenge.StoreId.IsZero() {
			return nil, sdkerrors.Wrapf(twofatypes.ErrInvalidRequest, "DynamicStoreChallenge requirement at index %d: StoreId cannot be nil or zero", idx)
		}
		// Check if dynamic store exists
		_, found := k.badgesKeeper.GetDynamicStoreFromStore(ctx, challenge.StoreId)
		if !found {
			return nil, sdkerrors.Wrapf(twofatypes.ErrInvalidRequest, "Dynamic store %s not found for challenge at index %d", challenge.StoreId.String(), idx)
		}
		// ownershipCheckParty validation is handled in the ante handler's determinePartyToCheckForDynamicStore
	}

	// Create the requirements object
	requirements := &twofatypes.User2FARequirements{
		MustOwnTokens:          msg.MustOwnTokens,
		DynamicStoreChallenges: msg.DynamicStoreChallenges,
	}

	// Set the 2FA requirements
	if err := k.SetUser2FARequirementsInStore(ctx, msg.Creator, requirements); err != nil {
		return nil, sdkerrors.Wrap(err, "Failed to store 2FA requirements")
	}

	// Emit event
	EmitMessageAndIndexerEvents(ctx,
		sdk.NewAttribute(sdk.AttributeKeyModule, "twofa"),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
		sdk.NewAttribute("msg_type", "set_user_2fa_requirements"),
		sdk.NewAttribute("must_own_tokens_count", fmt.Sprintf("%d", len(msg.MustOwnTokens))),
		sdk.NewAttribute("dynamic_store_challenges_count", fmt.Sprintf("%d", len(msg.DynamicStoreChallenges))),
	)

	return &twofatypes.MsgSetUser2FARequirementsResponse{}, nil
}
