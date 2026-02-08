package tokenization

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	tokenizationkeeper "github.com/bitbadges/bitbadgeschain/x/tokenization/keeper"
	tokenizationtypes "github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

// CreateCollection creates a new collection via the tokenization module.
func (p Precompile) CreateCollection(ctx sdk.Context, method *abi.Method, args []interface{}, contract *vm.Contract) ([]byte, error) {
	if len(args) < 6 {
		return nil, ErrInvalidInput(fmt.Sprintf("invalid number of arguments, expected at least 6, got %d", len(args)))
	}

	creatorCosmosAddr, err := p.GetCallerAddress(contract)
	if err != nil {
		return nil, err
	}

	// Extract arguments - simplified version accepting key fields
	// For full implementation, would need to handle all nested structs from TokenizationTypes
	// Type assertions with error checking to prevent silent data corruption
	defaultBalancesRaw, ok := args[0].(map[string]interface{})
	if !ok && args[0] != nil {
		return nil, ErrInvalidInput("invalid defaultBalances type, expected map[string]interface{}")
	}
	validTokenIdsRaw, ok := args[1].([]interface{})
	if !ok && args[1] != nil {
		return nil, ErrInvalidInput("invalid validTokenIds type, expected []interface{}")
	}
	collectionPermissionsRaw, ok := args[2].(map[string]interface{})
	if !ok && args[2] != nil {
		return nil, ErrInvalidInput("invalid collectionPermissions type, expected map[string]interface{}")
	}
	manager, ok := args[3].(string)
	if !ok && args[3] != nil {
		return nil, ErrInvalidInput("invalid manager type, expected string")
	}
	collectionMetadataRaw, ok := args[4].(map[string]interface{})
	if !ok && args[4] != nil {
		return nil, ErrInvalidInput("invalid collectionMetadata type, expected map[string]interface{}")
	}
	tokenMetadataRaw, ok := args[5].([]interface{})
	if !ok && args[5] != nil {
		return nil, ErrInvalidInput("invalid tokenMetadata type, expected []interface{}")
	}

	// Convert defaultBalances
	var defaultBalances *tokenizationtypes.UserBalanceStore
	if defaultBalancesRaw != nil {
		defaultBalances, err = ConvertUserBalanceStore(defaultBalancesRaw)
		if err != nil {
			return nil, ErrInvalidInput(fmt.Sprintf("defaultBalances: %v", err))
		}
	} else {
		defaultBalances = &tokenizationtypes.UserBalanceStore{
			UserPermissions: &tokenizationtypes.UserPermissions{},
		}
	}

	// Convert validTokenIds using the error-returning helper
	var convertedValidTokenIds []*tokenizationtypes.UintRange
	if validTokenIdsRaw != nil {
		ranges, err := convertUintRangeArrayFromInterfaceWithError(validTokenIdsRaw, "validTokenIds")
		if err != nil {
			return nil, ErrInvalidInput(err.Error())
		}
		convertedValidTokenIds, err = ConvertUintRangeArray(ranges)
		if err != nil {
			return nil, ErrInvalidInput(fmt.Sprintf("validTokenIds: %v", err))
		}
	}

	// Convert collectionPermissions
	var collectionPermissions *tokenizationtypes.CollectionPermissions
	if collectionPermissionsRaw != nil {
		collectionPermissions, err = ConvertCollectionPermissions(collectionPermissionsRaw)
		if err != nil {
			return nil, ErrInvalidInput(fmt.Sprintf("collectionPermissions: %v", err))
		}
	} else {
		collectionPermissions = &tokenizationtypes.CollectionPermissions{}
	}

	// Convert manager address
	managerCosmosAddr, err := ConvertManagerAddress(manager)
	if err != nil {
		return nil, ErrInvalidInput(fmt.Sprintf("manager: %v", err))
	}

	// Convert collectionMetadata
	var collectionMetadata *tokenizationtypes.CollectionMetadata
	if collectionMetadataRaw != nil {
		uri, _ := collectionMetadataRaw["uri"].(string)
		customData, _ := collectionMetadataRaw["customData"].(string)
		collectionMetadata = ConvertCollectionMetadata(uri, customData)
	} else {
		collectionMetadata = &tokenizationtypes.CollectionMetadata{}
	}

	// Convert tokenMetadata
	tokenMetadata := make([]*tokenizationtypes.TokenMetadata, 0)
	if tokenMetadataRaw != nil {
		for i, tmRaw := range tokenMetadataRaw {
			tmMap, ok := tmRaw.(map[string]interface{})
			if !ok {
				return nil, ErrInvalidInput(fmt.Sprintf("tokenMetadata[%d] must be a map", i))
			}
			uri, _ := tmMap["uri"].(string)
			customDataVal, _ := tmMap["customData"].(string)
			tokenIdsRaw, _ := tmMap["tokenIds"].([]interface{})
			tokenIds, err := convertUintRangeArrayFromInterfaceWithError(tokenIdsRaw, fmt.Sprintf("tokenMetadata[%d].tokenIds", i))
			if err != nil {
				return nil, ErrInvalidInput(err.Error())
			}
			tm, err := ConvertTokenMetadata(uri, customDataVal, tokenIds)
			if err != nil {
				return nil, ErrInvalidInput(fmt.Sprintf("tokenMetadata[%d]: %v", i, err))
			}
			tokenMetadata = append(tokenMetadata, tm)
		}
	}

	// Extract optional fields with bounds checking and explicit type validation
	var customData string
	var collectionApprovalsRaw []interface{}
	var standardsRaw []string
	var isArchived bool

	if len(args) > 6 && args[6] != nil {
		if val, ok := args[6].(string); ok {
			customData = val
		} else {
			return nil, ErrInvalidInput("invalid customData type, expected string")
		}
	}
	if len(args) > 7 && args[7] != nil {
		if val, ok := args[7].([]interface{}); ok {
			collectionApprovalsRaw = val
		} else {
			return nil, ErrInvalidInput("invalid collectionApprovals type, expected []interface{}")
		}
	}
	if len(args) > 8 && args[8] != nil {
		if val, ok := args[8].([]string); ok {
			standardsRaw = val
		} else {
			return nil, ErrInvalidInput("invalid standards type, expected []string")
		}
	}
	if len(args) > 9 && args[9] != nil {
		if val, ok := args[9].(bool); ok {
			isArchived = val
		} else {
			return nil, ErrInvalidInput("invalid isArchived type, expected bool")
		}
	}

	// Convert collectionApprovals using the dedicated array converter
	var collectionApprovals []*tokenizationtypes.CollectionApproval
	if collectionApprovalsRaw != nil {
		var err error
		collectionApprovals, err = ConvertCollectionApprovalArray(collectionApprovalsRaw)
		if err != nil {
			return nil, ErrInvalidInput(fmt.Sprintf("collectionApprovals: %v", err))
		}
	}

	// Convert standards
	standards := make([]string, 0)
	if standardsRaw != nil {
		standards = standardsRaw
	}

	// Extract invariants and paths (simplified - empty for now)
	// NOTE: invariants and paths are initialized as empty and not settable through the precompile
	// This is an intentional limitation. To set these, use the native Cosmos SDK interface.
	// See production readiness evaluation for details.
	invariants := &tokenizationtypes.InvariantsAddObject{}
	cosmosCoinWrapperPaths := []*tokenizationtypes.CosmosCoinWrapperPathAddObject{}
	aliasPaths := []*tokenizationtypes.AliasPathAddObject{}

	msg := &tokenizationtypes.MsgCreateCollection{
		Creator:                     creatorCosmosAddr,
		DefaultBalances:             defaultBalances,
		ValidTokenIds:               convertedValidTokenIds,
		CollectionPermissions:       collectionPermissions,
		Manager:                     managerCosmosAddr,
		CollectionMetadata:          collectionMetadata,
		TokenMetadata:               tokenMetadata,
		CustomData:                  customData,
		CollectionApprovals:         collectionApprovals,
		Standards:                   standards,
		IsArchived:                  isArchived,
		Invariants:                  invariants,
		CosmosCoinWrapperPathsToAdd: cosmosCoinWrapperPaths,
		AliasPathsToAdd:             aliasPaths,
	}

	msgServer := tokenizationkeeper.NewMsgServerImpl(p.tokenizationKeeper)
	resp, err := msgServer.CreateCollection(ctx, msg)
	if err != nil {
		return nil, WrapError(err, ErrorCodeInternalError, "create collection failed")
	}

	return method.Outputs.Pack(resp.CollectionId.BigInt())
}

