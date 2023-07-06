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

func TestMsgUpdateMetadata_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgUpdateMetadata
		err  error
	}{
		{
			name: "invalid address",
			msg: types.MsgUpdateMetadata{
				Creator:            "invalid_address",
				CollectionId:       sdkmath.NewUint(1),
				CollectionMetadataTimeline: []*types.CollectionMetadataTimeline{
					{
						CollectionMetadata: &types.CollectionMetadata{
							Uri: "https://example.com/{id}",
						},
					},
				},
				BadgeMetadataTimeline: []*types.BadgeMetadataTimeline{
					{
						BadgeMetadata: []*types.BadgeMetadata{
							{
								Uri: "https://example.com/{id}",
								BadgeIds: []*types.UintRange{
									{
										Start: sdkmath.NewUint(1),

										End:   sdkmath.NewUint(math.MaxUint64),
									},
								},
							},
						},
					},
				},
			},
			err: types.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: types.MsgUpdateMetadata{
				Creator:            sample.AccAddress(),
				CollectionId:       sdkmath.NewUint(1),
				CollectionMetadataTimeline: []*types.CollectionMetadataTimeline{
					{
						CollectionMetadata: &types.CollectionMetadata{
							Uri: "https://example.com/{id}",
						},
						TimelineTimes: []*types.UintRange{
							{
								Start: sdkmath.NewUint(1),
								End:  sdkmath.NewUint(math.MaxUint64),
							},
						},
					},
				},
				BadgeMetadataTimeline: []*types.BadgeMetadataTimeline{
					{
						BadgeMetadata: []*types.BadgeMetadata{
							{
								Uri: "https://example.com/{id}",
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
								End:  sdkmath.NewUint(math.MaxUint64),
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
