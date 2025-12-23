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

func TestMsgNewBadge_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgUniversalUpdateCollection
		err  error
	}{
		{
			name: "invalid address",
			msg: types.MsgUniversalUpdateCollection{
				Creator:                     "invalid_address",
				CollectionId:                sdkmath.NewUint(0),
				UpdateCollectionMetadata:    true,
				CollectionMetadata:          GetValidCollectionMetadata(),
				UpdateTokenMetadata:         true,
				TokenMetadata:               GetValidTokenMetadata(),
				UpdateCollectionPermissions: true,
				CollectionPermissions:       &types.CollectionPermissions{},
			},
			err: types.ErrInvalidAddress,
		}, {
			name: "valid state",
			msg: types.MsgUniversalUpdateCollection{
				Creator:                     sample.AccAddress(),
				CollectionId:                sdkmath.NewUint(0),
				UpdateCollectionMetadata:    true,
				UpdateCollectionPermissions: true,
				UpdateTokenMetadata:         true,
				CollectionMetadata:          GetValidCollectionMetadata(),
				TokenMetadata:               GetValidTokenMetadata(),
				CollectionPermissions:       &types.CollectionPermissions{},
			},
		}, {
			name: "invalid URI",
			msg: types.MsgUniversalUpdateCollection{
				Creator:                     sample.AccAddress(),
				CollectionId:                sdkmath.NewUint(0),
				UpdateCollectionMetadata:    true,
				UpdateCollectionPermissions: true,
				UpdateTokenMetadata:         true,
				CollectionMetadata: &types.CollectionMetadata{
					Uri: "asdfasdfasdf",
				},
				TokenMetadata:         GetValidTokenMetadata(),
				CollectionPermissions: &types.CollectionPermissions{},
			},

			err: types.ErrInvalidURI,
		},
		{
			name: "invalid Token URI",
			msg: types.MsgUniversalUpdateCollection{
				Creator:                     sample.AccAddress(),
				CollectionId:                sdkmath.NewUint(0),
				CollectionMetadata:          GetValidCollectionMetadata(),
				UpdateCollectionMetadata:    true,
				UpdateCollectionPermissions: true,
				UpdateTokenMetadata:         true,
				TokenMetadata: []*types.TokenMetadata{
					{
						Uri: "asdfasdfas",
						TokenIds: []*types.UintRange{
							{
								Start: sdkmath.NewUint(1),
								End:   sdkmath.NewUint(math.MaxUint64),
							},
						},
					},
				},
				CollectionPermissions: &types.CollectionPermissions{},
			},
			err: types.ErrInvalidURI,
		},
		// {
		// 	name: "invalid Permissions",
		// 	msg: types.MsgUniversalUpdateCollection{
		// 		Creator:                    sample.AccAddress(),
		// 		CollectionId: 						 sdkmath.NewUint(0),
		// 		BalancesType:               "Standard",
		// 		UpdateCollectionMetadataTimeline: true,
		// 		UpdateCollectionPermissions: true,
		// 		UpdateTokenMetadataTimeline: true,
		// 		CollectionMetadataTimeline: GetValidCollectionMetadataTimeline(),
		// 		TokenMetadataTimeline:      GetValidTokenMetadataTimeline(),
		// 	},
		// 	err: types.ErrPermissionsIsNil,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.Error(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
