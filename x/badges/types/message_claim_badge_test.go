package types

import (
	"testing"

	"github.com/bitbadges/bitbadgeschain/testutil/sample"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgClaimBadge_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgClaimBadge
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgClaimBadge{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgClaimBadge{
				Creator: sample.AccAddress(),
				Leaf: []byte("hello"),
				Proof: &Proof{
					LeafHash: []byte("hello"),
					Aunts: [][]byte{[]byte("hello")},
					Total: 1,
					Index: 0,
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
