package tokenization

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	tokenizationtypes "github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

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
	var gas uint64 = GasTransferTokensBase
	gas += uint64(len(toAddresses)) * GasPerRecipient
	gas += uint64(len(tokenIdsRanges)) * GasPerTokenIdRange
	gas += uint64(len(ownershipTimesRanges)) * GasPerOwnershipTimeRange
	return gas
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
	var gas uint64 = GasSetIncomingApprovalBase
	gas += uint64(len(transferTimes)) * GasPerApprovalField
	gas += uint64(len(tokenIds)) * GasPerTokenIdRange
	gas += uint64(len(ownershipTimes)) * GasPerOwnershipTimeRange
	return gas
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

