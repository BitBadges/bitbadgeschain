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

// createMintForbiddenPermission creates a collection approval permission that forbids
// updates to approvals with fromList matching Mint when cosmosCoinBackedPath is set.
// This should be prepended to CanUpdateCollectionApprovals for first-match behavior.
func createMintForbiddenPermission() *types.CollectionApprovalPermission {
	fullRanges := []*types.UintRange{
		{
			Start: sdkmath.NewUint(1),
			End:   sdkmath.NewUint(MaxUint64Value),
		},
	}

	return &types.CollectionApprovalPermission{
		FromListId:                types.MintAddress,
		ToListId:                  "All",
		InitiatedByListId:         "All",
		ApprovalId:                "All",
		TransferTimes:             fullRanges,
		TokenIds:                  fullRanges,
		OwnershipTimes:            fullRanges,
		PermanentlyPermittedTimes: []*types.UintRange{},
		PermanentlyForbiddenTimes: fullRanges,
	}
}

// ensureMintForbiddenPermission ensures that the Mint forbidden permission is prepended
// to CanUpdateCollectionApprovals when cosmosCoinBackedPath is set.
func ensureMintForbiddenPermission(permissions *types.CollectionPermissions, hasCosmosCoinBackedPath bool) {
	if !hasCosmosCoinBackedPath {
		return
	}

	mintForbiddenPerm := createMintForbiddenPermission()

	// Check if the permission already exists (compare key fields)
	alreadyExists := false
	for _, perm := range permissions.CanUpdateCollectionApprovals {
		if perm.FromListId == types.MintAddress &&
			perm.ToListId == "All" &&
			perm.InitiatedByListId == "All" &&
			perm.ApprovalId == "All" {
			alreadyExists = true
			break
		}
	}

	// If it doesn't exist, prepend it
	if !alreadyExists {
		permissions.CanUpdateCollectionApprovals = append([]*types.CollectionApprovalPermission{mintForbiddenPerm}, permissions.CanUpdateCollectionApprovals...)
	}
}

// generatePathAddress generates an address from a path string using the given prefix
func generatePathAddress(pathString string, prefix []byte) (sdk.AccAddress, error) {
	fullPathBytes := []byte(pathString)
	ac, err := authtypes.NewModuleCredential(types.ModuleName, prefix, fullPathBytes)
	if err != nil {
		return nil, err
	}
	return sdk.AccAddress(ac.Address()), nil
}

// setReservedProtocolAddressForPath sets a path address as a reserved protocol address
func (k msgServer) setReservedProtocolAddressForPath(ctx sdk.Context, address string, pathType string) error {
	err := k.SetReservedProtocolAddressInStore(ctx, address, true)
	if err != nil {
		return fmt.Errorf("failed to set %s path address as reserved protocol: %w", pathType, err)
	}
	return nil
}

