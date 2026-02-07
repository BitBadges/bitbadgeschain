package tokenization

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	tokenizationtypes "github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

// ConvertBalance converts a Solidity Balance struct to proto Balance
func ConvertBalance(balanceMap map[string]interface{}) (*tokenizationtypes.Balance, error) {
	amountBig, ok := balanceMap["amount"].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("balance.amount must be *big.Int")
	}
	if err := CheckOverflow(amountBig, "amount"); err != nil {
		return nil, err
	}
	amount := sdkmath.NewUintFromBigInt(amountBig)

	ownershipTimesRaw, ok := balanceMap["ownershipTimes"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("balance.ownershipTimes must be array")
	}
	ownershipTimes := make([]struct {
		Start *big.Int `json:"start"`
		End   *big.Int `json:"end"`
	}, 0, len(ownershipTimesRaw))
	for i, otRaw := range ownershipTimesRaw {
		otMap, ok := otRaw.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("ownershipTimes[%d] must be struct", i)
		}
		start, startOk := otMap["start"].(*big.Int)
		end, endOk := otMap["end"].(*big.Int)
		if !startOk || start == nil {
			return nil, fmt.Errorf("ownershipTimes[%d].start must be a valid *big.Int", i)
		}
		if !endOk || end == nil {
			return nil, fmt.Errorf("ownershipTimes[%d].end must be a valid *big.Int", i)
		}
		ownershipTimes = append(ownershipTimes, struct {
			Start *big.Int `json:"start"`
			End   *big.Int `json:"end"`
		}{Start: start, End: end})
	}
	convertedOwnershipTimes, err := ConvertUintRangeArray(ownershipTimes)
	if err != nil {
		return nil, fmt.Errorf("ownershipTimes: %w", err)
	}

	tokenIdsRaw, ok := balanceMap["tokenIds"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("balance.tokenIds must be array")
	}
	tokenIds := make([]struct {
		Start *big.Int `json:"start"`
		End   *big.Int `json:"end"`
	}, 0, len(tokenIdsRaw))
	for i, tidRaw := range tokenIdsRaw {
		tidMap, ok := tidRaw.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("tokenIds[%d] must be struct", i)
		}
		start, startOk := tidMap["start"].(*big.Int)
		end, endOk := tidMap["end"].(*big.Int)
		if !startOk || start == nil {
			return nil, fmt.Errorf("tokenIds[%d].start must be a valid *big.Int", i)
		}
		if !endOk || end == nil {
			return nil, fmt.Errorf("tokenIds[%d].end must be a valid *big.Int", i)
		}
		tokenIds = append(tokenIds, struct {
			Start *big.Int `json:"start"`
			End   *big.Int `json:"end"`
		}{Start: start, End: end})
	}
	convertedTokenIds, err := ConvertUintRangeArray(tokenIds)
	if err != nil {
		return nil, fmt.Errorf("tokenIds: %w", err)
	}

	return &tokenizationtypes.Balance{
		Amount:         amount,
		OwnershipTimes: convertedOwnershipTimes,
		TokenIds:       convertedTokenIds,
	}, nil
}

// ConvertBalanceArray converts an array of Solidity Balance structs
func ConvertBalanceArray(balancesRaw []interface{}) ([]*tokenizationtypes.Balance, error) {
	balances := make([]*tokenizationtypes.Balance, len(balancesRaw))
	for i, balRaw := range balancesRaw {
		balMap, ok := balRaw.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("balances[%d] must be struct", i)
		}
		balance, err := ConvertBalance(balMap)
		if err != nil {
			return nil, fmt.Errorf("balances[%d]: %w", i, err)
		}
		balances[i] = balance
	}
	return balances, nil
}

// ConvertUintRangeArray converts an array of Solidity UintRange structs to proto UintRange array
func ConvertUintRangeArray(ranges []struct {
	Start *big.Int `json:"start"`
	End   *big.Int `json:"end"`
},
) ([]*tokenizationtypes.UintRange, error) {
	return ConvertAndValidateBigIntRanges(ranges, "ranges")
}

// ConvertCollectionMetadata converts Solidity CollectionMetadata to proto CollectionMetadata
func ConvertCollectionMetadata(uri, customData string) *tokenizationtypes.CollectionMetadata {
	return &tokenizationtypes.CollectionMetadata{
		Uri:        uri,
		CustomData: customData,
	}
}

// ConvertTokenMetadata converts Solidity TokenMetadata to proto TokenMetadata
func ConvertTokenMetadata(uri, customData string, tokenIds []struct {
	Start *big.Int `json:"start"`
	End   *big.Int `json:"end"`
},
) (*tokenizationtypes.TokenMetadata, error) {
	convertedTokenIds, err := ConvertUintRangeArray(tokenIds)
	if err != nil {
		return nil, err
	}
	return &tokenizationtypes.TokenMetadata{
		Uri:        uri,
		CustomData: customData,
		TokenIds:   convertedTokenIds,
	}, nil
}

// ConvertActionPermission converts Solidity ActionPermission to proto ActionPermission
// Allows empty time arrays (empty means no time restrictions)
func ConvertActionPermission(permittedTimes, forbiddenTimes []struct {
	Start *big.Int `json:"start"`
	End   *big.Int `json:"end"`
},
) (*tokenizationtypes.ActionPermission, error) {
	permTimes, err := ConvertAndValidateBigIntRangesAllowEmpty(permittedTimes, "permanentlyPermittedTimes")
	if err != nil {
		return nil, err
	}
	forbTimes, err := ConvertAndValidateBigIntRangesAllowEmpty(forbiddenTimes, "permanentlyForbiddenTimes")
	if err != nil {
		return nil, err
	}
	return &tokenizationtypes.ActionPermission{
		PermanentlyPermittedTimes: permTimes,
		PermanentlyForbiddenTimes: forbTimes,
	}, nil
}

// ConvertActionPermissionArray converts an array of Solidity ActionPermissions
func ConvertActionPermissionArray(permissions []struct {
	PermanentlyPermittedTimes []struct {
		Start *big.Int `json:"start"`
		End   *big.Int `json:"end"`
	} `json:"permanentlyPermittedTimes"`
	PermanentlyForbiddenTimes []struct {
		Start *big.Int `json:"start"`
		End   *big.Int `json:"end"`
	} `json:"permanentlyForbiddenTimes"`
},
) ([]*tokenizationtypes.ActionPermission, error) {
	result := make([]*tokenizationtypes.ActionPermission, len(permissions))
	for i, perm := range permissions {
		converted, err := ConvertActionPermission(perm.PermanentlyPermittedTimes, perm.PermanentlyForbiddenTimes)
		if err != nil {
			return nil, fmt.Errorf("permission[%d]: %w", i, err)
		}
		result[i] = converted
	}
	return result, nil
}

// ConvertUserBalanceStore converts Solidity UserBalanceStore to proto UserBalanceStore
func ConvertUserBalanceStore(storeMap map[string]interface{}) (*tokenizationtypes.UserBalanceStore, error) {
	store := &tokenizationtypes.UserBalanceStore{
		UserPermissions: &tokenizationtypes.UserPermissions{},
	}

	if balancesRaw, ok := storeMap["balances"].([]interface{}); ok {
		balances, err := ConvertBalanceArray(balancesRaw)
		if err != nil {
			return nil, fmt.Errorf("balances: %w", err)
		}
		store.Balances = balances
	}

	if outgoingRaw, ok := storeMap["outgoingApprovals"].([]interface{}); ok {
		outgoing, err := ConvertUserOutgoingApprovalArray(outgoingRaw)
		if err != nil {
			return nil, fmt.Errorf("outgoingApprovals: %w", err)
		}
		store.OutgoingApprovals = outgoing
	}

	if incomingRaw, ok := storeMap["incomingApprovals"].([]interface{}); ok {
		incoming, err := ConvertUserIncomingApprovalArray(incomingRaw)
		if err != nil {
			return nil, fmt.Errorf("incomingApprovals: %w", err)
		}
		store.IncomingApprovals = incoming
	}

	if val, ok := storeMap["autoApproveSelfInitiatedOutgoingTransfers"].(bool); ok {
		store.AutoApproveSelfInitiatedOutgoingTransfers = val
	}
	if val, ok := storeMap["autoApproveSelfInitiatedIncomingTransfers"].(bool); ok {
		store.AutoApproveSelfInitiatedIncomingTransfers = val
	}
	if val, ok := storeMap["autoApproveAllIncomingTransfers"].(bool); ok {
		store.AutoApproveAllIncomingTransfers = val
	}

	if permsRaw, ok := storeMap["userPermissions"].(map[string]interface{}); ok {
		perms, err := ConvertUserPermissions(permsRaw)
		if err != nil {
			return nil, fmt.Errorf("userPermissions: %w", err)
		}
		store.UserPermissions = perms
	}

	return store, nil
}

// ConvertCollectionPermissions converts Solidity CollectionPermissions to proto CollectionPermissions
func ConvertCollectionPermissions(permsMap map[string]interface{}) (*tokenizationtypes.CollectionPermissions, error) {
	perms := &tokenizationtypes.CollectionPermissions{}

	if val, ok := permsMap["canDeleteCollection"].([]interface{}); ok {
		converted, err := ConvertActionPermissionArrayFromInterface(val)
		if err != nil {
			return nil, fmt.Errorf("canDeleteCollection: %w", err)
		}
		perms.CanDeleteCollection = converted
	}
	if val, ok := permsMap["canArchiveCollection"].([]interface{}); ok {
		converted, err := ConvertActionPermissionArrayFromInterface(val)
		if err != nil {
			return nil, fmt.Errorf("canArchiveCollection: %w", err)
		}
		perms.CanArchiveCollection = converted
	}
	if val, ok := permsMap["canUpdateStandards"].([]interface{}); ok {
		converted, err := ConvertActionPermissionArrayFromInterface(val)
		if err != nil {
			return nil, fmt.Errorf("canUpdateStandards: %w", err)
		}
		perms.CanUpdateStandards = converted
	}
	if val, ok := permsMap["canUpdateCustomData"].([]interface{}); ok {
		converted, err := ConvertActionPermissionArrayFromInterface(val)
		if err != nil {
			return nil, fmt.Errorf("canUpdateCustomData: %w", err)
		}
		perms.CanUpdateCustomData = converted
	}
	if val, ok := permsMap["canUpdateManager"].([]interface{}); ok {
		converted, err := ConvertActionPermissionArrayFromInterface(val)
		if err != nil {
			return nil, fmt.Errorf("canUpdateManager: %w", err)
		}
		perms.CanUpdateManager = converted
	}
	if val, ok := permsMap["canUpdateCollectionMetadata"].([]interface{}); ok {
		converted, err := ConvertActionPermissionArrayFromInterface(val)
		if err != nil {
			return nil, fmt.Errorf("canUpdateCollectionMetadata: %w", err)
		}
		perms.CanUpdateCollectionMetadata = converted
	}
	if val, ok := permsMap["canAddMoreAliasPaths"].([]interface{}); ok {
		converted, err := ConvertActionPermissionArrayFromInterface(val)
		if err != nil {
			return nil, fmt.Errorf("canAddMoreAliasPaths: %w", err)
		}
		perms.CanAddMoreAliasPaths = converted
	}
	if val, ok := permsMap["canAddMoreCosmosCoinWrapperPaths"].([]interface{}); ok {
		converted, err := ConvertActionPermissionArrayFromInterface(val)
		if err != nil {
			return nil, fmt.Errorf("canAddMoreCosmosCoinWrapperPaths: %w", err)
		}
		perms.CanAddMoreCosmosCoinWrapperPaths = converted
	}

	if val, ok := permsMap["canUpdateValidTokenIds"].([]interface{}); ok {
		converted, err := ConvertTokenIdsActionPermissionArrayFromInterface(val)
		if err != nil {
			return nil, fmt.Errorf("canUpdateValidTokenIds: %w", err)
		}
		perms.CanUpdateValidTokenIds = converted
	}
	if val, ok := permsMap["canUpdateTokenMetadata"].([]interface{}); ok {
		converted, err := ConvertTokenIdsActionPermissionArrayFromInterface(val)
		if err != nil {
			return nil, fmt.Errorf("canUpdateTokenMetadata: %w", err)
		}
		perms.CanUpdateTokenMetadata = converted
	}

	if val, ok := permsMap["canUpdateCollectionApprovals"].([]interface{}); ok {
		converted, err := ConvertCollectionApprovalPermissionArrayFromInterface(val)
		if err != nil {
			return nil, fmt.Errorf("canUpdateCollectionApprovals: %w", err)
		}
		perms.CanUpdateCollectionApprovals = converted
	}

	return perms, nil
}