// UpdateCollection updates an existing collection.
func (p Precompile) UpdateCollection(ctx sdk.Context, method *abi.Method, args []interface{}, contract *vm.Contract) ([]byte, error) {
	if len(args) < 2 {
		return nil, ErrInvalidInput(fmt.Sprintf("invalid number of arguments, expected at least 2, got %d", len(args)))
	}

	creatorCosmosAddr, err := p.GetCallerAddress(contract)
	if err != nil {
		return nil, err
	}

	collectionIdBig, ok := args[0].(*big.Int)
	if !ok {
		return nil, ErrInvalidInput("invalid collectionId type, expected *big.Int")
	}

	if err := ValidateCollectionId(collectionIdBig); err != nil {
		return nil, err
	}

	collectionId := sdkmath.NewUintFromBigInt(collectionIdBig)

	// Extract update flags and values with type validation
	updateValidTokenIds, ok := args[1].(bool)
	if !ok {
		return nil, ErrInvalidInput("invalid updateValidTokenIds type, expected bool")
	}
	validTokenIdsRaw, ok := args[2].([]interface{})
	if !ok && args[2] != nil {
		return nil, ErrInvalidInput("invalid validTokenIds type, expected []interface{}")
	}
	updateCollectionPermissions, ok := args[3].(bool)
	if !ok {
		return nil, ErrInvalidInput("invalid updateCollectionPermissions type, expected bool")
	}
	collectionPermissionsRaw, ok := args[4].(map[string]interface{})
	if !ok && args[4] != nil {
		return nil, ErrInvalidInput("invalid collectionPermissions type, expected map[string]interface{}")
	}
	updateManager, ok := args[5].(bool)
	if !ok {
		return nil, ErrInvalidInput("invalid updateManager type, expected bool")
	}
	manager, ok := args[6].(string)
	if !ok && args[6] != nil {
		return nil, ErrInvalidInput("invalid manager type, expected string")
	}
	updateCollectionMetadata, ok := args[7].(bool)
	if !ok {
		return nil, ErrInvalidInput("invalid updateCollectionMetadata type, expected bool")
	}
	collectionMetadataRaw, ok := args[8].(map[string]interface{})
	if !ok && args[8] != nil {
		return nil, ErrInvalidInput("invalid collectionMetadata type, expected map[string]interface{}")
	}
	updateTokenMetadata, ok := args[9].(bool)
	if !ok {
		return nil, ErrInvalidInput("invalid updateTokenMetadata type, expected bool")
	}
	tokenMetadataRaw, ok := args[10].([]interface{})
	if !ok && args[10] != nil {
		return nil, ErrInvalidInput("invalid tokenMetadata type, expected []interface{}")
	}
	updateCustomData, ok := args[11].(bool)
	if !ok {
		return nil, ErrInvalidInput("invalid updateCustomData type, expected bool")
	}
	customData, ok := args[12].(string)
	if !ok && args[12] != nil {
		return nil, ErrInvalidInput("invalid customData type, expected string")
	}
	updateCollectionApprovals, ok := args[13].(bool)
	if !ok {
		return nil, ErrInvalidInput("invalid updateCollectionApprovals type, expected bool")
	}
	collectionApprovalsRaw, ok := args[14].([]interface{})
	if !ok && args[14] != nil {
		return nil, ErrInvalidInput("invalid collectionApprovals type, expected []interface{}")
	}
	updateStandards, ok := args[15].(bool)
	if !ok {
		return nil, ErrInvalidInput("invalid updateStandards type, expected bool")
	}
	standardsRaw, ok := args[16].([]string)
	if !ok && args[16] != nil {
		return nil, ErrInvalidInput("invalid standards type, expected []string")
	}
	updateIsArchived, ok := args[17].(bool)
	if !ok {
		return nil, ErrInvalidInput("invalid updateIsArchived type, expected bool")
	}
	isArchived, ok := args[18].(bool)
	if !ok && args[18] != nil {
		return nil, ErrInvalidInput("invalid isArchived type, expected bool")
	}

	msg := &tokenizationtypes.MsgUpdateCollection{
		Creator:      creatorCosmosAddr,
		CollectionId: collectionId,
	}

	// Set update flags
	msg.UpdateValidTokenIds = updateValidTokenIds
	msg.UpdateCollectionPermissions = updateCollectionPermissions
	msg.UpdateManager = updateManager
	msg.UpdateCollectionMetadata = updateCollectionMetadata
	msg.UpdateTokenMetadata = updateTokenMetadata
	msg.UpdateCustomData = updateCustomData
	msg.UpdateCollectionApprovals = updateCollectionApprovals
	msg.UpdateStandards = updateStandards
	msg.UpdateIsArchived = updateIsArchived

	// Convert and set values if update flags are true
	if updateValidTokenIds && validTokenIdsRaw != nil {
		validTokenIds, err := convertUintRangeArrayFromInterfaceWithError(validTokenIdsRaw, "validTokenIds")
		if err != nil {
			return nil, ErrInvalidInput(err.Error())
		}
		msg.ValidTokenIds, err = ConvertUintRangeArray(validTokenIds)
		if err != nil {
			return nil, ErrInvalidInput(fmt.Sprintf("validTokenIds: %v", err))
		}
	}

	if updateCollectionPermissions && collectionPermissionsRaw != nil {
		var err error
		msg.CollectionPermissions, err = ConvertCollectionPermissions(collectionPermissionsRaw)
		if err != nil {
			return nil, ErrInvalidInput(fmt.Sprintf("collectionPermissions: %v", err))
		}
	}

	if updateManager {
		managerCosmosAddr, err := ConvertManagerAddress(manager)
		if err != nil {
			return nil, ErrInvalidInput(fmt.Sprintf("manager: %v", err))
		}
		msg.Manager = managerCosmosAddr
	}

	if updateCollectionMetadata && collectionMetadataRaw != nil {
		uri, _ := collectionMetadataRaw["uri"].(string)
		customData, _ := collectionMetadataRaw["customData"].(string)
		msg.CollectionMetadata = ConvertCollectionMetadata(uri, customData)
	}

	if updateTokenMetadata && tokenMetadataRaw != nil {
		tokenMetadata := make([]*tokenizationtypes.TokenMetadata, 0)
		for i, tmRaw := range tokenMetadataRaw {
			tmMap, ok := tmRaw.(map[string]interface{})
			if !ok {
				return nil, ErrInvalidInput(fmt.Sprintf("tokenMetadata[%d] must be a map", i))
			}
			uri, _ := tmMap["uri"].(string)
			customDataVal, _ := tmMap["customData"].(string)
			tokenIdsRaw, _ := tmMap["tokenIds"].([]interface{})
			tokenIds, err := convertUintRangeArrayFromInterfaceWithError(tokenIdsRaw, fmt.Sprintf("tokenMetadata[%d].tokenIds", i))
			if err != nil {
				return nil, ErrInvalidInput(err.Error())
			}
			tm, err := ConvertTokenMetadata(uri, customDataVal, tokenIds)
			if err != nil {
				return nil, ErrInvalidInput(fmt.Sprintf("tokenMetadata[%d]: %v", i, err))
			}
			tokenMetadata = append(tokenMetadata, tm)
		}
		msg.TokenMetadata = tokenMetadata
	}

	if updateCustomData {
		msg.CustomData = customData
	}

	if updateCollectionApprovals && collectionApprovalsRaw != nil {
		collectionApprovals, err := ConvertCollectionApprovalArray(collectionApprovalsRaw)
		if err != nil {
			return nil, ErrInvalidInput(fmt.Sprintf("collectionApprovals: %v", err))
		}
		msg.CollectionApprovals = collectionApprovals
	}

	if updateStandards && standardsRaw != nil {
		msg.Standards = standardsRaw
	}

	if updateIsArchived {
		msg.IsArchived = isArchived
	}

	// Set empty invariants and paths (can be expanded)
	msg.Invariants = &tokenizationtypes.InvariantsAddObject{}
	msg.CosmosCoinWrapperPathsToAdd = []*tokenizationtypes.CosmosCoinWrapperPathAddObject{}
	msg.AliasPathsToAdd = []*tokenizationtypes.AliasPathAddObject{}

	msgServer := tokenizationkeeper.NewMsgServerImpl(p.tokenizationKeeper)
	resp, err := msgServer.UpdateCollection(ctx, msg)
	if err != nil {
		return nil, WrapError(err, ErrorCodeCollectionNotFound, "update collection failed")
	}

	return method.Outputs.Pack(resp.CollectionId.BigInt())
}

