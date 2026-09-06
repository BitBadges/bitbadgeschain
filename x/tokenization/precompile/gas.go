package tokenization

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	sdk "github.com/cosmos/cosmos-sdk/types"

	tokenizationtypes "github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

// MeterMessage validates the size of a parsed message's variable-length parts
// against the caps in security.go and charges their per-element gas on the
// SDK meter, which RunNativeAction bills back to the EVM.
//
// RequiredGas runs before the JSON is parsed and can only price the raw
// input; this is the other half, taken once the message shape is known and
// before the keeper does any work.
func MeterMessage(ctx sdk.Context, msg sdk.Msg) error {
	gas, err := messageGas(msg)
	if err != nil {
		return err
	}
	ctx.GasMeter().ConsumeGas(gas, "precompile message elements")
	return nil
}

func messageGas(msg sdk.Msg) (uint64, error) {
	switch m := msg.(type) {
	case *tokenizationtypes.MsgTransferTokens:
		return transfersGas(m.Transfers)
	case *tokenizationtypes.MsgSetIncomingApproval:
		return incomingApprovalsGas([]*tokenizationtypes.UserIncomingApproval{m.Approval})
	case *tokenizationtypes.MsgSetOutgoingApproval:
		return outgoingApprovalsGas([]*tokenizationtypes.UserOutgoingApproval{m.Approval})
	case *tokenizationtypes.MsgUpdateUserApprovals:
		return sumGas(
			func() (uint64, error) { return incomingApprovalsGas(m.IncomingApprovals) },
			func() (uint64, error) { return outgoingApprovalsGas(m.OutgoingApprovals) },
		)
	case *tokenizationtypes.MsgSetCollectionApprovals:
		return collectionApprovalsGas(m.CollectionApprovals)
	case *tokenizationtypes.MsgCreateCollection:
		return collectionFieldsGas(m.CollectionMetadata, m.TokenMetadata, m.CollectionApprovals, m.CustomData)
	case *tokenizationtypes.MsgUpdateCollection:
		return collectionFieldsGas(m.CollectionMetadata, m.TokenMetadata, m.CollectionApprovals, m.CustomData)
	case *tokenizationtypes.MsgUniversalUpdateCollection:
		return collectionFieldsGas(m.CollectionMetadata, m.TokenMetadata, m.CollectionApprovals, m.CustomData)
	case *tokenizationtypes.MsgSetCollectionMetadata:
		return 0, collectionMetadataSize(m.CollectionMetadata)
	case *tokenizationtypes.MsgSetTokenMetadata:
		return tokenMetadataGas(m.TokenMetadata)
	case *tokenizationtypes.MsgSetCustomData:
		return 0, ValidateMetadataLength(m.CustomData, "customData")
	case *tokenizationtypes.MsgCreateAddressLists:
		return addressListsGas(m.AddressLists)
	}
	return 0, nil
}

func sumGas(parts ...func() (uint64, error)) (uint64, error) {
	var total uint64
	for _, part := range parts {
		gas, err := part()
		if err != nil {
			return 0, err
		}
		total += gas
	}
	return total, nil
}

func transferElementsGas(recipients, tokenIdRanges, ownershipTimeRanges int) uint64 {
	return uint64(recipients)*GasPerRecipient +
		uint64(tokenIdRanges)*GasPerTokenIdRange +
		uint64(ownershipTimeRanges)*GasPerOwnershipTimeRange
}

func approvalElementsGas(transferTimes, tokenIds, ownershipTimes int) uint64 {
	return uint64(transferTimes)*GasPerApprovalField +
		uint64(tokenIds)*GasPerTokenIdRange +
		uint64(ownershipTimes)*GasPerOwnershipTimeRange
}

func transfersGas(transfers []*tokenizationtypes.Transfer) (uint64, error) {
	var gas uint64
	for _, t := range transfers {
		if t == nil {
			continue
		}
		if err := ValidateArraySizeAllowEmpty(len(t.ToAddresses), MaxRecipients, "toAddresses"); err != nil {
			return 0, err
		}
		tokenIds, ownershipTimes := 0, 0
		for _, b := range t.Balances {
			if b == nil {
				continue
			}
			tokenIds += len(b.TokenIds)
			ownershipTimes += len(b.OwnershipTimes)
		}
		if err := ValidateArraySizeAllowEmpty(tokenIds, MaxTokenIdRanges, "tokenIds"); err != nil {
			return 0, err
		}
		if err := ValidateArraySizeAllowEmpty(ownershipTimes, MaxOwnershipTimeRanges, "ownershipTimes"); err != nil {
			return 0, err
		}
		if err := ValidateMetadataLength(t.Memo, "memo"); err != nil {
			return 0, err
		}
		gas += transferElementsGas(len(t.ToAddresses), tokenIds, ownershipTimes)
	}
	return gas, nil
}

