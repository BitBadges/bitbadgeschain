package types

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"
)

func TestMsgDeleteDynamicStore_ValidateBasic(t *testing.T) {
	msg := NewMsgDeleteDynamicStore("bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q", sdkmath.NewUint(1))
	require.NoError(t, msg.ValidateBasic())

	msg = NewMsgDeleteDynamicStore("", sdkmath.NewUint(1))
	err := msg.ValidateBasic()
	require.Error(t, err)

	msg = NewMsgDeleteDynamicStore("bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q", sdkmath.NewUint(0))
	err = msg.ValidateBasic()
	require.Error(t, err)
}
