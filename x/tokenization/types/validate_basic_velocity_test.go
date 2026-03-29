package types

import (
	"math"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

// makeValidApprovalWithApprovalAmounts creates a minimal valid CollectionApproval
// with the given ApprovalAmounts. Keeps other fields valid to isolate velocity
// limit validation logic.
func makeValidApprovalWithApprovalAmounts(amounts *ApprovalAmounts) *CollectionApproval {
	return &CollectionApproval{
		ApprovalId:        "velocity-test",
		FromListId:        "All",
		ToListId:          "All",
		InitiatedByListId: "All",
		TokenIds: []*UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
		},
		TransferTimes: []*UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
		},
		OwnershipTimes: []*UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
		},
		ApprovalCriteria: &ApprovalCriteria{
			ApprovalAmounts: amounts,
		},
		Version: sdkmath.NewUint(0),
	}
}

// makeValidApprovalWithMaxNumTransfers creates a minimal valid CollectionApproval
// with the given MaxNumTransfers.
func makeValidApprovalWithMaxNumTransfers(maxNum *MaxNumTransfers) *CollectionApproval {
	return &CollectionApproval{
		ApprovalId:        "velocity-test-num",
		FromListId:        "All",
		ToListId:          "All",
		InitiatedByListId: "All",
		TokenIds: []*UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
		},
		TransferTimes: []*UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
		},
		OwnershipTimes: []*UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
		},
		ApprovalCriteria: &ApprovalCriteria{
			MaxNumTransfers: maxNum,
		},
		Version: sdkmath.NewUint(0),
	}
}

func testCtx() sdk.Context {
	return sdk.Context{}
}

// --- ApprovalAmounts: mutual exclusivity ---

func TestValidateBasic_ApprovalAmounts_BothResetAndVelocity_Error(t *testing.T) {
	approval := makeValidApprovalWithApprovalAmounts(&ApprovalAmounts{
		OverallApprovalAmount: sdkmath.NewUint(100),
		AmountTrackerId:       "tracker1",
		ResetTimeIntervals: &ResetTimeIntervals{
			StartTime:      sdkmath.NewUint(1000),
			IntervalLength: sdkmath.NewUint(86400000),
		},
		VelocityLimit: &VelocityLimit{
			WindowDuration: sdkmath.NewUint(86400000),
		},
	})

	err := ValidateCollectionApprovals(testCtx(), []*CollectionApproval{approval}, false)
	require.Error(t, err, "setting both resetTimeIntervals and velocityLimit should be rejected")
	require.Contains(t, err.Error(), "mutually exclusive")
}

func TestValidateBasic_ApprovalAmounts_OnlyVelocityLimit_Valid(t *testing.T) {
	approval := makeValidApprovalWithApprovalAmounts(&ApprovalAmounts{
		OverallApprovalAmount: sdkmath.NewUint(100),
		AmountTrackerId:       "tracker1",
		VelocityLimit: &VelocityLimit{
			WindowDuration: sdkmath.NewUint(86400000),
		},
	})

	err := ValidateCollectionApprovals(testCtx(), []*CollectionApproval{approval}, false)
	require.NoError(t, err, "velocityLimit alone should be valid")
}

func TestValidateBasic_ApprovalAmounts_OnlyResetTimeIntervals_Valid(t *testing.T) {
	approval := makeValidApprovalWithApprovalAmounts(&ApprovalAmounts{
		OverallApprovalAmount: sdkmath.NewUint(100),
		AmountTrackerId:       "tracker1",
		ResetTimeIntervals: &ResetTimeIntervals{
			StartTime:      sdkmath.NewUint(1000),
			IntervalLength: sdkmath.NewUint(86400000),
		},
	})

	err := ValidateCollectionApprovals(testCtx(), []*CollectionApproval{approval}, false)
	require.NoError(t, err, "resetTimeIntervals alone should be valid (existing behavior)")
}

// --- ApprovalAmounts: window duration minimum ---

func TestValidateBasic_ApprovalAmounts_VelocityWindowTooSmall_Error(t *testing.T) {
	approval := makeValidApprovalWithApprovalAmounts(&ApprovalAmounts{
		OverallApprovalAmount: sdkmath.NewUint(100),
		AmountTrackerId:       "tracker1",
		VelocityLimit: &VelocityLimit{
			WindowDuration: sdkmath.NewUint(23), // less than MaxVelocityBuckets (24)
		},
	})

	err := ValidateCollectionApprovals(testCtx(), []*CollectionApproval{approval}, false)
	require.Error(t, err, "windowDuration < 24ms should be rejected")
	require.Contains(t, err.Error(), "windowDuration must be at least")
}

func TestValidateBasic_ApprovalAmounts_VelocityWindowExactMinimum_Valid(t *testing.T) {
	approval := makeValidApprovalWithApprovalAmounts(&ApprovalAmounts{
		OverallApprovalAmount: sdkmath.NewUint(100),
		AmountTrackerId:       "tracker1",
		VelocityLimit: &VelocityLimit{
			WindowDuration: sdkmath.NewUint(24), // exactly MaxVelocityBuckets
		},
	})

	err := ValidateCollectionApprovals(testCtx(), []*CollectionApproval{approval}, false)
	require.NoError(t, err, "windowDuration = 24ms should be valid (minimum)")
}

// --- ApprovalAmounts: missing amountTrackerId ---