// ConvertAddressListInput converts Solidity AddressListInput to proto AddressListInput
func ConvertAddressListInput(listId string, addresses []string, whitelist bool, uri, customData string) *tokenizationtypes.AddressListInput {
	cosmosAddrs := make([]string, len(addresses))
	for i, addr := range addresses {
		if common.IsHexAddress(addr) {
			evmAddr := common.HexToAddress(addr)
			cosmosAddrs[i] = sdk.AccAddress(evmAddr.Bytes()).String()
		} else {
			cosmosAddrs[i] = addr
		}
	}
	return &tokenizationtypes.AddressListInput{
		ListId:     listId,
		Addresses:  cosmosAddrs,
		Whitelist:  whitelist,
		Uri:        uri,
		CustomData: customData,
	}
}

// ConvertManagerAddress converts an EVM address or string to Cosmos address string
func ConvertManagerAddress(manager interface{}) (string, error) {
	switch m := manager.(type) {
	case common.Address:
		return sdk.AccAddress(m.Bytes()).String(), nil
	case string:
		if common.IsHexAddress(m) {
			evmAddr := common.HexToAddress(m)
			return sdk.AccAddress(evmAddr.Bytes()).String(), nil
		}
		return m, nil
	default:
		return "", fmt.Errorf("invalid manager type: expected common.Address or string, got %T", manager)
	}
}

// ConvertInvariantsAddObject converts Solidity InvariantsAddObject to proto InvariantsAddObject
func ConvertInvariantsAddObject(invariants interface{}) (*tokenizationtypes.InvariantsAddObject, error) {
	invariantsMap, ok := invariants.(map[string]interface{})
	if !ok {
		return &tokenizationtypes.InvariantsAddObject{
			NoCustomOwnershipTimes:      false,
			MaxSupplyPerId:              sdkmath.NewUint(0),
			CosmosCoinBackedPath:        nil,
			NoForcefulPostMintTransfers: false,
			DisablePoolCreation:         false,
		}, nil
	}

	result := &tokenizationtypes.InvariantsAddObject{}

	if val, ok := invariantsMap["noCustomOwnershipTimes"].(bool); ok {
		result.NoCustomOwnershipTimes = val
	}
	if val, ok := invariantsMap["noForcefulPostMintTransfers"].(bool); ok {
		result.NoForcefulPostMintTransfers = val
	}
	if val, ok := invariantsMap["disablePoolCreation"].(bool); ok {
		result.DisablePoolCreation = val
	}

	if maxSupplyBig, ok := invariantsMap["maxSupplyPerId"].(*big.Int); ok {
		if err := CheckOverflow(maxSupplyBig, "maxSupplyPerId"); err != nil {
			return nil, err
		}
		result.MaxSupplyPerId = sdkmath.NewUintFromBigInt(maxSupplyBig)
	} else {
		result.MaxSupplyPerId = sdkmath.NewUint(0)
	}

	if _, ok := invariantsMap["cosmosCoinBackedPath"].(map[string]interface{}); ok {
		result.CosmosCoinBackedPath = &tokenizationtypes.CosmosCoinBackedPathAddObject{
			Conversion: nil,
		}
	}

	return result, nil
}

// ConvertCosmosCoinWrapperPathAddObjectArray converts array of Solidity CosmosCoinWrapperPathAddObject
func ConvertCosmosCoinWrapperPathAddObjectArray(paths interface{}) ([]*tokenizationtypes.CosmosCoinWrapperPathAddObject, error) {
	pathsRaw, ok := paths.([]interface{})
	if !ok {
		return []*tokenizationtypes.CosmosCoinWrapperPathAddObject{}, nil
	}

	result := make([]*tokenizationtypes.CosmosCoinWrapperPathAddObject, 0, len(pathsRaw))
	for i, pathRaw := range pathsRaw {
		pathMap, ok := pathRaw.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("cosmosCoinWrapperPaths[%d] must be a map", i)
		}

		path := &tokenizationtypes.CosmosCoinWrapperPathAddObject{}

		if val, ok := pathMap["denom"].(string); ok {
			path.Denom = val
		}
		if val, ok := pathMap["symbol"].(string); ok {
			path.Symbol = val
		}
		if val, ok := pathMap["allowOverrideWithAnyValidToken"].(bool); ok {
			path.AllowOverrideWithAnyValidToken = val
		}

		if conversionRaw, ok := pathMap["conversion"].(map[string]interface{}); ok && conversionRaw != nil {
			path.Conversion = &tokenizationtypes.ConversionWithoutDenom{}
		}

		if denomUnitsRaw, ok := pathMap["denomUnits"].([]interface{}); ok {
			denomUnits := make([]*tokenizationtypes.DenomUnit, 0, len(denomUnitsRaw))
			for _, duRaw := range denomUnitsRaw {
				if duMap, ok := duRaw.(map[string]interface{}); ok {
					denomUnit := &tokenizationtypes.DenomUnit{}
					if decimalsBig, ok := duMap["decimals"].(*big.Int); ok {
						if err := CheckOverflow(decimalsBig, "decimals"); err == nil {
							denomUnit.Decimals = sdkmath.NewUintFromBigInt(decimalsBig)
						}
					}
					if val, ok := duMap["symbol"].(string); ok {
						denomUnit.Symbol = val
					}
					if val, ok := duMap["isDefaultDisplay"].(bool); ok {
						denomUnit.IsDefaultDisplay = val
					}
					if metadataRaw, ok := duMap["metadata"].(map[string]interface{}); ok {
						metadata := &tokenizationtypes.PathMetadata{}
						if val, ok := metadataRaw["uri"].(string); ok {
							metadata.Uri = val
						}
						if val, ok := metadataRaw["customData"].(string); ok {
							metadata.CustomData = val
						}
						denomUnit.Metadata = metadata
					}
					denomUnits = append(denomUnits, denomUnit)
				}
			}
			path.DenomUnits = denomUnits
		}

		if metadataRaw, ok := pathMap["metadata"].(map[string]interface{}); ok {
			metadata := &tokenizationtypes.PathMetadata{}
			if val, ok := metadataRaw["uri"].(string); ok {
				metadata.Uri = val
			}
			if val, ok := metadataRaw["customData"].(string); ok {
				metadata.CustomData = val
			}
			path.Metadata = metadata
		}

		result = append(result, path)
	}

	return result, nil
}

// ConvertAliasPathAddObjectArray converts array of Solidity AliasPathAddObject
func ConvertAliasPathAddObjectArray(paths interface{}) ([]*tokenizationtypes.AliasPathAddObject, error) {
	pathsRaw, ok := paths.([]interface{})
	if !ok {
		return []*tokenizationtypes.AliasPathAddObject{}, nil
	}

	result := make([]*tokenizationtypes.AliasPathAddObject, 0, len(pathsRaw))
	for i, pathRaw := range pathsRaw {
		pathMap, ok := pathRaw.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("aliasPaths[%d] must be a map", i)
		}

		path := &tokenizationtypes.AliasPathAddObject{}

		if val, ok := pathMap["denom"].(string); ok {
			path.Denom = val
		}
		if val, ok := pathMap["symbol"].(string); ok {
			path.Symbol = val
		}

		if conversionRaw, ok := pathMap["conversion"].(map[string]interface{}); ok && conversionRaw != nil {
			path.Conversion = &tokenizationtypes.ConversionWithoutDenom{}
		}

		if denomUnitsRaw, ok := pathMap["denomUnits"].([]interface{}); ok {
			denomUnits := make([]*tokenizationtypes.DenomUnit, 0, len(denomUnitsRaw))
			for _, duRaw := range denomUnitsRaw {
				if duMap, ok := duRaw.(map[string]interface{}); ok {
					denomUnit := &tokenizationtypes.DenomUnit{}
					if decimalsBig, ok := duMap["decimals"].(*big.Int); ok {
						if err := CheckOverflow(decimalsBig, "decimals"); err == nil {
							denomUnit.Decimals = sdkmath.NewUintFromBigInt(decimalsBig)
						}
					}
					if val, ok := duMap["symbol"].(string); ok {
						denomUnit.Symbol = val
					}
					if val, ok := duMap["isDefaultDisplay"].(bool); ok {
						denomUnit.IsDefaultDisplay = val
					}
					if metadataRaw, ok := duMap["metadata"].(map[string]interface{}); ok {
						metadata := &tokenizationtypes.PathMetadata{}
						if val, ok := metadataRaw["uri"].(string); ok {
							metadata.Uri = val
						}
						if val, ok := metadataRaw["customData"].(string); ok {
							metadata.CustomData = val
						}
						denomUnit.Metadata = metadata
					}
					denomUnits = append(denomUnits, denomUnit)
				}
			}
			path.DenomUnits = denomUnits
		}

		if metadataRaw, ok := pathMap["metadata"].(map[string]interface{}); ok {
			metadata := &tokenizationtypes.PathMetadata{}
			if val, ok := metadataRaw["uri"].(string); ok {
				metadata.Uri = val
			}
			if val, ok := metadataRaw["customData"].(string); ok {
				metadata.CustomData = val
			}
			path.Metadata = metadata
		}

		result = append(result, path)
	}

	return result, nil
}

