package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMsgCreateDynamicStore_ValidateBasic(t *testing.T) {
	msg := NewMsgCreateDynamicStore("bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q", false)
	require.NoError(t, msg.ValidateBasic())

	msg = NewMsgCreateDynamicStore("", false)
	err := msg.ValidateBasic()
	require.Error(t, err)
}
