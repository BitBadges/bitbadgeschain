package types_test

import (
	"testing"

	// "github.com/bitbadges/bitbadgeschain/testutil/sample"
	// "github.com/bitbadges/bitbadgeschain/x/badges/types"
	// sdk "github.com/cosmos/cosmos-sdk/types"
	// sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	// "github.com/stretchr/testify/require"
)

func TestMsgUpdateCollectionApprovedTransfers_ValidateBasic(t *testing.T) {
	// tests := []struct {
	// 	name string
	// 	msg  types.MsgUpdateCollectionApprovedTransfers
	// 	err  error
	// }{
	// 	{
	// 		name: "invalid address",
	// 		msg: types.MsgUpdateCollectionApprovedTransfers{
	// 			Creator:      "invalid_address",
	// 			CollectionId: sdk.NewUint(1),
	// 			ApprovedTransfers: []*types.CollectionApprovedTransfer{
	// 				{
	// 					From: &types.AddressMapping{
	// 						Addresses: []string{
	// 							"invalid_address",
	// 						},
	// 						ManagerOptions:       sdk.NewUint(uint64(types.AddressOptions_None)),
	// 						IncludeOnlySpecified: true,
	// 					},
	// 					To: &types.AddressMapping{
	// 						Addresses: []string{
	// 							"invalid_address",
	// 						},
	// 						ManagerOptions:       sdk.NewUint(uint64(types.AddressOptions_None)),
	// 						IncludeOnlySpecified: true,
	// 					},
	// 				},
	// 			},
	// 		},
	// 		err: sdkerrors.ErrInvalidAddress,
	// 	}, {
	// 		name: "valid address",
	// 		msg: types.MsgUpdateCollectionApprovedTransfers{
	// 			Creator:      sample.AccAddress(),
	// 			CollectionId: sdk.NewUint(1),
	// 			ApprovedTransfers: []*types.CollectionApprovedTransfer{
	// 				{
	// 					From: &types.AddressMapping{
	// 						Addresses: []string{
	// 							sample.AccAddress(),
	// 						},
	// 						ManagerOptions:       sdk.NewUint(uint64(types.AddressOptions_None)),
	// 						IncludeOnlySpecified: true,
	// 					},
	// 					To: &types.AddressMapping{
	// 						Addresses: []string{
	// 							sample.AccAddress(),
	// 						},
	// 						ManagerOptions:       sdk.NewUint(uint64(types.AddressOptions_None)),
	// 						IncludeOnlySpecified: true,
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
