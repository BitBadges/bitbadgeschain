package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMsgUpdateDynamicStore_ValidateBasic(t *testing.T) {
	msg := NewMsgUpdateDynamicStore("bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q", NewUintFromString("1"), false)
	require.NoError(t, msg.ValidateBasic())

	msg = NewMsgUpdateDynamicStore("", NewUintFromString("1"), false)
	err := msg.ValidateBasic()
	require.Error(t, err)

	msg = NewMsgUpdateDynamicStore("bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q", NewUintFromString("0"), false)
	err = msg.ValidateBasic()
	require.Error(t, err)
}