// convertUintRangeArrayFromInterfaceWithError converts []interface{} to []struct{Start, End *big.Int}
// Returns an error if any range has invalid Start/End values
func convertUintRangeArrayFromInterfaceWithError(rangesRaw []interface{}, fieldName string) ([]struct {
	Start *big.Int `json:"start"`
	End   *big.Int `json:"end"`
}, error) {
	if rangesRaw == nil {
		return nil, nil
	}
	ranges := make([]struct {
		Start *big.Int `json:"start"`
		End   *big.Int `json:"end"`
	}, 0, len(rangesRaw))
	for i, rRaw := range rangesRaw {
		if rMap, ok := rRaw.(map[string]interface{}); ok {
			start, startOk := rMap["start"].(*big.Int)
			end, endOk := rMap["end"].(*big.Int)
			if !startOk || start == nil {
				return nil, fmt.Errorf("%s[%d].start must be a valid *big.Int", fieldName, i)
			}
			if !endOk || end == nil {
				return nil, fmt.Errorf("%s[%d].end must be a valid *big.Int", fieldName, i)
			}
			ranges = append(ranges, struct {
				Start *big.Int `json:"start"`
				End   *big.Int `json:"end"`
			}{Start: start, End: end})
		} else {
			return nil, fmt.Errorf("%s[%d] must be a map with start and end fields", fieldName, i)
		}
	}
	return ranges, nil
}

// ConvertActionPermissionArrayFromInterface converts []interface{} to []*ActionPermission
func ConvertActionPermissionArrayFromInterface(permsRaw []interface{}) ([]*tokenizationtypes.ActionPermission, error) {
	perms := make([]*tokenizationtypes.ActionPermission, 0, len(permsRaw))
	for i, permRaw := range permsRaw {
		permMap, ok := permRaw.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("permission[%d] must be a map", i)
		}
		permTimesRaw, _ := permMap["permanentlyPermittedTimes"].([]interface{})
		forbTimesRaw, _ := permMap["permanentlyForbiddenTimes"].([]interface{})
		permTimes, err := convertUintRangeArrayFromInterfaceWithError(permTimesRaw, fmt.Sprintf("permission[%d].permanentlyPermittedTimes", i))
		if err != nil {
			return nil, err
		}
		forbTimes, err := convertUintRangeArrayFromInterfaceWithError(forbTimesRaw, fmt.Sprintf("permission[%d].permanentlyForbiddenTimes", i))
		if err != nil {
			return nil, err
		}
		perm, err := ConvertActionPermission(permTimes, forbTimes)
		if err != nil {
			return nil, fmt.Errorf("permission[%d]: %w", i, err)
		}
		perms = append(perms, perm)
	}
	return perms, nil
}

// ConvertTokenIdsActionPermissionArrayFromInterface converts []interface{} to []*TokenIdsActionPermission
func ConvertTokenIdsActionPermissionArrayFromInterface(permsRaw []interface{}) ([]*tokenizationtypes.TokenIdsActionPermission, error) {
	perms := make([]*tokenizationtypes.TokenIdsActionPermission, 0, len(permsRaw))
	for i, permRaw := range permsRaw {
		permMap, ok := permRaw.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("permission[%d] must be struct", i)
		}
		tokenIdsRaw, _ := permMap["tokenIds"].([]interface{})
		permTimesRaw, _ := permMap["permanentlyPermittedTimes"].([]interface{})
		forbTimesRaw, _ := permMap["permanentlyForbiddenTimes"].([]interface{})
		// TokenIds are optional (empty means all token IDs)
		tokenIdsConverted, err := convertUintRangeArrayFromInterfaceWithError(tokenIdsRaw, fmt.Sprintf("permission[%d].tokenIds", i))
		if err != nil {
			return nil, err
		}
		tokenIds, err := ConvertAndValidateBigIntRangesAllowEmpty(tokenIdsConverted, fmt.Sprintf("permission[%d].tokenIds", i))
		if err != nil {
			return nil, err
		}
		// Time arrays are optional (empty means no restrictions)
		permTimesConverted, err := convertUintRangeArrayFromInterfaceWithError(permTimesRaw, fmt.Sprintf("permission[%d].permanentlyPermittedTimes", i))
		if err != nil {
			return nil, err
		}
		permTimes, err := ConvertAndValidateBigIntRangesAllowEmpty(permTimesConverted, fmt.Sprintf("permission[%d].permanentlyPermittedTimes", i))
		if err != nil {
			return nil, err
		}
		forbTimesConverted, err := convertUintRangeArrayFromInterfaceWithError(forbTimesRaw, fmt.Sprintf("permission[%d].permanentlyForbiddenTimes", i))
		if err != nil {
			return nil, err
		}
		forbTimes, err := ConvertAndValidateBigIntRangesAllowEmpty(forbTimesConverted, fmt.Sprintf("permission[%d].permanentlyForbiddenTimes", i))
		if err != nil {
			return nil, err
		}
		perms = append(perms, &tokenizationtypes.TokenIdsActionPermission{
			TokenIds:                  tokenIds,
			PermanentlyPermittedTimes: permTimes,
			PermanentlyForbiddenTimes: forbTimes,
		})
	}
	return perms, nil
}

// ConvertCollectionApprovalPermissionArrayFromInterface converts []interface{} to []*CollectionApprovalPermission
func ConvertCollectionApprovalPermissionArrayFromInterface(permsRaw []interface{}) ([]*tokenizationtypes.CollectionApprovalPermission, error) {
	perms := make([]*tokenizationtypes.CollectionApprovalPermission, 0, len(permsRaw))
	for i, permRaw := range permsRaw {
		permMap, ok := permRaw.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("permission[%d] must be struct", i)
		}
		fromListId, _ := permMap["fromListId"].(string)
		toListId, _ := permMap["toListId"].(string)
		initiatedByListId, _ := permMap["initiatedByListId"].(string)
		approvalId, _ := permMap["approvalId"].(string)
		transferTimesRaw, _ := permMap["transferTimes"].([]interface{})
		tokenIdsRaw, _ := permMap["tokenIds"].([]interface{})
		ownershipTimesRaw, _ := permMap["ownershipTimes"].([]interface{})
		permTimesRaw, _ := permMap["permanentlyPermittedTimes"].([]interface{})
		forbTimesRaw, _ := permMap["permanentlyForbiddenTimes"].([]interface{})
		// All range arrays are optional (empty means all/no restrictions)
		transferTimesConverted, err := convertUintRangeArrayFromInterfaceWithError(transferTimesRaw, fmt.Sprintf("permission[%d].transferTimes", i))
		if err != nil {
			return nil, err
		}
		transferTimes, err := ConvertAndValidateBigIntRangesAllowEmpty(transferTimesConverted, fmt.Sprintf("permission[%d].transferTimes", i))
		if err != nil {
			return nil, err
		}
		tokenIdsConverted, err := convertUintRangeArrayFromInterfaceWithError(tokenIdsRaw, fmt.Sprintf("permission[%d].tokenIds", i))
		if err != nil {
			return nil, err
		}
		tokenIds, err := ConvertAndValidateBigIntRangesAllowEmpty(tokenIdsConverted, fmt.Sprintf("permission[%d].tokenIds", i))
		if err != nil {
			return nil, err
		}
		ownershipTimesConverted, err := convertUintRangeArrayFromInterfaceWithError(ownershipTimesRaw, fmt.Sprintf("permission[%d].ownershipTimes", i))
		if err != nil {
			return nil, err
		}
		ownershipTimes, err := ConvertAndValidateBigIntRangesAllowEmpty(ownershipTimesConverted, fmt.Sprintf("permission[%d].ownershipTimes", i))
		if err != nil {
			return nil, err
		}
		permTimesConverted, err := convertUintRangeArrayFromInterfaceWithError(permTimesRaw, fmt.Sprintf("permission[%d].permanentlyPermittedTimes", i))
		if err != nil {
			return nil, err
		}
		permTimes, err := ConvertAndValidateBigIntRangesAllowEmpty(permTimesConverted, fmt.Sprintf("permission[%d].permanentlyPermittedTimes", i))
		if err != nil {
			return nil, err
		}
		forbTimesConverted, err := convertUintRangeArrayFromInterfaceWithError(forbTimesRaw, fmt.Sprintf("permission[%d].permanentlyForbiddenTimes", i))
		if err != nil {
			return nil, err
		}
		forbTimes, err := ConvertAndValidateBigIntRangesAllowEmpty(forbTimesConverted, fmt.Sprintf("permission[%d].permanentlyForbiddenTimes", i))
		if err != nil {
			return nil, err
		}
		perms = append(perms, &tokenizationtypes.CollectionApprovalPermission{
			FromListId:                fromListId,
			ToListId:                  toListId,
			InitiatedByListId:         initiatedByListId,
			TransferTimes:             transferTimes,
			TokenIds:                  tokenIds,
			OwnershipTimes:            ownershipTimes,
			ApprovalId:                approvalId,
			PermanentlyPermittedTimes: permTimes,
			PermanentlyForbiddenTimes: forbTimes,
		})
	}
	return perms, nil
}

// ConvertUserOutgoingApprovalArray converts array of UserOutgoingApproval
func ConvertUserOutgoingApprovalArray(approvalsRaw []interface{}) ([]*tokenizationtypes.UserOutgoingApproval, error) {
	approvals := make([]*tokenizationtypes.UserOutgoingApproval, 0, len(approvalsRaw))
	for i, appRaw := range approvalsRaw {
		appMap, ok := appRaw.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("approval[%d] must be struct", i)
		}
		approval, err := ConvertUserOutgoingApproval(appMap)
		if err != nil {
			return nil, fmt.Errorf("approval[%d]: %w", i, err)
		}
		approvals = append(approvals, approval)
	}
	return approvals, nil
}

// ConvertUserIncomingApprovalArray converts array of UserIncomingApproval
func ConvertUserIncomingApprovalArray(approvalsRaw []interface{}) ([]*tokenizationtypes.UserIncomingApproval, error) {
	approvals := make([]*tokenizationtypes.UserIncomingApproval, 0, len(approvalsRaw))
	for i, appRaw := range approvalsRaw {
		appMap, ok := appRaw.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("approval[%d] must be struct", i)
		}
		approval, err := ConvertUserIncomingApproval(appMap)
		if err != nil {
			return nil, fmt.Errorf("approval[%d]: %w", i, err)
		}
		approvals = append(approvals, approval)
	}
	return approvals, nil
}

