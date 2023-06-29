package types

import (
	"testing"

	"github.com/bitbadges/bitbadgeschain/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgForkCollection_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgForkCollection
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgForkCollection{
				Creator: "invalid_address",
			},
			err: ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgForkCollection{
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
