package types

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"
)

func TestMsgUpdateDynamicStore_ValidateBasic(t *testing.T) {
	msg := NewMsgUpdateDynamicStore("bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q", sdkmath.NewUint(1), false)
	require.NoError(t, msg.ValidateBasic())

	msg = NewMsgUpdateDynamicStore("", sdkmath.NewUint(1), false)
	err := msg.ValidateBasic()
	require.Error(t, err)

	msg = NewMsgUpdateDynamicStore("bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q", sdkmath.NewUint(0), false)
	err = msg.ValidateBasic()
	require.Error(t, err)
}
