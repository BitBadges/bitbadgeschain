package types_test

import (
	"testing"

	"github.com/bitbadges/bitbadgeschain/testutil/sample"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgClaimBadge_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgClaimBadge
		err  error
	}{
		{
			name: "invalid address",
			msg: types.MsgClaimBadge{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: types.MsgClaimBadge{
				Creator: sample.AccAddress(),
				Proof: &types.Proof{
					Leaf: "hello",
					Aunts: []*types.ProofItem{
						{
							Aunt: "hello",
							OnRight: true,
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
