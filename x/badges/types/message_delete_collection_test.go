package types

import (
	"testing"

	"bitbadgeschain/x/badges/testutil/sample"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"
)

func TestMsgDeleteCollection_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgDeleteCollection
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgDeleteCollection{
				Creator:      "invalid_address",
				CollectionId: sdkmath.NewUint(1),
			},
			err: ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgDeleteCollection{
				Creator:      sample.AccAddress(),
				CollectionId: sdkmath.NewUint(1),
			},
		}, {
			name: "invalid collection id",
			msg: MsgDeleteCollection{
				Creator:      sample.AccAddress(),
				CollectionId: sdkmath.NewUint(0),
			},
			err: ErrInvalidCollectionID,
		},
		{
			name: "invalid collection id 2",
			msg: MsgDeleteCollection{
				Creator: sample.AccAddress(),
			},
			err: ErrInvalidCollectionID,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.Error(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
