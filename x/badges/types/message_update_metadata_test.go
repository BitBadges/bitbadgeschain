package types_test

import (
	// math "math"
	math "math"
	"testing"

	"github.com/bitbadges/bitbadgeschain/x/badges/testutil/sample"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"
)

func TestMsgUpdateMetadata_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgUniversalUpdateCollection
		err  error
	}{
		{
			name: "invalid address",
			msg: types.MsgUniversalUpdateCollection{
				Creator:                  "invalid_address",
				CollectionId:             sdkmath.NewUint(1),
				UpdateCollectionMetadata: true,
				UpdateTokenMetadata:      true,
				CollectionMetadata: &types.CollectionMetadata{
					Uri: "https://example.com/{id}",
				},
				TokenMetadata: []*types.TokenMetadata{
					{
						Uri: "https://example.com/{id}",
						TokenIds: []*types.UintRange{
							{
								Start: sdkmath.NewUint(1),
								End:   sdkmath.NewUint(math.MaxUint64),
							},
						},
					},
				},
			},
			err: types.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: types.MsgUniversalUpdateCollection{
				Creator:                  sample.AccAddress(),
				CollectionId:             sdkmath.NewUint(1),
				UpdateCollectionMetadata: true,
				UpdateTokenMetadata:      true,
				CollectionMetadata: &types.CollectionMetadata{
					Uri: "https://example.com/{id}",
				},
				TokenMetadata: []*types.TokenMetadata{
					{
						Uri: "https://example.com/{id}",
						TokenIds: []*types.UintRange{
							{
								Start: sdkmath.NewUint(1),
								End:   sdkmath.NewUint(math.MaxUint64),
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