// approvalGas prices the ranges shared by every approval type and checks the
// caps on its metadata and criteria lists.
func approvalGas(transferTimes, tokenIds, ownershipTimes []*tokenizationtypes.UintRange, uri, customData string,
	merkleChallenges, coinTransfers, mustOwnTokens, dynamicStoreChallenges, ethSignatureChallenges int,
) (uint64, error) {
	if err := ValidateArraySizeAllowEmpty(len(transferTimes), MaxApprovalRanges, "transferTimes"); err != nil {
		return 0, err
	}
	if err := ValidateArraySizeAllowEmpty(len(tokenIds), MaxTokenIdRanges, "tokenIds"); err != nil {
		return 0, err
	}
	if err := ValidateArraySizeAllowEmpty(len(ownershipTimes), MaxOwnershipTimeRanges, "ownershipTimes"); err != nil {
		return 0, err
	}
	if err := ValidateMetadataLength(uri, "uri"); err != nil {
		return 0, err
	}
	if err := ValidateMetadataLength(customData, "customData"); err != nil {
		return 0, err
	}
	if err := ValidateMerkleChallengesSize(merkleChallenges); err != nil {
		return 0, err
	}
	if err := ValidateCoinTransfersSize(coinTransfers); err != nil {
		return 0, err
	}
	if err := ValidateArraySizeAllowEmpty(mustOwnTokens, MaxMustOwnTokens, "mustOwnTokens"); err != nil {
		return 0, err
	}
	if err := ValidateArraySizeAllowEmpty(dynamicStoreChallenges, MaxDynamicStoreChallenges, "dynamicStoreChallenges"); err != nil {
		return 0, err
	}
	if err := ValidateArraySizeAllowEmpty(ethSignatureChallenges, MaxETHSignatureChallenges, "ethSignatureChallenges"); err != nil {
		return 0, err
	}
	return approvalElementsGas(len(transferTimes), len(tokenIds), len(ownershipTimes)), nil
}

func incomingApprovalsGas(approvals []*tokenizationtypes.UserIncomingApproval) (uint64, error) {
	var gas uint64
	for _, a := range approvals {
		if a == nil {
			continue
		}
		c := a.ApprovalCriteria
		if c == nil {
			c = &tokenizationtypes.IncomingApprovalCriteria{}
		}
		g, err := approvalGas(a.TransferTimes, a.TokenIds, a.OwnershipTimes, a.Uri, a.CustomData,
			len(c.MerkleChallenges), len(c.CoinTransfers), len(c.MustOwnTokens), len(c.DynamicStoreChallenges), len(c.EthSignatureChallenges))
		if err != nil {
			return 0, err
		}
		gas += g
	}
	return gas, nil
}

func outgoingApprovalsGas(approvals []*tokenizationtypes.UserOutgoingApproval) (uint64, error) {
	var gas uint64
	for _, a := range approvals {
		if a == nil {
			continue
		}
		c := a.ApprovalCriteria
		if c == nil {
			c = &tokenizationtypes.OutgoingApprovalCriteria{}
		}
		g, err := approvalGas(a.TransferTimes, a.TokenIds, a.OwnershipTimes, a.Uri, a.CustomData,
			len(c.MerkleChallenges), len(c.CoinTransfers), len(c.MustOwnTokens), len(c.DynamicStoreChallenges), len(c.EthSignatureChallenges))
		if err != nil {
			return 0, err
		}
		gas += g
	}
	return gas, nil
}

func collectionApprovalsGas(approvals []*tokenizationtypes.CollectionApproval) (uint64, error) {
	var gas uint64
	for _, a := range approvals {
		if a == nil {
			continue
		}
		c := a.ApprovalCriteria
		if c == nil {
			c = &tokenizationtypes.ApprovalCriteria{}
		}
		g, err := approvalGas(a.TransferTimes, a.TokenIds, a.OwnershipTimes, a.Uri, a.CustomData,
			len(c.MerkleChallenges), len(c.CoinTransfers), len(c.MustOwnTokens), len(c.DynamicStoreChallenges), len(c.EthSignatureChallenges))
		if err != nil {
			return 0, err
		}
		gas += g
	}
	return gas, nil
}

