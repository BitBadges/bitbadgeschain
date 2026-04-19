package types

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"
)

func TestIsFullOwnershipTimesRange(t *testing.T) {
	tests := []struct {
		name           string
		ownershipTimes []*UintRange
		expected       bool
	}{
		{
			name: "full range",
			ownershipTimes: []*UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(18446744073709551615), // math.MaxUint64
				},
			},
			expected: true,
		},
		{
			name: "not full range - different start",
			ownershipTimes: []*UintRange{
				{
					Start: sdkmath.NewUint(2),
					End:   sdkmath.NewUint(18446744073709551615),
				},
			},
			expected: false,
		},
		{
			name: "not full range - different end",
			ownershipTimes: []*UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(1000),
				},
			},
			expected: false,
		},
		{
			name: "multiple ranges",
			ownershipTimes: []*UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(200),
				},
			},
			expected: false,
		},
		{
			name:           "empty ranges",
			ownershipTimes: []*UintRange{},
			expected:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsFullOwnershipTimesRange(tt.ownershipTimes)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestValidateNoCustomOwnershipTimesInvariant(t *testing.T) {
	tests := []struct {
		name             string
		ownershipTimes   []*UintRange
		invariantEnabled bool
		expectError      bool
	}{
		{
			name: "invariant disabled - should pass",
			ownershipTimes: []*UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(1000),
				},
			},
			invariantEnabled: false,
			expectError:      false,
		},
		{
			name: "invariant enabled - full range - should pass",
			ownershipTimes: []*UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(18446744073709551615),
				},
			},
			invariantEnabled: true,
			expectError:      false,
		},
		{
			name: "invariant enabled - not full range - should fail",
			ownershipTimes: []*UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(1000),
				},
			},
			invariantEnabled: true,
			expectError:      true,
		},
		{
			name: "invariant enabled - multiple ranges - should fail",
			ownershipTimes: []*UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(100),
				},
				{
					Start: sdkmath.NewUint(101),
					End:   sdkmath.NewUint(200),
				},
			},
			invariantEnabled: true,
			expectError:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateNoCustomOwnershipTimesInvariant(tt.ownershipTimes, tt.invariantEnabled)
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateAltTimeChecks(t *testing.T) {
	tests := []struct {
		name        string
		altTimeChecks *AltTimeChecks
		expectError bool
	}{
		{
			name:        "nil altTimeChecks - should pass",
			altTimeChecks: nil,
			expectError:  false,
		},
		{
			name: "empty ranges - should pass",
			altTimeChecks: &AltTimeChecks{
				OfflineHours: []*UintRange{},
				OfflineDays:  []*UintRange{},
			},
			expectError: false,
		},
		{
			name: "valid single hour range - should pass",
			altTimeChecks: &AltTimeChecks{
				OfflineHours: []*UintRange{
					{
						Start: sdkmath.NewUint(9),
						End:   sdkmath.NewUint(17),
					},
				},
				OfflineDays: []*UintRange{},
			},
			expectError: false,
		},
		{
			name: "valid single day range - should pass",
			altTimeChecks: &AltTimeChecks{
				OfflineHours: []*UintRange{},
				OfflineDays: []*UintRange{
					{
						Start: sdkmath.NewUint(1),
						End:   sdkmath.NewUint(5),
					},
				},
			},
			expectError: false,
		},
		{
			name: "valid multiple hour ranges - should pass",
			altTimeChecks: &AltTimeChecks{
				OfflineHours: []*UintRange{
					{
						Start: sdkmath.NewUint(9),
						End:   sdkmath.NewUint(12),
					},
					{
						Start: sdkmath.NewUint(14),
						End:   sdkmath.NewUint(17),
					},
				},
				OfflineDays: []*UintRange{},
			},
			expectError: false,
		},
		{
			name: "valid ranges including zero - should pass",
			altTimeChecks: &AltTimeChecks{
				OfflineHours: []*UintRange{
					{
						Start: sdkmath.NewUint(0),
						End:   sdkmath.NewUint(5),
					},
				},
				OfflineDays: []*UintRange{
					{
						Start: sdkmath.NewUint(0),
						End:   sdkmath.NewUint(0),
					},
				},
			},
			expectError: false,
		},
		{
			name: "valid ranges at boundaries - should pass",
			altTimeChecks: &AltTimeChecks{
				OfflineHours: []*UintRange{
					{
						Start: sdkmath.NewUint(0),
						End:   sdkmath.NewUint(23),
					},
				},
				OfflineDays: []*UintRange{
					{
						Start: sdkmath.NewUint(0),
						End:   sdkmath.NewUint(6),
					},
				},
			},
			expectError: false,
		},
		{
			name: "hour range exceeds maximum - should fail",
			altTimeChecks: &AltTimeChecks{
				OfflineHours: []*UintRange{
					{
						Start: sdkmath.NewUint(9),
						End:   sdkmath.NewUint(24), // Invalid: max is 23
					},
				},
				OfflineDays: []*UintRange{},
			},
			expectError: true,
		},
		{
			name: "day range exceeds maximum - should fail",
			altTimeChecks: &AltTimeChecks{
				OfflineHours: []*UintRange{},
				OfflineDays: []*UintRange{
					{
						Start: sdkmath.NewUint(1),
						End:   sdkmath.NewUint(7), // Invalid: max is 6
					},
				},
			},
			expectError: true,
		},
		{
			name: "duplicate hour values across ranges - should fail",
			altTimeChecks: &AltTimeChecks{
				OfflineHours: []*UintRange{
					{
						Start: sdkmath.NewUint(9),
						End:   sdkmath.NewUint(12),
					},
					{
						Start: sdkmath.NewUint(11), // Overlaps with previous range (11-12)
						End:   sdkmath.NewUint(15),
					},
				},
				OfflineDays: []*UintRange{},
			},
			expectError: true,
		},
		{
			name: "duplicate day values across ranges - should fail",
			altTimeChecks: &AltTimeChecks{
				OfflineHours: []*UintRange{},
				OfflineDays: []*UintRange{
					{
						Start: sdkmath.NewUint(1),
						End:   sdkmath.NewUint(3),
					},
					{
						Start: sdkmath.NewUint(2), // Overlaps with previous range (2-3)
						End:   sdkmath.NewUint(5),
					},
				},
			},
			expectError: true,
		},
		{
			name: "wrapping hour range (start > end) - should fail",
			altTimeChecks: &AltTimeChecks{
				OfflineHours: []*UintRange{
					{
						Start: sdkmath.NewUint(23),
						End:   sdkmath.NewUint(0), // Invalid: wrapping not allowed
					},
				},
				OfflineDays: []*UintRange{},
			},
			expectError: true,
		},
		{
			name: "wrapping day range (start > end) - should fail",
			altTimeChecks: &AltTimeChecks{
				OfflineHours: []*UintRange{},
				OfflineDays: []*UintRange{
					{
						Start: sdkmath.NewUint(6),
						End:   sdkmath.NewUint(0), // Invalid: wrapping not allowed
					},
				},
			},
			expectError: true,
		},
		{
			name: "nil range in hours - should fail",
			altTimeChecks: &AltTimeChecks{
				OfflineHours: []*UintRange{
					nil,
				},
				OfflineDays: []*UintRange{},
			},
			expectError: true,
		},
		{
			name: "nil range in days - should fail",
			altTimeChecks: &AltTimeChecks{
				OfflineHours: []*UintRange{},
				OfflineDays: []*UintRange{
					nil,
				},
			},
			expectError: true,
		},
		{
			name: "valid separate ranges for wrapping scenario - should pass",
			altTimeChecks: &AltTimeChecks{
				OfflineHours: []*UintRange{
					{
						Start: sdkmath.NewUint(23),
						End:   sdkmath.NewUint(23), // Hour 23
					},
					{
						Start: sdkmath.NewUint(0),
						End:   sdkmath.NewUint(0), // Hour 0 (separate range)
					},
				},
				OfflineDays: []*UintRange{},
			},
			expectError: false,
		},
		// --- v29 field coverage: OfflineMonths (1-12) ---
		{
			name: "valid offline months boundary - should pass",
			altTimeChecks: &AltTimeChecks{
				OfflineMonths: []*UintRange{
					{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(12)},
				},
			},
			expectError: false,
		},
		{
			name: "offline months start below minimum (0) - should fail",
			altTimeChecks: &AltTimeChecks{
				OfflineMonths: []*UintRange{
					{Start: sdkmath.NewUint(0), End: sdkmath.NewUint(5)},
				},
			},
			expectError: true,
		},
		{
			name: "offline months end exceeds 12 - should fail",
			altTimeChecks: &AltTimeChecks{
				OfflineMonths: []*UintRange{
					{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(13)},
				},
			},
			expectError: true,
		},
		// --- v29 field coverage: OfflineDaysOfMonth (1-31) ---
		{
			name: "valid offline days of month boundary - should pass",
			altTimeChecks: &AltTimeChecks{
				OfflineDaysOfMonth: []*UintRange{
					{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(31)},
				},
			},
			expectError: false,
		},
		{
			name: "offline days of month start 0 - should fail",
			altTimeChecks: &AltTimeChecks{
				OfflineDaysOfMonth: []*UintRange{
					{Start: sdkmath.NewUint(0), End: sdkmath.NewUint(15)},
				},
			},
			expectError: true,
		},
		{
			name: "offline days of month end 32 - should fail",
			altTimeChecks: &AltTimeChecks{
				OfflineDaysOfMonth: []*UintRange{
					{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(32)},
				},
			},
			expectError: true,
		},
		// --- v29 field coverage: OfflineWeeksOfYear (1-53, ISO-week upper bound) ---
		{
			name: "valid offline weeks of year boundary 53 - should pass",
			altTimeChecks: &AltTimeChecks{
				OfflineWeeksOfYear: []*UintRange{
					{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(53)},
				},
			},
			expectError: false,
		},
		{
			name: "offline weeks of year start 0 - should fail",
			altTimeChecks: &AltTimeChecks{
				OfflineWeeksOfYear: []*UintRange{
					{Start: sdkmath.NewUint(0), End: sdkmath.NewUint(10)},
				},
			},
			expectError: true,
		},
		{
			name: "offline weeks of year end 54 - should fail",
			altTimeChecks: &AltTimeChecks{
				OfflineWeeksOfYear: []*UintRange{
					{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(54)},
				},
			},
			expectError: true,
		},
		// --- v29 field coverage: TimezoneOffsetMinutes (0..840) ---
		{
			name: "timezone offset zero (UTC) - should pass",
			altTimeChecks: &AltTimeChecks{
				TimezoneOffsetMinutes: sdkmath.NewUint(0),
			},
			expectError: false,
		},
		{
			name: "timezone offset 300 (EST) - should pass",
			altTimeChecks: &AltTimeChecks{
				TimezoneOffsetMinutes: sdkmath.NewUint(300),
			},
			expectError: false,
		},
		{
			name: "timezone offset 840 (UTC+14:00 Kiribati - max real-world) - should pass",
			altTimeChecks: &AltTimeChecks{
				TimezoneOffsetMinutes: sdkmath.NewUint(840),
			},
			expectError: false,
		},
		{
			name: "timezone offset 841 (just over max real-world) - should fail",
			altTimeChecks: &AltTimeChecks{
				TimezoneOffsetMinutes: sdkmath.NewUint(841),
			},
			expectError: true,
		},
		{
			name: "timezone offset in silent-wrap range [2^63+1] - should fail",
			altTimeChecks: &AltTimeChecks{
				// 2^63 + 1 — fits in uint64 but time.Duration(u64) * time.Minute
				// silently wraps through int64 overflow at the consumer.
				TimezoneOffsetMinutes: sdkmath.NewUintFromString("9223372036854775809"),
			},
			expectError: true,
		},
		{
			name: "timezone offset at uint64 max - should fail",
			altTimeChecks: &AltTimeChecks{
				TimezoneOffsetMinutes: sdkmath.NewUintFromString("18446744073709551615"),
			},
			expectError: true,
		},
		{
			name: "timezone offset exceeds uint64 (2^64) - should fail without panic",
			altTimeChecks: &AltTimeChecks{
				// 2^64 — would panic in .Uint64() at the consumer. Must be
				// caught by ValidateBasic with a regular error, not a panic.
				TimezoneOffsetMinutes: sdkmath.NewUintFromString("18446744073709551616"),
			},
			expectError: true,
		},
		{
			name: "timezone offset 200-bit value - should fail without panic",
			altTimeChecks: &AltTimeChecks{
				// Larger than 2^64 and well into sdkmath.Uint territory.
				TimezoneOffsetMinutes: sdkmath.NewUintFromString("1606938044258990275541962092341162602522202993782792835301376"),
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// The fix must never panic — attacker-controlled values going
			// into validate_basic should yield errors, not crashes.
			require.NotPanics(t, func() {
				err := ValidateAltTimeChecks(tt.altTimeChecks)
				if tt.expectError {
					require.Error(t, err, "expected error but got none")
				} else {
					require.NoError(t, err, "expected no error but got: %v", err)
				}
			})
		})
	}
}