// ConvertUserOutgoingApproval converts a single UserOutgoingApproval
func ConvertUserOutgoingApproval(appMap map[string]interface{}) (*tokenizationtypes.UserOutgoingApproval, error) {
	approval := &tokenizationtypes.UserOutgoingApproval{
		Version: sdkmath.NewUint(0),
	}

	if val, ok := appMap["approvalId"].(string); ok {
		approval.ApprovalId = val
	}
	if val, ok := appMap["toListId"].(string); ok {
		approval.ToListId = val
	}
	if val, ok := appMap["initiatedByListId"].(string); ok {
		approval.InitiatedByListId = val
	}
	if val, ok := appMap["uri"].(string); ok {
		approval.Uri = val
	}
	if val, ok := appMap["customData"].(string); ok {
		approval.CustomData = val
	}

	if transferTimesRaw, ok := appMap["transferTimes"].([]interface{}); ok {
		transferTimesConverted, err := convertUintRangeArrayFromInterfaceWithError(transferTimesRaw, "transferTimes")
		if err != nil {
			return nil, err
		}
		transferTimes, err := ConvertUintRangeArray(transferTimesConverted)
		if err != nil {
			return nil, fmt.Errorf("transferTimes: %w", err)
		}
		approval.TransferTimes = transferTimes
	}
	if tokenIdsRaw, ok := appMap["tokenIds"].([]interface{}); ok {
		tokenIdsConverted, err := convertUintRangeArrayFromInterfaceWithError(tokenIdsRaw, "tokenIds")
		if err != nil {
			return nil, err
		}
		tokenIds, err := ConvertUintRangeArray(tokenIdsConverted)
		if err != nil {
			return nil, fmt.Errorf("tokenIds: %w", err)
		}
		approval.TokenIds = tokenIds
	}
	if ownershipTimesRaw, ok := appMap["ownershipTimes"].([]interface{}); ok {
		ownershipTimesConverted, err := convertUintRangeArrayFromInterfaceWithError(ownershipTimesRaw, "ownershipTimes")
		if err != nil {
			return nil, err
		}
		ownershipTimes, err := ConvertUintRangeArray(ownershipTimesConverted)
		if err != nil {
			return nil, fmt.Errorf("ownershipTimes: %w", err)
		}
		approval.OwnershipTimes = ownershipTimes
	}

	// Convert ApprovalCriteria if provided, otherwise use empty default
	if approvalCriteriaRaw, ok := appMap["approvalCriteria"].(map[string]interface{}); ok {
		criteria, err := ConvertOutgoingApprovalCriteria(approvalCriteriaRaw)
		if err != nil {
			return nil, fmt.Errorf("approvalCriteria: %w", err)
		}
		approval.ApprovalCriteria = criteria
	} else {
		approval.ApprovalCriteria = &tokenizationtypes.OutgoingApprovalCriteria{}
	}

	return approval, nil
}

// ConvertUserIncomingApproval converts a single UserIncomingApproval
func ConvertUserIncomingApproval(appMap map[string]interface{}) (*tokenizationtypes.UserIncomingApproval, error) {
	approval := &tokenizationtypes.UserIncomingApproval{
		Version: sdkmath.NewUint(0),
	}

	if val, ok := appMap["approvalId"].(string); ok {
		approval.ApprovalId = val
	}
	if val, ok := appMap["fromListId"].(string); ok {
		approval.FromListId = val
	}
	if val, ok := appMap["initiatedByListId"].(string); ok {
		approval.InitiatedByListId = val
	}
	if val, ok := appMap["uri"].(string); ok {
		approval.Uri = val
	}
	if val, ok := appMap["customData"].(string); ok {
		approval.CustomData = val
	}

	if transferTimesRaw, ok := appMap["transferTimes"].([]interface{}); ok {
		transferTimesConverted, err := convertUintRangeArrayFromInterfaceWithError(transferTimesRaw, "transferTimes")
		if err != nil {
			return nil, err
		}
		transferTimes, err := ConvertUintRangeArray(transferTimesConverted)
		if err != nil {
			return nil, fmt.Errorf("transferTimes: %w", err)
		}
		approval.TransferTimes = transferTimes
	}
	if tokenIdsRaw, ok := appMap["tokenIds"].([]interface{}); ok {
		tokenIdsConverted, err := convertUintRangeArrayFromInterfaceWithError(tokenIdsRaw, "tokenIds")
		if err != nil {
			return nil, err
		}
		tokenIds, err := ConvertUintRangeArray(tokenIdsConverted)
		if err != nil {
			return nil, fmt.Errorf("tokenIds: %w", err)
		}
		approval.TokenIds = tokenIds
	}
	if ownershipTimesRaw, ok := appMap["ownershipTimes"].([]interface{}); ok {
		ownershipTimesConverted, err := convertUintRangeArrayFromInterfaceWithError(ownershipTimesRaw, "ownershipTimes")
		if err != nil {
			return nil, err
		}
		ownershipTimes, err := ConvertUintRangeArray(ownershipTimesConverted)
		if err != nil {
			return nil, fmt.Errorf("ownershipTimes: %w", err)
		}
		approval.OwnershipTimes = ownershipTimes
	}

	// Convert ApprovalCriteria if provided, otherwise use empty default
	if approvalCriteriaRaw, ok := appMap["approvalCriteria"].(map[string]interface{}); ok {
		criteria, err := ConvertIncomingApprovalCriteria(approvalCriteriaRaw)
		if err != nil {
			return nil, fmt.Errorf("approvalCriteria: %w", err)
		}
		approval.ApprovalCriteria = criteria
	} else {
		approval.ApprovalCriteria = &tokenizationtypes.IncomingApprovalCriteria{}
	}

	return approval, nil
}

// ConvertUserPermissions converts UserPermissions
func ConvertUserPermissions(permsMap map[string]interface{}) (*tokenizationtypes.UserPermissions, error) {
	perms := &tokenizationtypes.UserPermissions{}

	if val, ok := permsMap["canUpdateOutgoingApprovals"].([]interface{}); ok {
		converted, err := ConvertUserOutgoingApprovalPermissionArrayFromInterface(val)
		if err != nil {
			return nil, fmt.Errorf("canUpdateOutgoingApprovals: %w", err)
		}
		perms.CanUpdateOutgoingApprovals = converted
	}
	if val, ok := permsMap["canUpdateIncomingApprovals"].([]interface{}); ok {
		converted, err := ConvertUserIncomingApprovalPermissionArrayFromInterface(val)
		if err != nil {
			return nil, fmt.Errorf("canUpdateIncomingApprovals: %w", err)
		}
		perms.CanUpdateIncomingApprovals = converted
	}
	if val, ok := permsMap["canUpdateAutoApproveSelfInitiatedOutgoingTransfers"].([]interface{}); ok {
		converted, err := ConvertActionPermissionArrayFromInterface(val)
		if err != nil {
			return nil, fmt.Errorf("canUpdateAutoApproveSelfInitiatedOutgoingTransfers: %w", err)
		}
		perms.CanUpdateAutoApproveSelfInitiatedOutgoingTransfers = converted
	}
	if val, ok := permsMap["canUpdateAutoApproveSelfInitiatedIncomingTransfers"].([]interface{}); ok {
		converted, err := ConvertActionPermissionArrayFromInterface(val)
		if err != nil {
			return nil, fmt.Errorf("canUpdateAutoApproveSelfInitiatedIncomingTransfers: %w", err)
		}
		perms.CanUpdateAutoApproveSelfInitiatedIncomingTransfers = converted
	}
	if val, ok := permsMap["canUpdateAutoApproveAllIncomingTransfers"].([]interface{}); ok {
		converted, err := ConvertActionPermissionArrayFromInterface(val)
		if err != nil {
			return nil, fmt.Errorf("canUpdateAutoApproveAllIncomingTransfers: %w", err)
		}
		perms.CanUpdateAutoApproveAllIncomingTransfers = converted
	}

	return perms, nil
}

// ConvertUserOutgoingApprovalPermissionArrayFromInterface converts []interface{} to []*UserOutgoingApprovalPermission
func ConvertUserOutgoingApprovalPermissionArrayFromInterface(permsRaw []interface{}) ([]*tokenizationtypes.UserOutgoingApprovalPermission, error) {
	perms := make([]*tokenizationtypes.UserOutgoingApprovalPermission, 0, len(permsRaw))
	for i, permRaw := range permsRaw {
		permMap, ok := permRaw.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("permission[%d] must be struct", i)
		}
		toListId, _ := permMap["toListId"].(string)
		initiatedByListId, _ := permMap["initiatedByListId"].(string)
		approvalId, _ := permMap["approvalId"].(string)
		transferTimesRaw, _ := permMap["transferTimes"].([]interface{})
		tokenIdsRaw, _ := permMap["tokenIds"].([]interface{})
		ownershipTimesRaw, _ := permMap["ownershipTimes"].([]interface{})
		permTimesRaw, _ := permMap["permanentlyPermittedTimes"].([]interface{})
		forbTimesRaw, _ := permMap["permanentlyForbiddenTimes"].([]interface{})
		// All range arrays are optional (empty means all/no restrictions)
		transferTimesConverted, err := convertUintRangeArrayFromInterfaceWithError(transferTimesRaw, fmt.Sprintf("permission[%d].transferTimes", i))
		if err != nil {
			return nil, err
		}
		transferTimes, err := ConvertAndValidateBigIntRangesAllowEmpty(transferTimesConverted, fmt.Sprintf("permission[%d].transferTimes", i))
		if err != nil {
			return nil, err
		}
		tokenIdsConverted, err := convertUintRangeArrayFromInterfaceWithError(tokenIdsRaw, fmt.Sprintf("permission[%d].tokenIds", i))
		if err != nil {
			return nil, err
		}
		tokenIds, err := ConvertAndValidateBigIntRangesAllowEmpty(tokenIdsConverted, fmt.Sprintf("permission[%d].tokenIds", i))
		if err != nil {
			return nil, err
		}
		ownershipTimesConverted, err := convertUintRangeArrayFromInterfaceWithError(ownershipTimesRaw, fmt.Sprintf("permission[%d].ownershipTimes", i))
		if err != nil {
			return nil, err
		}
		ownershipTimes, err := ConvertAndValidateBigIntRangesAllowEmpty(ownershipTimesConverted, fmt.Sprintf("permission[%d].ownershipTimes", i))
		if err != nil {
			return nil, err
		}
		permTimesConverted, err := convertUintRangeArrayFromInterfaceWithError(permTimesRaw, fmt.Sprintf("permission[%d].permanentlyPermittedTimes", i))
		if err != nil {
			return nil, err
		}
		permTimes, err := ConvertAndValidateBigIntRangesAllowEmpty(permTimesConverted, fmt.Sprintf("permission[%d].permanentlyPermittedTimes", i))
		if err != nil {
			return nil, err
		}
		forbTimesConverted, err := convertUintRangeArrayFromInterfaceWithError(forbTimesRaw, fmt.Sprintf("permission[%d].permanentlyForbiddenTimes", i))
		if err != nil {
			return nil, err
		}
		forbTimes, err := ConvertAndValidateBigIntRangesAllowEmpty(forbTimesConverted, fmt.Sprintf("permission[%d].permanentlyForbiddenTimes", i))
		if err != nil {
			return nil, err
		}
		perms = append(perms, &tokenizationtypes.UserOutgoingApprovalPermission{
			ToListId:                  toListId,
			InitiatedByListId:         initiatedByListId,
			TransferTimes:             transferTimes,
			TokenIds:                  tokenIds,
			OwnershipTimes:            ownershipTimes,
			ApprovalId:                approvalId,
			PermanentlyPermittedTimes: permTimes,
			PermanentlyForbiddenTimes: forbTimes,
		})
	}
	return perms, nil
}

