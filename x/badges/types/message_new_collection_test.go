package types_test

import (
	// math "math"
	math "math"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/testutil/sample"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	"github.com/stretchr/testify/require"
	// "github.com/bitbadges/bitbadgeschain/testutil/sample"
	// sdk "github.com/cosmos/cosmos-sdk/types"
	// sdkerrors "cosmossdk.io/errors"
	// "github.com/stretchr/testify/require"
	// "github.com/bitbadges/bitbadgeschain/x/badges/types"
)

func TestMsgNewBadge_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgNewCollection
		err  error
	}{
		{
			name: "invalid address",
			msg: types.MsgNewCollection{
				Creator:            "invalid_address",
				BalancesType: sdkmath.NewUint(0),
				CollectionMetadataTimeline: GetValidCollectionMetadataTimeline(),
				BadgeMetadataTimeline: GetValidBadgeMetadataTimeline(),
				Permissions: &types.CollectionPermissions{},
			},
			err: types.ErrInvalidAddress,
		}, {
			name: "valid state",
			msg: types.MsgNewCollection{
				Creator:            sample.AccAddress(),
				BalancesType: sdkmath.NewUint(0),
				CollectionMetadataTimeline: GetValidCollectionMetadataTimeline(),
				BadgeMetadataTimeline: GetValidBadgeMetadataTimeline(),
				Permissions: &types.CollectionPermissions{},
			},
		}, {
			name: "invalid URI",
			msg: types.MsgNewCollection{
				Creator:            sample.AccAddress(),
				BalancesType: sdkmath.NewUint(0),
				CollectionMetadataTimeline: []*types.CollectionMetadataTimeline{
					{
						CollectionMetadata: &types.CollectionMetadata{
							Uri: "",
						},
						TimelineTimes: []*types.UintRange{
							{
								Start: sdkmath.NewUint(0),
								End:  sdkmath.NewUint(math.MaxUint64),
							},
						},
					},
				},
				BadgeMetadataTimeline: GetValidBadgeMetadataTimeline(),
				Permissions: &types.CollectionPermissions{},
			},

			err: types.ErrInvalidURI,
		},
		{
			name: "invalid Badge URI",
			msg: types.MsgNewCollection{
				Creator:            sample.AccAddress(),
				BalancesType: sdkmath.NewUint(0),
				CollectionMetadataTimeline: GetValidCollectionMetadataTimeline(),
				BadgeMetadataTimeline: []*types.BadgeMetadataTimeline{
					{
						BadgeMetadata: []*types.BadgeMetadata{
							{
								Uri: "",
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
								Start: sdkmath.NewUint(0),
								End:  sdkmath.NewUint(math.MaxUint64),
							},
						},
					},
				},
				Permissions: &types.CollectionPermissions{},
			},
			err: types.ErrInvalidURI,
		},
		{
			name: "invalid Permissions",
			msg: types.MsgNewCollection{
				Creator:            sample.AccAddress(),
				BalancesType: sdkmath.NewUint(0),
				CollectionMetadataTimeline: GetValidCollectionMetadataTimeline(),
				BadgeMetadataTimeline: GetValidBadgeMetadataTimeline(),
			},
			err: types.ErrPermissionsIsNil,
		},
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
