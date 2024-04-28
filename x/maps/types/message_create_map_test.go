package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMsgCreateMap_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgCreateMap
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgCreateMap{
				Creator: "invalid_address",
			},
			err: ErrInvalidAddress,
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