// ConvertUserIncomingApprovalPermissionArrayFromInterface converts []interface{} to []*UserIncomingApprovalPermission
func ConvertUserIncomingApprovalPermissionArrayFromInterface(permsRaw []interface{}) ([]*tokenizationtypes.UserIncomingApprovalPermission, error) {
	perms := make([]*tokenizationtypes.UserIncomingApprovalPermission, 0, len(permsRaw))
	for i, permRaw := range permsRaw {
		permMap, ok := permRaw.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("permission[%d] must be struct", i)
		}
		fromListId, _ := permMap["fromListId"].(string)
		initiatedByListId, _ := permMap["initiatedByListId"].(string)
		approvalId, _ := permMap["approvalId"].(string)
		transferTimesRaw, _ := permMap["transferTimes"].([]interface{})
		tokenIdsRaw, _ := permMap["tokenIds"].([]interface{})
		ownershipTimesRaw, _ := permMap["ownershipTimes"].([]interface{})
		permTimesRaw, _ := permMap["permanentlyPermittedTimes"].([]interface{})
		forbTimesRaw, _ := permMap["permanentlyForbiddenTimes"].([]interface{})
		// All range arrays are optional (empty means all/no restrictions)
		transferTimesConverted, err := convertUintRangeArrayFromInterfaceWithError(transferTimesRaw, fmt.Sprintf("permission[%d].transferTimes", i))
		if err != nil {
			return nil, err
		}
		transferTimes, err := ConvertAndValidateBigIntRangesAllowEmpty(transferTimesConverted, fmt.Sprintf("permission[%d].transferTimes", i))
		if err != nil {
			return nil, err
		}
		tokenIdsConverted, err := convertUintRangeArrayFromInterfaceWithError(tokenIdsRaw, fmt.Sprintf("permission[%d].tokenIds", i))
		if err != nil {
			return nil, err
		}
		tokenIds, err := ConvertAndValidateBigIntRangesAllowEmpty(tokenIdsConverted, fmt.Sprintf("permission[%d].tokenIds", i))
		if err != nil {
			return nil, err
		}
		ownershipTimesConverted, err := convertUintRangeArrayFromInterfaceWithError(ownershipTimesRaw, fmt.Sprintf("permission[%d].ownershipTimes", i))
		if err != nil {
			return nil, err
		}
		ownershipTimes, err := ConvertAndValidateBigIntRangesAllowEmpty(ownershipTimesConverted, fmt.Sprintf("permission[%d].ownershipTimes", i))
		if err != nil {
			return nil, err
		}
		permTimesConverted, err := convertUintRangeArrayFromInterfaceWithError(permTimesRaw, fmt.Sprintf("permission[%d].permanentlyPermittedTimes", i))
		if err != nil {
			return nil, err
		}
		permTimes, err := ConvertAndValidateBigIntRangesAllowEmpty(permTimesConverted, fmt.Sprintf("permission[%d].permanentlyPermittedTimes", i))
		if err != nil {
			return nil, err
		}
		forbTimesConverted, err := convertUintRangeArrayFromInterfaceWithError(forbTimesRaw, fmt.Sprintf("permission[%d].permanentlyForbiddenTimes", i))
		if err != nil {
			return nil, err
		}
		forbTimes, err := ConvertAndValidateBigIntRangesAllowEmpty(forbTimesConverted, fmt.Sprintf("permission[%d].permanentlyForbiddenTimes", i))
		if err != nil {
			return nil, err
		}
		perms = append(perms, &tokenizationtypes.UserIncomingApprovalPermission{
			FromListId:                fromListId,
			InitiatedByListId:         initiatedByListId,
			TransferTimes:             transferTimes,
			TokenIds:                  tokenIds,
			OwnershipTimes:            ownershipTimes,
			ApprovalId:                approvalId,
			PermanentlyPermittedTimes: permTimes,
			PermanentlyForbiddenTimes: forbTimes,
		})
	}
	return perms, nil
}

// ConvertCollectionApprovalArray converts array of CollectionApproval from interface{} slice
func ConvertCollectionApprovalArray(approvalsRaw []interface{}) ([]*tokenizationtypes.CollectionApproval, error) {
	collectionApprovals := make([]*tokenizationtypes.CollectionApproval, 0, len(approvalsRaw))
	for i, appRaw := range approvalsRaw {
		appMap, ok := appRaw.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("collectionApprovals[%d] must be struct", i)
		}

		approval := &tokenizationtypes.CollectionApproval{
			Version: sdkmath.NewUint(0),
		}

		// Extract basic string fields
		if val, ok := appMap["fromListId"].(string); ok {
			approval.FromListId = val
		}
		if val, ok := appMap["toListId"].(string); ok {
			approval.ToListId = val
		}
		if val, ok := appMap["initiatedByListId"].(string); ok {
			approval.InitiatedByListId = val
		}
		if val, ok := appMap["approvalId"].(string); ok {
			approval.ApprovalId = val
		}
		if val, ok := appMap["uri"].(string); ok {
			approval.Uri = val
		}
		if val, ok := appMap["customData"].(string); ok {
			approval.CustomData = val
		}

		// Convert ranges
		if transferTimesRaw, ok := appMap["transferTimes"].([]interface{}); ok {
			transferTimesConverted, err := convertUintRangeArrayFromInterfaceWithError(transferTimesRaw, fmt.Sprintf("collectionApprovals[%d].transferTimes", i))
			if err != nil {
				return nil, err
			}
			transferTimes, err := ConvertUintRangeArray(transferTimesConverted)
			if err != nil {
				return nil, fmt.Errorf("collectionApprovals[%d].transferTimes: %w", i, err)
			}
			approval.TransferTimes = transferTimes
		}
		if tokenIdsRaw, ok := appMap["tokenIds"].([]interface{}); ok {
			tokenIdsConverted, err := convertUintRangeArrayFromInterfaceWithError(tokenIdsRaw, fmt.Sprintf("collectionApprovals[%d].tokenIds", i))
			if err != nil {
				return nil, err
			}
			tokenIds, err := ConvertUintRangeArray(tokenIdsConverted)
			if err != nil {
				return nil, fmt.Errorf("collectionApprovals[%d].tokenIds: %w", i, err)
			}
			approval.TokenIds = tokenIds
		}
		if ownershipTimesRaw, ok := appMap["ownershipTimes"].([]interface{}); ok {
			ownershipTimesConverted, err := convertUintRangeArrayFromInterfaceWithError(ownershipTimesRaw, fmt.Sprintf("collectionApprovals[%d].ownershipTimes", i))
			if err != nil {
				return nil, err
			}
			ownershipTimes, err := ConvertUintRangeArray(ownershipTimesConverted)
			if err != nil {
				return nil, fmt.Errorf("collectionApprovals[%d].ownershipTimes: %w", i, err)
			}
			approval.OwnershipTimes = ownershipTimes
		}

		// Convert ApprovalCriteria if provided
		if criteriaRaw, ok := appMap["approvalCriteria"].(map[string]interface{}); ok {
			criteria, err := ConvertApprovalCriteria(criteriaRaw)
			if err != nil {
				return nil, fmt.Errorf("collectionApprovals[%d].approvalCriteria: %w", i, err)
			}
			approval.ApprovalCriteria = criteria
		} else {
			// Initialize empty criteria if not provided
			approval.ApprovalCriteria = &tokenizationtypes.ApprovalCriteria{}
		}

		collectionApprovals = append(collectionApprovals, approval)
	}
	return collectionApprovals, nil
}

// ConvertApprovalCriteria converts Solidity ApprovalCriteria to proto ApprovalCriteria
// Returns an error if any nested conversion fails.
func ConvertApprovalCriteria(criteriaMap map[string]interface{}) (*tokenizationtypes.ApprovalCriteria, error) {
	criteria := &tokenizationtypes.ApprovalCriteria{}

	// Boolean flags
	if val, ok := criteriaMap["requireToEqualsInitiatedBy"].(bool); ok {
		criteria.RequireToEqualsInitiatedBy = val
	}
	if val, ok := criteriaMap["requireFromEqualsInitiatedBy"].(bool); ok {
		criteria.RequireFromEqualsInitiatedBy = val
	}
	if val, ok := criteriaMap["requireToDoesNotEqualInitiatedBy"].(bool); ok {
		criteria.RequireToDoesNotEqualInitiatedBy = val
	}
	if val, ok := criteriaMap["requireFromDoesNotEqualInitiatedBy"].(bool); ok {
		criteria.RequireFromDoesNotEqualInitiatedBy = val
	}
	if val, ok := criteriaMap["overridesFromOutgoingApprovals"].(bool); ok {
		criteria.OverridesFromOutgoingApprovals = val
	}
	if val, ok := criteriaMap["overridesToIncomingApprovals"].(bool); ok {
		criteria.OverridesToIncomingApprovals = val
	}
	if val, ok := criteriaMap["mustPrioritize"].(bool); ok {
		criteria.MustPrioritize = val
	}
	if val, ok := criteriaMap["allowBackedMinting"].(bool); ok {
		criteria.AllowBackedMinting = val
	}
	if val, ok := criteriaMap["allowSpecialWrapping"].(bool); ok {
		criteria.AllowSpecialWrapping = val
	}

	// Convert MerkleChallenges
	if merkleChallengesRaw, ok := criteriaMap["merkleChallenges"].([]interface{}); ok {
		merkleChallenges, err := convertMerkleChallenges(merkleChallengesRaw)
		if err != nil {
			return nil, fmt.Errorf("approvalCriteria.merkleChallenges: %w", err)
		}
		criteria.MerkleChallenges = merkleChallenges
	}

	// Convert PredeterminedBalances
	if predeterminedBalancesRaw, ok := criteriaMap["predeterminedBalances"].(map[string]interface{}); ok {
		predeterminedBalances, err := convertPredeterminedBalances(predeterminedBalancesRaw)
		if err != nil {
			return nil, fmt.Errorf("approvalCriteria.predeterminedBalances: %w", err)
		}
		criteria.PredeterminedBalances = predeterminedBalances
	}

	// Convert ApprovalAmounts
	if approvalAmountsRaw, ok := criteriaMap["approvalAmounts"].(map[string]interface{}); ok {
		criteria.ApprovalAmounts = convertApprovalAmounts(approvalAmountsRaw)
	}

	// Convert MaxNumTransfers
	if maxNumTransfersRaw, ok := criteriaMap["maxNumTransfers"].(map[string]interface{}); ok {
		criteria.MaxNumTransfers = convertMaxNumTransfers(maxNumTransfersRaw)
	}

	// Convert CoinTransfers
	if coinTransfersRaw, ok := criteriaMap["coinTransfers"].([]interface{}); ok {
		coinTransfers, err := convertCoinTransfers(coinTransfersRaw)
		if err != nil {
			return nil, fmt.Errorf("approvalCriteria.coinTransfers: %w", err)
		}
		criteria.CoinTransfers = coinTransfers
	}

	// Convert AutoDeletionOptions
	if autoDeletionRaw, ok := criteriaMap["autoDeletionOptions"].(map[string]interface{}); ok {
		criteria.AutoDeletionOptions = convertAutoDeletionOptions(autoDeletionRaw)
	}

	// Convert UserRoyalties
	if userRoyaltiesRaw, ok := criteriaMap["userRoyalties"].(map[string]interface{}); ok {
		criteria.UserRoyalties = convertUserRoyalties(userRoyaltiesRaw)
	}

	// Convert MustOwnTokens
	if mustOwnTokensRaw, ok := criteriaMap["mustOwnTokens"].([]interface{}); ok {
		mustOwnTokens, err := convertMustOwnTokens(mustOwnTokensRaw)
		if err != nil {
			return nil, fmt.Errorf("approvalCriteria.mustOwnTokens: %w", err)
		}
		criteria.MustOwnTokens = mustOwnTokens
	}

	// Convert DynamicStoreChallenges
	if dynamicStoreChallengesRaw, ok := criteriaMap["dynamicStoreChallenges"].([]interface{}); ok {
		dynamicStoreChallenges, err := convertDynamicStoreChallenges(dynamicStoreChallengesRaw)
		if err != nil {
			return nil, fmt.Errorf("approvalCriteria.dynamicStoreChallenges: %w", err)
		}
		criteria.DynamicStoreChallenges = dynamicStoreChallenges
	}

	// Convert EthSignatureChallenges
	if ethSignatureChallengesRaw, ok := criteriaMap["ethSignatureChallenges"].([]interface{}); ok {
		ethSignatureChallenges, err := convertETHSignatureChallenges(ethSignatureChallengesRaw)
		if err != nil {
			return nil, fmt.Errorf("approvalCriteria.ethSignatureChallenges: %w", err)
		}
		criteria.EthSignatureChallenges = ethSignatureChallenges
	}

	// Convert address check types
	if senderChecksRaw, ok := criteriaMap["senderChecks"].(map[string]interface{}); ok {
		criteria.SenderChecks = convertAddressChecks(senderChecksRaw)
	}
	if recipientChecksRaw, ok := criteriaMap["recipientChecks"].(map[string]interface{}); ok {
		criteria.RecipientChecks = convertAddressChecks(recipientChecksRaw)
	}
	if initiatorChecksRaw, ok := criteriaMap["initiatorChecks"].(map[string]interface{}); ok {
		criteria.InitiatorChecks = convertAddressChecks(initiatorChecksRaw)
	}

	// Convert AltTimeChecks
	if altTimeChecksRaw, ok := criteriaMap["altTimeChecks"].(map[string]interface{}); ok {
		altTimeChecks, err := convertAltTimeChecks(altTimeChecksRaw)
		if err != nil {
			return nil, fmt.Errorf("approvalCriteria.altTimeChecks: %w", err)
		}
		criteria.AltTimeChecks = altTimeChecks
	}

	// Convert VotingChallenges
	if votingChallengesRaw, ok := criteriaMap["votingChallenges"].([]interface{}); ok {
		votingChallenges, err := convertVotingChallenges(votingChallengesRaw)
		if err != nil {
			return nil, fmt.Errorf("approvalCriteria.votingChallenges: %w", err)
		}
		criteria.VotingChallenges = votingChallenges
	}

	return criteria, nil
}