// DeleteCollection deletes a collection.
func (p Precompile) DeleteCollection(ctx sdk.Context, method *abi.Method, args []interface{}, contract *vm.Contract) ([]byte, error) {
	if len(args) != 1 {
		return nil, ErrInvalidInput(fmt.Sprintf("invalid number of arguments, expected 1, got %d", len(args)))
	}

	collectionIdBig, ok := args[0].(*big.Int)
	if !ok {
		return nil, ErrInvalidInput("invalid collectionId type, expected *big.Int")
	}

	if err := ValidateCollectionId(collectionIdBig); err != nil {
		return nil, err
	}

	creatorCosmosAddr, err := p.GetCallerAddress(contract)
	if err != nil {
		return nil, err
	}

	collectionId := sdkmath.NewUintFromBigInt(collectionIdBig)

	msg := &tokenizationtypes.MsgDeleteCollection{
		Creator:      creatorCosmosAddr,
		CollectionId: collectionId,
	}

	msgServer := tokenizationkeeper.NewMsgServerImpl(p.tokenizationKeeper)
	_, err = msgServer.DeleteCollection(ctx, msg)
	if err != nil {
		return nil, WrapError(err, ErrorCodeCollectionNotFound, "delete collection failed")
	}

	return method.Outputs.Pack(true)
}

// CreateAddressLists creates address lists.
func (p Precompile) CreateAddressLists(ctx sdk.Context, method *abi.Method, args []interface{}, contract *vm.Contract) ([]byte, error) {
	if len(args) != 1 {
		return nil, ErrInvalidInput(fmt.Sprintf("invalid number of arguments, expected 1, got %d", len(args)))
	}

	// Extract address lists array from args
	// The ABI will pass this as an array of structs
	addressListsRaw, ok := args[0].([]interface{})
	if !ok {
		return nil, ErrInvalidInput("invalid addressLists type, expected array")
	}

	if len(addressListsRaw) == 0 {
		return nil, ErrInvalidInput("addressLists cannot be empty")
	}

	creatorCosmosAddr, err := p.GetCallerAddress(contract)
	if err != nil {
		return nil, err
	}

	// Convert each address list
	addressLists := make([]*tokenizationtypes.AddressListInput, 0, len(addressListsRaw))
	for i, listRaw := range addressListsRaw {
		// The ABI will pass structs as maps or specific struct types
		// For now, we'll handle it as a map[string]interface{}
		listMap, ok := listRaw.(map[string]interface{})
		if !ok {
			return nil, ErrInvalidInput(fmt.Sprintf("addressLists[%d] must be a struct", i))
		}

		listId, _ := listMap["listId"].(string)
		addressesRaw, _ := listMap["addresses"].([]interface{})
		whitelist, _ := listMap["whitelist"].(bool)
		uri, _ := listMap["uri"].(string)
		customData, _ := listMap["customData"].(string)

		// Convert addresses
		addresses := make([]string, len(addressesRaw))
		for j, addrRaw := range addressesRaw {
			if addrStr, ok := addrRaw.(string); ok {
				addresses[j] = addrStr
			} else if addrAddr, ok := addrRaw.(common.Address); ok {
				addresses[j] = sdk.AccAddress(addrAddr.Bytes()).String()
			} else {
				return nil, ErrInvalidInput(fmt.Sprintf("addressLists[%d].addresses[%d] must be string or address", i, j))
			}
		}

		addressList := ConvertAddressListInput(listId, addresses, whitelist, uri, customData)
		addressLists = append(addressLists, addressList)
	}

	msg := &tokenizationtypes.MsgCreateAddressLists{
		Creator:      creatorCosmosAddr,
		AddressLists: addressLists,
	}

	msgServer := tokenizationkeeper.NewMsgServerImpl(p.tokenizationKeeper)
	_, err = msgServer.CreateAddressLists(ctx, msg)
	if err != nil {
		return nil, WrapError(err, ErrorCodeInternalError, "create address lists failed")
	}

	return method.Outputs.Pack(true)
}

// UpdateUserApprovals updates user approvals.
func (p Precompile) UpdateUserApprovals(ctx sdk.Context, method *abi.Method, args []interface{}, contract *vm.Contract) ([]byte, error) {
	if len(args) < 3 {
		return nil, ErrInvalidInput(fmt.Sprintf("invalid number of arguments, expected at least 3, got %d", len(args)))
	}

	creatorCosmosAddr, err := p.GetCallerAddress(contract)
	if err != nil {
		return nil, err
	}

	collectionIdBig, ok := args[0].(*big.Int)
	if !ok {
		return nil, ErrInvalidInput("invalid collectionId type, expected *big.Int")
	}

	if err := ValidateCollectionId(collectionIdBig); err != nil {
		return nil, err
	}

	collectionId := sdkmath.NewUintFromBigInt(collectionIdBig)

	// Extract update flags and values with type validation
	updateOutgoingApprovals, ok := args[1].(bool)
	if !ok {
		return nil, ErrInvalidInput("invalid updateOutgoingApprovals type, expected bool")
	}
	outgoingApprovalsRaw, ok := args[2].([]interface{})
	if !ok && args[2] != nil {
		return nil, ErrInvalidInput("invalid outgoingApprovals type, expected []interface{}")
	}
	updateIncomingApprovals, ok := args[3].(bool)
	if !ok {
		return nil, ErrInvalidInput("invalid updateIncomingApprovals type, expected bool")
	}
	incomingApprovalsRaw, ok := args[4].([]interface{})
	if !ok && args[4] != nil {
		return nil, ErrInvalidInput("invalid incomingApprovals type, expected []interface{}")
	}
	updateAutoApproveSelfInitiatedOutgoingTransfers, ok := args[5].(bool)
	if !ok {
		return nil, ErrInvalidInput("invalid updateAutoApproveSelfInitiatedOutgoingTransfers type, expected bool")
	}
	autoApproveSelfInitiatedOutgoingTransfers, ok := args[6].(bool)
	if !ok && args[6] != nil {
		return nil, ErrInvalidInput("invalid autoApproveSelfInitiatedOutgoingTransfers type, expected bool")
	}
	updateAutoApproveSelfInitiatedIncomingTransfers, ok := args[7].(bool)
	if !ok {
		return nil, ErrInvalidInput("invalid updateAutoApproveSelfInitiatedIncomingTransfers type, expected bool")
	}
	autoApproveSelfInitiatedIncomingTransfers, ok := args[8].(bool)
	if !ok && args[8] != nil {
		return nil, ErrInvalidInput("invalid autoApproveSelfInitiatedIncomingTransfers type, expected bool")
	}
	updateAutoApproveAllIncomingTransfers, ok := args[9].(bool)
	if !ok {
		return nil, ErrInvalidInput("invalid updateAutoApproveAllIncomingTransfers type, expected bool")
	}
	autoApproveAllIncomingTransfers, ok := args[10].(bool)
	if !ok && args[10] != nil {
		return nil, ErrInvalidInput("invalid autoApproveAllIncomingTransfers type, expected bool")
	}
	updateUserPermissions, ok := args[11].(bool)
	if !ok {
		return nil, ErrInvalidInput("invalid updateUserPermissions type, expected bool")
	}
	userPermissionsRaw, ok := args[12].(map[string]interface{})
	if !ok && args[12] != nil {
		return nil, ErrInvalidInput("invalid userPermissions type, expected map[string]interface{}")
	}

	msg := &tokenizationtypes.MsgUpdateUserApprovals{
		Creator:                 creatorCosmosAddr,
		CollectionId:            collectionId,
		UpdateOutgoingApprovals: updateOutgoingApprovals,
		UpdateIncomingApprovals: updateIncomingApprovals,
		UpdateAutoApproveSelfInitiatedOutgoingTransfers: updateAutoApproveSelfInitiatedOutgoingTransfers,
		AutoApproveSelfInitiatedOutgoingTransfers:       autoApproveSelfInitiatedOutgoingTransfers,
		UpdateAutoApproveSelfInitiatedIncomingTransfers: updateAutoApproveSelfInitiatedIncomingTransfers,
		AutoApproveSelfInitiatedIncomingTransfers:       autoApproveSelfInitiatedIncomingTransfers,
		UpdateAutoApproveAllIncomingTransfers:           updateAutoApproveAllIncomingTransfers,
		AutoApproveAllIncomingTransfers:                 autoApproveAllIncomingTransfers,
		UpdateUserPermissions:                           updateUserPermissions,
	}

	// Convert outgoing approvals
	if updateOutgoingApprovals && outgoingApprovalsRaw != nil {
		outgoing, err := ConvertUserOutgoingApprovalArray(outgoingApprovalsRaw)
		if err != nil {
			return nil, ErrInvalidInput(fmt.Sprintf("outgoingApprovals: %v", err))
		}
		msg.OutgoingApprovals = outgoing
	}

	// Convert incoming approvals
	if updateIncomingApprovals && incomingApprovalsRaw != nil {
		incoming, err := ConvertUserIncomingApprovalArray(incomingApprovalsRaw)
		if err != nil {
			return nil, ErrInvalidInput(fmt.Sprintf("incomingApprovals: %v", err))
		}
		msg.IncomingApprovals = incoming
	}

	// Convert user permissions
	if updateUserPermissions && userPermissionsRaw != nil {
		perms, err := ConvertUserPermissions(userPermissionsRaw)
		if err != nil {
			return nil, ErrInvalidInput(fmt.Sprintf("userPermissions: %v", err))
		}
		msg.UserPermissions = perms
	}

	msgServer := tokenizationkeeper.NewMsgServerImpl(p.tokenizationKeeper)
	_, err = msgServer.UpdateUserApprovals(ctx, msg)
	if err != nil {
		return nil, WrapError(err, ErrorCodeApprovalFailed, "update user approvals failed")
	}

	return method.Outputs.Pack(true)
}

