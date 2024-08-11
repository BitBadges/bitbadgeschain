package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgCreateProposal_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgCreateProposal
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgCreateProposal{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		// {
		// 	name: "valid address",
		// 	msg: MsgCreateProposal{
		// 		Creator: sample.AccAddress(),
		// 		Parties: []*Parties{
		// 			{
		// 				Creator: sample.AccAddress(),
		// 				// MsgsToExecute: GetM
		// 			},
		// 		},
		// 	},
		// },
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
