package types_test

import (
	// math "math"
	"testing"

	// "github.com/bitbadges/bitbadgeschain/testutil/sample"
	// sdk "github.com/cosmos/cosmos-sdk/types"
	// sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	// "github.com/stretchr/testify/require"

	// "github.com/bitbadges/bitbadgeschain/x/badges/types"
)

func TestMsgNewBadge_ValidateBasic(t *testing.T) {
	// tests := []struct {
	// 	name string
	// 	msg  types.MsgNewCollection
	// 	err  error
	// }{
	// 	{
	// 		name: "invalid address",
	// 		msg: types.MsgNewCollection{
	// 			Creator:            "invalid_address",
	// 			Standard:           sdk.NewUint(0),
	// 			CollectionMetadata: "https://example.com",
	// 			BadgeMetadata: []*types.BadgeMetadata{
	// 				{
	// 					Uri: "https://example.com/{id}",
	// 					BadgeIds: []*types.IdRange{
	// 						{
	// 							Start: sdk.NewUint(1),
	// 							End:   sdk.NewUint(math.MaxUint64),
	// 						},
	// 					},
	// 				},
	// 			},
	// 			Permissions:       sdk.NewUint(15),
	// 			ApprovedTransfers: []*types.CollectionApprovedTransfer{},
	// 		},
	// 		err: sdkerrors.ErrInvalidAddress,
	// 	}, {
	// 		name: "valid state",
	// 		msg: types.MsgNewCollection{
	// 			Creator:            sample.AccAddress(),
	// 			Standard:           sdk.NewUint(0),
	// 			CollectionMetadata: "https://example.com",
	// 			BadgeMetadata: []*types.BadgeMetadata{
	// 				{
	// 					Uri: "https://example.com/{id}",
	// 					BadgeIds: []*types.IdRange{
	// 						{
	// 							Start: sdk.NewUint(1),
	// 							End:   sdk.NewUint(math.MaxUint64),
	// 						},
	// 					},
	// 				},
	// 			},
	// 			Permissions:       sdk.NewUint(15),
	// 			ApprovedTransfers: []*types.CollectionApprovedTransfer{},
	// 		},
	// 	}, {
	// 		name: "invalid URI",
	// 		msg: types.MsgNewCollection{
	// 			Creator:            sample.AccAddress(),
	// 			Standard:           sdk.NewUint(0),
	// 			CollectionMetadata: "",
	// 			BadgeMetadata: []*types.BadgeMetadata{
	// 				{
	// 					Uri: "https://example.com/{id}",
	// 					BadgeIds: []*types.IdRange{
	// 						{
	// 							Start: sdk.NewUint(1),
	// 							End:   sdk.NewUint(math.MaxUint64),
	// 						},
	// 					},
	// 				},
	// 			},
	// 			Permissions:       sdk.NewUint(15),
	// 			ApprovedTransfers: []*types.CollectionApprovedTransfer{},
	// 		},

	// 		err: types.ErrInvalidBadgeURI,
	// 	},
	// 	{
	// 		name: "invalid Badge URI",
	// 		msg: types.MsgNewCollection{
	// 			Creator:            sample.AccAddress(),
	// 			Standard:           sdk.NewUint(0),
	// 			CollectionMetadata: "https://example.com",
	// 			BadgeMetadata: []*types.BadgeMetadata{
	// 				{
	// 					Uri: "",
	// 					BadgeIds: []*types.IdRange{
	// 						{
	// 							Start: sdk.NewUint(1),
	// 							End:   sdk.NewUint(math.MaxUint64),
	// 						},
	// 					},
	// 				},
	// 			},
	// 			Permissions:       sdk.NewUint(15),
	// 			ApprovedTransfers: []*types.CollectionApprovedTransfer{},
	// 		},
	// 		err: types.ErrInvalidBadgeURI,
	// 	},
	// 	{
	// 		name: "invalid Permissions",
	// 		msg: types.MsgNewCollection{
	// 			Creator:            sample.AccAddress(),
	// 			Standard:           sdk.NewUint(0),
	// 			CollectionMetadata: "https://example.com",
	// 			BadgeMetadata: []*types.BadgeMetadata{
	// 				{
	// 					Uri: "https://example.com/{id}",
	// 					BadgeIds: []*types.IdRange{
	// 						{
	// 							Start: sdk.NewUint(1),
	// 							End:   sdk.NewUint(math.MaxUint64),
	// 						},
	// 					},
	// 				},
	// 			},
	// 			Permissions:       sdk.NewUint(100000),
	// 			ApprovedTransfers: []*types.CollectionApprovedTransfer{},
	// 		},
	// 		err: types.ErrInvalidPermissionsLeadingZeroes,
	// 	},
	// }
	// for _, tt := range tests {
	// 	t.Run(tt.name, func(t *testing.T) {
	// 		err := tt.msg.ValidateBasic()
	// 		if tt.err != nil {
	// 			require.ErrorIs(t, err, tt.err)
	// 			return
	// 		}
	// 		require.NoError(t, err)
	// 	})
	// }
}
