package types_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/testutil/sample"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	"github.com/stretchr/testify/require"
)

func TestMsgMintAndDistributeBadges_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgUniversalUpdateCollection
		err  error
	}{
		{
			name: "invalid address",
			msg: types.MsgUniversalUpdateCollection{
				Creator:      "invalid_address",
				CollectionId: sdkmath.NewUint(1),
			},
			err: types.ErrInvalidAddress,
		}, {
			name: "valid state",
			msg: types.MsgUniversalUpdateCollection{
				Creator:        sample.AccAddress(),
				CollectionId:   sdkmath.NewUint(1),
				BadgesToCreate: []*types.Balance{},
			},
		}, {
			name: "invalid amount",
			msg: types.MsgUniversalUpdateCollection{
				Creator:      sample.AccAddress(),
				CollectionId: sdkmath.NewUint(1),
				BadgesToCreate: []*types.Balance{
					{
						Amount: sdkmath.NewUint(0),
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
				require.Error(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
