package keeper

import (
	"context"
	"fmt"
	"strconv"

	"cosmossdk.io/log"
	"github.com/bitbadges/bitbadgeschain/third_party/osmomath"
	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v8/modules/core/24-host"

	"github.com/bitbadges/bitbadgeschain/x/custom-hooks/types"
	poolmanagertypes "github.com/bitbadges/bitbadgeschain/x/poolmanager/types"
)

type (
	// PoolManagerKeeper interface for poolmanager module
	PoolManagerKeeper interface {
		RouteExactAmountIn(ctx sdk.Context, sender sdk.AccAddress, routes []poolmanagertypes.SwapAmountInRoute, tokenIn sdk.Coin, tokenOutMinAmount osmomath.Int) (osmomath.Int, error)
	}

	// BankKeeper interface for bank module
	BankKeeper interface {
		SendCoins(ctx context.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	}

	// ICS4Wrapper interface for sending IBC packets
	ICS4Wrapper interface {
		SendPacket(
			ctx sdk.Context,
			channelCap *capabilitytypes.Capability,
			sourcePort string,
			sourceChannel string,
			timeoutHeight clienttypes.Height,
			timeoutTimestamp uint64,
			data []byte,
		) (uint64, error)
	}

	// ChannelKeeper interface for getting channel capabilities
	ChannelKeeper interface {
		GetChannel(ctx sdk.Context, portID, channelID string) (channeltypes.Channel, bool)
	}

	// ScopedKeeper interface for getting channel capabilities
	ScopedKeeper interface {
		GetCapability(ctx sdk.Context, name string) (*capabilitytypes.Capability, bool)
	}

	Keeper struct {
		logger            log.Logger
		poolManagerKeeper PoolManagerKeeper
		bankKeeper        BankKeeper
		ics4Wrapper       ICS4Wrapper
		channelKeeper     ChannelKeeper
		scopedKeeper      ScopedKeeper
	}
)

func NewKeeper(
	logger log.Logger,
	poolManagerKeeper PoolManagerKeeper,
	bankKeeper BankKeeper,
	ics4Wrapper ICS4Wrapper,
	channelKeeper ChannelKeeper,
	scopedKeeper ScopedKeeper,
) Keeper {
	return Keeper{
		logger:            logger,
		poolManagerKeeper: poolManagerKeeper,
		bankKeeper:        bankKeeper,
		ics4Wrapper:       ics4Wrapper,
		channelKeeper:     channelKeeper,
		scopedKeeper:      scopedKeeper,
	}
}

// Logger returns a logger for the x/custom-hooks module
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

// ExecuteHook executes a hook based on the hook data
func (k Keeper) ExecuteHook(ctx sdk.Context, sender sdk.AccAddress, hookData *types.HookData, tokenIn sdk.Coin) error {
	if hookData == nil {
		return nil
	}

	// Execute SwapAndAction
	if hookData.SwapAndAction != nil {
		return k.ExecuteSwapAndAction(ctx, sender, hookData.SwapAndAction, tokenIn)
	}

	return nil
}

// ExecuteSwapAndAction executes a SwapAndAction hook
func (k Keeper) ExecuteSwapAndAction(ctx sdk.Context, sender sdk.AccAddress, swapAndAction *types.SwapAndAction, tokenIn sdk.Coin) error {
	if swapAndAction == nil {
		return nil
	}

	// Validate: if post_swap_action is defined, a swap must also be defined
	// If they only want to do an IBC transfer without a swap, they should use packet-forward-middleware
	if swapAndAction.PostSwapAction != nil &&
		(swapAndAction.UserSwap == nil || swapAndAction.UserSwap.SwapExactAssetIn == nil) {
		return fmt.Errorf("post_swap_action requires a swap to be defined; use packet-forward-middleware for transfers without swaps")
	}

	// Execute the swap first
	if swapAndAction.UserSwap != nil && swapAndAction.UserSwap.SwapExactAssetIn != nil {
		// Convert operations to routes
		routes, err := k.ConvertOperationsToRoutes(swapAndAction.UserSwap.SwapExactAssetIn.Operations)
		if err != nil {
			return err
		}

		// Get minimum amount from min_asset
		var tokenOut osmomath.Int

		if swapAndAction.MinAsset != nil && swapAndAction.MinAsset.Native != nil {
			tokenOutMinAmount, ok := osmomath.NewIntFromString(swapAndAction.MinAsset.Native.Amount)
			if !ok {
				return fmt.Errorf("invalid min_asset amount: %s", swapAndAction.MinAsset.Native.Amount)
			}
			// Execute swap
			tokenOut, err = k.poolManagerKeeper.RouteExactAmountIn(ctx, sender, routes, tokenIn, tokenOutMinAmount)
			if err != nil {
				return err
			}

			// Update tokenIn to the output token for post-swap action
			tokenIn = sdk.NewCoin(swapAndAction.MinAsset.Native.Denom, tokenOut)
		} else {
			// No min_asset specified, use zero
			tokenOut, err = k.poolManagerKeeper.RouteExactAmountIn(ctx, sender, routes, tokenIn, osmomath.ZeroInt())
			if err != nil {
				return err
			}
			// Update tokenIn to the output token (use last route's output denom)
			if len(routes) > 0 {
				tokenIn = sdk.NewCoin(routes[len(routes)-1].TokenOutDenom, tokenOut)
			}
		}
	}

	// Execute post-swap action (IBC transfer)
	if swapAndAction.PostSwapAction != nil && swapAndAction.PostSwapAction.IBCTransfer != nil {
		// Use timeout_timestamp if provided, otherwise use default
		var timeoutTimestamp uint64
		if swapAndAction.TimeoutTimestamp != nil {
			timeoutTimestamp = *swapAndAction.TimeoutTimestamp
		} else {
			// Default: current time + 1 hour
			timeoutTimestamp = uint64(ctx.BlockTime().UnixNano()) + uint64(3600*1e9)
		}

		if err := k.ExecuteIBCTransfer(ctx, sender, swapAndAction.PostSwapAction.IBCTransfer, tokenIn, timeoutTimestamp); err != nil {
			return err
		}
	}

	return nil
}

// ConvertOperationsToRoutes converts operations array to routes format
func (k Keeper) ConvertOperationsToRoutes(operations []types.Operation) ([]poolmanagertypes.SwapAmountInRoute, error) {
	routes := make([]poolmanagertypes.SwapAmountInRoute, len(operations))
	for i, op := range operations {
		poolId, err := strconv.ParseUint(op.Pool, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid pool ID: %s", op.Pool)
		}

		routes[i] = poolmanagertypes.SwapAmountInRoute{
			PoolId:        poolId,
			TokenOutDenom: op.DenomOut,
		}
	}
	return routes, nil
}

// ExecuteIBCTransfer executes an IBC transfer
func (k Keeper) ExecuteIBCTransfer(ctx sdk.Context, sender sdk.AccAddress, ibcTransfer *types.IBCTransferInfo, token sdk.Coin, timeoutTimestamp uint64) error {
	if ibcTransfer == nil || ibcTransfer.IBCInfo == nil {
		return nil
	}

	ibcInfo := ibcTransfer.IBCInfo

	// Get channel capability
	capPath := host.ChannelCapabilityPath(transfertypes.PortID, ibcInfo.SourceChannel)
	channelCap, ok := k.scopedKeeper.GetCapability(ctx, capPath)
	if !ok {
		return fmt.Errorf("channel capability not found for channel %s", ibcInfo.SourceChannel)
	}

	// Use zero height for timeout (no timeout height)
	timeoutHeight := clienttypes.ZeroHeight()

	// Create transfer packet data
	packetData := transfertypes.FungibleTokenPacketData{
		Denom:    token.Denom,
		Amount:   token.Amount.String(),
		Sender:   sender.String(),
		Receiver: ibcInfo.Receiver,
		Memo:     ibcInfo.Memo,
	}

	// Marshal packet data
	data, err := transfertypes.ModuleCdc.MarshalJSON(&packetData)
	if err != nil {
		return err
	}

	// Send IBC packet
	_, err = k.ics4Wrapper.SendPacket(
		ctx,
		channelCap,
		transfertypes.PortID,
		ibcInfo.SourceChannel,
		timeoutHeight,
		timeoutTimestamp,
		data,
	)
	if err != nil {
		return err
	}

	return nil
}
