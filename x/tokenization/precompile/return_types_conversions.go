package tokenization

import (
	"fmt"
	"math/big"

	tokenizationtypes "github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

// uintRangesToSolidity converts []*UintRange to Solidity UintRange[] format (each element is [start, end]).
func uintRangesToSolidity(ranges []*tokenizationtypes.UintRange) []interface{} {
	if len(ranges) == 0 {
		return make([]interface{}, 0)
	}
	out := make([]interface{}, 0, len(ranges))
	for _, r := range ranges {
		if r == nil {
			continue
		}
		start := big.NewInt(0)
		end := big.NewInt(0)
		if !r.Start.IsNil() {
			start = r.Start.BigInt()
		}
		if !r.End.IsNil() {
			end = r.End.BigInt()
		}
		out = append(out, []interface{}{start, end})
	}
	return out
}

// conversionWithoutDenomToSolidity converts ConversionWithoutDenom to Solidity struct format.
// Solidity: ConversionSideA sideA (amount), Balance[] sideB.
func conversionWithoutDenomToSolidity(c *tokenizationtypes.ConversionWithoutDenom) []interface{} {
	if c == nil {
		return []interface{}{
			[]interface{}{big.NewInt(0)},
			make([]interface{}, 0), // sideB Balance[]
		}
	}
	var sideA []interface{}
	if c.SideA != nil {
		amt := big.NewInt(0)
		if !c.SideA.Amount.IsNil() {
			amt = c.SideA.Amount.BigInt()
		}
		sideA = []interface{}{amt}
	} else {
		sideA = []interface{}{big.NewInt(0)}
	}
	sideB := make([]interface{}, 0, len(c.SideB))
	for _, b := range c.SideB {
		if b == nil {
			continue
		}
		conv, err := ConvertBalanceToSolidityStruct(b)
		if err != nil {
			continue
		}
		sideB = append(sideB, conv)
	}
	return []interface{}{sideA, sideB}
}

// conversionToSolidity converts Conversion (with denom) to Solidity struct format.
// Solidity: ConversionSideAWithDenom sideA (amount, denom), Balance[] sideB.
func conversionToSolidity(c *tokenizationtypes.Conversion) []interface{} {
	if c == nil {
		return []interface{}{
			[]interface{}{big.NewInt(0), ""},
			make([]interface{}, 0), // sideB Balance[]
		}
	}
	var sideA []interface{}
	if c.SideA != nil {
		amt := big.NewInt(0)
		if !c.SideA.Amount.IsNil() {
			amt = c.SideA.Amount.BigInt()
		}
		denom := ""
		if c.SideA.Denom != "" {
			denom = c.SideA.Denom
		}
		sideA = []interface{}{amt, denom}
	} else {
		sideA = []interface{}{big.NewInt(0), ""}
	}
	sideB := make([]interface{}, 0, len(c.SideB))
	for _, b := range c.SideB {
		if b == nil {
			continue
		}
		conv, err := ConvertBalanceToSolidityStruct(b)
		if err != nil {
			continue
		}
		sideB = append(sideB, conv)
	}
	return []interface{}{sideA, sideB}
}

// cosmosCoinBackedPathToSolidity converts CosmosCoinBackedPath to Solidity struct format.
// Solidity: string addr, Conversion conversion.
func cosmosCoinBackedPathToSolidity(p *tokenizationtypes.CosmosCoinBackedPath) []interface{} {
	if p == nil {
		return []interface{}{"", conversionToSolidity(nil)}
	}
	addr := p.Address
	return []interface{}{addr, conversionToSolidity(p.Conversion)}
}

// actionPermissionToSolidity converts ActionPermission to Solidity struct format.
// Solidity: UintRange[] permanentlyPermittedTimes, UintRange[] permanentlyForbiddenTimes.
func actionPermissionToSolidity(p *tokenizationtypes.ActionPermission) []interface{} {
	if p == nil {
		return []interface{}{uintRangesToSolidity(nil), uintRangesToSolidity(nil)}
	}
	return []interface{}{
		uintRangesToSolidity(p.PermanentlyPermittedTimes),
		uintRangesToSolidity(p.PermanentlyForbiddenTimes),
	}
}

// tokenIdsActionPermissionToSolidity converts TokenIdsActionPermission to Solidity struct format.
func tokenIdsActionPermissionToSolidity(p *tokenizationtypes.TokenIdsActionPermission) []interface{} {
	if p == nil {
		return []interface{}{uintRangesToSolidity(nil), uintRangesToSolidity(nil), uintRangesToSolidity(nil)}
	}
	return []interface{}{
		uintRangesToSolidity(p.TokenIds),
		uintRangesToSolidity(p.PermanentlyPermittedTimes),
		uintRangesToSolidity(p.PermanentlyForbiddenTimes),
	}
}

// collectionApprovalPermissionToSolidity converts CollectionApprovalPermission to Solidity struct format.
func collectionApprovalPermissionToSolidity(p *tokenizationtypes.CollectionApprovalPermission) []interface{} {
	if p == nil {
		return []interface{}{
			"", "", "", uintRangesToSolidity(nil), uintRangesToSolidity(nil), uintRangesToSolidity(nil),
			"", uintRangesToSolidity(nil), uintRangesToSolidity(nil),
		}
	}
	return []interface{}{
		p.FromListId,
		p.ToListId,
		p.InitiatedByListId,
		uintRangesToSolidity(p.TransferTimes),
		uintRangesToSolidity(p.TokenIds),
		uintRangesToSolidity(p.OwnershipTimes),
		p.ApprovalId,
		uintRangesToSolidity(p.PermanentlyPermittedTimes),
		uintRangesToSolidity(p.PermanentlyForbiddenTimes),
	}
}

// userOutgoingApprovalPermissionToSolidity converts UserOutgoingApprovalPermission to Solidity struct format.
func userOutgoingApprovalPermissionToSolidity(p *tokenizationtypes.UserOutgoingApprovalPermission) []interface{} {
	if p == nil {
		return []interface{}{
			"", "", uintRangesToSolidity(nil), uintRangesToSolidity(nil), uintRangesToSolidity(nil),
			"", uintRangesToSolidity(nil), uintRangesToSolidity(nil),
		}
	}
	return []interface{}{
		p.ToListId,
		p.InitiatedByListId,
		uintRangesToSolidity(p.TransferTimes),
		uintRangesToSolidity(p.TokenIds),
		uintRangesToSolidity(p.OwnershipTimes),
		p.ApprovalId,
		uintRangesToSolidity(p.PermanentlyPermittedTimes),
		uintRangesToSolidity(p.PermanentlyForbiddenTimes),
	}
}

// userIncomingApprovalPermissionToSolidity converts UserIncomingApprovalPermission to Solidity struct format.
func userIncomingApprovalPermissionToSolidity(p *tokenizationtypes.UserIncomingApprovalPermission) []interface{} {
	if p == nil {
		return []interface{}{
			"", "", uintRangesToSolidity(nil), uintRangesToSolidity(nil), uintRangesToSolidity(nil),
			"", uintRangesToSolidity(nil), uintRangesToSolidity(nil),
		}
	}
	return []interface{}{
		p.FromListId,
		p.InitiatedByListId,
		uintRangesToSolidity(p.TransferTimes),
		uintRangesToSolidity(p.TokenIds),
		uintRangesToSolidity(p.OwnershipTimes),
		p.ApprovalId,
		uintRangesToSolidity(p.PermanentlyPermittedTimes),
		uintRangesToSolidity(p.PermanentlyForbiddenTimes),
	}
}

// collectionPermissionsToSolidity converts CollectionPermissions to Solidity struct format.
// Returns tuple of 11 arrays: canDeleteCollection ... canAddMoreCosmosCoinWrapperPaths.
func collectionPermissionsToSolidity(p *tokenizationtypes.CollectionPermissions) []interface{} {
	empty := make([]interface{}, 0)
	if p == nil {
		return []interface{}{
			empty, empty, empty, empty, empty,
			empty, empty, empty, empty,
			empty, empty,
		}
	}
	canDelete := make([]interface{}, 0, len(p.CanDeleteCollection))
	for _, x := range p.CanDeleteCollection {
		canDelete = append(canDelete, actionPermissionToSolidity(x))
	}
	canArchive := make([]interface{}, 0, len(p.CanArchiveCollection))
	for _, x := range p.CanArchiveCollection {
		canArchive = append(canArchive, actionPermissionToSolidity(x))
	}
	canUpdateStandards := make([]interface{}, 0, len(p.CanUpdateStandards))
	for _, x := range p.CanUpdateStandards {
		canUpdateStandards = append(canUpdateStandards, actionPermissionToSolidity(x))
	}
	canUpdateCustomData := make([]interface{}, 0, len(p.CanUpdateCustomData))
	for _, x := range p.CanUpdateCustomData {
		canUpdateCustomData = append(canUpdateCustomData, actionPermissionToSolidity(x))
	}
	canUpdateManager := make([]interface{}, 0, len(p.CanUpdateManager))
	for _, x := range p.CanUpdateManager {
		canUpdateManager = append(canUpdateManager, actionPermissionToSolidity(x))
	}
	canUpdateCollMeta := make([]interface{}, 0, len(p.CanUpdateCollectionMetadata))
	for _, x := range p.CanUpdateCollectionMetadata {
		canUpdateCollMeta = append(canUpdateCollMeta, actionPermissionToSolidity(x))
	}
	canUpdateValidIds := make([]interface{}, 0, len(p.CanUpdateValidTokenIds))
	for _, x := range p.CanUpdateValidTokenIds {
		canUpdateValidIds = append(canUpdateValidIds, tokenIdsActionPermissionToSolidity(x))
	}
	canUpdateTokenMeta := make([]interface{}, 0, len(p.CanUpdateTokenMetadata))
	for _, x := range p.CanUpdateTokenMetadata {
		canUpdateTokenMeta = append(canUpdateTokenMeta, tokenIdsActionPermissionToSolidity(x))
	}
	canUpdateCollApprovals := make([]interface{}, 0, len(p.CanUpdateCollectionApprovals))
	for _, x := range p.CanUpdateCollectionApprovals {
		canUpdateCollApprovals = append(canUpdateCollApprovals, collectionApprovalPermissionToSolidity(x))
	}
	canAddAlias := make([]interface{}, 0, len(p.CanAddMoreAliasPaths))
	for _, x := range p.CanAddMoreAliasPaths {
		canAddAlias = append(canAddAlias, actionPermissionToSolidity(x))
	}
	canAddCosmos := make([]interface{}, 0, len(p.CanAddMoreCosmosCoinWrapperPaths))
	for _, x := range p.CanAddMoreCosmosCoinWrapperPaths {
		canAddCosmos = append(canAddCosmos, actionPermissionToSolidity(x))
	}
	return []interface{}{
		canDelete, canArchive, canUpdateStandards, canUpdateCustomData, canUpdateManager,
		canUpdateCollMeta, canUpdateValidIds, canUpdateTokenMeta, canUpdateCollApprovals,
		canAddAlias, canAddCosmos,
	}
}

// userPermissionsToSolidity converts UserPermissions to Solidity struct format.
// Returns tuple of 5 arrays.
func userPermissionsToSolidity(p *tokenizationtypes.UserPermissions) []interface{} {
	empty := make([]interface{}, 0)
	if p == nil {
		return []interface{}{empty, empty, empty, empty, empty}
	}
	outgoing := make([]interface{}, 0, len(p.CanUpdateOutgoingApprovals))
	for _, x := range p.CanUpdateOutgoingApprovals {
		outgoing = append(outgoing, userOutgoingApprovalPermissionToSolidity(x))
	}
	incoming := make([]interface{}, 0, len(p.CanUpdateIncomingApprovals))
	for _, x := range p.CanUpdateIncomingApprovals {
		incoming = append(incoming, userIncomingApprovalPermissionToSolidity(x))
	}
	autoOut := make([]interface{}, 0, len(p.CanUpdateAutoApproveSelfInitiatedOutgoingTransfers))
	for _, x := range p.CanUpdateAutoApproveSelfInitiatedOutgoingTransfers {
		autoOut = append(autoOut, actionPermissionToSolidity(x))
	}
	autoIn := make([]interface{}, 0, len(p.CanUpdateAutoApproveSelfInitiatedIncomingTransfers))
	for _, x := range p.CanUpdateAutoApproveSelfInitiatedIncomingTransfers {
		autoIn = append(autoIn, actionPermissionToSolidity(x))
	}
	allIn := make([]interface{}, 0, len(p.CanUpdateAutoApproveAllIncomingTransfers))
	for _, x := range p.CanUpdateAutoApproveAllIncomingTransfers {
		allIn = append(allIn, actionPermissionToSolidity(x))
	}
	return []interface{}{outgoing, incoming, autoOut, autoIn, allIn}
}

// ConvertApprovalTrackerToSolidityStruct converts ApprovalTracker to Solidity struct format.
// Solidity: numTransfers, Balance[] amounts, lastUpdatedAt.
func ConvertApprovalTrackerToSolidityStruct(t *tokenizationtypes.ApprovalTracker) ([]interface{}, error) {
	if t == nil {
		return nil, fmt.Errorf("approval tracker cannot be nil")
	}
	amounts := make([]interface{}, 0, len(t.Amounts))
	for _, b := range t.Amounts {
		if b == nil {
			continue
		}
		conv, err := ConvertBalanceToSolidityStruct(b)
		if err != nil {
			continue
		}
		amounts = append(amounts, conv)
	}
	numTransfers := big.NewInt(0)
	lastUpdatedAt := big.NewInt(0)
	if !t.NumTransfers.IsNil() {
		numTransfers = t.NumTransfers.BigInt()
	}
	if !t.LastUpdatedAt.IsNil() {
		lastUpdatedAt = t.LastUpdatedAt.BigInt()
	}
	return []interface{}{numTransfers, amounts, lastUpdatedAt}, nil
}

// ConvertDynamicStoreToSolidityStruct converts DynamicStore to Solidity struct format.
func ConvertDynamicStoreToSolidityStruct(d *tokenizationtypes.DynamicStore) ([]interface{}, error) {
	if d == nil {
		return nil, fmt.Errorf("dynamic store cannot be nil")
	}
	storeId := big.NewInt(0)
	if !d.StoreId.IsNil() {
		storeId = d.StoreId.BigInt()
	}
	return []interface{}{
		storeId,
		d.CreatedBy,
		d.DefaultValue,
		d.GlobalEnabled,
		d.Uri,
		d.CustomData,
	}, nil
}

// ConvertDynamicStoreValueToSolidityStruct converts DynamicStoreValue to Solidity struct format.
func ConvertDynamicStoreValueToSolidityStruct(v *tokenizationtypes.DynamicStoreValue) ([]interface{}, error) {
	if v == nil {
		return nil, fmt.Errorf("dynamic store value cannot be nil")
	}
	storeId := big.NewInt(0)
	if !v.StoreId.IsNil() {
		storeId = v.StoreId.BigInt()
	}
	return []interface{}{storeId, v.Address, v.Value}, nil
}

// ConvertVoteProofToSolidityStruct converts VoteProof to Solidity struct format.
func ConvertVoteProofToSolidityStruct(v *tokenizationtypes.VoteProof) ([]interface{}, error) {
	if v == nil {
		return nil, fmt.Errorf("vote proof cannot be nil")
	}
	yesWeight := big.NewInt(0)
	if !v.YesWeight.IsNil() {
		yesWeight = v.YesWeight.BigInt()
	}
	return []interface{}{v.ProposalId, v.Voter, yesWeight}, nil
}

// ConvertParamsToSolidityStruct converts Params to Solidity struct format.
func ConvertParamsToSolidityStruct(p *tokenizationtypes.Params) ([]interface{}, error) {
	if p == nil {
		return nil, fmt.Errorf("params cannot be nil")
	}
	denoms := make([]interface{}, len(p.AllowedDenoms))
	for i, d := range p.AllowedDenoms {
		denoms[i] = d
	}
	affiliate := big.NewInt(0)
	if !p.AffiliatePercentage.IsNil() {
		affiliate = p.AffiliatePercentage.BigInt()
	}
	return []interface{}{denoms, affiliate}, nil
}