func TestValidateBasic_ApprovalAmounts_VelocityWithoutTrackerId_Error(t *testing.T) {
	approval := makeValidApprovalWithApprovalAmounts(&ApprovalAmounts{
		OverallApprovalAmount: sdkmath.NewUint(100),
		AmountTrackerId:       "", // missing
		VelocityLimit: &VelocityLimit{
			WindowDuration: sdkmath.NewUint(86400000),
		},
	})

	err := ValidateCollectionApprovals(testCtx(), []*CollectionApproval{approval}, false)
	require.Error(t, err, "velocityLimit without amountTrackerId should be rejected")
	require.Contains(t, err.Error(), "amountTrackerId")
}

// --- MaxNumTransfers: mutual exclusivity ---

func TestValidateBasic_MaxNumTransfers_BothResetAndVelocity_Error(t *testing.T) {
	approval := makeValidApprovalWithMaxNumTransfers(&MaxNumTransfers{
		OverallMaxNumTransfers: sdkmath.NewUint(10),
		AmountTrackerId:        "tracker1",
		ResetTimeIntervals: &ResetTimeIntervals{
			StartTime:      sdkmath.NewUint(1000),
			IntervalLength: sdkmath.NewUint(86400000),
		},
		VelocityLimit: &VelocityLimit{
			WindowDuration: sdkmath.NewUint(86400000),
		},
	})

	err := ValidateCollectionApprovals(testCtx(), []*CollectionApproval{approval}, false)
	require.Error(t, err, "setting both resetTimeIntervals and velocityLimit on MaxNumTransfers should be rejected")
	require.Contains(t, err.Error(), "mutually exclusive")
}

func TestValidateBasic_MaxNumTransfers_OnlyVelocityLimit_Valid(t *testing.T) {
	approval := makeValidApprovalWithMaxNumTransfers(&MaxNumTransfers{
		OverallMaxNumTransfers: sdkmath.NewUint(10),
		AmountTrackerId:        "tracker1",
		VelocityLimit: &VelocityLimit{
			WindowDuration: sdkmath.NewUint(86400000),
		},
	})

	err := ValidateCollectionApprovals(testCtx(), []*CollectionApproval{approval}, false)
	require.NoError(t, err, "velocityLimit alone on MaxNumTransfers should be valid")
}

func TestValidateBasic_MaxNumTransfers_OnlyResetTimeIntervals_Valid(t *testing.T) {
	approval := makeValidApprovalWithMaxNumTransfers(&MaxNumTransfers{
		OverallMaxNumTransfers: sdkmath.NewUint(10),
		AmountTrackerId:        "tracker1",
		ResetTimeIntervals: &ResetTimeIntervals{
			StartTime:      sdkmath.NewUint(1000),
			IntervalLength: sdkmath.NewUint(86400000),
		},
	})

	err := ValidateCollectionApprovals(testCtx(), []*CollectionApproval{approval}, false)
	require.NoError(t, err, "resetTimeIntervals alone on MaxNumTransfers should be valid (existing behavior)")
}

// --- MaxNumTransfers: window duration minimum ---

func TestValidateBasic_MaxNumTransfers_VelocityWindowTooSmall_Error(t *testing.T) {
	approval := makeValidApprovalWithMaxNumTransfers(&MaxNumTransfers{
		OverallMaxNumTransfers: sdkmath.NewUint(10),
		AmountTrackerId:        "tracker1",
		VelocityLimit: &VelocityLimit{
			WindowDuration: sdkmath.NewUint(1), // way below minimum
		},
	})

	err := ValidateCollectionApprovals(testCtx(), []*CollectionApproval{approval}, false)
	require.Error(t, err, "windowDuration < 24ms on MaxNumTransfers should be rejected")
}

// --- MaxNumTransfers: missing amountTrackerId ---

func TestValidateBasic_MaxNumTransfers_VelocityWithoutTrackerId_Error(t *testing.T) {
	approval := makeValidApprovalWithMaxNumTransfers(&MaxNumTransfers{
		OverallMaxNumTransfers: sdkmath.NewUint(10),
		AmountTrackerId:        "", // missing
		VelocityLimit: &VelocityLimit{
			WindowDuration: sdkmath.NewUint(86400000),
		},
	})

	err := ValidateCollectionApprovals(testCtx(), []*CollectionApproval{approval}, false)
	require.Error(t, err, "velocityLimit without amountTrackerId on MaxNumTransfers should be rejected")
}

// --- Neither set: existing behavior unchanged ---

func TestValidateBasic_NeitherResetNorVelocity_Valid(t *testing.T) {
	approval := makeValidApprovalWithApprovalAmounts(&ApprovalAmounts{
		OverallApprovalAmount: sdkmath.NewUint(100),
		AmountTrackerId:       "tracker1",
	})

	err := ValidateCollectionApprovals(testCtx(), []*CollectionApproval{approval}, false)
	require.NoError(t, err, "no resetTimeIntervals or velocityLimit should still be valid")
}

// --- VelocityLimit nil/zero treated as not set ---

func TestValidateBasic_VelocityLimitZeroWindowDuration_TreatedAsNotSet(t *testing.T) {
	// A zero-window VelocityLimit should be treated as nil (not set),
	// so no mutual exclusivity error even if resetTimeIntervals is also set.
	approval := makeValidApprovalWithApprovalAmounts(&ApprovalAmounts{
		OverallApprovalAmount: sdkmath.NewUint(100),
		AmountTrackerId:       "tracker1",
		ResetTimeIntervals: &ResetTimeIntervals{
			StartTime:      sdkmath.NewUint(1000),
			IntervalLength: sdkmath.NewUint(86400000),
		},
		VelocityLimit: &VelocityLimit{
			WindowDuration: sdkmath.NewUint(0), // zero = basically nil
		},
	})

	err := ValidateCollectionApprovals(testCtx(), []*CollectionApproval{approval}, false)
	require.NoError(t, err, "zero-window velocityLimit should be treated as not set")
}
