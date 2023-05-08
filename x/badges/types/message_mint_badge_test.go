package types_test

import (
	"testing"

	"github.com/bitbadges/bitbadgeschain/testutil/sample"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
)

func TestMsgMintAndDistributeBadges_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgMintAndDistributeBadges
		err  error
	}{
		{
			name: "invalid address",
			msg: types.MsgMintAndDistributeBadges{
				Creator: "invalid_address",
				CollectionId: sdk.NewUint(1),
				BadgeSupplys: []*types.BadgeSupplyAndAmount{
					{
						Supply: sdk.NewUint(10),
						Amount: sdk.NewUint(1),
					},
				},
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid state",
			msg: types.MsgMintAndDistributeBadges{
				Creator: sample.AccAddress(),
				CollectionId: sdk.NewUint(1),
				BadgeSupplys: []*types.BadgeSupplyAndAmount{
					{
						Supply: sdk.NewUint(10),
						Amount: sdk.NewUint(1),
					},
				},
			},
		}, {
			name: "invalid amount",
			msg: types.MsgMintAndDistributeBadges{
				Creator: sample.AccAddress(),
				CollectionId: sdk.NewUint(1),
				BadgeSupplys: []*types.BadgeSupplyAndAmount{
					{
						Supply: sdk.NewUint(10),
						Amount: sdk.NewUint(0),
					},
				},
			},
			err: types.ErrElementCantEqualThis,
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
