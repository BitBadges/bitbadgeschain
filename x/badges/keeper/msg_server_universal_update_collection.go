package keeper

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"slices"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

const (
	// DerivationKeyLength represents the length of the derivation key for account generation
	DerivationKeyLength = 8
	// MaxUint64Value represents the maximum value for uint64
	MaxUint64Value = math.MaxUint64
)

// Legacy function that is all-inclusive (creates and updates)
func (k msgServer) UniversalUpdateCollection(goCtx context.Context, msg *types.MsgUniversalUpdateCollection) (*types.MsgUniversalUpdateCollectionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	err := msg.CheckAndCleanMsg(ctx, true)
	if err != nil {
		return nil, err
	}

	collection := &types.TokenCollection{}
	if msg.CollectionId.Equal(sdkmath.NewUint(NewCollectionId)) {
		//Creation case
		nextCollectionId := k.GetNextCollectionId(ctx)
		if err := k.IncrementNextCollectionId(ctx); err != nil {
			return nil, err
		}

		// From cosmos SDK x/group module
		// Generate account address of collection
		var accountAddr sdk.AccAddress
		for {
			derivationKey := make([]byte, DerivationKeyLength)
			binary.BigEndian.PutUint64(derivationKey, nextCollectionId.Uint64())

			ac, err := authtypes.NewModuleCredential(types.ModuleName, AccountGenerationPrefix, derivationKey)
			if err != nil {
				return nil, err
			}
			//generate the address from the credential
			accountAddr = sdk.AccAddress(ac.Address())

			break
		}

		collection = &types.TokenCollection{
			CollectionId:          nextCollectionId,
			CollectionPermissions: &types.CollectionPermissions{},
			DefaultBalances:       msg.DefaultBalances,
			CreatedBy:             msg.Creator,
			// Default manager is the creator if not updating manager timeline
			ManagerTimeline: []*types.ManagerTimeline{
				{
					Manager: msg.Creator,
					TimelineTimes: []*types.UintRange{
						{
							Start: sdkmath.NewUint(1),
							End:   sdkmath.NewUint(MaxUint64Value),
						},
					},
				},
			},
			MintEscrowAddress: accountAddr.String(),
			Invariants:        msg.Invariants,
		}
	} else {
		//Update case
		found := false
		collection, found = k.GetCollectionFromStore(ctx, msg.CollectionId)
		if !found {
			return nil, ErrCollectionNotExists
		}
	}

	//Check must be manager
	err = k.UniversalValidate(ctx, collection, UniversalValidationParams{
		Creator:       msg.Creator,
		MustBeManager: true,
	})
	if err != nil {
		return nil, err
	}

	//Other cases:
	//previouslyArchived && not stillArchived - we have just unarchived the collection
	//not previouslyArchived && stillArchived - we have just archived the collection (all TXs moving forward will fail, but we allow this one)
	//not previouslyArchived && not stillArchived - unarchived before and now so we allow
	previouslyArchived := types.GetIsArchived(ctx, collection)
	if msg.UpdateIsArchivedTimeline {
		if err := k.ValidateIsArchivedUpdate(ctx, collection.IsArchivedTimeline, msg.IsArchivedTimeline, collection.CollectionPermissions.CanArchiveCollection); err != nil {
			return nil, err
		}
		collection.IsArchivedTimeline = msg.IsArchivedTimeline
	}
	stillArchived := types.GetIsArchived(ctx, collection)
	if previouslyArchived && stillArchived {
		return nil, ErrCollectionIsArchived
	}

	if msg.UpdateCollectionApprovals {
		if err := k.ValidateCollectionApprovalsUpdate(ctx, collection, collection.CollectionApprovals, msg.CollectionApprovals, collection.CollectionPermissions.CanUpdateCollectionApprovals); err != nil {
			return nil, err
		}

		// Create a map of existing approvals for quick lookup
		existingApprovals := make(map[string]*types.CollectionApproval)
		for _, approval := range collection.CollectionApprovals {
			existingApprovals[approval.ApprovalId] = approval
		}

		// Only increment versions for approvals that have changed
		newApprovalsWithVersion := []*types.CollectionApproval{}
		for _, newApproval := range msg.CollectionApprovals {
			existingApproval, exists := existingApprovals[newApproval.ApprovalId]
			// Only increment version if approval is new or changed
			if !exists || existingApproval.String() != newApproval.String() {
				newVersion := k.IncrementApprovalVersion(ctx, collection.CollectionId, "collection", "", newApproval.ApprovalId)
				newApproval.Version = newVersion
			} else {
				// Keep existing version if approval hasn't changed
				newApproval.Version = existingApproval.Version
			}
			newApprovalsWithVersion = append(newApprovalsWithVersion, newApproval)
		}
		collection.CollectionApprovals = newApprovalsWithVersion
	}

	if msg.UpdateCollectionMetadataTimeline {
		if err := k.ValidateCollectionMetadataUpdate(ctx, collection.CollectionMetadataTimeline, msg.CollectionMetadataTimeline, collection.CollectionPermissions.CanUpdateCollectionMetadata); err != nil {
			return nil, err
		}
		collection.CollectionMetadataTimeline = msg.CollectionMetadataTimeline
	}

	if msg.UpdateTokenMetadataTimeline {
		if err := k.ValidateTokenMetadataUpdate(ctx, collection.TokenMetadataTimeline, msg.TokenMetadataTimeline, collection.CollectionPermissions.CanUpdateTokenMetadata); err != nil {
			return nil, err
		}
		collection.TokenMetadataTimeline = msg.TokenMetadataTimeline
	}

	if msg.UpdateManagerTimeline {
		if err := k.ValidateManagerUpdate(ctx, collection.ManagerTimeline, msg.ManagerTimeline, collection.CollectionPermissions.CanUpdateManager); err != nil {
			return nil, err
		}
		collection.ManagerTimeline = msg.ManagerTimeline
	}

	if msg.UpdateStandardsTimeline {
		if err := k.ValidateStandardsUpdate(ctx, collection.StandardsTimeline, msg.StandardsTimeline, collection.CollectionPermissions.CanUpdateStandards); err != nil {
			return nil, err
		}
		collection.StandardsTimeline = msg.StandardsTimeline
	}

	if msg.UpdateCustomDataTimeline {
		if err := k.ValidateCustomDataUpdate(ctx, collection.CustomDataTimeline, msg.CustomDataTimeline, collection.CollectionPermissions.CanUpdateCustomData); err != nil {
			return nil, err
		}
		collection.CustomDataTimeline = msg.CustomDataTimeline
	}

	if msg.UpdateValidTokenIds {
		collection, err = k.CreateTokens(ctx, collection, msg.ValidTokenIds)
		if err != nil {
			return nil, err
		}
	}

	if msg.UpdateCollectionPermissions {
		err = k.ValidatePermissionsUpdate(ctx, collection.CollectionPermissions, msg.CollectionPermissions)
		if err != nil {
			return nil, err
		}
		collection.CollectionPermissions = msg.CollectionPermissions
	}

	if len(msg.MintEscrowCoinsToTransfer) > 0 {
		from := sdk.MustAccAddressFromBech32(msg.Creator)
		to := sdk.MustAccAddressFromBech32(collection.MintEscrowAddress)

		for _, coin := range msg.MintEscrowCoinsToTransfer {
			allowedDenoms := k.GetParams(ctx).AllowedDenoms
			if !slices.Contains(allowedDenoms, coin.Denom) {
				return nil, fmt.Errorf("denom %s is not allowed", coin.Denom)
			}

			err = k.bankKeeper.SendCoins(ctx, from, to, sdk.NewCoins(*coin))
			if err != nil {
				return nil, err
			}
		}
	}

	if len(msg.CosmosCoinWrapperPathsToAdd) > 0 {
		pathsToAdd := make([]*types.CosmosCoinWrapperPath, len(msg.CosmosCoinWrapperPathsToAdd))
		for i, path := range msg.CosmosCoinWrapperPathsToAdd {
			var accountAddr sdk.AccAddress
			for {
				fullPath := path.Denom
				fullPathBytes := []byte(fullPath)

				ac, err := authtypes.NewModuleCredential(types.ModuleName, WrapperPathGenerationPrefix, fullPathBytes)
				if err != nil {
					return nil, err
				}
				//generate the address from the credential
				accountAddr = sdk.AccAddress(ac.Address())

				break
			}

			pathsToAdd[i] = &types.CosmosCoinWrapperPath{
				Address:                        accountAddr.String(),
				Denom:                          path.Denom,
				Balances:                       path.Balances,
				Symbol:                         path.Symbol,
				DenomUnits:                     path.DenomUnits,
				AllowOverrideWithAnyValidToken: path.AllowOverrideWithAnyValidToken,
				AllowCosmosWrapping:            path.AllowCosmosWrapping,
			}
		}

		collection.CosmosCoinWrapperPaths = append(collection.CosmosCoinWrapperPaths, pathsToAdd...)
	}

	// Ensure no duplicate denom paths
	denomPaths := make(map[string]bool)
	for _, path := range collection.CosmosCoinWrapperPaths {
		if _, ok := denomPaths[path.Denom]; ok {
			return nil, fmt.Errorf("duplicate ibc wrapper path denom: %s", path.Denom)
		}
		denomPaths[path.Denom] = true
	}

	// Ensure no duplicate symbols (including base symbol and denom unit symbols)
	symbolPaths := make(map[string]bool)
	for _, path := range collection.CosmosCoinWrapperPaths {
		// Check the main path symbol
		if path.Symbol != "" {
			if _, ok := symbolPaths[path.Symbol]; ok {
				return nil, fmt.Errorf("duplicate ibc wrapper path symbol: %s", path.Symbol)
			}
			symbolPaths[path.Symbol] = true
		}

		// Check denom unit symbols
		for _, denomUnit := range path.DenomUnits {
			if denomUnit.Symbol != "" {
				if _, ok := symbolPaths[denomUnit.Symbol]; ok {
					return nil, fmt.Errorf("duplicate denom unit symbol: %s", denomUnit.Symbol)
				}
				symbolPaths[denomUnit.Symbol] = true
			}
		}

		// Validate that only one denom unit per path has isDefaultDisplay set to true
		defaultDisplayCount := 0
		decimalsSet := make(map[string]bool)
		for _, denomUnit := range path.DenomUnits {
			if denomUnit.IsDefaultDisplay {
				defaultDisplayCount++
			}

			// Check that decimals is not 0
			if denomUnit.Decimals.IsZero() {
				return nil, fmt.Errorf("denom unit decimals cannot be 0")
			}

			// Check for duplicate decimals
			decimalsStr := denomUnit.Decimals.String()
			if _, ok := decimalsSet[decimalsStr]; ok {
				return nil, fmt.Errorf("duplicate denom unit decimals: %s", decimalsStr)
			}
			decimalsSet[decimalsStr] = true
		}

		if defaultDisplayCount > 1 {
			return nil, fmt.Errorf("only one denom unit per path can have isDefaultDisplay set to true, found %d", defaultDisplayCount)
		}
	}

	if err := k.SetCollectionInStore(ctx, collection); err != nil {
		return nil, err
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	EmitMessageAndIndexerEvents(ctx,
		sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
		sdk.NewAttribute("msg_type", "universal_update_collection"),
		sdk.NewAttribute("msg", string(msgBytes)),
		sdk.NewAttribute("collectionId", fmt.Sprint(collection.CollectionId)),
	)

	return &types.MsgUniversalUpdateCollectionResponse{
		CollectionId: collection.CollectionId,
	}, nil
}