func TestEnforceMustPrioritizeForNonAutoScannable(t *testing.T) {
	tests := []struct {
		name     string
		criteria *ApprovalCriteria
		expected bool
	}{
		{
			name:     "nil criteria - no panic",
			criteria: nil,
			expected: false, // not checked
		},
		{
			name:     "auto-scannable + mustPrioritize false - stays false",
			criteria: &ApprovalCriteria{MustPrioritize: false},
			expected: false,
		},
		{
			name:     "auto-scannable + mustPrioritize true - stays true",
			criteria: &ApprovalCriteria{MustPrioritize: true},
			expected: true,
		},
		{
			name: "non-auto-scannable (coinTransfers) + mustPrioritize false - forced true",
			criteria: &ApprovalCriteria{
				MustPrioritize: false,
				CoinTransfers:  []*CoinTransfer{{}},
			},
			expected: true,
		},
		{
			name: "non-auto-scannable (merkleChallenges) + mustPrioritize false - forced true",
			criteria: &ApprovalCriteria{
				MustPrioritize:  false,
				MerkleChallenges: []*MerkleChallenge{{}},
			},
			expected: true,
		},
		{
			name: "non-auto-scannable (ethSignatureChallenges) + mustPrioritize false - forced true",
			criteria: &ApprovalCriteria{
				MustPrioritize:         false,
				EthSignatureChallenges: []*ETHSignatureChallenge{{}},
			},
			expected: true,
		},
		{
			name: "non-auto-scannable (predeterminedBalances) + mustPrioritize false - forced true",
			criteria: &ApprovalCriteria{
				MustPrioritize: false,
				PredeterminedBalances: &PredeterminedBalances{
					ManualBalances: []*ManualBalances{{}},
				},
			},
			expected: true,
		},
		{
			name: "non-auto-scannable + mustPrioritize true - stays true",
			criteria: &ApprovalCriteria{
				MustPrioritize: true,
				CoinTransfers:  []*CoinTransfer{{}},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			EnforceMustPrioritizeForNonAutoScannable(tt.criteria)
			if tt.criteria != nil {
				require.Equal(t, tt.expected, tt.criteria.MustPrioritize)
			}
		})
	}
}