// DeleteIncomingApproval deletes an incoming approval.
func (p Precompile) DeleteIncomingApproval(ctx sdk.Context, method *abi.Method, args []interface{}, contract *vm.Contract) ([]byte, error) {
	if len(args) != 2 {
		return nil, ErrInvalidInput(fmt.Sprintf("invalid number of arguments, expected 2, got %d", len(args)))
	}

	collectionIdBig, ok := args[0].(*big.Int)
	if !ok {
		return nil, ErrInvalidInput("invalid collectionId type, expected *big.Int")
	}

	approvalId, ok := args[1].(string)
	if !ok {
		return nil, ErrInvalidInput("invalid approvalId type, expected string")
	}

	if err := ValidateCollectionId(collectionIdBig); err != nil {
		return nil, err
	}
	if err := ValidateString(approvalId, "approvalId"); err != nil {
		return nil, err
	}

	creatorCosmosAddr, err := p.GetCallerAddress(contract)
	if err != nil {
		return nil, err
	}

	collectionId := sdkmath.NewUintFromBigInt(collectionIdBig)

	msg := &tokenizationtypes.MsgDeleteIncomingApproval{
		Creator:      creatorCosmosAddr,
		CollectionId: collectionId,
		ApprovalId:   approvalId,
	}

	msgServer := tokenizationkeeper.NewMsgServerImpl(p.tokenizationKeeper)
	_, err = msgServer.DeleteIncomingApproval(ctx, msg)
	if err != nil {
		return nil, WrapError(err, ErrorCodeApprovalFailed, "delete incoming approval failed")
	}

	return method.Outputs.Pack(true)
}

// DeleteOutgoingApproval deletes an outgoing approval.
func (p Precompile) DeleteOutgoingApproval(ctx sdk.Context, method *abi.Method, args []interface{}, contract *vm.Contract) ([]byte, error) {
	if len(args) != 2 {
		return nil, ErrInvalidInput(fmt.Sprintf("invalid number of arguments, expected 2, got %d", len(args)))
	}

	collectionIdBig, ok := args[0].(*big.Int)
	if !ok {
		return nil, ErrInvalidInput("invalid collectionId type, expected *big.Int")
	}

	approvalId, ok := args[1].(string)
	if !ok {
		return nil, ErrInvalidInput("invalid approvalId type, expected string")
	}

	if err := ValidateCollectionId(collectionIdBig); err != nil {
		return nil, err
	}
	if err := ValidateString(approvalId, "approvalId"); err != nil {
		return nil, err
	}

	creatorCosmosAddr, err := p.GetCallerAddress(contract)
	if err != nil {
		return nil, err
	}

	collectionId := sdkmath.NewUintFromBigInt(collectionIdBig)

	msg := &tokenizationtypes.MsgDeleteOutgoingApproval{
		Creator:      creatorCosmosAddr,
		CollectionId: collectionId,
		ApprovalId:   approvalId,
	}

	msgServer := tokenizationkeeper.NewMsgServerImpl(p.tokenizationKeeper)
	_, err = msgServer.DeleteOutgoingApproval(ctx, msg)
	if err != nil {
		return nil, WrapError(err, ErrorCodeApprovalFailed, "delete outgoing approval failed")
	}

	return method.Outputs.Pack(true)
}

// PurgeApprovals purges expired approvals.
func (p Precompile) PurgeApprovals(ctx sdk.Context, method *abi.Method, args []interface{}, contract *vm.Contract) ([]byte, error) {
	if len(args) < 2 {
		return nil, ErrInvalidInput(fmt.Sprintf("invalid number of arguments, expected at least 2, got %d", len(args)))
	}

	creatorCosmosAddr, err := p.GetCallerAddress(contract)
	if err != nil {
		return nil, err
	}

	collectionIdBig, ok := args[0].(*big.Int)
	if !ok {
		return nil, ErrInvalidInput("invalid collectionId type, expected *big.Int")
	}

	if err := ValidateCollectionId(collectionIdBig); err != nil {
		return nil, err
	}

	collectionId := sdkmath.NewUintFromBigInt(collectionIdBig)

	purgeExpired, ok := args[1].(bool)
	if !ok {
		return nil, ErrInvalidInput("invalid purgeExpired type, expected bool")
	}
	approverAddress, ok := args[2].(string)
	if !ok && args[2] != nil {
		return nil, ErrInvalidInput("invalid approverAddress type, expected string")
	}
	purgeCounterpartyApprovals, ok := args[3].(bool)
	if !ok {
		return nil, ErrInvalidInput("invalid purgeCounterpartyApprovals type, expected bool")
	}
	approvalsToPurgeRaw, ok := args[4].([]interface{})
	if !ok && args[4] != nil {
		return nil, ErrInvalidInput("invalid approvalsToPurge type, expected []interface{}")
	}

	// Convert approver address if provided
	approverCosmosAddr := ""
	if approverAddress != "" {
		approverCosmosAddr, err = ConvertManagerAddress(approverAddress)
		if err != nil {
			return nil, ErrInvalidInput(fmt.Sprintf("approverAddress: %v", err))
		}
	} else {
		approverCosmosAddr = creatorCosmosAddr
	}

	// Convert approvals to purge
	approvalsToPurge := make([]*tokenizationtypes.ApprovalIdentifierDetails, 0)
	if approvalsToPurgeRaw != nil {
		for _, appRaw := range approvalsToPurgeRaw {
			if appMap, ok := appRaw.(map[string]interface{}); ok {
				approvalId, _ := appMap["approvalId"].(string)
				approvalLevel, _ := appMap["approvalLevel"].(string)
				approverAddr, _ := appMap["approverAddress"].(string)
				versionBig, _ := appMap["version"].(*big.Int)
				version := sdkmath.NewUint(0)
				if versionBig != nil {
					version = sdkmath.NewUintFromBigInt(versionBig)
				}
				approvalsToPurge = append(approvalsToPurge, &tokenizationtypes.ApprovalIdentifierDetails{
					ApprovalId:      approvalId,
					ApprovalLevel:   approvalLevel,
					ApproverAddress: approverAddr,
					Version:         version,
				})
			}
		}
	}

	msg := &tokenizationtypes.MsgPurgeApprovals{
		Creator:                    creatorCosmosAddr,
		CollectionId:               collectionId,
		PurgeExpired:               purgeExpired,
		ApproverAddress:            approverCosmosAddr,
		PurgeCounterpartyApprovals: purgeCounterpartyApprovals,
		ApprovalsToPurge:           approvalsToPurge,
	}

	msgServer := tokenizationkeeper.NewMsgServerImpl(p.tokenizationKeeper)
	resp, err := msgServer.PurgeApprovals(ctx, msg)
	if err != nil {
		return nil, WrapError(err, ErrorCodeApprovalFailed, "purge approvals failed")
	}

	return method.Outputs.Pack(resp.NumPurged.BigInt())
}

// CreateDynamicStore creates a dynamic store.
func (p Precompile) CreateDynamicStore(ctx sdk.Context, method *abi.Method, args []interface{}, contract *vm.Contract) ([]byte, error) {
	if len(args) != 3 {
		return nil, ErrInvalidInput(fmt.Sprintf("invalid number of arguments, expected 3, got %d", len(args)))
	}

	defaultValue, ok := args[0].(bool)
	if !ok {
		return nil, ErrInvalidInput("invalid defaultValue type, expected bool")
	}

	uri, ok := args[1].(string)
	if !ok {
		return nil, ErrInvalidInput("invalid uri type, expected string")
	}

	customData, ok := args[2].(string)
	if !ok {
		return nil, ErrInvalidInput("invalid customData type, expected string")
	}

	creatorCosmosAddr, err := p.GetCallerAddress(contract)
	if err != nil {
		return nil, err
	}

	msg := &tokenizationtypes.MsgCreateDynamicStore{
		Creator:      creatorCosmosAddr,
		DefaultValue: defaultValue,
		Uri:          uri,
		CustomData:   customData,
	}

	msgServer := tokenizationkeeper.NewMsgServerImpl(p.tokenizationKeeper)
	resp, err := msgServer.CreateDynamicStore(ctx, msg)
	if err != nil {
		return nil, WrapError(err, ErrorCodeInternalError, "create dynamic store failed")
	}

	return method.Outputs.Pack(resp.StoreId.BigInt())
}