func collectionMetadataSize(md *tokenizationtypes.CollectionMetadata) error {
	if md == nil {
		return nil
	}
	if err := ValidateMetadataLength(md.Uri, "collectionMetadata.uri"); err != nil {
		return err
	}
	return ValidateMetadataLength(md.CustomData, "collectionMetadata.customData")
}

func tokenMetadataGas(entries []*tokenizationtypes.TokenMetadata) (uint64, error) {
	var gas uint64
	for _, md := range entries {
		if md == nil {
			continue
		}
		if err := ValidateMetadataLength(md.Uri, "tokenMetadata.uri"); err != nil {
			return 0, err
		}
		if err := ValidateMetadataLength(md.CustomData, "tokenMetadata.customData"); err != nil {
			return 0, err
		}
		if err := ValidateArraySizeAllowEmpty(len(md.TokenIds), MaxTokenIdRanges, "tokenMetadata.tokenIds"); err != nil {
			return 0, err
		}
		gas += uint64(len(md.TokenIds)) * GasPerTokenIdRange
	}
	return gas, nil
}

func collectionFieldsGas(md *tokenizationtypes.CollectionMetadata, tokenMetadata []*tokenizationtypes.TokenMetadata,
	approvals []*tokenizationtypes.CollectionApproval, customData string,
) (uint64, error) {
	if err := collectionMetadataSize(md); err != nil {
		return 0, err
	}
	if err := ValidateMetadataLength(customData, "customData"); err != nil {
		return 0, err
	}
	return sumGas(
		func() (uint64, error) { return tokenMetadataGas(tokenMetadata) },
		func() (uint64, error) { return collectionApprovalsGas(approvals) },
	)
}

func addressListsGas(lists []*tokenizationtypes.AddressListInput) (uint64, error) {
	var gas uint64
	for _, l := range lists {
		if l == nil {
			continue
		}
		if err := ValidateArraySizeAllowEmpty(len(l.Addresses), MaxAddressListEntries, "addresses"); err != nil {
			return 0, err
		}
		if err := ValidateMetadataLength(l.Uri, "uri"); err != nil {
			return 0, err
		}
		if err := ValidateMetadataLength(l.CustomData, "customData"); err != nil {
			return 0, err
		}
		gas += uint64(len(l.Addresses)) * GasPerAddressListEntry
	}
	return gas, nil
}

// CalculateTransferGas calculates dynamic gas for transfer operations
func CalculateTransferGas(
	toAddresses []common.Address,
	tokenIdsRanges []struct {
		Start *big.Int `json:"start"`
		End   *big.Int `json:"end"`
	},
	ownershipTimesRanges []struct {
		Start *big.Int `json:"start"`
		End   *big.Int `json:"end"`
	},
) uint64 {
	return GasTransferTokensBase + transferElementsGas(len(toAddresses), len(tokenIdsRanges), len(ownershipTimesRanges))
}

// CalculateApprovalGas calculates dynamic gas for approval operations
func CalculateApprovalGas(
	transferTimes []struct {
		Start *big.Int `json:"start"`
		End   *big.Int `json:"end"`
	},
	tokenIds []struct {
		Start *big.Int `json:"start"`
		End   *big.Int `json:"end"`
	},
	ownershipTimes []struct {
		Start *big.Int `json:"start"`
		End   *big.Int `json:"end"`
	},
) uint64 {
	return GasSetIncomingApprovalBase + approvalElementsGas(len(transferTimes), len(tokenIds), len(ownershipTimes))
}

// CalculateQueryGas calculates dynamic gas for query operations with ranges
func CalculateQueryGas(
	tokenIdsRanges []struct {
		Start *big.Int `json:"start"`
		End   *big.Int `json:"end"`
	},
	ownershipTimesRanges []struct {
		Start *big.Int `json:"start"`
		End   *big.Int `json:"end"`
	},
	baseGas uint64,
) uint64 {
	var gas uint64 = baseGas
	gas += uint64(len(tokenIdsRanges)) * GasPerQueryRange
	gas += uint64(len(ownershipTimesRanges)) * GasPerQueryRange
	return gas
}

// CalculateQueryGasFromUintRanges calculates dynamic gas from UintRange slices
func CalculateQueryGasFromUintRanges(
	tokenIds []*tokenizationtypes.UintRange,
	ownershipTimes []*tokenizationtypes.UintRange,
	baseGas uint64,
) uint64 {
	var gas uint64 = baseGas
	gas += uint64(len(tokenIds)) * GasPerQueryRange
	gas += uint64(len(ownershipTimes)) * GasPerQueryRange
	return gas
}

