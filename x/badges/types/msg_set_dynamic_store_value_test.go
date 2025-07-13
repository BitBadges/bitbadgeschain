package types

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"
)

func TestMsgSetDynamicStoreValue_ValidateBasic(t *testing.T) {
	msg := NewMsgSetDynamicStoreValue("bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q", sdkmath.NewUint(1), "bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q", true)
	require.NoError(t, msg.ValidateBasic())

	msg = NewMsgSetDynamicStoreValue("", sdkmath.NewUint(1), "bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q", true)
	err := msg.ValidateBasic()
	require.Error(t, err)

	msg = NewMsgSetDynamicStoreValue("bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q", sdkmath.NewUint(0), "bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q", true)
	err = msg.ValidateBasic()
	require.Error(t, err)

	msg = NewMsgSetDynamicStoreValue("bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q", sdkmath.NewUint(1), "", true)
	err = msg.ValidateBasic()
	require.Error(t, err)

	msg = NewMsgSetDynamicStoreValue("bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q", sdkmath.NewUint(1), "invalid-address", true)
	err = msg.ValidateBasic()
	require.Error(t, err)
}