// UpdateDynamicStore updates a dynamic store.
func (p Precompile) UpdateDynamicStore(ctx sdk.Context, method *abi.Method, args []interface{}, contract *vm.Contract) ([]byte, error) {
	if len(args) != 5 {
		return nil, ErrInvalidInput(fmt.Sprintf("invalid number of arguments, expected 5, got %d", len(args)))
	}

	storeIdBig, ok := args[0].(*big.Int)
	if !ok {
		return nil, ErrInvalidInput("invalid storeId type, expected *big.Int")
	}

	defaultValue, ok := args[1].(bool)
	if !ok {
		return nil, ErrInvalidInput("invalid defaultValue type, expected bool")
	}

	globalEnabled, ok := args[2].(bool)
	if !ok {
		return nil, ErrInvalidInput("invalid globalEnabled type, expected bool")
	}

	uri, ok := args[3].(string)
	if !ok {
		return nil, ErrInvalidInput("invalid uri type, expected string")
	}

	customData, ok := args[4].(string)
	if !ok {
		return nil, ErrInvalidInput("invalid customData type, expected string")
	}

	if err := CheckOverflow(storeIdBig, "storeId"); err != nil {
		return nil, err
	}

	creatorCosmosAddr, err := p.GetCallerAddress(contract)
	if err != nil {
		return nil, err
	}

	storeId := sdkmath.NewUintFromBigInt(storeIdBig)

	msg := &tokenizationtypes.MsgUpdateDynamicStore{
		Creator:       creatorCosmosAddr,
		StoreId:       storeId,
		DefaultValue:  defaultValue,
		GlobalEnabled: globalEnabled,
		Uri:           uri,
		CustomData:    customData,
	}

	msgServer := tokenizationkeeper.NewMsgServerImpl(p.tokenizationKeeper)
	_, err = msgServer.UpdateDynamicStore(ctx, msg)
	if err != nil {
		return nil, WrapError(err, ErrorCodeInternalError, "update dynamic store failed")
	}

	return method.Outputs.Pack(true)
}

// DeleteDynamicStore deletes a dynamic store.
func (p Precompile) DeleteDynamicStore(ctx sdk.Context, method *abi.Method, args []interface{}, contract *vm.Contract) ([]byte, error) {
	if len(args) != 1 {
		return nil, ErrInvalidInput(fmt.Sprintf("invalid number of arguments, expected 1, got %d", len(args)))
	}

	storeIdBig, ok := args[0].(*big.Int)
	if !ok {
		return nil, ErrInvalidInput("invalid storeId type, expected *big.Int")
	}

	if err := CheckOverflow(storeIdBig, "storeId"); err != nil {
		return nil, err
	}

	creatorCosmosAddr, err := p.GetCallerAddress(contract)
	if err != nil {
		return nil, err
	}

	storeId := sdkmath.NewUintFromBigInt(storeIdBig)

	msg := &tokenizationtypes.MsgDeleteDynamicStore{
		Creator: creatorCosmosAddr,
		StoreId: storeId,
	}

	msgServer := tokenizationkeeper.NewMsgServerImpl(p.tokenizationKeeper)
	_, err = msgServer.DeleteDynamicStore(ctx, msg)
	if err != nil {
		return nil, WrapError(err, ErrorCodeInternalError, "delete dynamic store failed")
	}

	return method.Outputs.Pack(true)
}

// SetDynamicStoreValue sets a dynamic store value.
func (p Precompile) SetDynamicStoreValue(ctx sdk.Context, method *abi.Method, args []interface{}, contract *vm.Contract) ([]byte, error) {
	if len(args) != 3 {
		return nil, ErrInvalidInput(fmt.Sprintf("invalid number of arguments, expected 3, got %d", len(args)))
	}

	storeIdBig, ok := args[0].(*big.Int)
	if !ok {
		return nil, ErrInvalidInput("invalid storeId type, expected *big.Int")
	}

	address, ok := args[1].(common.Address)
	if !ok {
		return nil, ErrInvalidInput("invalid address type, expected common.Address")
	}

	value, ok := args[2].(bool)
	if !ok {
		return nil, ErrInvalidInput("invalid value type, expected bool")
	}

	if err := CheckOverflow(storeIdBig, "storeId"); err != nil {
		return nil, err
	}
	if err := ValidateAddress(address, "address"); err != nil {
		return nil, err
	}

	creatorCosmosAddr, err := p.GetCallerAddress(contract)
	if err != nil {
		return nil, err
	}

	storeId := sdkmath.NewUintFromBigInt(storeIdBig)
	addressCosmosAddr := sdk.AccAddress(address.Bytes()).String()

	msg := &tokenizationtypes.MsgSetDynamicStoreValue{
		Creator: creatorCosmosAddr,
		StoreId: storeId,
		Address: addressCosmosAddr,
		Value:   value,
	}

	msgServer := tokenizationkeeper.NewMsgServerImpl(p.tokenizationKeeper)
	_, err = msgServer.SetDynamicStoreValue(ctx, msg)
	if err != nil {
		return nil, WrapError(err, ErrorCodeInternalError, "set dynamic store value failed")
	}

	return method.Outputs.Pack(true)
}

// SetValidTokenIds sets valid token IDs for a collection.
func (p Precompile) SetValidTokenIds(ctx sdk.Context, method *abi.Method, args []interface{}, contract *vm.Contract) ([]byte, error) {
	if len(args) < 2 {
		return nil, ErrInvalidInput(fmt.Sprintf("invalid number of arguments, expected at least 2, got %d", len(args)))
	}

	creatorCosmosAddr, err := p.GetCallerAddress(contract)
	if err != nil {
		return nil, err
	}

	collectionIdBig, ok := args[0].(*big.Int)
	if !ok {
		return nil, ErrInvalidInput("invalid collectionId type, expected *big.Int")
	}

	if err := ValidateCollectionId(collectionIdBig); err != nil {
		return nil, err
	}

	collectionId := sdkmath.NewUintFromBigInt(collectionIdBig)

	validTokenIdsRaw, ok := args[1].([]interface{})
	if !ok && args[1] != nil {
		return nil, ErrInvalidInput("invalid validTokenIds type, expected []interface{}")
	}
	canUpdateValidTokenIdsRaw, ok := args[2].([]interface{})
	if !ok && args[2] != nil {
		return nil, ErrInvalidInput("invalid canUpdateValidTokenIds type, expected []interface{}")
	}

	// Convert validTokenIds
	validTokenIds, err := convertUintRangeArrayFromInterfaceWithError(validTokenIdsRaw, "validTokenIds")
	if err != nil {
		return nil, ErrInvalidInput(err.Error())
	}
	convertedValidTokenIds, err := ConvertUintRangeArray(validTokenIds)
	if err != nil {
		return nil, ErrInvalidInput(fmt.Sprintf("validTokenIds: %v", err))
	}

	// Convert canUpdateValidTokenIds permissions
	canUpdateValidTokenIds, err := ConvertTokenIdsActionPermissionArrayFromInterface(canUpdateValidTokenIdsRaw)
	if err != nil {
		return nil, ErrInvalidInput(fmt.Sprintf("canUpdateValidTokenIds: %v", err))
	}

	msg := &tokenizationtypes.MsgSetValidTokenIds{
		Creator:                creatorCosmosAddr,
		CollectionId:           collectionId,
		ValidTokenIds:          convertedValidTokenIds,
		CanUpdateValidTokenIds: canUpdateValidTokenIds,
	}

	msgServer := tokenizationkeeper.NewMsgServerImpl(p.tokenizationKeeper)
	resp, err := msgServer.SetValidTokenIds(ctx, msg)
	if err != nil {
		return nil, WrapError(err, ErrorCodeCollectionNotFound, "set valid token IDs failed")
	}

	return method.Outputs.Pack(resp.CollectionId.BigInt())
}