// convertMerkleChallenges converts array of MerkleChallenge
// Returns an error if any item in the array has an invalid type.
func convertMerkleChallenges(challengesRaw []interface{}) ([]*tokenizationtypes.MerkleChallenge, error) {
	challenges := make([]*tokenizationtypes.MerkleChallenge, 0, len(challengesRaw))
	for i, challengeRaw := range challengesRaw {
		challengeMap, ok := challengeRaw.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("merkleChallenges[%d] must be a map", i)
		}

		challenge := &tokenizationtypes.MerkleChallenge{}

		if val, ok := challengeMap["root"].(string); ok {
			challenge.Root = val
		}
		if val, ok := challengeMap["expectedProofLength"].(*big.Int); ok {
			challenge.ExpectedProofLength = sdkmath.NewUintFromBigInt(val)
		}
		if val, ok := challengeMap["useCreatorAddressAsLeaf"].(bool); ok {
			challenge.UseCreatorAddressAsLeaf = val
		}
		if val, ok := challengeMap["maxUsesPerLeaf"].(*big.Int); ok {
			challenge.MaxUsesPerLeaf = sdkmath.NewUintFromBigInt(val)
		}
		if val, ok := challengeMap["uri"].(string); ok {
			challenge.Uri = val
		}
		if val, ok := challengeMap["customData"].(string); ok {
			challenge.CustomData = val
		}
		if val, ok := challengeMap["challengeTrackerId"].(string); ok {
			challenge.ChallengeTrackerId = val
		}

		challenges = append(challenges, challenge)
	}
	return challenges, nil
}

// convertPredeterminedBalances converts PredeterminedBalances
// Returns an error if any field conversion fails.
func convertPredeterminedBalances(balancesRaw map[string]interface{}) (*tokenizationtypes.PredeterminedBalances, error) {
	balances := &tokenizationtypes.PredeterminedBalances{}

	if val, ok := balancesRaw["orderCalculationMethod"].(map[string]interface{}); ok {
		method := &tokenizationtypes.PredeterminedOrderCalculationMethod{}
		if useOverallNumTransfers, ok := val["useOverallNumTransfers"].(bool); ok {
			method.UseOverallNumTransfers = useOverallNumTransfers
		}
		if usePerToAddressNumTransfers, ok := val["usePerToAddressNumTransfers"].(bool); ok {
			method.UsePerToAddressNumTransfers = usePerToAddressNumTransfers
		}
		if usePerFromAddressNumTransfers, ok := val["usePerFromAddressNumTransfers"].(bool); ok {
			method.UsePerFromAddressNumTransfers = usePerFromAddressNumTransfers
		}
		if usePerInitiatedByAddressNumTransfers, ok := val["usePerInitiatedByAddressNumTransfers"].(bool); ok {
			method.UsePerInitiatedByAddressNumTransfers = usePerInitiatedByAddressNumTransfers
		}
		if useMerkleChallengeLeafIndex, ok := val["useMerkleChallengeLeafIndex"].(bool); ok {
			method.UseMerkleChallengeLeafIndex = useMerkleChallengeLeafIndex
		}
		if challengeTrackerId, ok := val["challengeTrackerId"].(string); ok {
			method.ChallengeTrackerId = challengeTrackerId
		}
		balances.OrderCalculationMethod = method
	}

	if incrementedBalancesRaw, ok := balancesRaw["incrementedBalances"].(map[string]interface{}); ok {
		incremented := &tokenizationtypes.IncrementedBalances{}
		if startBalancesRaw, ok := incrementedBalancesRaw["startBalances"].([]interface{}); ok {
			startBalances, err := ConvertBalanceArray(startBalancesRaw)
			if err != nil {
				return nil, fmt.Errorf("predeterminedBalances.incrementedBalances.startBalances: %w", err)
			}
			incremented.StartBalances = startBalances
		}
		if val, ok := incrementedBalancesRaw["incrementTokenIdsBy"].(*big.Int); ok {
			incremented.IncrementTokenIdsBy = sdkmath.NewUintFromBigInt(val)
		}
		if val, ok := incrementedBalancesRaw["incrementOwnershipTimesBy"].(*big.Int); ok {
			incremented.IncrementOwnershipTimesBy = sdkmath.NewUintFromBigInt(val)
		}
		if val, ok := incrementedBalancesRaw["durationFromTimestamp"].(*big.Int); ok {
			incremented.DurationFromTimestamp = sdkmath.NewUintFromBigInt(val)
		}
		if val, ok := incrementedBalancesRaw["allowOverrideTimestamp"].(bool); ok {
			incremented.AllowOverrideTimestamp = val
		}
		if val, ok := incrementedBalancesRaw["allowOverrideWithAnyValidToken"].(bool); ok {
			incremented.AllowOverrideWithAnyValidToken = val
		}
		balances.IncrementedBalances = incremented
	}

	return balances, nil
}

// convertApprovalAmounts converts ApprovalAmounts
func convertApprovalAmounts(amountsRaw map[string]interface{}) *tokenizationtypes.ApprovalAmounts {
	amounts := &tokenizationtypes.ApprovalAmounts{}

	if val, ok := amountsRaw["overallApprovalAmount"].(*big.Int); ok {
		amounts.OverallApprovalAmount = sdkmath.NewUintFromBigInt(val)
	}
	if val, ok := amountsRaw["perToAddressApprovalAmount"].(*big.Int); ok {
		amounts.PerToAddressApprovalAmount = sdkmath.NewUintFromBigInt(val)
	}
	if val, ok := amountsRaw["perFromAddressApprovalAmount"].(*big.Int); ok {
		amounts.PerFromAddressApprovalAmount = sdkmath.NewUintFromBigInt(val)
	}
	if val, ok := amountsRaw["perInitiatedByAddressApprovalAmount"].(*big.Int); ok {
		amounts.PerInitiatedByAddressApprovalAmount = sdkmath.NewUintFromBigInt(val)
	}
	if val, ok := amountsRaw["amountTrackerId"].(string); ok {
		amounts.AmountTrackerId = val
	}

	return amounts
}

// convertMaxNumTransfers converts MaxNumTransfers
func convertMaxNumTransfers(transfersRaw map[string]interface{}) *tokenizationtypes.MaxNumTransfers {
	transfers := &tokenizationtypes.MaxNumTransfers{}

	if val, ok := transfersRaw["overallMaxNumTransfers"].(*big.Int); ok {
		transfers.OverallMaxNumTransfers = sdkmath.NewUintFromBigInt(val)
	}
	if val, ok := transfersRaw["perToAddressMaxNumTransfers"].(*big.Int); ok {
		transfers.PerToAddressMaxNumTransfers = sdkmath.NewUintFromBigInt(val)
	}
	if val, ok := transfersRaw["perFromAddressMaxNumTransfers"].(*big.Int); ok {
		transfers.PerFromAddressMaxNumTransfers = sdkmath.NewUintFromBigInt(val)
	}
	if val, ok := transfersRaw["perInitiatedByAddressMaxNumTransfers"].(*big.Int); ok {
		transfers.PerInitiatedByAddressMaxNumTransfers = sdkmath.NewUintFromBigInt(val)
	}
	if val, ok := transfersRaw["amountTrackerId"].(string); ok {
		transfers.AmountTrackerId = val
	}

	return transfers
}

// convertCoinTransfers converts array of CoinTransfer
// Returns an error if any item in the array has an invalid type.
func convertCoinTransfers(transfersRaw []interface{}) ([]*tokenizationtypes.CoinTransfer, error) {
	transfers := make([]*tokenizationtypes.CoinTransfer, 0, len(transfersRaw))
	for i, transferRaw := range transfersRaw {
		transferMap, ok := transferRaw.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("coinTransfers[%d] must be a map", i)
		}

		transfer := &tokenizationtypes.CoinTransfer{}

		if val, ok := transferMap["to"].(string); ok {
			// Convert address if it's an EVM address
			if common.IsHexAddress(val) {
				evmAddr := common.HexToAddress(val)
				transfer.To = sdk.AccAddress(evmAddr.Bytes()).String()
			} else {
				transfer.To = val
			}
		}
		if coinsRaw, ok := transferMap["coins"].([]interface{}); ok {
			coins, err := convertCoinsToPointers(coinsRaw)
			if err != nil {
				return nil, fmt.Errorf("coinTransfers[%d].coins: %w", i, err)
			}
			transfer.Coins = coins
		}

		transfers = append(transfers, transfer)
	}
	return transfers, nil
}

