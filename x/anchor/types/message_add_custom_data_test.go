package types

import (
	"testing"

	"bitbadgeschain/testutil/sample"

	"github.com/stretchr/testify/require"
)

func TestMsgAddCustomData_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgAddCustomData
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgAddCustomData{
				Creator: "invalid_address",
			},
			err: ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgAddCustomData{
				Creator: sample.AccAddress(),
				Data:    "test",
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