// SetManager sets the manager for a collection.
func (p Precompile) SetManager(ctx sdk.Context, method *abi.Method, args []interface{}, contract *vm.Contract) ([]byte, error) {
	if len(args) != 2 {
		return nil, ErrInvalidInput(fmt.Sprintf("invalid number of arguments, expected 2, got %d", len(args)))
	}

	collectionIdBig, ok := args[0].(*big.Int)
	if !ok {
		return nil, ErrInvalidInput("invalid collectionId type, expected *big.Int")
	}

	manager, ok := args[1].(string)
	if !ok {
		// Try to convert from common.Address if it's an EVM address
		if managerAddr, ok2 := args[1].(common.Address); ok2 {
			manager = sdk.AccAddress(managerAddr.Bytes()).String()
		} else {
			return nil, ErrInvalidInput("invalid manager type, expected string or address")
		}
	}

	if err := ValidateCollectionId(collectionIdBig); err != nil {
		return nil, err
	}

	creatorCosmosAddr, err := p.GetCallerAddress(contract)
	if err != nil {
		return nil, err
	}

	collectionId := sdkmath.NewUintFromBigInt(collectionIdBig)

	// Convert manager address if it's an EVM address
	managerCosmosAddr, err := ConvertManagerAddress(manager)
	if err != nil {
		return nil, ErrInvalidInput(fmt.Sprintf("invalid manager address: %v", err))
	}

	msg := &tokenizationtypes.MsgSetManager{
		Creator:      creatorCosmosAddr,
		CollectionId: collectionId,
		Manager:      managerCosmosAddr,
		// Empty permissions array - can be extended later
		CanUpdateManager: []*tokenizationtypes.ActionPermission{},
	}

	msgServer := tokenizationkeeper.NewMsgServerImpl(p.tokenizationKeeper)
	_, err = msgServer.SetManager(ctx, msg)
	if err != nil {
		return nil, WrapError(err, ErrorCodeCollectionNotFound, "set manager failed")
	}

	return method.Outputs.Pack(collectionId.BigInt())
}

// SetCollectionMetadata sets collection metadata.
func (p Precompile) SetCollectionMetadata(ctx sdk.Context, method *abi.Method, args []interface{}, contract *vm.Contract) ([]byte, error) {
	if len(args) != 3 {
		return nil, ErrInvalidInput(fmt.Sprintf("invalid number of arguments, expected 3, got %d", len(args)))
	}

	collectionIdBig, ok := args[0].(*big.Int)
	if !ok {
		return nil, ErrInvalidInput("invalid collectionId type, expected *big.Int")
	}

	uri, ok := args[1].(string)
	if !ok {
		return nil, ErrInvalidInput("invalid uri type, expected string")
	}

	customData, ok := args[2].(string)
	if !ok {
		return nil, ErrInvalidInput("invalid customData type, expected string")
	}

	if err := ValidateCollectionId(collectionIdBig); err != nil {
		return nil, err
	}

	creatorCosmosAddr, err := p.GetCallerAddress(contract)
	if err != nil {
		return nil, err
	}

	collectionId := sdkmath.NewUintFromBigInt(collectionIdBig)

	metadata := ConvertCollectionMetadata(uri, customData)

	msg := &tokenizationtypes.MsgSetCollectionMetadata{
		Creator:            creatorCosmosAddr,
		CollectionId:       collectionId,
		CollectionMetadata: metadata,
		// Empty permissions array - can be extended later
		CanUpdateCollectionMetadata: []*tokenizationtypes.ActionPermission{},
	}

	msgServer := tokenizationkeeper.NewMsgServerImpl(p.tokenizationKeeper)
	_, err = msgServer.SetCollectionMetadata(ctx, msg)
	if err != nil {
		return nil, WrapError(err, ErrorCodeCollectionNotFound, "set collection metadata failed")
	}

	return method.Outputs.Pack(collectionId.BigInt())
}

// SetTokenMetadata sets token metadata.
func (p Precompile) SetTokenMetadata(ctx sdk.Context, method *abi.Method, args []interface{}, contract *vm.Contract) ([]byte, error) {
	if len(args) < 2 {
		return nil, ErrInvalidInput(fmt.Sprintf("invalid number of arguments, expected at least 2, got %d", len(args)))
	}

	creatorCosmosAddr, err := p.GetCallerAddress(contract)
	if err != nil {
		return nil, err
	}

	collectionIdBig, ok := args[0].(*big.Int)
	if !ok {
		return nil, ErrInvalidInput("invalid collectionId type, expected *big.Int")
	}

	if err := ValidateCollectionId(collectionIdBig); err != nil {
		return nil, err
	}

	collectionId := sdkmath.NewUintFromBigInt(collectionIdBig)

	tokenMetadataRaw, ok := args[1].([]interface{})
	if !ok && args[1] != nil {
		return nil, ErrInvalidInput("invalid tokenMetadata type, expected []interface{}")
	}
	canUpdateTokenMetadataRaw, ok := args[2].([]interface{})
	if !ok && args[2] != nil {
		return nil, ErrInvalidInput("invalid canUpdateTokenMetadata type, expected []interface{}")
	}

	// Convert tokenMetadata
	tokenMetadata := make([]*tokenizationtypes.TokenMetadata, 0)
	if tokenMetadataRaw != nil {
		for i, tmRaw := range tokenMetadataRaw {
			tmMap, ok := tmRaw.(map[string]interface{})
			if !ok {
				return nil, ErrInvalidInput(fmt.Sprintf("tokenMetadata[%d] must be a map", i))
			}
			uri, _ := tmMap["uri"].(string)
			customDataVal, _ := tmMap["customData"].(string)
			tokenIdsRaw, _ := tmMap["tokenIds"].([]interface{})
			tokenIds, err := convertUintRangeArrayFromInterfaceWithError(tokenIdsRaw, fmt.Sprintf("tokenMetadata[%d].tokenIds", i))
			if err != nil {
				return nil, ErrInvalidInput(err.Error())
			}
			tm, err := ConvertTokenMetadata(uri, customDataVal, tokenIds)
			if err != nil {
				return nil, ErrInvalidInput(fmt.Sprintf("tokenMetadata[%d]: %v", i, err))
			}
			tokenMetadata = append(tokenMetadata, tm)
		}
	}

	// Convert canUpdateTokenMetadata permissions
	canUpdateTokenMetadata, err := ConvertTokenIdsActionPermissionArrayFromInterface(canUpdateTokenMetadataRaw)
	if err != nil {
		return nil, ErrInvalidInput(fmt.Sprintf("canUpdateTokenMetadata: %v", err))
	}

	msg := &tokenizationtypes.MsgSetTokenMetadata{
		Creator:                creatorCosmosAddr,
		CollectionId:           collectionId,
		TokenMetadata:          tokenMetadata,
		CanUpdateTokenMetadata: canUpdateTokenMetadata,
	}

	msgServer := tokenizationkeeper.NewMsgServerImpl(p.tokenizationKeeper)
	resp, err := msgServer.SetTokenMetadata(ctx, msg)
	if err != nil {
		return nil, WrapError(err, ErrorCodeCollectionNotFound, "set token metadata failed")
	}

	return method.Outputs.Pack(resp.CollectionId.BigInt())
}

// SetCustomData sets custom data for a collection.
func (p Precompile) SetCustomData(ctx sdk.Context, method *abi.Method, args []interface{}, contract *vm.Contract) ([]byte, error) {
	if len(args) != 2 {
		return nil, ErrInvalidInput(fmt.Sprintf("invalid number of arguments, expected 2, got %d", len(args)))
	}

	collectionIdBig, ok := args[0].(*big.Int)
	if !ok {
		return nil, ErrInvalidInput("invalid collectionId type, expected *big.Int")
	}

	customData, ok := args[1].(string)
	if !ok {
		return nil, ErrInvalidInput("invalid customData type, expected string")
	}

	if err := ValidateCollectionId(collectionIdBig); err != nil {
		return nil, err
	}

	creatorCosmosAddr, err := p.GetCallerAddress(contract)
	if err != nil {
		return nil, err
	}

	collectionId := sdkmath.NewUintFromBigInt(collectionIdBig)

	msg := &tokenizationtypes.MsgSetCustomData{
		Creator:      creatorCosmosAddr,
		CollectionId: collectionId,
		CustomData:   customData,
	}

	msgServer := tokenizationkeeper.NewMsgServerImpl(p.tokenizationKeeper)
	_, err = msgServer.SetCustomData(ctx, msg)
	if err != nil {
		return nil, WrapError(err, ErrorCodeCollectionNotFound, "set custom data failed")
	}

	return method.Outputs.Pack(collectionId.BigInt())
}

