package validation

import (
	"math"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

// GenerateValidUintRange generates a valid UintRange for standard contexts (TokenIds, TransferTimes, OwnershipTimes)
func GenerateValidUintRange(start, end uint64) *types.UintRange {
	return &types.UintRange{
		Start: sdkmath.NewUint(start),
		End:   sdkmath.NewUint(end),
	}
}

// GenerateValidUintRangeMax generates a valid UintRange covering full spectrum
func GenerateValidUintRangeMax() *types.UintRange {
	return &types.UintRange{
		Start: sdkmath.NewUint(1),
		End:   sdkmath.NewUint(math.MaxUint64),
	}
}

// GenerateInvalidUintRangeNilStart generates a UintRange with nil Start
func GenerateInvalidUintRangeNilStart() *types.UintRange {
	return &types.UintRange{
		Start: sdkmath.Uint{},
		End:   sdkmath.NewUint(100),
	}
}

// GenerateInvalidUintRangeNilEnd generates a UintRange with nil End
func GenerateInvalidUintRangeNilEnd() *types.UintRange {
	return &types.UintRange{
		Start: sdkmath.NewUint(1),
		End:   sdkmath.Uint{},
	}
}

// GenerateInvalidUintRangeZeroStart generates a UintRange with zero Start (invalid when allowAllUints=false)
func GenerateInvalidUintRangeZeroStart() *types.UintRange {
	return &types.UintRange{
		Start: sdkmath.NewUint(0),
		End:   sdkmath.NewUint(100),
	}
}

// GenerateInvalidUintRangeZeroEnd generates a UintRange with zero End (invalid when allowAllUints=false)
func GenerateInvalidUintRangeZeroEnd() *types.UintRange {
	return &types.UintRange{
		Start: sdkmath.NewUint(1),
		End:   sdkmath.NewUint(0),
	}
}

// GenerateInvalidUintRangeStartGreaterThanEnd generates a UintRange with Start > End
func GenerateInvalidUintRangeStartGreaterThanEnd() *types.UintRange {
	return &types.UintRange{
		Start: sdkmath.NewUint(100),
		End:   sdkmath.NewUint(1),
	}
}

// GenerateOverlappingUintRanges generates two overlapping UintRanges
func GenerateOverlappingUintRanges() ([]*types.UintRange, []*types.UintRange) {
	// Exact overlap
	exactOverlap := []*types.UintRange{
		GenerateValidUintRange(1, 10),
		GenerateValidUintRange(1, 10),
	}

	// Partial overlap
	partialOverlap := []*types.UintRange{
		GenerateValidUintRange(1, 10),
		GenerateValidUintRange(5, 15),
	}

	return exactOverlap, partialOverlap
}

// GenerateAdjacentUintRanges generates two adjacent UintRanges
func GenerateAdjacentUintRanges() []*types.UintRange {
	return []*types.UintRange{
		GenerateValidUintRange(1, 5),
		GenerateValidUintRange(6, 10),
	}
}

// GenerateGappedUintRanges generates two UintRanges with a gap
func GenerateGappedUintRanges() []*types.UintRange {
	return []*types.UintRange{
		GenerateValidUintRange(1, 5),
		GenerateValidUintRange(7, 10),
	}
}

// GenerateContainedUintRanges generates ranges where one is fully contained
func GenerateContainedUintRanges() []*types.UintRange {
	return []*types.UintRange{
		GenerateValidUintRange(1, 10),
		GenerateValidUintRange(3, 7),
	}
}

// GenerateAltTimeHoursRanges generates valid OfflineHours ranges (0-23)
func GenerateAltTimeHoursRanges() []*types.UintRange {
	return []*types.UintRange{
		{
			Start: sdkmath.NewUint(0),
			End:   sdkmath.NewUint(23),
		},
	}
}

// GenerateInvalidAltTimeHoursRanges generates invalid OfflineHours ranges
func GenerateInvalidAltTimeHoursRanges() []*types.UintRange {
	return []*types.UintRange{
		{
			Start: sdkmath.NewUint(0),
			End:   sdkmath.NewUint(24), // Invalid: exceeds max 23
		},
	}
}

// GenerateAltTimeDaysRanges generates valid OfflineDays ranges (0-6)
func GenerateAltTimeDaysRanges() []*types.UintRange {
	return []*types.UintRange{
		{
			Start: sdkmath.NewUint(0),
			End:   sdkmath.NewUint(6),
		},
	}
}

// GenerateInvalidAltTimeDaysRanges generates invalid OfflineDays ranges
func GenerateInvalidAltTimeDaysRanges() []*types.UintRange {
	return []*types.UintRange{
		{
			Start: sdkmath.NewUint(0),
			End:   sdkmath.NewUint(7), // Invalid: exceeds max 6
		},
	}
}

// GenerateValidBalance generates a valid Balance
func GenerateValidBalance(amount uint64, tokenStart, tokenEnd, timeStart, timeEnd uint64) *types.Balance {
	return &types.Balance{
		Amount: sdkmath.NewUint(amount),
		TokenIds: []*types.UintRange{
			GenerateValidUintRange(tokenStart, tokenEnd),
		},
		OwnershipTimes: []*types.UintRange{
			GenerateValidUintRange(timeStart, timeEnd),
		},
	}
}

// GenerateInvalidBalanceZeroAmount generates a Balance with zero amount
func GenerateInvalidBalanceZeroAmount() *types.Balance {
	return &types.Balance{
		Amount:         sdkmath.NewUint(0),
		TokenIds:       []*types.UintRange{GenerateValidUintRange(1, 10)},
		OwnershipTimes: []*types.UintRange{GenerateValidUintRange(1, 100)},
	}
}

// GenerateInvalidBalanceNilAmount generates a Balance with nil amount
func GenerateInvalidBalanceNilAmount() *types.Balance {
	return &types.Balance{
		Amount:         sdkmath.Uint{},
		TokenIds:       []*types.UintRange{GenerateValidUintRange(1, 10)},
		OwnershipTimes: []*types.UintRange{GenerateValidUintRange(1, 100)},
	}
}

// GenerateInvalidBalanceEmptyTokenIds generates a Balance with empty TokenIds
func GenerateInvalidBalanceEmptyTokenIds() *types.Balance {
	return &types.Balance{
		Amount:         sdkmath.NewUint(1),
		TokenIds:       []*types.UintRange{},
		OwnershipTimes: []*types.UintRange{GenerateValidUintRange(1, 100)},
	}
}

// GenerateInvalidBalanceEmptyOwnershipTimes generates a Balance with empty OwnershipTimes
func GenerateInvalidBalanceEmptyOwnershipTimes() *types.Balance {
	return &types.Balance{
		Amount:         sdkmath.NewUint(1),
		TokenIds:       []*types.UintRange{GenerateValidUintRange(1, 10)},
		OwnershipTimes: []*types.UintRange{},
	}
}

// GenerateDuplicateBalances generates two balances with exact same (TokenIds, OwnershipTimes)
func GenerateDuplicateBalances() ([]*types.Balance, []*types.Balance) {
	// Exact duplicate
	exactDuplicate := []*types.Balance{
		GenerateValidBalance(1, 1, 10, 1, 100),
		GenerateValidBalance(2, 1, 10, 1, 100), // Same TokenIds and OwnershipTimes, different amount
	}

	// Overlapping TokenIds, same OwnershipTimes
	overlappingTokenIds := []*types.Balance{
		GenerateValidBalance(1, 1, 10, 1, 100),
		GenerateValidBalance(2, 5, 15, 1, 100), // Overlapping TokenIds, same OwnershipTimes
	}

	return exactDuplicate, overlappingTokenIds
}

// GenerateBalanceWithOverlappingTokenIds generates a balance with overlapping TokenIds within itself
func GenerateBalanceWithOverlappingTokenIds() *types.Balance {
	return &types.Balance{
		Amount: sdkmath.NewUint(1),
		TokenIds: []*types.UintRange{
			GenerateValidUintRange(1, 10),
			GenerateValidUintRange(5, 20), // Overlaps with first range
		},
		OwnershipTimes: []*types.UintRange{
			GenerateValidUintRange(1, 100),
		},
	}
}

// GenerateBalanceWithOverlappingOwnershipTimes generates a balance with overlapping OwnershipTimes within itself
func GenerateBalanceWithOverlappingOwnershipTimes() *types.Balance {
	return &types.Balance{
		Amount: sdkmath.NewUint(1),
		TokenIds: []*types.UintRange{
			GenerateValidUintRange(1, 10),
		},
		OwnershipTimes: []*types.UintRange{
			GenerateValidUintRange(1, 100),
			GenerateValidUintRange(50, 200), // Overlaps with first range
		},
	}
}

// CreateTestContext creates a test context for validation tests
func CreateTestContext() sdk.Context {
	return sdk.Context{}
}