// convertCoinsToPointers converts array of sdk.Coin to []*sdk.Coin
// Returns an error if any coin has an invalid type or missing required fields.
func convertCoinsToPointers(coinsRaw []interface{}) ([]*sdk.Coin, error) {
	coins := make([]*sdk.Coin, 0, len(coinsRaw))
	for i, coinRaw := range coinsRaw {
		coinMap, ok := coinRaw.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("coins[%d] must be a map", i)
		}

		denom, denomOk := coinMap["denom"].(string)
		amountBig, amountOk := coinMap["amount"].(*big.Int)
		if !denomOk || denom == "" {
			return nil, fmt.Errorf("coins[%d].denom must be a non-empty string", i)
		}
		if !amountOk || amountBig == nil {
			return nil, fmt.Errorf("coins[%d].amount must be a valid *big.Int", i)
		}
		amount := sdkmath.NewIntFromBigInt(amountBig)
		coin := sdk.NewCoin(denom, amount)
		coins = append(coins, &coin)
	}
	return coins, nil
}

// convertAutoDeletionOptions converts AutoDeletionOptions
func convertAutoDeletionOptions(optionsRaw map[string]interface{}) *tokenizationtypes.AutoDeletionOptions {
	options := &tokenizationtypes.AutoDeletionOptions{}

	if val, ok := optionsRaw["afterOneUse"].(bool); ok {
		options.AfterOneUse = val
	}
	if val, ok := optionsRaw["afterOverallMaxNumTransfers"].(bool); ok {
		options.AfterOverallMaxNumTransfers = val
	}
	if val, ok := optionsRaw["allowCounterpartyPurge"].(bool); ok {
		options.AllowCounterpartyPurge = val
	}
	if val, ok := optionsRaw["allowPurgeIfExpired"].(bool); ok {
		options.AllowPurgeIfExpired = val
	}

	return options
}

// convertUserRoyalties converts UserRoyalties
func convertUserRoyalties(royaltiesRaw map[string]interface{}) *tokenizationtypes.UserRoyalties {
	royalties := &tokenizationtypes.UserRoyalties{}

	if val, ok := royaltiesRaw["percentage"].(*big.Int); ok {
		royalties.Percentage = sdkmath.NewUintFromBigInt(val)
	}
	if val, ok := royaltiesRaw["payoutAddress"].(string); ok {
		// Convert address if it's an EVM address
		if common.IsHexAddress(val) {
			evmAddr := common.HexToAddress(val)
			royalties.PayoutAddress = sdk.AccAddress(evmAddr.Bytes()).String()
		} else {
			royalties.PayoutAddress = val
		}
	}

	return royalties
}

// convertMustOwnTokens converts array of MustOwnTokens
// Returns an error if any item in the array has an invalid type.
func convertMustOwnTokens(tokensRaw []interface{}) ([]*tokenizationtypes.MustOwnTokens, error) {
	tokens := make([]*tokenizationtypes.MustOwnTokens, 0, len(tokensRaw))
	for i, tokenRaw := range tokensRaw {
		tokenMap, ok := tokenRaw.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("mustOwnTokens[%d] must be a map", i)
		}

		token := &tokenizationtypes.MustOwnTokens{}

		if val, ok := tokenMap["collectionId"].(*big.Int); ok {
			token.CollectionId = sdkmath.NewUintFromBigInt(val)
		}
		if val, ok := tokenMap["amountRange"].(map[string]interface{}); ok {
			start, startOk := val["start"].(*big.Int)
			end, endOk := val["end"].(*big.Int)
			if !startOk || start == nil {
				return nil, fmt.Errorf("mustOwnTokens[%d].amountRange.start must be a valid *big.Int", i)
			}
			if !endOk || end == nil {
				return nil, fmt.Errorf("mustOwnTokens[%d].amountRange.end must be a valid *big.Int", i)
			}
			token.AmountRange = &tokenizationtypes.UintRange{
				Start: sdkmath.NewUintFromBigInt(start),
				End:   sdkmath.NewUintFromBigInt(end),
			}
		}
		if tokenIdsRaw, ok := tokenMap["tokenIds"].([]interface{}); ok {
			ranges, err := convertUintRangeArrayFromInterfaceWithError(tokenIdsRaw, fmt.Sprintf("mustOwnTokens[%d].tokenIds", i))
			if err != nil {
				return nil, err
			}
			token.TokenIds, err = ConvertUintRangeArray(ranges)
			if err != nil {
				return nil, fmt.Errorf("mustOwnTokens[%d].tokenIds: %w", i, err)
			}
		}
		if ownershipTimesRaw, ok := tokenMap["ownershipTimes"].([]interface{}); ok {
			ranges, err := convertUintRangeArrayFromInterfaceWithError(ownershipTimesRaw, fmt.Sprintf("mustOwnTokens[%d].ownershipTimes", i))
			if err != nil {
				return nil, err
			}
			token.OwnershipTimes, err = ConvertUintRangeArray(ranges)
			if err != nil {
				return nil, fmt.Errorf("mustOwnTokens[%d].ownershipTimes: %w", i, err)
			}
		}
		if val, ok := tokenMap["overrideWithCurrentTime"].(bool); ok {
			token.OverrideWithCurrentTime = val
		}
		if val, ok := tokenMap["mustSatisfyForAllAssets"].(bool); ok {
			token.MustSatisfyForAllAssets = val
		}
		if val, ok := tokenMap["ownershipCheckParty"].(string); ok {
			token.OwnershipCheckParty = val
		}

		tokens = append(tokens, token)
	}
	return tokens, nil
}

// convertDynamicStoreChallenges converts array of DynamicStoreChallenge
// Returns an error if any item in the array has an invalid type.
func convertDynamicStoreChallenges(challengesRaw []interface{}) ([]*tokenizationtypes.DynamicStoreChallenge, error) {
	challenges := make([]*tokenizationtypes.DynamicStoreChallenge, 0, len(challengesRaw))
	for i, challengeRaw := range challengesRaw {
		challengeMap, ok := challengeRaw.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("dynamicStoreChallenges[%d] must be a map", i)
		}

		challenge := &tokenizationtypes.DynamicStoreChallenge{}

		if val, ok := challengeMap["storeId"].(*big.Int); ok {
			challenge.StoreId = sdkmath.NewUintFromBigInt(val)
		}
		if val, ok := challengeMap["ownershipCheckParty"].(string); ok {
			challenge.OwnershipCheckParty = val
		}

		challenges = append(challenges, challenge)
	}
	return challenges, nil
}

// convertETHSignatureChallenges converts array of ETHSignatureChallenge
// Returns an error if any item in the array has an invalid type.
func convertETHSignatureChallenges(challengesRaw []interface{}) ([]*tokenizationtypes.ETHSignatureChallenge, error) {
	challenges := make([]*tokenizationtypes.ETHSignatureChallenge, 0, len(challengesRaw))
	for i, challengeRaw := range challengesRaw {
		challengeMap, ok := challengeRaw.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("ethSignatureChallenges[%d] must be a map", i)
		}

		challenge := &tokenizationtypes.ETHSignatureChallenge{}

		if val, ok := challengeMap["signer"].(string); ok {
			challenge.Signer = val
		}
		if val, ok := challengeMap["challengeTrackerId"].(string); ok {
			challenge.ChallengeTrackerId = val
		}
		if val, ok := challengeMap["uri"].(string); ok {
			challenge.Uri = val
		}
		if val, ok := challengeMap["customData"].(string); ok {
			challenge.CustomData = val
		}

		challenges = append(challenges, challenge)
	}
	return challenges, nil
}

// convertAddressChecks converts AddressChecks
func convertAddressChecks(checksRaw map[string]interface{}) *tokenizationtypes.AddressChecks {
	checks := &tokenizationtypes.AddressChecks{}

	if val, ok := checksRaw["mustBeEvmContract"].(bool); ok {
		checks.MustBeEvmContract = val
	}
	if val, ok := checksRaw["mustNotBeEvmContract"].(bool); ok {
		checks.MustNotBeEvmContract = val
	}
	if val, ok := checksRaw["mustBeLiquidityPool"].(bool); ok {
		checks.MustBeLiquidityPool = val
	}
	if val, ok := checksRaw["mustNotBeLiquidityPool"].(bool); ok {
		checks.MustNotBeLiquidityPool = val
	}

	return checks
}

// convertAltTimeChecks converts AltTimeChecks
// Returns an error if any range has invalid values.
func convertAltTimeChecks(checksRaw map[string]interface{}) (*tokenizationtypes.AltTimeChecks, error) {
	checks := &tokenizationtypes.AltTimeChecks{}

	if offlineHoursRaw, ok := checksRaw["offlineHours"].([]interface{}); ok {
		ranges, err := convertUintRangeArrayFromInterfaceWithError(offlineHoursRaw, "altTimeChecks.offlineHours")
		if err != nil {
			return nil, err
		}
		checks.OfflineHours, err = ConvertUintRangeArray(ranges)
		if err != nil {
			return nil, fmt.Errorf("altTimeChecks.offlineHours: %w", err)
		}
	}
	if offlineDaysRaw, ok := checksRaw["offlineDays"].([]interface{}); ok {
		ranges, err := convertUintRangeArrayFromInterfaceWithError(offlineDaysRaw, "altTimeChecks.offlineDays")
		if err != nil {
			return nil, err
		}
		checks.OfflineDays, err = ConvertUintRangeArray(ranges)
		if err != nil {
			return nil, fmt.Errorf("altTimeChecks.offlineDays: %w", err)
		}
	}

	return checks, nil
}

// convertVotingChallenges converts array of VotingChallenge
// Returns an error if any item in the array has an invalid type.
func convertVotingChallenges(challengesRaw []interface{}) ([]*tokenizationtypes.VotingChallenge, error) {
	challenges := make([]*tokenizationtypes.VotingChallenge, 0, len(challengesRaw))
	for i, challengeRaw := range challengesRaw {
		challengeMap, ok := challengeRaw.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("votingChallenges[%d] must be a map", i)
		}

		challenge := &tokenizationtypes.VotingChallenge{}

		if val, ok := challengeMap["proposalId"].(string); ok {
			challenge.ProposalId = val
		}
		if val, ok := challengeMap["quorumThreshold"].(*big.Int); ok {
			challenge.QuorumThreshold = sdkmath.NewUintFromBigInt(val)
		}
		// Voters would need their own conversion function - simplified for now
		if val, ok := challengeMap["uri"].(string); ok {
			challenge.Uri = val
		}
		if val, ok := challengeMap["customData"].(string); ok {
			challenge.CustomData = val
		}

		challenges = append(challenges, challenge)
	}
	return challenges, nil
}

