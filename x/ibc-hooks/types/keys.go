package types

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

const (
	ModuleName = "ibchooks"
	RouterKey  = ModuleName
	StoreKey   = "hooks-for-ibc" // not using the module name because of collisions with key "ibc"

	SenderPrefix = "ibc-hook-intermediary"
)

// DeriveIntermediateSender derives the sender address to be used for custom hooks
func DeriveIntermediateSender(channel, originalSender, bech32Prefix string) (string, error) {
	senderStr := fmt.Sprintf("%s/%s", channel, originalSender)
	senderHash32 := address.Hash(SenderPrefix, []byte(senderStr))
	sender := types.AccAddress(senderHash32[:])
	return types.Bech32ifyAddressBytes(bech32Prefix, sender)
}
