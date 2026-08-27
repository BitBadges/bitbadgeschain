package keeper

import (
	"bytes"
	"fmt"
	"strings"

	sdkmath "cosmossdk.io/math"
	storetypes "github.com/cosmos/cosmos-sdk/store/v2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/ibc-rate-limit/types"
)

// denomSeparator is the byte the key builders in types/keys.go join components
// with. Neither an IBC channel id ("channel-N") nor a bech32 address can
// contain it, so splitting on it recovers the components unambiguously.
const denomSeparator = "|"

// RedenominateFlows moves this module's in-flight accounting from one denom to
// another, scaling the amounts it carries.
//
// Four record families are keyed by denom:
//
//	KeyPrefixChannelFlow           channelID|denom[|tfType|tfDuration] -> ChannelFlow
//	KeyPrefixChannelFlowWindow     channelID|denom[|tfType|tfDuration] -> ChannelFlowWindow
//	KeyPrefixAddressTransferData   address|channelID|denom|tfType|tfDuration -> AddressTransferData
//	KeyPrefixAddressTransferWindow address|channelID|denom|tfType|tfDuration -> ChannelFlowWindow
//
// The two flow families also carry an amount in the old scale — ChannelFlow's
// NetFlow and AddressTransferData's TotalAmount — so they need both a re-key
// and a rescale. The two window families carry only timestamps, so they need
// the re-key alone; without it the window looks unopened and the accumulated
// flow it belongs to would be discarded on the next transfer.
//
// Leaving these behind is not a value bug — nothing is spendable from them —
// but it silently restarts every rate-limit window at zero at the upgrade
// height, which is the permissive direction and briefly allows up to twice the
// configured cap across the boundary.
//
// A key whose shape is not recognised is logged and skipped rather than
// failing. These records are transient window state that resets on its own
// timer; halting a chain over one is a worse outcome than one reset window.
func (k Keeper) RedenominateFlows(
	ctx sdk.Context,
	legacyDenom, newDenom string,
	factor sdkmath.Int,
) (int, error) {
	store := ctx.KVStore(k.storeKey)

	type rewrite struct {
		oldKey, newKey []byte
		value          []byte
	}
	var pending []rewrite

	// denomIndex says which "|"-separated component after the prefix byte holds
	// the denom, and how many components the key must have.
	families := []struct {
		prefix     []byte
		denomIndex int
		// scale names the amount field to rescale; "" means the record carries
		// no amount.
		scale string
	}{
		{types.KeyPrefixChannelFlow, 1, "channel_flow"},
		{types.KeyPrefixChannelFlowWindow, 1, ""},
		{types.KeyPrefixAddressTransferData, 2, "address_transfer"},
		{types.KeyPrefixAddressTransferWindow, 2, ""},
	}

	for _, fam := range families {
		iter := storetypes.KVStorePrefixIterator(store, fam.prefix)
		for ; iter.Valid(); iter.Next() {
			key := append([]byte(nil), iter.Key()...)
			body := string(bytes.TrimPrefix(key, fam.prefix))
			parts := strings.Split(body, denomSeparator)

			if len(parts) <= fam.denomIndex {
				ctx.Logger().Error("v35: skipping rate-limit record with an unrecognised key shape",
					"prefix", fmt.Sprintf("%x", fam.prefix), "key", body)
				continue
			}
			if parts[fam.denomIndex] != legacyDenom {
				continue
			}

			parts[fam.denomIndex] = newDenom
			newKey := append(append([]byte(nil), fam.prefix...), []byte(strings.Join(parts, denomSeparator))...)

			value := append([]byte(nil), iter.Value()...)
			switch fam.scale {
			case "channel_flow":
				var flow types.ChannelFlow
				if err := k.cdc.Unmarshal(value, &flow); err != nil {
					iter.Close()
					return 0, fmt.Errorf("unmarshalling channel flow %s: %w", body, err)
				}
				flow.NetFlow = flow.NetFlow.Mul(factor)
				scaled, err := k.cdc.Marshal(&flow)
				if err != nil {
					iter.Close()
					return 0, fmt.Errorf("marshalling channel flow %s: %w", body, err)
				}
				value = scaled
			case "address_transfer":
				var data types.AddressTransferData
				if err := k.cdc.Unmarshal(value, &data); err != nil {
					iter.Close()
					return 0, fmt.Errorf("unmarshalling address transfer data %s: %w", body, err)
				}
				data.TotalAmount = data.TotalAmount.Mul(factor)
				scaled, err := k.cdc.Marshal(&data)
				if err != nil {
					iter.Close()
					return 0, fmt.Errorf("marshalling address transfer data %s: %w", body, err)
				}
				value = scaled
			}

			pending = append(pending, rewrite{oldKey: key, newKey: newKey, value: value})
		}
		iter.Close()
	}

	// Written after every iterator is closed: these writes land in the same
	// prefixes being walked.
	for _, r := range pending {
		store.Delete(r.oldKey)
		store.Set(r.newKey, r.value)
	}

	if len(pending) > 0 {
		ctx.Logger().Info("v35: re-keyed ibc rate limit flows",
			"from", legacyDenom, "to", newDenom, "records", len(pending))
	}

	return len(pending), nil
}