// ConvertIncomingApprovalCriteria converts Solidity IncomingApprovalCriteria to proto IncomingApprovalCriteria
// Returns an error if any nested conversion fails.
func ConvertIncomingApprovalCriteria(criteriaMap map[string]interface{}) (*tokenizationtypes.IncomingApprovalCriteria, error) {
	if criteriaMap == nil {
		return nil, nil
	}

	criteria := &tokenizationtypes.IncomingApprovalCriteria{}

	// Boolean flags
	if val, ok := criteriaMap["requireFromEqualsInitiatedBy"].(bool); ok {
		criteria.RequireFromEqualsInitiatedBy = val
	}
	if val, ok := criteriaMap["requireFromDoesNotEqualInitiatedBy"].(bool); ok {
		criteria.RequireFromDoesNotEqualInitiatedBy = val
	}
	if val, ok := criteriaMap["mustPrioritize"].(bool); ok {
		criteria.MustPrioritize = val
	}

	// Convert MerkleChallenges
	if merkleChallengesRaw, ok := criteriaMap["merkleChallenges"].([]interface{}); ok {
		merkleChallenges, err := convertMerkleChallenges(merkleChallengesRaw)
		if err != nil {
			return nil, fmt.Errorf("incomingApprovalCriteria.merkleChallenges: %w", err)
		}
		criteria.MerkleChallenges = merkleChallenges
	}

	// Convert PredeterminedBalances
	if predeterminedBalancesRaw, ok := criteriaMap["predeterminedBalances"].(map[string]interface{}); ok {
		predeterminedBalances, err := convertPredeterminedBalances(predeterminedBalancesRaw)
		if err != nil {
			return nil, fmt.Errorf("incomingApprovalCriteria.predeterminedBalances: %w", err)
		}
		criteria.PredeterminedBalances = predeterminedBalances
	}

	// Convert ApprovalAmounts
	if approvalAmountsRaw, ok := criteriaMap["approvalAmounts"].(map[string]interface{}); ok {
		criteria.ApprovalAmounts = convertApprovalAmounts(approvalAmountsRaw)
	}

	// Convert MaxNumTransfers
	if maxNumTransfersRaw, ok := criteriaMap["maxNumTransfers"].(map[string]interface{}); ok {
		criteria.MaxNumTransfers = convertMaxNumTransfers(maxNumTransfersRaw)
	}

	// Convert CoinTransfers
	if coinTransfersRaw, ok := criteriaMap["coinTransfers"].([]interface{}); ok {
		coinTransfers, err := convertCoinTransfers(coinTransfersRaw)
		if err != nil {
			return nil, fmt.Errorf("incomingApprovalCriteria.coinTransfers: %w", err)
		}
		criteria.CoinTransfers = coinTransfers
	}

	// Convert AutoDeletionOptions
	if autoDeletionRaw, ok := criteriaMap["autoDeletionOptions"].(map[string]interface{}); ok {
		criteria.AutoDeletionOptions = convertAutoDeletionOptions(autoDeletionRaw)
	}

	// Convert MustOwnTokens
	if mustOwnTokensRaw, ok := criteriaMap["mustOwnTokens"].([]interface{}); ok {
		mustOwnTokens, err := convertMustOwnTokens(mustOwnTokensRaw)
		if err != nil {
			return nil, fmt.Errorf("incomingApprovalCriteria.mustOwnTokens: %w", err)
		}
		criteria.MustOwnTokens = mustOwnTokens
	}

	// Convert DynamicStoreChallenges
	if dynamicStoreChallengesRaw, ok := criteriaMap["dynamicStoreChallenges"].([]interface{}); ok {
		dynamicStoreChallenges, err := convertDynamicStoreChallenges(dynamicStoreChallengesRaw)
		if err != nil {
			return nil, fmt.Errorf("incomingApprovalCriteria.dynamicStoreChallenges: %w", err)
		}
		criteria.DynamicStoreChallenges = dynamicStoreChallenges
	}

	// Convert EthSignatureChallenges
	if ethSignatureChallengesRaw, ok := criteriaMap["ethSignatureChallenges"].([]interface{}); ok {
		ethSignatureChallenges, err := convertETHSignatureChallenges(ethSignatureChallengesRaw)
		if err != nil {
			return nil, fmt.Errorf("incomingApprovalCriteria.ethSignatureChallenges: %w", err)
		}
		criteria.EthSignatureChallenges = ethSignatureChallenges
	}

	// Convert address check types (IncomingApprovalCriteria has SenderChecks and InitiatorChecks, no RecipientChecks)
	if senderChecksRaw, ok := criteriaMap["senderChecks"].(map[string]interface{}); ok {
		criteria.SenderChecks = convertAddressChecks(senderChecksRaw)
	}
	if initiatorChecksRaw, ok := criteriaMap["initiatorChecks"].(map[string]interface{}); ok {
		criteria.InitiatorChecks = convertAddressChecks(initiatorChecksRaw)
	}

	// Convert AltTimeChecks
	if altTimeChecksRaw, ok := criteriaMap["altTimeChecks"].(map[string]interface{}); ok {
		altTimeChecks, err := convertAltTimeChecks(altTimeChecksRaw)
		if err != nil {
			return nil, fmt.Errorf("incomingApprovalCriteria.altTimeChecks: %w", err)
		}
		criteria.AltTimeChecks = altTimeChecks
	}

	// Convert VotingChallenges
	if votingChallengesRaw, ok := criteriaMap["votingChallenges"].([]interface{}); ok {
		votingChallenges, err := convertVotingChallenges(votingChallengesRaw)
		if err != nil {
			return nil, fmt.Errorf("incomingApprovalCriteria.votingChallenges: %w", err)
		}
		criteria.VotingChallenges = votingChallenges
	}

	return criteria, nil
}

// ConvertOutgoingApprovalCriteria converts Solidity OutgoingApprovalCriteria to proto OutgoingApprovalCriteria
// Returns an error if any nested conversion fails.
func ConvertOutgoingApprovalCriteria(criteriaMap map[string]interface{}) (*tokenizationtypes.OutgoingApprovalCriteria, error) {
	if criteriaMap == nil {
		return nil, nil
	}

	criteria := &tokenizationtypes.OutgoingApprovalCriteria{}

	// Boolean flags
	if val, ok := criteriaMap["requireToEqualsInitiatedBy"].(bool); ok {
		criteria.RequireToEqualsInitiatedBy = val
	}
	if val, ok := criteriaMap["requireToDoesNotEqualInitiatedBy"].(bool); ok {
		criteria.RequireToDoesNotEqualInitiatedBy = val
	}
	if val, ok := criteriaMap["mustPrioritize"].(bool); ok {
		criteria.MustPrioritize = val
	}

	// Convert MerkleChallenges
	if merkleChallengesRaw, ok := criteriaMap["merkleChallenges"].([]interface{}); ok {
		merkleChallenges, err := convertMerkleChallenges(merkleChallengesRaw)
		if err != nil {
			return nil, fmt.Errorf("outgoingApprovalCriteria.merkleChallenges: %w", err)
		}
		criteria.MerkleChallenges = merkleChallenges
	}

	// Convert PredeterminedBalances
	if predeterminedBalancesRaw, ok := criteriaMap["predeterminedBalances"].(map[string]interface{}); ok {
		predeterminedBalances, err := convertPredeterminedBalances(predeterminedBalancesRaw)
		if err != nil {
			return nil, fmt.Errorf("outgoingApprovalCriteria.predeterminedBalances: %w", err)
		}
		criteria.PredeterminedBalances = predeterminedBalances
	}

	// Convert ApprovalAmounts
	if approvalAmountsRaw, ok := criteriaMap["approvalAmounts"].(map[string]interface{}); ok {
		criteria.ApprovalAmounts = convertApprovalAmounts(approvalAmountsRaw)
	}

	// Convert MaxNumTransfers
	if maxNumTransfersRaw, ok := criteriaMap["maxNumTransfers"].(map[string]interface{}); ok {
		criteria.MaxNumTransfers = convertMaxNumTransfers(maxNumTransfersRaw)
	}

	// Convert CoinTransfers
	if coinTransfersRaw, ok := criteriaMap["coinTransfers"].([]interface{}); ok {
		coinTransfers, err := convertCoinTransfers(coinTransfersRaw)
		if err != nil {
			return nil, fmt.Errorf("outgoingApprovalCriteria.coinTransfers: %w", err)
		}
		criteria.CoinTransfers = coinTransfers
	}

	// Convert AutoDeletionOptions
	if autoDeletionRaw, ok := criteriaMap["autoDeletionOptions"].(map[string]interface{}); ok {
		criteria.AutoDeletionOptions = convertAutoDeletionOptions(autoDeletionRaw)
	}

	// Convert MustOwnTokens
	if mustOwnTokensRaw, ok := criteriaMap["mustOwnTokens"].([]interface{}); ok {
		mustOwnTokens, err := convertMustOwnTokens(mustOwnTokensRaw)
		if err != nil {
			return nil, fmt.Errorf("outgoingApprovalCriteria.mustOwnTokens: %w", err)
		}
		criteria.MustOwnTokens = mustOwnTokens
	}

	// Convert DynamicStoreChallenges
	if dynamicStoreChallengesRaw, ok := criteriaMap["dynamicStoreChallenges"].([]interface{}); ok {
		dynamicStoreChallenges, err := convertDynamicStoreChallenges(dynamicStoreChallengesRaw)
		if err != nil {
			return nil, fmt.Errorf("outgoingApprovalCriteria.dynamicStoreChallenges: %w", err)
		}
		criteria.DynamicStoreChallenges = dynamicStoreChallenges
	}

	// Convert EthSignatureChallenges
	if ethSignatureChallengesRaw, ok := criteriaMap["ethSignatureChallenges"].([]interface{}); ok {
		ethSignatureChallenges, err := convertETHSignatureChallenges(ethSignatureChallengesRaw)
		if err != nil {
			return nil, fmt.Errorf("outgoingApprovalCriteria.ethSignatureChallenges: %w", err)
		}
		criteria.EthSignatureChallenges = ethSignatureChallenges
	}

	// Convert address check types (OutgoingApprovalCriteria has RecipientChecks and InitiatorChecks, no SenderChecks)
	if recipientChecksRaw, ok := criteriaMap["recipientChecks"].(map[string]interface{}); ok {
		criteria.RecipientChecks = convertAddressChecks(recipientChecksRaw)
	}
	if initiatorChecksRaw, ok := criteriaMap["initiatorChecks"].(map[string]interface{}); ok {
		criteria.InitiatorChecks = convertAddressChecks(initiatorChecksRaw)
	}

	// Convert AltTimeChecks
	if altTimeChecksRaw, ok := criteriaMap["altTimeChecks"].(map[string]interface{}); ok {
		altTimeChecks, err := convertAltTimeChecks(altTimeChecksRaw)
		if err != nil {
			return nil, fmt.Errorf("outgoingApprovalCriteria.altTimeChecks: %w", err)
		}
		criteria.AltTimeChecks = altTimeChecks
	}

	// Convert VotingChallenges
	if votingChallengesRaw, ok := criteriaMap["votingChallenges"].([]interface{}); ok {
		votingChallenges, err := convertVotingChallenges(votingChallengesRaw)
		if err != nil {
			return nil, fmt.Errorf("outgoingApprovalCriteria.votingChallenges: %w", err)
		}
		criteria.VotingChallenges = votingChallenges
	}

	return criteria, nil
}