// SetStandards sets standards for a collection.
func (p Precompile) SetStandards(ctx sdk.Context, method *abi.Method, args []interface{}, contract *vm.Contract) ([]byte, error) {
	if len(args) != 2 {
		return nil, ErrInvalidInput(fmt.Sprintf("invalid number of arguments, expected 2, got %d", len(args)))
	}

	collectionIdBig, ok := args[0].(*big.Int)
	if !ok {
		return nil, ErrInvalidInput("invalid collectionId type, expected *big.Int")
	}

	standards, ok := args[1].([]string)
	if !ok {
		return nil, ErrInvalidInput("invalid standards type, expected []string")
	}

	if err := ValidateCollectionId(collectionIdBig); err != nil {
		return nil, err
	}

	creatorCosmosAddr, err := p.GetCallerAddress(contract)
	if err != nil {
		return nil, err
	}

	collectionId := sdkmath.NewUintFromBigInt(collectionIdBig)

	msg := &tokenizationtypes.MsgSetStandards{
		Creator:      creatorCosmosAddr,
		CollectionId: collectionId,
		Standards:    standards,
		// Empty permissions array - can be extended later
		CanUpdateStandards: []*tokenizationtypes.ActionPermission{},
	}

	msgServer := tokenizationkeeper.NewMsgServerImpl(p.tokenizationKeeper)
	_, err = msgServer.SetStandards(ctx, msg)
	if err != nil {
		return nil, WrapError(err, ErrorCodeCollectionNotFound, "set standards failed")
	}

	return method.Outputs.Pack(collectionId.BigInt())
}

// SetCollectionApprovals sets collection approvals.
func (p Precompile) SetCollectionApprovals(ctx sdk.Context, method *abi.Method, args []interface{}, contract *vm.Contract) ([]byte, error) {
	if len(args) < 2 {
		return nil, ErrInvalidInput(fmt.Sprintf("invalid number of arguments, expected at least 2, got %d", len(args)))
	}

	creatorCosmosAddr, err := p.GetCallerAddress(contract)
	if err != nil {
		return nil, err
	}

	collectionIdBig, ok := args[0].(*big.Int)
	if !ok {
		return nil, ErrInvalidInput("invalid collectionId type, expected *big.Int")
	}

	if err := ValidateCollectionId(collectionIdBig); err != nil {
		return nil, err
	}

	collectionId := sdkmath.NewUintFromBigInt(collectionIdBig)

	collectionApprovalsRaw, ok := args[1].([]interface{})
	if !ok && args[1] != nil {
		return nil, ErrInvalidInput("invalid collectionApprovals type, expected []interface{}")
	}
	canUpdateCollectionApprovalsRaw, ok := args[2].([]interface{})
	if !ok && args[2] != nil {
		return nil, ErrInvalidInput("invalid canUpdateCollectionApprovals type, expected []interface{}")
	}

	// Convert collectionApprovals using the dedicated array converter
	var collectionApprovals []*tokenizationtypes.CollectionApproval
	if collectionApprovalsRaw != nil {
		var err error
		collectionApprovals, err = ConvertCollectionApprovalArray(collectionApprovalsRaw)
		if err != nil {
			return nil, ErrInvalidInput(fmt.Sprintf("collectionApprovals: %v", err))
		}
	}

	// Convert canUpdateCollectionApprovals permissions
	canUpdateCollectionApprovals, err := ConvertCollectionApprovalPermissionArrayFromInterface(canUpdateCollectionApprovalsRaw)
	if err != nil {
		return nil, ErrInvalidInput(fmt.Sprintf("canUpdateCollectionApprovals: %v", err))
	}

	msg := &tokenizationtypes.MsgSetCollectionApprovals{
		Creator:                      creatorCosmosAddr,
		CollectionId:                 collectionId,
		CollectionApprovals:          collectionApprovals,
		CanUpdateCollectionApprovals: canUpdateCollectionApprovals,
	}

	msgServer := tokenizationkeeper.NewMsgServerImpl(p.tokenizationKeeper)
	resp, err := msgServer.SetCollectionApprovals(ctx, msg)
	if err != nil {
		return nil, WrapError(err, ErrorCodeCollectionNotFound, "set collection approvals failed")
	}

	return method.Outputs.Pack(resp.CollectionId.BigInt())
}

// SetIsArchived sets the archived status for a collection.
func (p Precompile) SetIsArchived(ctx sdk.Context, method *abi.Method, args []interface{}, contract *vm.Contract) ([]byte, error) {
	if len(args) != 2 {
		return nil, ErrInvalidInput(fmt.Sprintf("invalid number of arguments, expected 2, got %d", len(args)))
	}

	collectionIdBig, ok := args[0].(*big.Int)
	if !ok {
		return nil, ErrInvalidInput("invalid collectionId type, expected *big.Int")
	}

	isArchived, ok := args[1].(bool)
	if !ok {
		return nil, ErrInvalidInput("invalid isArchived type, expected bool")
	}

	if err := ValidateCollectionId(collectionIdBig); err != nil {
		return nil, err
	}

	creatorCosmosAddr, err := p.GetCallerAddress(contract)
	if err != nil {
		return nil, err
	}

	collectionId := sdkmath.NewUintFromBigInt(collectionIdBig)

	msg := &tokenizationtypes.MsgSetIsArchived{
		Creator:      creatorCosmosAddr,
		CollectionId: collectionId,
		IsArchived:   isArchived,
	}

	msgServer := tokenizationkeeper.NewMsgServerImpl(p.tokenizationKeeper)
	_, err = msgServer.SetIsArchived(ctx, msg)
	if err != nil {
		return nil, WrapError(err, ErrorCodeCollectionNotFound, "set is archived failed")
	}

	return method.Outputs.Pack(collectionId.BigInt())
}

// CastVote casts a vote on a voting challenge.
func (p Precompile) CastVote(ctx sdk.Context, method *abi.Method, args []interface{}, contract *vm.Contract) ([]byte, error) {
	if len(args) != 6 {
		return nil, ErrInvalidInput(fmt.Sprintf("invalid number of arguments, expected 6, got %d", len(args)))
	}

	collectionIdBig, ok := args[0].(*big.Int)
	if !ok {
		return nil, ErrInvalidInput("invalid collectionId type, expected *big.Int")
	}

	approvalLevel, ok := args[1].(string)
	if !ok {
		return nil, ErrInvalidInput("invalid approvalLevel type, expected string")
	}

	approverAddress, ok := args[2].(string)
	if !ok {
		return nil, ErrInvalidInput("invalid approverAddress type, expected string")
	}

	approvalId, ok := args[3].(string)
	if !ok {
		return nil, ErrInvalidInput("invalid approvalId type, expected string")
	}

	proposalId, ok := args[4].(string)
	if !ok {
		return nil, ErrInvalidInput("invalid proposalId type, expected string")
	}

	yesWeightBig, ok := args[5].(*big.Int)
	if !ok {
		return nil, ErrInvalidInput("invalid yesWeight type, expected *big.Int")
	}

	if err := ValidateCollectionId(collectionIdBig); err != nil {
		return nil, err
	}
	if err := ValidateString(approvalLevel, "approvalLevel"); err != nil {
		return nil, err
	}
	if err := ValidateString(approverAddress, "approverAddress"); err != nil {
		return nil, err
	}
	if err := ValidateString(approvalId, "approvalId"); err != nil {
		return nil, err
	}
	if err := ValidateString(proposalId, "proposalId"); err != nil {
		return nil, err
	}
	if err := CheckOverflow(yesWeightBig, "yesWeight"); err != nil {
		return nil, err
	}

	creatorCosmosAddr, err := p.GetCallerAddress(contract)
	if err != nil {
		return nil, err
	}

	collectionId := sdkmath.NewUintFromBigInt(collectionIdBig)
	yesWeight := sdkmath.NewUintFromBigInt(yesWeightBig)

	msg := &tokenizationtypes.MsgCastVote{
		Creator:         creatorCosmosAddr,
		CollectionId:    collectionId,
		ApprovalLevel:   approvalLevel,
		ApproverAddress: approverAddress,
		ApprovalId:      approvalId,
		ProposalId:      proposalId,
		YesWeight:       yesWeight,
	}

	msgServer := tokenizationkeeper.NewMsgServerImpl(p.tokenizationKeeper)
	_, err = msgServer.CastVote(ctx, msg)
	if err != nil {
		return nil, WrapError(err, ErrorCodeInternalError, "cast vote failed")
	}

	return method.Outputs.Pack(true)
}

