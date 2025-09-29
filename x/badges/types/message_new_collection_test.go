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
				Creator:                          "invalid_address",
				CollectionId:                     sdkmath.NewUint(0),
				UpdateCollectionMetadataTimeline: true,
				CollectionMetadataTimeline:       GetValidCollectionMetadataTimeline(),
				UpdateBadgeMetadataTimeline:      true,
				BadgeMetadataTimeline:            GetValidBadgeMetadataTimeline(),
				UpdateCollectionPermissions:      true,
				CollectionPermissions:            &types.CollectionPermissions{},
			},
			err: types.ErrInvalidAddress,
		}, {
			name: "valid state",
			msg: types.MsgUniversalUpdateCollection{
				Creator:                          sample.AccAddress(),
				CollectionId:                     sdkmath.NewUint(0),
				UpdateCollectionMetadataTimeline: true,
				UpdateCollectionPermissions:      true,
				UpdateBadgeMetadataTimeline:      true,
				CollectionMetadataTimeline:       GetValidCollectionMetadataTimeline(),
				BadgeMetadataTimeline:            GetValidBadgeMetadataTimeline(),
				CollectionPermissions:            &types.CollectionPermissions{},
			},
		}, {
			name: "invalid URI",
			msg: types.MsgUniversalUpdateCollection{
				Creator:                          sample.AccAddress(),
				CollectionId:                     sdkmath.NewUint(0),
				UpdateCollectionMetadataTimeline: true,
				UpdateCollectionPermissions:      true,
				UpdateBadgeMetadataTimeline:      true,
				CollectionMetadataTimeline: []*types.CollectionMetadataTimeline{
					{
						CollectionMetadata: &types.CollectionMetadata{
							Uri: "asdfasdfasdf",
						},
						TimelineTimes: []*types.UintRange{
							{
								Start: sdkmath.NewUint(1),
								End:   sdkmath.NewUint(math.MaxUint64),
							},
						},
					},
				},
				BadgeMetadataTimeline: GetValidBadgeMetadataTimeline(),
				CollectionPermissions: &types.CollectionPermissions{},
			},

			err: types.ErrInvalidURI,
		},
		{
			name: "invalid Token URI",
			msg: types.MsgUniversalUpdateCollection{
				Creator:                          sample.AccAddress(),
				CollectionId:                     sdkmath.NewUint(0),
				CollectionMetadataTimeline:       GetValidCollectionMetadataTimeline(),
				UpdateCollectionMetadataTimeline: true,
				UpdateCollectionPermissions:      true,
				UpdateBadgeMetadataTimeline:      true,
				BadgeMetadataTimeline: []*types.BadgeMetadataTimeline{
					{
						BadgeMetadata: []*types.BadgeMetadata{
							{
								Uri: "asdfasdfas",
								BadgeIds: []*types.UintRange{
									{
										Start: sdkmath.NewUint(1),
										End:   sdkmath.NewUint(math.MaxUint64),
									},
								},
							},
						},
						TimelineTimes: []*types.UintRange{
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
		// 		UpdateBadgeMetadataTimeline: true,
		// 		CollectionMetadataTimeline: GetValidCollectionMetadataTimeline(),
		// 		BadgeMetadataTimeline:      GetValidBadgeMetadataTimeline(),
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
