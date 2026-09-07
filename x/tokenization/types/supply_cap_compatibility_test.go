package types_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/testutil/sample"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	"github.com/stretchr/testify/require"
)

func TestSupplyCapCompatibilityMessages(t *testing.T) {
	for _, tc := range []struct {
		name    string
		cap     sdkmath.Uint
		backed  bool
		invalid bool
	}{
		{"backed capped", sdkmath.OneUint(), true, true},
		{"backed unlimited", sdkmath.ZeroUint(), true, false},
		{"backed unset", sdkmath.Uint{}, true, false},
		{"unbacked capped", sdkmath.OneUint(), false, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			invariants := &types.InvariantsAddObject{MaxSupplyPerId: tc.cap}
			if tc.backed {
				invariants.CosmosCoinBackedPath = &types.CosmosCoinBackedPathAddObject{Conversion: &types.Conversion{
					SideA: &types.ConversionSideAWithDenom{Amount: sdkmath.OneUint(), Denom: "ubadge"},
					SideB: []*types.Balance{{Amount: sdkmath.OneUint(), TokenIds: []*types.UintRange{{Start: sdkmath.OneUint(), End: sdkmath.OneUint()}}, OwnershipTimes: []*types.UintRange{{Start: sdkmath.OneUint(), End: sdkmath.OneUint()}}}},
				}}
			}
			creator := sample.AccAddress()
			messages := []interface{ ValidateBasic() error }{
				&types.MsgCreateCollection{Creator: creator, Invariants: invariants},
				&types.MsgUniversalUpdateCollection{Creator: creator, CollectionId: sdkmath.ZeroUint(), Invariants: invariants},
				&types.MsgUniversalUpdateCollection{Creator: creator, CollectionId: sdkmath.OneUint(), Invariants: invariants},
			}
			for _, msg := range messages {
				err := msg.ValidateBasic()
				if tc.invalid {
					require.ErrorContains(t, err, "maxSupplyPerId cannot be used with cosmosCoinBackedPath")
					require.ErrorIs(t, err, types.ErrInvalidRequest)
				} else {
					require.NoError(t, err)
				}
			}
		})
	}
}
