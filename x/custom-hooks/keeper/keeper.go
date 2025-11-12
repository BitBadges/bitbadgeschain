package keeper

import (
	"context"
	"fmt"
	"strconv"

	"cosmossdk.io/log"
	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/third_party/osmomath"
	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v8/modules/core/24-host"

	"github.com/bitbadges/bitbadgeschain/x/custom-hooks/types"
	gammtypes "github.com/bitbadges/bitbadgeschain/x/gamm/types"
	poolmanagertypes "github.com/bitbadges/bitbadgeschain/x/poolmanager/types"
)

type (
	// GammKeeper interface for gamm module
	// Note: This should match the gamm.Keeper methods
	GammKeeper interface {
		GetCFMMPool(ctx sdk.Context, poolId uint64) (gammtypes.CFMMPoolI, error)
		SwapExactAmountIn(ctx sdk.Context, sender sdk.AccAddress, pool poolmanagertypes.PoolI, tokenIn sdk.Coin, tokenOutDenom string, tokenOutMinAmount osmomath.Int, spreadFactor osmomath.Dec) (osmomath.Int, error)
		CheckIsWrappedDenom(ctx sdk.Context, denom string) bool
		SendNativeTokensFromPool(ctx sdk.Context, poolAddress string, recipientAddress string, denom string, amount sdkmath.Uint) error
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
		logger        log.Logger
		gammKeeper    GammKeeper
		bankKeeper    BankKeeper
		ics4Wrapper   ICS4Wrapper
		channelKeeper ChannelKeeper
		scopedKeeper  ScopedKeeper
	}
)

func NewKeeper(
	logger log.Logger,
	gammKeeper GammKeeper,
	bankKeeper BankKeeper,
	ics4Wrapper ICS4Wrapper,
	channelKeeper ChannelKeeper,
	scopedKeeper ScopedKeeper,
) Keeper {
	return Keeper{
		logger:        logger,
		gammKeeper:    gammKeeper,
		bankKeeper:    bankKeeper,
		ics4Wrapper:   ics4Wrapper,
		channelKeeper: channelKeeper,
		scopedKeeper:  scopedKeeper,
	}
}

// Logger returns a logger for the x/custom-hooks module
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

// BankKeeper returns the bank keeper for token transfers
func (k Keeper) BankKeeper() BankKeeper {
	return k.bankKeeper
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

	// Validate: post_swap_action is required and must have exactly one of IBCTransfer or Transfer
	if swapAndAction.PostSwapAction == nil {
		return fmt.Errorf("post_swap_action is required and must have either ibc_transfer or transfer")
	}

	// Validate that exactly one of IBCTransfer or Transfer is set
	hasIBCTransfer := swapAndAction.PostSwapAction.IBCTransfer != nil
	hasTransfer := swapAndAction.PostSwapAction.Transfer != nil
	if hasIBCTransfer && hasTransfer {
		return fmt.Errorf("post_swap_action cannot have both ibc_transfer and transfer; must specify exactly one")
	}
	if !hasIBCTransfer && !hasTransfer {
		return fmt.Errorf("post_swap_action must have either ibc_transfer or transfer")
	}

	// Validate: if post_swap_action is defined, a swap must also be defined
	// If they only want to do a transfer without a swap, they should use packet-forward-middleware
	if swapAndAction.UserSwap == nil || swapAndAction.UserSwap.SwapExactAssetIn == nil {
		return fmt.Errorf("post_swap_action requires a swap to be defined; use packet-forward-middleware for transfers without swaps")
	}

	// Validate post-swap action format and sanity checks BEFORE executing swap
	// This prevents wasting gas on a swap if the post-swap action is malformed
	if err := k.ValidatePostSwapAction(ctx, swapAndAction.PostSwapAction); err != nil {
		return fmt.Errorf("post_swap_action validation failed: %w", err)
	}

	// Execute the swap first
	if swapAndAction.UserSwap != nil && swapAndAction.UserSwap.SwapExactAssetIn != nil {
		operations := swapAndAction.UserSwap.SwapExactAssetIn.Operations
		if len(operations) == 0 {
			return fmt.Errorf("no operations provided for swap")
		}

		if len(operations) > 1 {
			return fmt.Errorf("multi-hop swaps are not supported")
		}

		// For now, we only support single-hop swaps (one operation)
		// Multi-hop swaps would require chaining multiple swaps
		operation := operations[0]
		poolId, err := strconv.ParseUint(operation.Pool, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid pool ID: %s", operation.Pool)
		}

		if operation.DenomIn != tokenIn.Denom {
			return fmt.Errorf("denom_in does not match token_in: %s != %s", operation.DenomIn, tokenIn.Denom)
		}

		// Get the pool
		pool, err := k.gammKeeper.GetCFMMPool(ctx, poolId)
		if err != nil {
			return fmt.Errorf("failed to get pool %d: %w", poolId, err)
		}

		// Get minimum amount from min_asset
		var tokenOut osmomath.Int
		tokenOutMinAmount := osmomath.ZeroInt()

		if swapAndAction.MinAsset != nil && swapAndAction.MinAsset.Native != nil {
			var ok bool
			tokenOutMinAmount, ok = osmomath.NewIntFromString(swapAndAction.MinAsset.Native.Amount)
			if !ok {
				return fmt.Errorf("invalid min_asset amount: %s", swapAndAction.MinAsset.Native.Amount)
			}
		}

		// Get spread factor from pool
		spreadFactor := pool.GetSpreadFactor(ctx)

		// Execute swap using Gamm keeper directly
		tokenOut, err = k.gammKeeper.SwapExactAmountIn(
			ctx,
			sender,
			pool,
			tokenIn,
			operation.DenomOut,
			tokenOutMinAmount,
			spreadFactor,
		)
		if err != nil {
			return fmt.Errorf("swap failed: %w", err)
		}

		//for post-swap, we need to use the exact amount
		tokenIn = sdk.NewCoin(operation.DenomOut, tokenOut)
	}

	// Execute post-swap action (IBC transfer or local transfer)
	if swapAndAction.PostSwapAction.IBCTransfer != nil {
		// Use timeout_timestamp if provided, otherwise use default
		var timeoutTimestamp uint64
		if swapAndAction.TimeoutTimestamp != nil {
			timeoutTimestamp = *swapAndAction.TimeoutTimestamp
		} else {
			// Default: current time + 5 minutes
			timeoutTimestamp = uint64(ctx.BlockTime().UnixNano()) + uint64(5*60*1e9)
		}

		recoverAddress := swapAndAction.PostSwapAction.IBCTransfer.IBCInfo.RecoverAddress
		recoverAddressAddr, err := sdk.AccAddressFromBech32(recoverAddress)
		if err != nil {
			return fmt.Errorf("invalid recover address: %w", err)
		}

		// Send to recover address then IBC it out
		if err := k.ExecuteLocalTransfer(ctx, sender, &types.TransferInfo{ToAddress: recoverAddress}, tokenIn); err != nil {
			return err
		}

		if err := k.ExecuteIBCTransfer(ctx, recoverAddressAddr, swapAndAction.PostSwapAction.IBCTransfer, tokenIn, timeoutTimestamp); err != nil {
			return err
		}
	} else if swapAndAction.PostSwapAction.Transfer != nil {
		if err := k.ExecuteLocalTransfer(ctx, sender, swapAndAction.PostSwapAction.Transfer, tokenIn); err != nil {
			return err
		}
	}

	return nil
}

