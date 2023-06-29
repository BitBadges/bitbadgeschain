package types

import (
	"testing"

	"github.com/bitbadges/bitbadgeschain/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgUpdateUserApprovedTransfers_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgUpdateUserApprovedTransfers
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgUpdateUserApprovedTransfers{
				Creator: "invalid_address",
			},
			err: ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgUpdateUserApprovedTransfers{
				Creator: sample.AccAddress(),
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
