package keeper

import (
	"fmt"

	"cosmossdk.io/log"
	sdkmath "cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/ibc-rate-limit/types"
)

type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey storetypes.StoreKey

	bankKeeper types.BankKeeper
	authority  string
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	bankKeeper types.BankKeeper,
	authority string,
) Keeper {
	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address: %s", authority))
	}

	return Keeper{
		cdc:        cdc,
		storeKey:   storeKey,
		bankKeeper: bankKeeper,
		authority:  authority,
	}
}

// Logger returns a module-specific logger
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

// GetAuthority returns the module's authority
func (k Keeper) GetAuthority() string {
	return k.authority
}

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ParamsKey)
	if bz == nil {
		return types.DefaultParams()
	}

	var params types.Params
	k.cdc.MustUnmarshal(bz, &params)
	return params
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) error {
	if err := params.Validate(); err != nil {
		return err
	}

	store := ctx.KVStore(k.storeKey)
	bz, err := k.cdc.Marshal(&params)
	if err != nil {
		return err
	}
	store.Set(types.ParamsKey, bz)

	return nil
}

// GetChannelFlow gets the current flow state for a channel and denom (backward compatibility - uses legacy key)
func (k Keeper) GetChannelFlow(ctx sdk.Context, channelID, denom string) (types.ChannelFlow, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.ChannelFlowKeyLegacy(channelID, denom)
	bz := store.Get(key)
	if bz == nil {
		return types.ChannelFlow{NetFlow: sdkmath.ZeroInt()}, false
	}

	var flow types.ChannelFlow
	k.cdc.MustUnmarshal(bz, &flow)
	return flow, true
}

// SetChannelFlow sets the flow state for a channel and denom (backward compatibility - uses legacy key)
func (k Keeper) SetChannelFlow(ctx sdk.Context, channelID, denom string, flow types.ChannelFlow) {
	store := ctx.KVStore(k.storeKey)
	key := types.ChannelFlowKeyLegacy(channelID, denom)
	bz := k.cdc.MustMarshal(&flow)
	store.Set(key, bz)
}

// GetChannelFlowWindow gets the time window for a channel and denom (backward compatibility - uses legacy key)
func (k Keeper) GetChannelFlowWindow(ctx sdk.Context, channelID, denom string) (types.ChannelFlowWindow, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.ChannelFlowWindowKeyLegacy(channelID, denom)
	bz := store.Get(key)
	if bz == nil {
		return types.ChannelFlowWindow{}, false
	}

	var window types.ChannelFlowWindow
	k.cdc.MustUnmarshal(bz, &window)
	return window, true
}

// SetChannelFlowWindow sets the time window for a channel and denom (backward compatibility - uses legacy key)
func (k Keeper) SetChannelFlowWindow(ctx sdk.Context, channelID, denom string, window types.ChannelFlowWindow) {
	store := ctx.KVStore(k.storeKey)
	key := types.ChannelFlowWindowKeyLegacy(channelID, denom)
	bz := k.cdc.MustMarshal(&window)
	store.Set(key, bz)
}

// ResetChannelFlowWindow resets the flow window for a channel and denom if it has expired
func (k Keeper) ResetChannelFlowWindow(ctx sdk.Context, channelID, denom string, windowDuration int64) {
	currentHeight := ctx.BlockHeight()
	window, found := k.GetChannelFlowWindow(ctx, channelID, denom)

	// If no window exists or window has expired, create/reset it
	if !found || currentHeight >= window.WindowStart+window.WindowDuration {
		newWindow := types.ChannelFlowWindow{
			WindowStart:    currentHeight,
			WindowDuration: windowDuration,
		}
		k.SetChannelFlowWindow(ctx, channelID, denom, newWindow)
		// Reset flow to zero
		k.SetChannelFlow(ctx, channelID, denom, types.ChannelFlow{NetFlow: sdkmath.ZeroInt()})
	}
}

// UpdateChannelFlow updates the net flow for a channel and denom
// positive amount = inflow, negative amount = outflow
func (k Keeper) UpdateChannelFlow(ctx sdk.Context, channelID, denom string, amount sdkmath.Int) {
	params := k.GetParams(ctx)
	config := params.FindMatchingConfig(channelID, denom)

	// If no config, don't track flow (no rate limit)
	if config == nil {
		return
	}

	// Reset window if needed
	k.ResetChannelFlowWindow(ctx, channelID, denom, config.WindowDuration)

	flow, _ := k.GetChannelFlow(ctx, channelID, denom)
	flow.NetFlow = flow.NetFlow.Add(amount)
	k.SetChannelFlow(ctx, channelID, denom, flow)
}