// ValidatePostSwapAction performs sanity checks on the post-swap action before executing the swap
// This prevents wasting gas on a swap if the post-swap action is malformed
func (k Keeper) ValidatePostSwapAction(ctx sdk.Context, postSwapAction *types.PostSwapAction) error {
	if postSwapAction == nil {
		return fmt.Errorf("post_swap_action cannot be nil")
	}

	// Validate IBC transfer if present
	if postSwapAction.IBCTransfer != nil {
		if postSwapAction.IBCTransfer.IBCInfo == nil {
			return fmt.Errorf("ibc_transfer.ibc_info is required")
		}

		ibcInfo := postSwapAction.IBCTransfer.IBCInfo

		// Validate source channel is not empty
		if ibcInfo.SourceChannel == "" {
			return fmt.Errorf("ibc_transfer.ibc_info.source_channel is required")
		}

		// Validate channel exists
		_, found := k.channelKeeper.GetChannel(ctx, transfertypes.PortID, ibcInfo.SourceChannel)
		if !found {
			return fmt.Errorf("IBC channel %s does not exist", ibcInfo.SourceChannel)
		}

		// Validate receiver is not empty
		if ibcInfo.Receiver == "" {
			return fmt.Errorf("ibc_transfer.ibc_info.receiver is required")
		}

		// Validate recover address if provided
		if ibcInfo.RecoverAddress != "" {
			_, err := sdk.AccAddressFromBech32(ibcInfo.RecoverAddress)
			if err != nil {
				return fmt.Errorf("invalid ibc_transfer.ibc_info.recover_address: %w", err)
			}
		}

		// Validate channel capability exists (this is what we'll need later)
		capPath := host.ChannelCapabilityPath(transfertypes.PortID, ibcInfo.SourceChannel)
		_, ok := k.scopedKeeper.GetCapability(ctx, capPath)
		if !ok {
			return fmt.Errorf("channel capability not found for channel %s", ibcInfo.SourceChannel)
		}
	}

	// Validate local transfer if present
	if postSwapAction.Transfer != nil {
		if postSwapAction.Transfer.ToAddress == "" {
			return fmt.Errorf("transfer.to_address is required")
		}

		// Validate to_address can be parsed as a valid address
		_, err := sdk.AccAddressFromBech32(postSwapAction.Transfer.ToAddress)
		if err != nil {
			return fmt.Errorf("invalid transfer.to_address: %w", err)
		}
	}

	return nil
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

// ExecuteLocalTransfer executes a local bank transfer
func (k Keeper) ExecuteLocalTransfer(ctx sdk.Context, sender sdk.AccAddress, transfer *types.TransferInfo, token sdk.Coin) error {
	if transfer == nil {
		return fmt.Errorf("transfer info is nil")
	}

	if transfer.ToAddress == "" {
		return fmt.Errorf("to_address is required for local transfer")
	}

	// Parse the recipient address
	toAddr, err := sdk.AccAddressFromBech32(transfer.ToAddress)
	if err != nil {
		return fmt.Errorf("invalid to_address: %w", err)
	}

	// Execute bank transfer
	if err := k.SendCoinsFromIntermediateAddress(ctx, sender, toAddr, sdk.Coins{token}); err != nil {
		return fmt.Errorf("failed to send coins: %w", err)
	}

	return nil
}

// IMPORTANT: Should ONLY be called when from address is a pool address
func (k Keeper) SendCoinsFromIntermediateAddress(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, coins sdk.Coins) error {
	for _, coin := range coins {
		if k.gammKeeper.CheckIsWrappedDenom(ctx, coin.Denom) {
			err := k.gammKeeper.SendNativeTokensFromPool(ctx, from.String(), to.String(), coin.Denom, sdkmath.NewUintFromBigInt(coin.Amount.BigInt()))
			if err != nil {
				return err
			}
		} else {
			err := k.bankKeeper.SendCoins(ctx, from, to, sdk.NewCoins(coin))
			if err != nil {
				return err
			}
		}
	}

	return nil
}