func TestEnforceMustPrioritizeForIncoming(t *testing.T) {
	// auto-scannable - stays false
	ac := &IncomingApprovalCriteria{MustPrioritize: false}
	EnforceMustPrioritizeForIncoming(ac)
	require.False(t, ac.MustPrioritize)

	// non-auto-scannable (coinTransfers) - forced true
	ac2 := &IncomingApprovalCriteria{
		MustPrioritize: false,
		CoinTransfers:  []*CoinTransfer{{}},
	}
	EnforceMustPrioritizeForIncoming(ac2)
	require.True(t, ac2.MustPrioritize)

	// nil - no panic
	EnforceMustPrioritizeForIncoming(nil)
}

func TestEnforceMustPrioritizeForOutgoing(t *testing.T) {
	// auto-scannable - stays false
	ac := &OutgoingApprovalCriteria{MustPrioritize: false}
	EnforceMustPrioritizeForOutgoing(ac)
	require.False(t, ac.MustPrioritize)

	// non-auto-scannable (merkleChallenges) - forced true
	ac2 := &OutgoingApprovalCriteria{
		MustPrioritize:  false,
		MerkleChallenges: []*MerkleChallenge{{}},
	}
	EnforceMustPrioritizeForOutgoing(ac2)
	require.True(t, ac2.MustPrioritize)

	// nil - no panic
	EnforceMustPrioritizeForOutgoing(nil)
}