// CheckRateLimit checks if a transfer would exceed the rate limit
// Returns error if rate limit would be exceeded
// senderAddr is the address of the sender (empty string if not available)
func (k Keeper) CheckRateLimit(ctx sdk.Context, channelID string, denom string, amount sdkmath.Int, isInflow bool, senderAddr string) error {
	params := k.GetParams(ctx)

	// Find matching config for this channel and denom
	config := params.FindMatchingConfig(channelID, denom)

	// If no config matches, allow the transfer (no rate limit)
	if config == nil {
		return nil
	}

	// Check backward compatibility: if old fields are set, use them
	if !config.MaxSupplyShift.IsZero() && config.WindowDuration > 0 {
		// Use deprecated single timeframe limit
		k.ResetChannelFlowWindow(ctx, channelID, denom, config.WindowDuration)
		flow, _ := k.GetChannelFlow(ctx, channelID, denom)
		newFlow := flow.NetFlow
		if isInflow {
			newFlow = newFlow.Add(amount)
		} else {
			newFlow = newFlow.Sub(amount)
		}
		absNewFlow := newFlow.Abs()
		if absNewFlow.GT(config.MaxSupplyShift) {
			return types.ErrRateLimitExceeded
		}
	}

	// Check multiple timeframe supply shift limits
	for _, limit := range config.SupplyShiftLimits {
		if limit.MaxAmount.IsZero() {
			continue // Limit disabled
		}

		// Reset window if needed
		k.ResetChannelFlowWindowWithTimeframe(ctx, channelID, denom, limit.TimeframeType, limit.TimeframeDuration)

		// Get current flow
		flow, _ := k.GetChannelFlowWithTimeframe(ctx, channelID, denom, limit.TimeframeType, limit.TimeframeDuration)

		// Calculate new flow after this transfer
		newFlow := flow.NetFlow
		if isInflow {
			newFlow = newFlow.Add(amount)
		} else {
			newFlow = newFlow.Sub(amount)
		}

		// Check if the absolute value of net flow exceeds the limit
		absNewFlow := newFlow.Abs()
		if absNewFlow.GT(limit.MaxAmount) {
			return types.ErrRateLimitExceeded
		}
	}

	// Check unique sender limits (only for inflows with sender address)
	if isInflow && senderAddr != "" {
		for _, limit := range config.UniqueSenderLimits {
			if limit.MaxUniqueSenders == 0 {
				continue // Limit disabled
			}

			// Reset window if needed
			k.ResetUniqueSendersWindow(ctx, channelID, limit.TimeframeType, limit.TimeframeDuration)

			// Get current unique senders
			senders, _ := k.GetUniqueSenders(ctx, channelID, limit.TimeframeType, limit.TimeframeDuration)

			// Check if sender is already in the list
			senderExists := false
			for _, addr := range senders.Senders {
				if addr == senderAddr {
					senderExists = true
					break
				}
			}

			// If sender doesn't exist, check if adding them would exceed limit
			if !senderExists {
				if int64(len(senders.Senders)) >= limit.MaxUniqueSenders {
					return types.ErrRateLimitExceeded
				}
			}
		}
	}

	// Check per-address limits (only if sender address is provided)
	if senderAddr != "" {
		for _, limit := range config.AddressLimits {
			// Reset window if needed
			k.ResetAddressTransferWindow(ctx, senderAddr, channelID, denom, limit.TimeframeType, limit.TimeframeDuration)

			// Get current transfer data
			data, _ := k.GetAddressTransferData(ctx, senderAddr, channelID, denom, limit.TimeframeType, limit.TimeframeDuration)

			// Check transfer count limit
			if limit.MaxTransfers > 0 {
				newCount := data.TransferCount + 1
				if newCount > limit.MaxTransfers {
					return types.ErrRateLimitExceeded
				}
			}

			// Check amount limit
			if !limit.MaxAmount.IsZero() {
				newAmount := data.TotalAmount.Add(amount)
				if newAmount.GT(limit.MaxAmount) {
					return types.ErrRateLimitExceeded
				}
			}
		}
	}

	return nil
}