// setAutoApproveFlagsForPathAddress sets auto-approve flags for a path address
func (k msgServer) setAutoApproveFlagsForPathAddress(ctx sdk.Context, collection *types.TokenCollection, pathAddress string, pathType string) error {
	currBalances, _ := k.GetBalanceOrApplyDefault(ctx, collection, pathAddress)

	alreadyAutoApprovedAllIncomingTransfers := currBalances.AutoApproveAllIncomingTransfers
	alreadyAutoApprovedSelfInitiatedOutgoingTransfers := currBalances.AutoApproveSelfInitiatedOutgoingTransfers
	alreadyAutoApprovedSelfInitiatedIncomingTransfers := currBalances.AutoApproveSelfInitiatedIncomingTransfers

	autoApprovedAll := alreadyAutoApprovedAllIncomingTransfers && alreadyAutoApprovedSelfInitiatedOutgoingTransfers && alreadyAutoApprovedSelfInitiatedIncomingTransfers

	if !autoApprovedAll {
		// We override all approvals to be default allowed
		// Incoming - All, no matter what
		// Outgoing - Self-initiated
		// Incoming - Self-initiated
		currBalances.AutoApproveAllIncomingTransfers = true
		currBalances.AutoApproveSelfInitiatedOutgoingTransfers = true
		currBalances.AutoApproveSelfInitiatedIncomingTransfers = true
		err := k.SetBalanceForAddress(ctx, collection, pathAddress, currBalances)
		if err != nil {
			return fmt.Errorf("failed to set auto-approve flags for %s path address: %w", pathType, err)
		}
	}
	return nil
}

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
		}

		// Convert InvariantsAddObject to CollectionInvariants if present
		if msg.Invariants != nil {
			collection.Invariants = &types.CollectionInvariants{
				NoCustomOwnershipTimes:      msg.Invariants.NoCustomOwnershipTimes,
				MaxSupplyPerId:              msg.Invariants.MaxSupplyPerId,
				NoForcefulPostMintTransfers: msg.Invariants.NoForcefulPostMintTransfers,
				DisablePoolCreation:         msg.Invariants.DisablePoolCreation,
			}

			// Handle cosmos coin backed path - generate address
			if msg.Invariants.CosmosCoinBackedPath != nil {
				pathAddObject := msg.Invariants.CosmosCoinBackedPath
				backedAccountAddr, err := generatePathAddress(pathAddObject.IbcDenom, BackedPathGenerationPrefix)
				if err != nil {
					return nil, err
				}

				backedPath := &types.CosmosCoinBackedPath{
					Address:   backedAccountAddr.String(),
					IbcDenom:  pathAddObject.IbcDenom,
					Balances:  pathAddObject.Balances,
					IbcAmount: pathAddObject.IbcAmount,
				}

				// Auto-set the cosmoscoinbacked path address as a reserved protocol address
				if err := k.setReservedProtocolAddressForPath(ctx, backedAccountAddr.String(), "cosmoscoinbacked"); err != nil {
					return nil, err
				}

				collection.Invariants.CosmosCoinBackedPath = backedPath
			}
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
		// Create a temporary collection with invariants from the message for validation
		// This ensures invariants are checked even if they haven't been set on the collection yet
		tempCollection := *collection
		if msg.Invariants != nil {
			if tempCollection.Invariants == nil {
				tempCollection.Invariants = &types.CollectionInvariants{}
			}
			tempCollection.Invariants.NoCustomOwnershipTimes = msg.Invariants.NoCustomOwnershipTimes
			tempCollection.Invariants.MaxSupplyPerId = msg.Invariants.MaxSupplyPerId
			tempCollection.Invariants.NoForcefulPostMintTransfers = msg.Invariants.NoForcefulPostMintTransfers
			// Note: CosmosCoinBackedPath requires address generation, so we skip it here
			// It will be validated separately in validateCollectionBeforeStore
		}

		if err := k.ValidateCollectionApprovalsUpdate(ctx, &tempCollection, collection.CollectionApprovals, msg.CollectionApprovals, collection.CollectionPermissions.CanUpdateCollectionApprovals); err != nil {
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
			accountAddr, err := generatePathAddress(path.Denom, WrapperPathGenerationPrefix)
			if err != nil {
				return nil, err
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

			// Auto-set the cosmoscoinwrapper path address as a reserved protocol address
			if err := k.setReservedProtocolAddressForPath(ctx, accountAddr.String(), "cosmoscoinwrapper"); err != nil {
				return nil, err
			}
		}

		collection.CosmosCoinWrapperPaths = append(collection.CosmosCoinWrapperPaths, pathsToAdd...)
	}

	// Handle invariants - convert InvariantsAddObject to CollectionInvariants
	if msg.Invariants != nil {
		if collection.Invariants == nil {
			collection.Invariants = &types.CollectionInvariants{}
		}

		// Set basic invariant fields
		collection.Invariants.NoCustomOwnershipTimes = msg.Invariants.NoCustomOwnershipTimes
		collection.Invariants.MaxSupplyPerId = msg.Invariants.MaxSupplyPerId
		collection.Invariants.NoForcefulPostMintTransfers = msg.Invariants.NoForcefulPostMintTransfers
		collection.Invariants.DisablePoolCreation = msg.Invariants.DisablePoolCreation

		// Handle cosmos coin backed path - generate address
		if msg.Invariants.CosmosCoinBackedPath != nil {
			pathAddObject := msg.Invariants.CosmosCoinBackedPath
			accountAddr, err := generatePathAddress(pathAddObject.IbcDenom, BackedPathGenerationPrefix)
			if err != nil {
				return nil, err
			}

			// ibcAmount is validated in ValidateBasic to be non-zero
			ibcAmount := pathAddObject.IbcAmount

			backedPath := &types.CosmosCoinBackedPath{
				Address:   accountAddr.String(),
				IbcDenom:  pathAddObject.IbcDenom,
				Balances:  pathAddObject.Balances,
				IbcAmount: ibcAmount,
			}

			// Auto-set the cosmoscoinbacked path address as a reserved protocol address
			if err := k.setReservedProtocolAddressForPath(ctx, accountAddr.String(), "cosmoscoinbacked"); err != nil {
				return nil, err
			}

			// Set the backed path in invariants
			collection.Invariants.CosmosCoinBackedPath = backedPath
		}
	}

	// Ensure no duplicate denom paths for wrapper paths
	denomPaths := make(map[string]bool)
	for _, path := range collection.CosmosCoinWrapperPaths {
		if _, ok := denomPaths[path.Denom]; ok {
			return nil, fmt.Errorf("duplicate ibc wrapper path denom: %s", path.Denom)
		}
		denomPaths[path.Denom] = true
	}

	// No need to check for duplicate ibc denom paths for backed paths since only one is allowed

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

	// Auto-set collection permission to forbid Mint address approvals if cosmosCoinBackedPath is set
	hasCosmosCoinBackedPath := collection.Invariants != nil && collection.Invariants.CosmosCoinBackedPath != nil
	if hasCosmosCoinBackedPath {
		if collection.CollectionPermissions == nil {
			collection.CollectionPermissions = &types.CollectionPermissions{}
		}
		ensureMintForbiddenPermission(collection.CollectionPermissions, true)
	}

	if err := k.SetCollectionInStore(ctx, collection); err != nil {
		return nil, err
	}

	// Set auto-approve flags for path addresses (must be after SetCollectionInStore)
	// This needs to happen after the collection is stored because GetBalanceOrApplyDefault requires the collection to exist
	for _, path := range collection.CosmosCoinWrapperPaths {
		if err := k.setAutoApproveFlagsForPathAddress(ctx, collection, path.Address, "cosmoscoinwrapper"); err != nil {
			return nil, err
		}
	}

	if collection.Invariants != nil && collection.Invariants.CosmosCoinBackedPath != nil {
		if err := k.setAutoApproveFlagsForPathAddress(ctx, collection, collection.Invariants.CosmosCoinBackedPath.Address, "cosmoscoinbacked"); err != nil {
			return nil, err
		}
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
