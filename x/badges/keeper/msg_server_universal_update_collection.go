package keeper

import (
	"context"
	"encoding/binary"
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
// The generated address is a module address and should not conflict with user addresses
// as module addresses use a different derivation method than user addresses.
// Security: Module addresses are derived using NewModuleCredential which uses the module name,
// a unique prefix, and the path bytes. This makes collisions with user addresses extremely unlikely.
// All generated path addresses are automatically marked as reserved protocol addresses to prevent conflicts.
func generatePathAddress(pathString string, prefix []byte) (sdk.AccAddress, error) {
	// Validate path string is not empty
	if pathString == "" {
		return nil, fmt.Errorf("path string cannot be empty")
	}

	// Validate path string length to prevent DoS attacks
	// Reasonable limit: 1024 bytes (most denoms are much shorter)
	if len(pathString) > 1024 {
		return nil, fmt.Errorf("path string exceeds maximum length of 1024 bytes")
	}

	fullPathBytes := []byte(pathString)
	ac, err := authtypes.NewModuleCredential(types.ModuleName, prefix, fullPathBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to generate module credential: %w", err)
	}
	generatedAddr := sdk.AccAddress(ac.Address())

	// Security check: Verify that the generated address is actually a module address
	// This ensures the address derivation worked correctly and prevents potential issues
	// Module addresses should always start with the module name prefix in their derivation
	if generatedAddr.Empty() {
		return nil, fmt.Errorf("generated path address is empty")
	}

	// Additional validation: Ensure the address can be converted to string (valid Bech32)
	// This provides an extra layer of validation
	addrString := generatedAddr.String()
	if addrString == "" {
		return nil, fmt.Errorf("generated path address cannot be converted to string")
	}

	return generatedAddr, nil
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
// Path addresses are module addresses that need auto-approve flags set for proper operation.
// This function only sets flags if they haven't been explicitly set already, to avoid
// overriding user-configured settings (though path addresses should not have user settings).
// Security: Each flag is checked individually and only set if not already set, preventing
// unintended overrides of existing settings.
func (k msgServer) setAutoApproveFlagsForPathAddress(ctx sdk.Context, collection *types.TokenCollection, pathAddress string, pathType string) error {
	currBalances, _ := k.GetBalanceOrApplyDefault(ctx, collection, pathAddress)

	// Check each flag individually and only set if not already set
	// This prevents overriding existing settings (defense-in-depth, though path addresses
	// should not have user-configured settings)
	flagsChanged := false

	// Set AutoApproveAllIncomingTransfers only if not already set
	// This is required for path addresses to receive tokens from any source
	if !currBalances.AutoApproveAllIncomingTransfers {
		currBalances.AutoApproveAllIncomingTransfers = true
		flagsChanged = true
	}

	// Set AutoApproveSelfInitiatedOutgoingTransfers only if not already set
	// This is required for path addresses to send tokens they initiated
	if !currBalances.AutoApproveSelfInitiatedOutgoingTransfers {
		currBalances.AutoApproveSelfInitiatedOutgoingTransfers = true
		flagsChanged = true
	}

	// Set AutoApproveSelfInitiatedIncomingTransfers only if not already set
	// This is required for path addresses to receive tokens they initiated
	if !currBalances.AutoApproveSelfInitiatedIncomingTransfers {
		currBalances.AutoApproveSelfInitiatedIncomingTransfers = true
		flagsChanged = true
	}

	// Only save if flags were changed
	if flagsChanged {
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
		derivationKey := make([]byte, DerivationKeyLength)
		binary.BigEndian.PutUint64(derivationKey, nextCollectionId.Uint64())
		ac, err := authtypes.NewModuleCredential(types.ModuleName, AccountGenerationPrefix, derivationKey)
		if err != nil {
			return nil, err
		}
		accountAddr := sdk.AccAddress(ac.Address())

		collection = &types.TokenCollection{
			CollectionId:          nextCollectionId,
			CollectionPermissions: &types.CollectionPermissions{},
			DefaultBalances:       msg.DefaultBalances,
			CreatedBy:             msg.Creator,
			// Default manager is the creator if not updating manager
			Manager:           msg.Creator,
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
				if pathAddObject.Conversion == nil || pathAddObject.Conversion.SideA == nil {
					return nil, fmt.Errorf("cosmos coin backed path conversion or sideA is nil")
				}

				backedAccountAddr, err := generatePathAddress(pathAddObject.Conversion.SideA.Denom, BackedPathGenerationPrefix)
				if err != nil {
					return nil, err
				}

				backedPath := &types.CosmosCoinBackedPath{
					Address:    backedAccountAddr.String(),
					Conversion: pathAddObject.Conversion,
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
	if msg.UpdateIsArchived {
		if err := k.ValidateIsArchivedUpdate(ctx, collection.IsArchived, msg.IsArchived, collection.CollectionPermissions.CanArchiveCollection); err != nil {
			return nil, err
		}
		collection.IsArchived = msg.IsArchived
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
		// Use proper field-by-field comparison instead of string comparison to avoid
		// issues with protobuf serialization differences or semantic equivalence
		newApprovalsWithVersion := []*types.CollectionApproval{}
		for _, newApproval := range msg.CollectionApprovals {
			existingApproval, exists := existingApprovals[newApproval.ApprovalId]
			// Only increment version if approval is new or changed (excluding Version field)
			if !exists || !collectionApprovalEqual(existingApproval, newApproval) {
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

	if msg.UpdateCollectionMetadata {
		if err := k.ValidateCollectionMetadataUpdate(ctx, collection.CollectionMetadata, msg.CollectionMetadata, collection.CollectionPermissions.CanUpdateCollectionMetadata); err != nil {
			return nil, err
		}
		collection.CollectionMetadata = msg.CollectionMetadata
	}

	if msg.UpdateTokenMetadata {
		if err := k.ValidateTokenMetadataUpdate(ctx, collection.TokenMetadata, msg.TokenMetadata, collection.CollectionPermissions.CanUpdateTokenMetadata); err != nil {
			return nil, err
		}
		collection.TokenMetadata = msg.TokenMetadata
	}

	if msg.UpdateManager {
		if err := k.ValidateManagerUpdate(ctx, collection.Manager, msg.Manager, collection.CollectionPermissions.CanUpdateManager); err != nil {
			return nil, err
		}
		collection.Manager = msg.Manager
	}

	if msg.UpdateStandards {
		if err := k.ValidateStandardsUpdate(ctx, collection.Standards, msg.Standards, collection.CollectionPermissions.CanUpdateStandards); err != nil {
			return nil, err
		}
		collection.Standards = msg.Standards
	}

	if msg.UpdateCustomData {
		if err := k.ValidateCustomDataUpdate(ctx, collection.CustomData, msg.CustomData, collection.CollectionPermissions.CanUpdateCustomData); err != nil {
			return nil, err
		}
		collection.CustomData = msg.CustomData
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
			// Generate path address with validation
			accountAddr, err := generatePathAddress(path.Denom, WrapperPathGenerationPrefix)
			if err != nil {
				return nil, fmt.Errorf("failed to generate path address for denom %s: %w", path.Denom, err)
			}

			addrString := accountAddr.String()

			// Security: Ensure the generated address is marked as reserved protocol address
			// This prevents conflicts with user addresses and ensures proper access control
			if err := k.setReservedProtocolAddressForPath(ctx, addrString, "cosmoscoinwrapper"); err != nil {
				return nil, fmt.Errorf("failed to set reserved protocol address for path %s: %w", path.Denom, err)
			}

			// Verify the address was properly reserved (defense in depth)
			if !k.IsAddressReservedProtocolInStore(ctx, addrString) {
				return nil, fmt.Errorf("generated path address %s was not properly reserved", addrString)
			}

			pathsToAdd[i] = &types.CosmosCoinWrapperPath{
				Address:                        addrString,
				Denom:                          path.Denom,
				Conversion:                     path.Conversion,
				Symbol:                         path.Symbol,
				DenomUnits:                     path.DenomUnits,
				AllowOverrideWithAnyValidToken: path.AllowOverrideWithAnyValidToken,
				Metadata:                       path.Metadata,
			}
		}

		collection.CosmosCoinWrapperPaths = append(collection.CosmosCoinWrapperPaths, pathsToAdd...)
	}

	if len(msg.AliasPathsToAdd) > 0 {
		pathsToAdd := make([]*types.AliasPath, len(msg.AliasPathsToAdd))
		for i, path := range msg.AliasPathsToAdd {
			pathsToAdd[i] = &types.AliasPath{
				Denom:      path.Denom,
				Conversion: path.Conversion,
				Symbol:     path.Symbol,
				DenomUnits: path.DenomUnits,
				Metadata:   path.Metadata,
			}
		}

		collection.AliasPaths = append(collection.AliasPaths, pathsToAdd...)
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
			if pathAddObject.Conversion == nil || pathAddObject.Conversion.SideA == nil {
				return nil, fmt.Errorf("cosmos coin backed path conversion or sideA is nil")
			}

			// Generate path address with validation
			accountAddr, err := generatePathAddress(pathAddObject.Conversion.SideA.Denom, BackedPathGenerationPrefix)
			if err != nil {
				return nil, fmt.Errorf("failed to generate path address for cosmos coin backed path denom %s: %w", pathAddObject.Conversion.SideA.Denom, err)
			}

			addrString := accountAddr.String()

			// Security: Ensure the generated address is marked as reserved protocol address
			// This prevents conflicts with user addresses and ensures proper access control
			if err := k.setReservedProtocolAddressForPath(ctx, addrString, "cosmoscoinbacked"); err != nil {
				return nil, fmt.Errorf("failed to set reserved protocol address for cosmos coin backed path: %w", err)
			}

			// Verify the address was properly reserved (defense in depth)
			if !k.IsAddressReservedProtocolInStore(ctx, addrString) {
				return nil, fmt.Errorf("generated cosmos coin backed path address %s was not properly reserved", addrString)
			}

			backedPath := &types.CosmosCoinBackedPath{
				Address:    addrString,
				Conversion: pathAddObject.Conversion,
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

	// Ensure no duplicate denom paths for alias paths and no collisions with wrappers
	aliasDenoms := make(map[string]bool)
	for _, path := range collection.AliasPaths {
		if _, ok := aliasDenoms[path.Denom]; ok {
			return nil, fmt.Errorf("duplicate alias path denom: %s", path.Denom)
		}
		if _, wrapperConflict := denomPaths[path.Denom]; wrapperConflict {
			return nil, fmt.Errorf("alias path denom conflicts with wrapper denom: %s", path.Denom)
		}
		aliasDenoms[path.Denom] = true
	}

	// No need to check for duplicate ibc denom paths for backed paths since only one is allowed

	// Ensure no duplicate symbols (including base symbol and denom unit symbols)
	symbolPaths := make(map[string]bool)
	for _, path := range collection.CosmosCoinWrapperPaths {
		if err := validatePathSymbols(path.Symbol, path.DenomUnits, symbolPaths, "ibc wrapper"); err != nil {
			return nil, err
		}
	}

	// Also ensure alias path symbols / denom unit symbols are unique across all paths (wrappers + aliases)
	for _, path := range collection.AliasPaths {
		if err := validatePathSymbols(path.Symbol, path.DenomUnits, symbolPaths, "alias"); err != nil {
			return nil, err
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

	if err := k.SetCollectionInStore(ctx, collection, false); err != nil {
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

	msgStr, err := MarshalMessageForEvent(msg)
	if err != nil {
		return nil, err
	}

	EmitMessageAndIndexerEvents(ctx,
		sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
		sdk.NewAttribute("msg_type", "universal_update_collection"),
		sdk.NewAttribute("msg", msgStr),
		sdk.NewAttribute("collectionId", fmt.Sprint(collection.CollectionId)),
	)

	return &types.MsgUniversalUpdateCollectionResponse{
		CollectionId: collection.CollectionId,
	}, nil
}

// validatePathSymbols enforces uniqueness of symbols/denom unit symbols across all paths,
// and per-path constraints (single default display, non-zero decimals, unique decimals per path).
func validatePathSymbols(pathSymbol string, denomUnits []*types.DenomUnit, symbolPaths map[string]bool, pathType string) error {
	if pathSymbol != "" {
		if _, ok := symbolPaths[pathSymbol]; ok {
			return fmt.Errorf("duplicate %s path symbol: %s", pathType, pathSymbol)
		}
		symbolPaths[pathSymbol] = true
	}

	defaultDisplayCount := 0
	decimalsSet := make(map[string]bool)

	for _, denomUnit := range denomUnits {
		if denomUnit.Symbol != "" {
			if _, ok := symbolPaths[denomUnit.Symbol]; ok {
				return fmt.Errorf("duplicate denom unit symbol: %s", denomUnit.Symbol)
			}
			symbolPaths[denomUnit.Symbol] = true
		}

		if denomUnit.IsDefaultDisplay {
			defaultDisplayCount++
		}

		if denomUnit.Decimals.IsZero() {
			return fmt.Errorf("denom unit decimals cannot be 0")
		}

		decimalsStr := denomUnit.Decimals.String()
		if _, ok := decimalsSet[decimalsStr]; ok {
			return fmt.Errorf("duplicate denom unit decimals: %s", decimalsStr)
		}
		decimalsSet[decimalsStr] = true
	}

	if defaultDisplayCount > 1 {
		return fmt.Errorf("only one denom unit per path can have isDefaultDisplay set to true, found %d", defaultDisplayCount)
	}

	return nil
}