// UniversalUpdateCollection performs a universal update on a collection (legacy method).
func (p Precompile) UniversalUpdateCollection(ctx sdk.Context, method *abi.Method, args []interface{}, contract *vm.Contract) ([]byte, error) {
	if len(args) < 2 {
		return nil, ErrInvalidInput(fmt.Sprintf("invalid number of arguments, expected at least 2, got %d", len(args)))
	}

	creatorCosmosAddr, err := p.GetCallerAddress(contract)
	if err != nil {
		return nil, err
	}

	collectionIdBig, ok := args[0].(*big.Int)
	if !ok {
		return nil, ErrInvalidInput("invalid collectionId type, expected *big.Int")
	}

	if err := ValidateCollectionId(collectionIdBig); err != nil {
		return nil, err
	}

	collectionId := sdkmath.NewUintFromBigInt(collectionIdBig)

	// Extract all fields (similar to UpdateCollection but with defaultBalances)
	// Type assertions with error checking to prevent silent data corruption
	defaultBalancesRaw, ok := args[1].(map[string]interface{})
	if !ok && args[1] != nil {
		return nil, ErrInvalidInput("invalid defaultBalances type, expected map[string]interface{}")
	}
	updateValidTokenIds, ok := args[2].(bool)
	if !ok {
		return nil, ErrInvalidInput("invalid updateValidTokenIds type, expected bool")
	}
	validTokenIdsRaw, ok := args[3].([]interface{})
	if !ok && args[3] != nil {
		return nil, ErrInvalidInput("invalid validTokenIds type, expected []interface{}")
	}
	updateCollectionPermissions, ok := args[4].(bool)
	if !ok {
		return nil, ErrInvalidInput("invalid updateCollectionPermissions type, expected bool")
	}
	collectionPermissionsRaw, ok := args[5].(map[string]interface{})
	if !ok && args[5] != nil {
		return nil, ErrInvalidInput("invalid collectionPermissions type, expected map[string]interface{}")
	}
	updateManager, ok := args[6].(bool)
	if !ok {
		return nil, ErrInvalidInput("invalid updateManager type, expected bool")
	}
	manager, ok := args[7].(string)
	if !ok && args[7] != nil {
		return nil, ErrInvalidInput("invalid manager type, expected string")
	}
	updateCollectionMetadata, ok := args[8].(bool)
	if !ok {
		return nil, ErrInvalidInput("invalid updateCollectionMetadata type, expected bool")
	}
	collectionMetadataRaw, ok := args[9].(map[string]interface{})
	if !ok && args[9] != nil {
		return nil, ErrInvalidInput("invalid collectionMetadata type, expected map[string]interface{}")
	}
	updateTokenMetadata, ok := args[10].(bool)
	if !ok {
		return nil, ErrInvalidInput("invalid updateTokenMetadata type, expected bool")
	}
	tokenMetadataRaw, ok := args[11].([]interface{})
	if !ok && args[11] != nil {
		return nil, ErrInvalidInput("invalid tokenMetadata type, expected []interface{}")
	}
	updateCustomData, ok := args[12].(bool)
	if !ok {
		return nil, ErrInvalidInput("invalid updateCustomData type, expected bool")
	}
	customData, ok := args[13].(string)
	if !ok && args[13] != nil {
		return nil, ErrInvalidInput("invalid customData type, expected string")
	}
	updateCollectionApprovals, ok := args[14].(bool)
	if !ok {
		return nil, ErrInvalidInput("invalid updateCollectionApprovals type, expected bool")
	}
	collectionApprovalsRaw, ok := args[15].([]interface{})
	if !ok && args[15] != nil {
		return nil, ErrInvalidInput("invalid collectionApprovals type, expected []interface{}")
	}
	updateStandards, ok := args[16].(bool)
	if !ok {
		return nil, ErrInvalidInput("invalid updateStandards type, expected bool")
	}
	standardsRaw, ok := args[17].([]string)
	if !ok && args[17] != nil {
		return nil, ErrInvalidInput("invalid standards type, expected []string")
	}
	updateIsArchived, ok := args[18].(bool)
	if !ok {
		return nil, ErrInvalidInput("invalid updateIsArchived type, expected bool")
	}
	isArchived, ok := args[19].(bool)
	if !ok && args[19] != nil {
		return nil, ErrInvalidInput("invalid isArchived type, expected bool")
	}

	msg := &tokenizationtypes.MsgUniversalUpdateCollection{
		Creator:                     creatorCosmosAddr,
		CollectionId:                collectionId,
		UpdateValidTokenIds:         updateValidTokenIds,
		UpdateCollectionPermissions: updateCollectionPermissions,
		UpdateManager:               updateManager,
		UpdateCollectionMetadata:    updateCollectionMetadata,
		UpdateTokenMetadata:         updateTokenMetadata,
		UpdateCustomData:            updateCustomData,
		UpdateCollectionApprovals:   updateCollectionApprovals,
		UpdateStandards:             updateStandards,
		UpdateIsArchived:            updateIsArchived,
	}

	// Convert defaultBalances
	if defaultBalancesRaw != nil {
		defaultBalances, err := ConvertUserBalanceStore(defaultBalancesRaw)
		if err != nil {
			return nil, ErrInvalidInput(fmt.Sprintf("defaultBalances: %v", err))
		}
		msg.DefaultBalances = defaultBalances
	}

	// Convert other fields (similar to UpdateCollection)
	if updateValidTokenIds && validTokenIdsRaw != nil {
		validTokenIds, err := convertUintRangeArrayFromInterfaceWithError(validTokenIdsRaw, "validTokenIds")
		if err != nil {
			return nil, ErrInvalidInput(err.Error())
		}
		msg.ValidTokenIds, err = ConvertUintRangeArray(validTokenIds)
		if err != nil {
			return nil, ErrInvalidInput(fmt.Sprintf("validTokenIds: %v", err))
		}
	}

	if updateCollectionPermissions && collectionPermissionsRaw != nil {
		var err error
		msg.CollectionPermissions, err = ConvertCollectionPermissions(collectionPermissionsRaw)
		if err != nil {
			return nil, ErrInvalidInput(fmt.Sprintf("collectionPermissions: %v", err))
		}
	}

	if updateManager {
		managerCosmosAddr, err := ConvertManagerAddress(manager)
		if err != nil {
			return nil, ErrInvalidInput(fmt.Sprintf("manager: %v", err))
		}
		msg.Manager = managerCosmosAddr
	}

	if updateCollectionMetadata && collectionMetadataRaw != nil {
		uri, _ := collectionMetadataRaw["uri"].(string)
		customData, _ := collectionMetadataRaw["customData"].(string)
		msg.CollectionMetadata = ConvertCollectionMetadata(uri, customData)
	}

	if updateTokenMetadata && tokenMetadataRaw != nil {
		tokenMetadata := make([]*tokenizationtypes.TokenMetadata, 0)
		for i, tmRaw := range tokenMetadataRaw {
			tmMap, ok := tmRaw.(map[string]interface{})
			if !ok {
				return nil, ErrInvalidInput(fmt.Sprintf("tokenMetadata[%d] must be a map", i))
			}
			uri, _ := tmMap["uri"].(string)
			customDataVal, _ := tmMap["customData"].(string)
			tokenIdsRaw, _ := tmMap["tokenIds"].([]interface{})
			tokenIds, err := convertUintRangeArrayFromInterfaceWithError(tokenIdsRaw, fmt.Sprintf("tokenMetadata[%d].tokenIds", i))
			if err != nil {
				return nil, ErrInvalidInput(err.Error())
			}
			tm, err := ConvertTokenMetadata(uri, customDataVal, tokenIds)
			if err != nil {
				return nil, ErrInvalidInput(fmt.Sprintf("tokenMetadata[%d]: %v", i, err))
			}
			tokenMetadata = append(tokenMetadata, tm)
		}
		msg.TokenMetadata = tokenMetadata
	}

	if updateCustomData {
		msg.CustomData = customData
	}

	if updateCollectionApprovals && collectionApprovalsRaw != nil {
		collectionApprovals, err := ConvertCollectionApprovalArray(collectionApprovalsRaw)
		if err != nil {
			return nil, ErrInvalidInput(fmt.Sprintf("collectionApprovals: %v", err))
		}
		msg.CollectionApprovals = collectionApprovals
	}

	if updateStandards && standardsRaw != nil {
		msg.Standards = standardsRaw
	}

	if updateIsArchived {
		msg.IsArchived = isArchived
	}

	// Set empty invariants and paths
	msg.Invariants = &tokenizationtypes.InvariantsAddObject{}
	msg.CosmosCoinWrapperPathsToAdd = []*tokenizationtypes.CosmosCoinWrapperPathAddObject{}
	msg.AliasPathsToAdd = []*tokenizationtypes.AliasPathAddObject{}

	msgServer := tokenizationkeeper.NewMsgServerImpl(p.tokenizationKeeper)
	resp, err := msgServer.UniversalUpdateCollection(ctx, msg)
	if err != nil {
		return nil, WrapError(err, ErrorCodeCollectionNotFound, "universal update collection failed")
	}

	return method.Outputs.Pack(resp.CollectionId.BigInt())
}
