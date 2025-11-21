package keeper

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"cosmossdk.io/log"
	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/third_party/osmomath"
	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v8/modules/core/24-host"

	badgestypes "github.com/bitbadges/bitbadgeschain/x/badges/types"
	"github.com/bitbadges/bitbadgeschain/x/custom-hooks/types"
	gammtypes "github.com/bitbadges/bitbadgeschain/x/gamm/types"
	poolmanagertypes "github.com/bitbadges/bitbadgeschain/x/poolmanager/types"
)

type (
	// GammKeeper interface for gamm module
	// Note: This should match the gamm.Keeper methods
	GammKeeper interface {
		GetCFMMPool(ctx sdk.Context, poolId uint64) (gammtypes.CFMMPoolI, error)
		SwapExactAmountIn(ctx sdk.Context, sender sdk.AccAddress, pool poolmanagertypes.PoolI, tokenIn sdk.Coin, tokenOutDenom string, tokenOutMinAmount osmomath.Int, spreadFactor osmomath.Dec, affiliates []poolmanagertypes.Affiliate) (osmomath.Int, error)
		RouteExactAmountIn(ctx sdk.Context, sender sdk.AccAddress, routes []poolmanagertypes.SwapAmountInRoute, tokenIn sdk.Coin, tokenOutMinAmount osmomath.Int, affiliates []poolmanagertypes.Affiliate) (osmomath.Int, error)
		CheckIsWrappedDenom(ctx sdk.Context, denom string) bool
		SendNativeTokensFromPool(ctx sdk.Context, poolAddress string, recipientAddress string, denom string, amount sdkmath.Uint) error
	}

	// BankKeeper interface for bank module
	BankKeeper interface {
		SendCoins(ctx context.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	}

	// BadgesKeeper interface for badges module operations
	BadgesKeeper interface {
		GetBalanceOrApplyDefault(ctx sdk.Context, collection *badgestypes.TokenCollection, address string) (*badgestypes.UserBalanceStore, bool)
		SetBalanceForAddress(ctx sdk.Context, collection *badgestypes.TokenCollection, address string, balance *badgestypes.UserBalanceStore) error
		GetCollectionFromStore(ctx sdk.Context, collectionId sdkmath.Uint) (*badgestypes.TokenCollection, bool)
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
		badgesKeeper  BadgesKeeper
		ics4Wrapper   ICS4Wrapper
		channelKeeper ChannelKeeper
		scopedKeeper  ScopedKeeper
	}
)

func NewKeeper(
	logger log.Logger,
	gammKeeper GammKeeper,
	bankKeeper BankKeeper,
	badgesKeeper BadgesKeeper,
	ics4Wrapper ICS4Wrapper,
	channelKeeper ChannelKeeper,
	scopedKeeper ScopedKeeper,
) Keeper {
	return Keeper{
		logger:        logger,
		gammKeeper:    gammKeeper,
		bankKeeper:    bankKeeper,
		badgesKeeper:  badgesKeeper,
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

	// Validate destination recover address if provided
	if swapAndAction.DestinationRecoverAddress != "" {
		_, err := sdk.AccAddressFromBech32(swapAndAction.DestinationRecoverAddress)
		if err != nil {
			return fmt.Errorf("invalid destination_recover_address: %w", err)
		}
	}

	// Validate: min_asset is always required for swaps
	if swapAndAction.MinAsset == nil || swapAndAction.MinAsset.Native == nil {
		return fmt.Errorf("min_asset is required for swaps")
	}

	// Store original tokenIn for fallback handling
	originalTokenIn := sdk.Coin{
		Denom:  tokenIn.Denom,
		Amount: tokenIn.Amount,
	}

	// Execute the swap first in a cached context
	// This ensures we can rollback the swap if it fails and we use the fallback
	if swapAndAction.UserSwap != nil && swapAndAction.UserSwap.SwapExactAssetIn != nil {
		operations := swapAndAction.UserSwap.SwapExactAssetIn.Operations
		if len(operations) == 0 {
			return fmt.Errorf("no operations provided for swap")
		}

		// Convert operations to routes and validate pool IDs first
		routes := make([]poolmanagertypes.SwapAmountInRoute, len(operations))
		for i, operation := range operations {
			poolId, err := strconv.ParseUint(operation.Pool, 10, 64)
			if err != nil {
				return fmt.Errorf("invalid pool ID in operation %d: %s", i, operation.Pool)
			}
			routes[i] = poolmanagertypes.SwapAmountInRoute{
				PoolId:        poolId,
				TokenOutDenom: operation.DenomOut,
			}
		}

		// Validate first operation matches tokenIn
		if operations[0].DenomIn != tokenIn.Denom {
			return fmt.Errorf("first operation denom_in does not match token_in: %s != %s", operations[0].DenomIn, tokenIn.Denom)
		}

		// Validate operations chain correctly (each DenomOut matches next DenomIn)
		for i := 0; i < len(operations)-1; i++ {
			if operations[i].DenomOut != operations[i+1].DenomIn {
				return fmt.Errorf("operations do not chain correctly: operation %d denom_out %s does not match operation %d denom_in %s", i, operations[i].DenomOut, i+1, operations[i+1].DenomIn)
			}
		}

		// Get minimum amount from min_asset for swap validation
		var tokenOut osmomath.Int
		tokenOutMinAmount := osmomath.ZeroInt()

		if swapAndAction.MinAsset != nil && swapAndAction.MinAsset.Native != nil {
			var ok bool
			tokenOutMinAmount, ok = osmomath.NewIntFromString(swapAndAction.MinAsset.Native.Amount)
			if !ok {
				return fmt.Errorf("invalid min_asset amount: %s", swapAndAction.MinAsset.Native.Amount)
			}
		}

		// Validate and convert custom hooks Affiliate types to poolmanager Affiliate types
		var poolmanagerAffiliates []poolmanagertypes.Affiliate
		if len(swapAndAction.Affiliates) > 0 {
			// Validate affiliates before converting
			totalBasisPoints := uint64(0)
			for i, affiliate := range swapAndAction.Affiliates {
				// Validate address
				if affiliate.Address == "" {
					return fmt.Errorf("affiliate_%d.address is required", i)
				}
				_, err := sdk.AccAddressFromBech32(affiliate.Address)
				if err != nil {
					return fmt.Errorf("invalid affiliate_%d.address: %w", i, err)
				}

				// Validate basis_points_fee
				if affiliate.BasisPointsFee == "" {
					return fmt.Errorf("affiliate_%d.basis_points_fee is required", i)
				}
				basisPoints, err := strconv.ParseUint(affiliate.BasisPointsFee, 10, 64)
				if err != nil {
					return fmt.Errorf("invalid affiliate_%d.basis_points_fee: %w", i, err)
				}

				// Validate individual basis points don't exceed 10000
				if basisPoints > 10000 {
					return fmt.Errorf("affiliate_%d.basis_points_fee cannot exceed 10000", i)
				}

				totalBasisPoints += basisPoints
			}

			// Validate total basis points don't exceed 10000
			if totalBasisPoints > 10000 {
				return fmt.Errorf("total affiliate basis_points_fee cannot exceed 10000")
			}

			// Convert to poolmanager types
			poolmanagerAffiliates = make([]poolmanagertypes.Affiliate, len(swapAndAction.Affiliates))
			for i, affiliate := range swapAndAction.Affiliates {
				poolmanagerAffiliates[i] = poolmanagertypes.Affiliate{
					BasisPointsFee: affiliate.BasisPointsFee,
					Address:        affiliate.Address,
				}
			}
		}

		// Set up auto-approve for intermediate address if dealing with wrapped denoms
		// This must be done before the swap call to ensure tokens can be received
		// Check all denoms in the route (input and all intermediate/output denoms)
		allDenoms := []string{tokenIn.Denom}
		for _, operation := range operations {
			allDenoms = append(allDenoms, operation.DenomOut)
		}
		for _, denom := range allDenoms {
			if k.gammKeeper.CheckIsWrappedDenom(ctx, denom) {
				if err := k.setAutoApproveForIntermediateAddress(ctx, sender.String(), denom); err != nil {
					return fmt.Errorf("failed to set auto-approve for intermediate address (denom %s): %w", denom, err)
				}
			}
		}

		// Execute swap in a cached context so we can rollback if it fails
		swapCacheCtx, writeSwapCache := ctx.CacheContext()

		// Always use RouteExactAmountIn (works for both single-hop and multi-hop swaps)
		tokenOut, err := k.gammKeeper.RouteExactAmountIn(
			swapCacheCtx,
			sender,
			routes,
			tokenIn,
			tokenOutMinAmount,
			poolmanagerAffiliates,
		)
		if err != nil {
			// Swap failed - check if we have a destination recover address
			if swapAndAction.DestinationRecoverAddress != "" {
				// Cache context is automatically discarded - swap never happened
				// Log the fallback action
				k.Logger(ctx).Info("custom-hooks: swap failed, using destination recover address",
					"error", err,
					"destination_recover_address", swapAndAction.DestinationRecoverAddress,
					"original_token", originalTokenIn.String())

				// Send original token to destination recover address in the main context
				// This allows the IBC transfer to succeed without the swap
				if err := k.ExecuteLocalTransfer(ctx, sender, &types.TransferInfo{ToAddress: swapAndAction.DestinationRecoverAddress}, originalTokenIn); err != nil {
					return fmt.Errorf("failed to send tokens to destination recover address: %w", err)
				}

				// Emit event for fallback path
				event := sdk.NewEvent(
					"swap_and_action_fallback",
					sdk.NewAttribute("module", "custom-hooks"),
					sdk.NewAttribute("sender", sender.String()),
					sdk.NewAttribute("swap_error", err.Error()),
					sdk.NewAttribute("destination_recover_address", swapAndAction.DestinationRecoverAddress),
					sdk.NewAttribute("original_token", originalTokenIn.String()),
					sdk.NewAttribute("num_operations", strconv.Itoa(len(operations))),
				)
				// Add operation details for debugging
				for i, op := range operations {
					event = event.AppendAttributes(
						sdk.NewAttribute(fmt.Sprintf("operation_%d_pool", i), op.Pool),
						sdk.NewAttribute(fmt.Sprintf("operation_%d_denom_in", i), op.DenomIn),
						sdk.NewAttribute(fmt.Sprintf("operation_%d_denom_out", i), op.DenomOut),
					)
				}
				ctx.EventManager().EmitEvent(event)

				// Successfully handled fallback - return nil to allow IBC transfer to succeed
				// The swap cache context is automatically discarded, so no swap state changes
				return nil
			}

			// No fallback address - cache context is automatically discarded
			// Return error as before
			return fmt.Errorf("swap failed: %w", err)
		}

		// Swap succeeded - commit the swap cache context
		writeSwapCache()

		// Emit event for successful swap path
		firstOperation := operations[0]
		lastOperation := operations[len(operations)-1]
		event := sdk.NewEvent(
			"swap_and_action_success",
			sdk.NewAttribute("module", "custom-hooks"),
			sdk.NewAttribute("sender", sender.String()),
			sdk.NewAttribute("token_in", tokenIn.String()),
			sdk.NewAttribute("token_out", sdk.NewCoin(lastOperation.DenomOut, tokenOut).String()),
			sdk.NewAttribute("denom_in", firstOperation.DenomIn),
			sdk.NewAttribute("denom_out", lastOperation.DenomOut),
			sdk.NewAttribute("num_operations", strconv.Itoa(len(operations))),
		)
		// Add operation details for debugging
		for i, op := range operations {
			event = event.AppendAttributes(
				sdk.NewAttribute(fmt.Sprintf("operation_%d_pool", i), op.Pool),
				sdk.NewAttribute(fmt.Sprintf("operation_%d_denom_in", i), op.DenomIn),
				sdk.NewAttribute(fmt.Sprintf("operation_%d_denom_out", i), op.DenomOut),
			)
		}
		ctx.EventManager().EmitEvent(event)

		// For post-swap, we need to use the exact amount (affiliates already processed)
		tokenIn = sdk.NewCoin(lastOperation.DenomOut, tokenOut)
	}

	// Execute post-swap action (IBC transfer or local transfer)
	if swapAndAction.PostSwapAction.IBCTransfer != nil {
		// Check if attempting to IBC transfer a wrapped denomination
		if k.gammKeeper.CheckIsWrappedDenom(ctx, tokenIn.Denom) {
			return fmt.Errorf("cannot IBC transfer BitBadges denominations: %s", tokenIn.Denom)
		}

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

// IMPORTANT: Should ONLY be called when from address is an intermediate address
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

// setAutoApproveForIntermediateAddress sets auto-approve flags for an intermediate address
// This is similar to how it's done for pool addresses in gamm keeper
func (k Keeper) setAutoApproveForIntermediateAddress(ctx sdk.Context, intermediateAddress string, denom string) error {
	// Parse collection ID from denom (format: badges:COLL_ID:* or badgeslp:COLL_ID:*)
	parts := strings.Split(denom, ":")
	if len(parts) < 3 {
		return fmt.Errorf("invalid denom format: %s", denom)
	}

	collectionIdStr := parts[1]
	collectionId, err := strconv.ParseUint(collectionIdStr, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid collection ID in denom: %w", err)
	}

	// Get collection
	collection, found := k.badgesKeeper.GetCollectionFromStore(ctx, sdkmath.NewUint(collectionId))
	if !found {
		return fmt.Errorf("collection %s not found", collectionIdStr)
	}

	// Get current balance or apply default
	currBalances, _ := k.badgesKeeper.GetBalanceOrApplyDefault(ctx, collection, intermediateAddress)

	// Check if already auto-approved (all flags)
	alreadyAutoApprovedAllIncomingTransfers := currBalances.AutoApproveAllIncomingTransfers
	alreadyAutoApprovedSelfInitiatedOutgoingTransfers := currBalances.AutoApproveSelfInitiatedOutgoingTransfers
	alreadyAutoApprovedSelfInitiatedIncomingTransfers := currBalances.AutoApproveSelfInitiatedIncomingTransfers

	autoApprovedAll := alreadyAutoApprovedAllIncomingTransfers && alreadyAutoApprovedSelfInitiatedOutgoingTransfers && alreadyAutoApprovedSelfInitiatedIncomingTransfers

	if !autoApprovedAll {
		// Set all auto-approve flags to true
		// Incoming - All, no matter what
		// Outgoing - Self-initiated
		// Incoming - Self-initiated
		currBalances.AutoApproveAllIncomingTransfers = true
		currBalances.AutoApproveSelfInitiatedOutgoingTransfers = true
		currBalances.AutoApproveSelfInitiatedIncomingTransfers = true

		// Save the balance
		err = k.badgesKeeper.SetBalanceForAddress(ctx, collection, intermediateAddress, currBalances)
		if err != nil {
			return fmt.Errorf("failed to set auto-approve for intermediate address: %w", err)
		}
	}

	return nil
}
