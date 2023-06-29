package types_test

import (
	"testing"
	// "github.com/bitbadges/bitbadgeschain/testutil/sample"
	// sdk "github.com/cosmos/cosmos-sdk/types"
	// sdkerrors "cosmossdk.io/errors"
	// "github.com/stretchr/testify/require"
	// "github.com/bitbadges/bitbadgeschain/x/badges/types"
)

func TestMsgTransferBadge_ValidateBasic(t *testing.T) {
	// tests := []struct {
	// 	name string
	// 	msg  types.MsgTransferBadge
	// 	err  error
	// }{
	// 	{
	// 		name: "invalid address",
	// 		msg: types.MsgTransferBadge{
	// 			Creator:      "invalid_address",
	// 			From:         sample.AccAddress(),
	// 			CollectionId: sdkmath.NewUint(1),
	// 			Transfers: []*types.Transfer{
	// 				{
	// 					ToAddresses: []string{
	// 						sample.AccAddress(),
	// 					},
	// 					Balances: []*types.Balance{
	// 						{
	// 							Amount: sdkmath.NewUint(10),
	// 							BadgeIds: []*types.IdRange{
	// 								{
	// 									Start: sdkmath.NewUint(1),
	// 									End:   sdkmath.NewUint(1),
	// 								},
	// 							},
	// 						},
	// 					},
	// 				},
	// 			},
	// 		},
	// 		err: ErrInvalidAddress,
	// 	}, {
	// 		name: "valid state",
	// 		msg: types.MsgTransferBadge{
	// 			Creator: sample.AccAddress(),
	// 			From:    sample.AccAddress(),
	// 			Transfers: []*types.Transfer{
	// 				{
	// 					ToAddresses: []string{
	// 						sample.AccAddress(),
	// 					},
	// 					Balances: []*types.Balance{
	// 						{
	// 							Amount: sdkmath.NewUint(10),
	// 							BadgeIds: []*types.IdRange{
	// 								{
	// 									Start: sdkmath.NewUint(1),
	// 									End:   sdkmath.NewUint(1),
	// 								},
	// 							},
	// 						},
	// 					},
	// 				},
	// 			},
	// 		},
	// 	}, {
	// 		name: "invalid amounts",
	// 		msg: types.MsgTransferBadge{
	// 			Creator: sample.AccAddress(),
	// 			From:    sample.AccAddress(),
	// 			Transfers: []*types.Transfer{
	// 				{
	// 					ToAddresses: []string{
	// 						sample.AccAddress(),
	// 					},
	// 					Balances: []*types.Balance{
	// 						{
	// 							Amount: sdkmath.NewUint(0),
	// 							BadgeIds: []*types.IdRange{
	// 								{
	// 									Start: sdkmath.NewUint(1),
	// 									End:   sdkmath.NewUint(1),
	// 								},
	// 							},
	// 						},
	// 					},
	// 				},
	// 			},
	// 		},
	// 		err: types.ErrAmountEqualsZero,
	// 	},
	// 	{
	// 		name: "invalid badge range",
	// 		msg: types.MsgTransferBadge{
	// 			Creator: sample.AccAddress(),
	// 			From:    sample.AccAddress(),
	// 			Transfers: []*types.Transfer{
	// 				{
	// 					ToAddresses: []string{sample.AccAddress()},
	// 					Balances: []*types.Balance{
	// 						{
	// 							Amount: sdkmath.NewUint(10),
	// 							BadgeIds: []*types.IdRange{
	// 								{
	// 									Start: sdkmath.NewUint(10),
	// 									End:   sdkmath.NewUint(1),
	// 								},
	// 							},
	// 						},
	// 					},
	// 				},
	// 			},
	// 		},
	// 		err: types.ErrStartGreaterThanEnd,
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
