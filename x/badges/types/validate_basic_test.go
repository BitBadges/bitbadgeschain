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
