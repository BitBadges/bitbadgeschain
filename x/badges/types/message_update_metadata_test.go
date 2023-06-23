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

func TestMsgUpdateMetadata_ValidateBasic(t *testing.T) {
	// tests := []struct {
	// 	name string
	// 	msg  types.MsgUpdateMetadata
	// 	err  error
	// }{
	// 	{
	// 		name: "invalid address",
	// 		msg: types.MsgUpdateMetadata{
	// 			Creator:            "invalid_address",
	// 			CollectionId:       sdk.NewUint(1),
	// 			CollectionMetadata: "https://facebook.com",
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
	// 		},
	// 		err: sdkerrors.ErrInvalidAddress,
	// 	}, {
	// 		name: "valid address",
	// 		msg: types.MsgUpdateMetadata{
	// 			Creator:            sample.AccAddress(),
	// 			CollectionMetadata: "https://facebook.com",
	// 			CollectionId:       sdk.NewUint(1),
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
	// 		},
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
