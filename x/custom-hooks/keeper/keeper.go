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
	transfertypes "github.com/cosmos/ibc-go/v10/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v10/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v10/modules/core/04-channel/types"
	ibcexported "github.com/cosmos/ibc-go/v10/modules/core/exported"

	tokenizationtypes "github.com/bitbadges/bitbadgeschain/x/tokenization/types"
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
	}

	// BankKeeper interface for bank module
	BankKeeper interface {
		SendCoins(ctx context.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	}

	// TokenizationKeeper interface for tokenization module operations
	TokenizationKeeper interface {
		ParseCollectionFromDenom(ctx sdk.Context, denom string) (*tokenizationtypes.TokenCollection, error)
		SetAllAutoApprovalFlagsForIntermediateAddress(ctx sdk.Context, collection *tokenizationtypes.TokenCollection, address string) error
		GetBalanceOrApplyDefault(ctx sdk.Context, collection *tokenizationtypes.TokenCollection, address string) (*tokenizationtypes.UserBalanceStore, bool, error)
		SetBalanceForAddress(ctx sdk.Context, collection *tokenizationtypes.TokenCollection, address string, balance *tokenizationtypes.UserBalanceStore) error
		GetCollectionFromStore(ctx sdk.Context, collectionId sdkmath.Uint) (*tokenizationtypes.TokenCollection, bool)
		SendNativeTokensViaAliasDenom(ctx sdk.Context, fromAddress string, recipientAddress string, denom string, amount sdkmath.Uint) error
		CheckIsAliasDenom(ctx sdk.Context, denom string) bool
	}

	// SendManagerKeeper interface for sendmanager module operations
	SendManagerKeeper interface {
		SendCoinWithAliasRouting(ctx sdk.Context, fromAddressAcc sdk.AccAddress, toAddressAcc sdk.AccAddress, coin *sdk.Coin) error
		IsICS20Compatible(ctx sdk.Context, denom string) bool
		StandardName(ctx sdk.Context, denom string) string
	}

	// ICS4Wrapper interface for sending IBC packets (IBC v10: capabilities removed)
	ICS4Wrapper interface {
		SendPacket(
			ctx sdk.Context,
			sourcePort string,
			sourceChannel string,
			timeoutHeight clienttypes.Height,
			timeoutTimestamp uint64,
			data []byte,
		) (uint64, error)
	}

	// ChannelKeeper interface for getting channel information
	ChannelKeeper interface {
		GetChannel(ctx sdk.Context, portID, channelID string) (channeltypes.Channel, bool)
	}

	Keeper struct {
		logger            log.Logger
		gammKeeper        GammKeeper
		bankKeeper        BankKeeper
		tokenizationKeeper TokenizationKeeper
		sendManagerKeeper SendManagerKeeper
		transferKeeper    gammtypes.TransferKeeper
		ics4Wrapper       ICS4Wrapper
		channelKeeper     ChannelKeeper
	}
)

func NewKeeper(
	logger log.Logger,
	gammKeeper GammKeeper,
	bankKeeper BankKeeper,
	tokenizationKeeper TokenizationKeeper,
	sendManagerKeeper SendManagerKeeper,
	transferKeeper gammtypes.TransferKeeper,
	ics4Wrapper ICS4Wrapper,
	channelKeeper ChannelKeeper,
) Keeper {
	return Keeper{
		logger:            logger,
		gammKeeper:        gammKeeper,
		bankKeeper:        bankKeeper,
		tokenizationKeeper: tokenizationKeeper,
		sendManagerKeeper: sendManagerKeeper,
		transferKeeper:    transferKeeper,
		ics4Wrapper:       ics4Wrapper,
		channelKeeper:     channelKeeper,
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
func (k Keeper) ExecuteHook(ctx sdk.Context, sender sdk.AccAddress, hookData *types.HookData, tokenIn sdk.Coin, originalSender string) ibcexported.Acknowledgement {
	if hookData == nil {
		return types.NewSuccessAcknowledgement()
	}

	// Execute SwapAndAction
	if hookData.SwapAndAction != nil {
		return k.ExecuteSwapAndAction(ctx, sender, hookData.SwapAndAction, tokenIn, originalSender)
	}

	return types.NewSuccessAcknowledgement()
}

// ExecuteSwapAndAction executes a SwapAndAction hook
func (k Keeper) ExecuteSwapAndAction(ctx sdk.Context, sender sdk.AccAddress, swapAndAction *types.SwapAndAction, tokenIn sdk.Coin, originalSender string) ibcexported.Acknowledgement {
	if swapAndAction == nil {
		return types.NewSuccessAcknowledgement()
	}

	// Validate: post_swap_action is required and must have exactly one of IBCTransfer or Transfer
	if swapAndAction.PostSwapAction == nil {
		return types.NewCustomErrorAcknowledgement("post_swap_action is required and must have either ibc_transfer or transfer")
	}

	// Validate that exactly one of IBCTransfer or Transfer is set
	hasIBCTransfer := swapAndAction.PostSwapAction.IBCTransfer != nil
	hasTransfer := swapAndAction.PostSwapAction.Transfer != nil
	if hasIBCTransfer && hasTransfer {
		return types.NewCustomErrorAcknowledgement("post_swap_action cannot have both ibc_transfer and transfer; must specify exactly one")
	}
	if !hasIBCTransfer && !hasTransfer {
		return types.NewCustomErrorAcknowledgement("post_swap_action must have either ibc_transfer or transfer")
	}

	// Validate: if post_swap_action is defined, a swap must also be defined
	// If they only want to do a transfer without a swap, they should use packet-forward-middleware
	if swapAndAction.UserSwap == nil || swapAndAction.UserSwap.SwapExactAssetIn == nil {
		return types.NewCustomErrorAcknowledgement("post_swap_action requires a swap to be defined; use packet-forward-middleware for transfers without swaps")
	}

	// Validate post-swap action format and sanity checks BEFORE executing swap
	// This prevents wasting gas on a swap if the post-swap action is malformed
	ack := k.ValidatePostSwapAction(ctx, swapAndAction.PostSwapAction)
	if !ack.Success() {
		return ack
	}

	// Validate destination recover address if provided
	if swapAndAction.DestinationRecoverAddress != "" {
		_, err := sdk.AccAddressFromBech32(swapAndAction.DestinationRecoverAddress)
		if err != nil {
			return types.NewCustomErrorAcknowledgement(fmt.Sprintf("invalid destination_recover_address: %s", swapAndAction.DestinationRecoverAddress))
		}
	}

	// Validate: min_asset is always required for swaps
	if swapAndAction.MinAsset == nil || swapAndAction.MinAsset.Native == nil {
		return types.NewCustomErrorAcknowledgement("min_asset is required for swaps")
	}

	// Store original tokenIn for fallback handling
	// Use sdk.NewCoin to ensure proper deep copy of the coin
	originalTokenIn := sdk.NewCoin(tokenIn.Denom, tokenIn.Amount)

	// Execute the swap first in a cached context
	// This ensures we can rollback the swap if it fails and we use the fallback
	if swapAndAction.UserSwap != nil && swapAndAction.UserSwap.SwapExactAssetIn != nil {
		operations := swapAndAction.UserSwap.SwapExactAssetIn.Operations
		if len(operations) == 0 {
			return types.NewCustomErrorAcknowledgement("no operations provided for swap")
		}

		// Convert operations to routes and validate pool IDs first
		routes := make([]poolmanagertypes.SwapAmountInRoute, len(operations))
		for i, operation := range operations {
			poolId, err := strconv.ParseUint(operation.Pool, 10, 64)
			if err != nil {
				return types.NewCustomErrorAcknowledgement(fmt.Sprintf("invalid pool ID in operation %d: %s", i, operation.Pool))
			}
			routes[i] = poolmanagertypes.SwapAmountInRoute{
				PoolId:        poolId,
				TokenOutDenom: operation.DenomOut,
			}
		}

		// Validate first operation matches tokenIn
		if operations[0].DenomIn != tokenIn.Denom {
			return types.NewCustomErrorAcknowledgement(fmt.Sprintf("first operation denom_in %s does not match token_in %s", operations[0].DenomIn, tokenIn.Denom))
		}

		// Validate operations chain correctly (each DenomOut matches next DenomIn)
		for i := 0; i < len(operations)-1; i++ {
			if operations[i].DenomOut != operations[i+1].DenomIn {
				return types.NewCustomErrorAcknowledgement(fmt.Sprintf("operation %d denom_out %s does not match operation %d denom_in %s", i, operations[i].DenomOut, i+1, operations[i+1].DenomIn))
			}
		}

		// Validate that last operation's DenomOut matches MinAsset.Denom
		lastOperation := operations[len(operations)-1]
		if swapAndAction.MinAsset != nil && swapAndAction.MinAsset.Native != nil {
			if lastOperation.DenomOut != swapAndAction.MinAsset.Native.Denom {
				return types.NewCustomErrorAcknowledgement(fmt.Sprintf("last operation denom_out %s does not match min_asset denom %s", lastOperation.DenomOut, swapAndAction.MinAsset.Native.Denom))
			}
		}

		// Pre-validate denom expansion for IBC transfers BEFORE executing swap
		// This prevents swap from succeeding but IBC transfer failing due to denom expansion issues
		if swapAndAction.PostSwapAction != nil && swapAndAction.PostSwapAction.IBCTransfer != nil {
			expectedOutputDenom := lastOperation.DenomOut

			// Try to expand the denom to validate it can be expanded
			// Use a cache context to avoid side effects
			testCtx, _ := ctx.CacheContext()
			_, err := gammtypes.ExpandIBCDenomToFullPath(testCtx, expectedOutputDenom, k.transferKeeper)
			if err != nil {
				return types.NewCustomErrorAcknowledgement(fmt.Sprintf("cannot expand IBC denom %s for post-swap IBC transfer: %v", expectedOutputDenom, err))
			}
		}

		// Get minimum amount from min_asset for swap validation
		var tokenOut osmomath.Int
		tokenOutMinAmount := osmomath.ZeroInt()

		if swapAndAction.MinAsset != nil && swapAndAction.MinAsset.Native != nil {
			var ok bool
			tokenOutMinAmount, ok = osmomath.NewIntFromString(swapAndAction.MinAsset.Native.Amount)
			if !ok {
				return types.NewCustomErrorAcknowledgement(fmt.Sprintf("invalid min_asset amount: %s", swapAndAction.MinAsset.Native.Amount))
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
					return types.NewCustomErrorAcknowledgement(fmt.Sprintf("affiliate address is required for affiliate_%d", i))
				}
				_, err := sdk.AccAddressFromBech32(affiliate.Address)
				if err != nil {
					return types.NewCustomErrorAcknowledgement(fmt.Sprintf("invalid affiliate address for affiliate_%d: %s", i, affiliate.Address))
				}

				// Validate basis_points_fee
				if affiliate.BasisPointsFee == "" {
					return types.NewCustomErrorAcknowledgement(fmt.Sprintf("affiliate basis_points_fee is required for affiliate_%d", i))
				}
				basisPoints, err := strconv.ParseUint(affiliate.BasisPointsFee, 10, 64)
				if err != nil {
					return types.NewCustomErrorAcknowledgement(fmt.Sprintf("invalid affiliate basis_points_fee for affiliate_%d: %s", i, affiliate.BasisPointsFee))
				}

				// Validate individual basis points don't exceed 10000
				if basisPoints > 10000 {
					return types.NewCustomErrorAcknowledgement(fmt.Sprintf("affiliate basis_points_fee cannot exceed 10000 for affiliate_%d", i))
				}

				totalBasisPoints += basisPoints
			}

			// Validate total basis points don't exceed 10000
			if totalBasisPoints > 10000 {
				return types.NewCustomErrorAcknowledgement("total affiliate basis_points_fee cannot exceed 10000")
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
			if k.tokenizationKeeper.CheckIsAliasDenom(ctx, denom) {
				ack := k.setAutoApproveForIntermediateAddress(ctx, sender.String(), denom)
				if !ack.Success() {
					return types.NewCustomErrorAcknowledgement(fmt.Sprintf("failed to set auto-approve for intermediate address, denom: %s", denom))
				}
			}
		}

		// Execute swap in a cached context so we can rollback if it fails
		swapCacheCtx, writeSwapCache := ctx.CacheContext()

		// Clear any previous deterministic errors before starting the swap
		types.ClearDeterministicError(swapCacheCtx)

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
			k.Logger(ctx).Error("custom-hooks: swap failed", "error", err, "sender", sender.String(), "original_sender", originalSender, "token_in", tokenIn.String(), "num_operations", len(operations))

			swapErrMsg := fmt.Sprintf("swap failed: sender=%s (derived from %s), token_in=%s, num_operations=%d", sender.String(), originalSender, tokenIn.String(), len(operations))

			// Hacky way to get exact error messages without the lack of deterministic error messages
			// Try to get deterministic error message from context
			if detErrMsg, found := types.GetDeterministicError(swapCacheCtx); found {
				swapErrMsg = fmt.Sprintf("%s: %s", swapErrMsg, detErrMsg)
			}

			// Swap failed - check if we have a destination recover address
			if swapAndAction.DestinationRecoverAddress != "" {
				// Cache context is automatically discarded - swap never happened
				// Log the fallback action
				k.Logger(ctx).Info("custom-hooks: swap failed, using destination recover address",
					"error", err,
					"sender", sender.String(),
					"original_sender", originalSender,
					"destination_recover_address", swapAndAction.DestinationRecoverAddress,
					"original_token", originalTokenIn.String())

				// Execute fallback transfer in a cache context to ensure atomicity
				// If fallback transfer fails, we want to roll back the IBC transfer
				fallbackCacheCtx, writeFallbackCache := ctx.CacheContext()
				types.ClearDeterministicError(fallbackCacheCtx)
				ack := k.ExecuteLocalTransfer(fallbackCacheCtx, sender, &types.TransferInfo{ToAddress: swapAndAction.DestinationRecoverAddress}, originalTokenIn)
				if !ack.Success() {
					// Fallback transfer failed - cache is automatically discarded
					// Combine ack error message with deterministic error from context
					ackErrMsg, _ := types.GetAckError(ack)
					combinedMsg := "failed to send tokens to destination recover address: swap failed and fallback transfer failed"
					if ackErrMsg != "" {
						combinedMsg = fmt.Sprintf("%s: %s", combinedMsg, ackErrMsg)
					}
					if detErrMsg, found := types.GetDeterministicError(fallbackCacheCtx); found {
						combinedMsg = fmt.Sprintf("%s: %s", combinedMsg, detErrMsg)
					}
					return types.NewCustomErrorAcknowledgement(combinedMsg)
				}

				// Fallback transfer succeeded - commit the fallback state
				writeFallbackCache()

				// Emit event for fallback path
				event := sdk.NewEvent(
					"swap_and_action_fallback",
					sdk.NewAttribute("module", "custom-hooks"),
					sdk.NewAttribute("sender", sender.String()),
					sdk.NewAttribute("swap_error", swapErrMsg),
					sdk.NewAttribute("original_sender", originalSender),
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

				// Successfully handled fallback - return success to allow IBC transfer to succeed
				// The swap cache context was automatically discarded, fallback cache was committed
				return types.NewSuccessAcknowledgement()
			}

			// No fallback address - cache context is automatically discarded
			// Return error acknowledgement with deterministic error message
			return types.NewCustomErrorAcknowledgement(swapErrMsg)
		}

		// Swap succeeded - commit the swap cache context
		writeSwapCache()

		// Emit event for successful swap path
		firstOperation := operations[0]
		lastOperation = operations[len(operations)-1]
		event := sdk.NewEvent(
			"swap_and_action_success",
			sdk.NewAttribute("module", "custom-hooks"),
			sdk.NewAttribute("sender", sender.String()),
			sdk.NewAttribute("original_sender", originalSender),
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
		// Check if attempting to IBC transfer a non-ICS20 compatible denomination
		if !k.sendManagerKeeper.IsICS20Compatible(ctx, tokenIn.Denom) {
			standardName := k.sendManagerKeeper.StandardName(ctx, tokenIn.Denom)
			return types.NewCustomErrorAcknowledgement(fmt.Sprintf("cannot IBC transfer %s denominations: %s", standardName, tokenIn.Denom))
		}

		// Security: LOW-010 - Timeout timestamp bounds validation
		// Use timeout_timestamp if provided, otherwise use default
		var timeoutTimestamp uint64
		currentTime := uint64(ctx.BlockTime().UnixNano())

		if swapAndAction.TimeoutTimestamp != nil {
			timeoutTimestamp = *swapAndAction.TimeoutTimestamp

			// Validate timeout is not in the past
			if timeoutTimestamp <= currentTime {
				return types.NewCustomErrorAcknowledgement(fmt.Sprintf("timeout timestamp must be in the future: provided=%d, current=%d", timeoutTimestamp, currentTime))
			}

			// Validate timeout is not too far in the future (max 1 week)
			// 1 week is reasonable for IBC transfers which may take time to complete
			maxTimeout := currentTime + uint64(7*24*60*60*1e9) // 1 week in nanoseconds
			if timeoutTimestamp > maxTimeout {
				return types.NewCustomErrorAcknowledgement(fmt.Sprintf("timeout timestamp exceeds maximum allowed: provided=%d, max=%d", timeoutTimestamp, maxTimeout))
			}
		} else {
			// Default: current time + 5 minutes
			timeoutTimestamp = currentTime + uint64(5*60*1e9)
		}

		recoverAddress := swapAndAction.PostSwapAction.IBCTransfer.IBCInfo.RecoverAddress
		recoverAddressAddr, err := sdk.AccAddressFromBech32(recoverAddress)
		if err != nil {
			return types.NewCustomErrorAcknowledgement(fmt.Sprintf("invalid recover address: %s", recoverAddress))
		}

		// Send to recover address then IBC it out
		// Clear any previous deterministic errors before starting the transfer
		types.ClearDeterministicError(ctx)
		ack := k.ExecuteLocalTransfer(ctx, sender, &types.TransferInfo{ToAddress: recoverAddress}, tokenIn)
		if !ack.Success() {
			// Combine ack error message with deterministic error from context
			ackErrMsg, _ := types.GetAckError(ack)
			if detErrMsg, found := types.GetDeterministicError(ctx); found {
				combinedMsg := fmt.Sprintf("%s: %s", ackErrMsg, detErrMsg)
				return types.NewCustomErrorAcknowledgement(combinedMsg)
			}
			return ack
		}

		ack = k.ExecuteIBCTransfer(ctx, recoverAddressAddr, swapAndAction.PostSwapAction.IBCTransfer, tokenIn, timeoutTimestamp)
		if !ack.Success() {
			return ack
		}
	} else if swapAndAction.PostSwapAction.Transfer != nil {
		// Clear any previous deterministic errors before starting the transfer
		types.ClearDeterministicError(ctx)
		ack := k.ExecuteLocalTransfer(ctx, sender, swapAndAction.PostSwapAction.Transfer, tokenIn)
		if !ack.Success() {
			// Combine ack error message with deterministic error from context
			ackErrMsg, _ := types.GetAckError(ack)
			if detErrMsg, found := types.GetDeterministicError(ctx); found {
				combinedMsg := fmt.Sprintf("%s: %s", ackErrMsg, detErrMsg)
				return types.NewCustomErrorAcknowledgement(combinedMsg)
			}
			return ack
		}
	}

	return types.NewSuccessAcknowledgement()
}

// ValidatePostSwapAction performs sanity checks on the post-swap action before executing the swap
// This prevents wasting gas on a swap if the post-swap action is malformed
func (k Keeper) ValidatePostSwapAction(ctx sdk.Context, postSwapAction *types.PostSwapAction) ibcexported.Acknowledgement {
	if postSwapAction == nil {
		return types.NewCustomErrorAcknowledgement("post_swap_action cannot be nil")
	}

	// Validate IBC transfer if present
	if postSwapAction.IBCTransfer != nil {
		if postSwapAction.IBCTransfer.IBCInfo == nil {
			return types.NewCustomErrorAcknowledgement("ibc_transfer.ibc_info is required")
		}

		ibcInfo := postSwapAction.IBCTransfer.IBCInfo

		// Validate source channel is not empty
		if ibcInfo.SourceChannel == "" {
			return types.NewCustomErrorAcknowledgement("ibc_transfer.ibc_info.source_channel is required")
		}

		// Validate channel exists
		_, found := k.channelKeeper.GetChannel(ctx, transfertypes.PortID, ibcInfo.SourceChannel)
		if !found {
			return types.NewCustomErrorAcknowledgement(fmt.Sprintf("IBC channel does not exist: %s", ibcInfo.SourceChannel))
		}

		// Validate receiver is not empty
		if ibcInfo.Receiver == "" {
			return types.NewCustomErrorAcknowledgement("ibc_transfer.ibc_info.receiver is required")
		}

		// Validate recover address if provided
		if ibcInfo.RecoverAddress != "" {
			_, err := sdk.AccAddressFromBech32(ibcInfo.RecoverAddress)
			if err != nil {
				return types.NewCustomErrorAcknowledgement(fmt.Sprintf("invalid ibc_transfer.ibc_info.recover_address: %s", ibcInfo.RecoverAddress))
			}
		}

		// Security: MED-007 - Channel capability validation timing
		// Note: Capability is validated here for early failure, but it's also validated
		// IBC v10: Capabilities removed - channel validation is handled by IBC core
	}

	// Validate local transfer if present
	if postSwapAction.Transfer != nil {
		if postSwapAction.Transfer.ToAddress == "" {
			return types.NewCustomErrorAcknowledgement("transfer.to_address is required")
		}

		// Validate to_address can be parsed as a valid address
		_, err := sdk.AccAddressFromBech32(postSwapAction.Transfer.ToAddress)
		if err != nil {
			return types.NewCustomErrorAcknowledgement(fmt.Sprintf("invalid transfer.to_address: %s", postSwapAction.Transfer.ToAddress))
		}
	}

	return types.NewSuccessAcknowledgement()
}

// ExecuteIBCTransfer executes an IBC transfer
func (k Keeper) ExecuteIBCTransfer(ctx sdk.Context, sender sdk.AccAddress, ibcTransfer *types.IBCTransferInfo, token sdk.Coin, timeoutTimestamp uint64) ibcexported.Acknowledgement {
	if ibcTransfer == nil || ibcTransfer.IBCInfo == nil {
		return types.NewSuccessAcknowledgement()
	}

	ibcInfo := ibcTransfer.IBCInfo

	// IBC v10: Capabilities removed - channel validation is handled by IBC core

	// Use zero height for timeout (no timeout height)
	timeoutHeight := clienttypes.ZeroHeight()

	// Create transfer packet data
	denom, err := gammtypes.ExpandIBCDenomToFullPath(ctx, token.Denom, k.transferKeeper)
	if err != nil {
		return types.NewCustomErrorAcknowledgement(fmt.Sprintf("failed to expand IBC denom: %v", err))
	}

	packetData := transfertypes.FungibleTokenPacketData{
		Denom:    denom,
		Amount:   token.Amount.String(),
		Sender:   sender.String(),
		Receiver: ibcInfo.Receiver,
		Memo:     ibcInfo.Memo,
	}

	// Marshal packet data
	data, err := transfertypes.ModuleCdc.MarshalJSON(&packetData)
	if err != nil {
		return types.NewCustomErrorAcknowledgement("failed to marshal IBC packet data")
	}

	// Send IBC packet (IBC v10: capabilities removed)
	_, err = k.ics4Wrapper.SendPacket(
		ctx,
		transfertypes.PortID,
		ibcInfo.SourceChannel,
		timeoutHeight,
		timeoutTimestamp,
		data,
	)
	if err != nil {
		return types.NewCustomErrorAcknowledgement(fmt.Sprintf("failed to send IBC packet: channel=%s, denom=%s, amount=%s", ibcInfo.SourceChannel, token.Denom, token.Amount.String()))
	}

	return types.NewSuccessAcknowledgement()
}

// ExecuteLocalTransfer executes a local bank transfer
func (k Keeper) ExecuteLocalTransfer(ctx sdk.Context, sender sdk.AccAddress, transfer *types.TransferInfo, token sdk.Coin) ibcexported.Acknowledgement {
	if transfer == nil {
		return types.NewCustomErrorAcknowledgement("transfer info is nil")
	}

	if transfer.ToAddress == "" {
		return types.NewCustomErrorAcknowledgement("to_address is required for local transfer")
	}

	// Parse the recipient address
	toAddr, err := sdk.AccAddressFromBech32(transfer.ToAddress)
	if err != nil {
		return types.NewCustomErrorAcknowledgement(fmt.Sprintf("invalid transfer.to_address: %s", transfer.ToAddress))
	}

	// Only set auto-approve flags for wrapped tokenization denoms
	// For regular denoms, SendCoinsFromIntermediateAddress will handle them via bank keeper
	if k.tokenizationKeeper.CheckIsAliasDenom(ctx, token.Denom) {
		collection, err := k.tokenizationKeeper.ParseCollectionFromDenom(ctx, token.Denom)
		if err != nil {
			return types.NewCustomErrorAcknowledgement(fmt.Sprintf("failed to parse collection from denom: %s", token.Denom))
		}

		// This is setting the auto-approve flags for the intermediate sender address
		// Edge case, but this sets it in the case of default self initiated outgoing is not approved
		err = k.tokenizationKeeper.SetAllAutoApprovalFlagsForIntermediateAddress(ctx, collection, sender.String())
		if err != nil {
			return types.NewCustomErrorAcknowledgement(fmt.Sprintf("failed to set all auto approval flags for address: %s", toAddr.String()))
		}
	}

	// Execute bank transfer
	ack := k.SendCoinsFromIntermediateAddress(ctx, sender, toAddr, sdk.Coins{token})
	if !ack.Success() {
		return ack
	}

	return types.NewSuccessAcknowledgement()
}

// IMPORTANT: Should ONLY be called when from address is an intermediate address
func (k Keeper) SendCoinsFromIntermediateAddress(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, coins sdk.Coins) ibcexported.Acknowledgement {
	for _, coin := range coins {
		err := k.sendManagerKeeper.SendCoinWithAliasRouting(ctx, from, to, &coin)
		if err != nil {
			return types.NewCustomErrorAcknowledgement(fmt.Sprintf("failed to send wrapped tokens: denom=%s, from=%s, to=%s", coin.Denom, from.String(), to.String()))
		}
	}

	return types.NewSuccessAcknowledgement()
}

// setAutoApproveForIntermediateAddress sets auto-approve flags for an intermediate address
// This is similar to how it's done for pool addresses in gamm keeper
func (k Keeper) setAutoApproveForIntermediateAddress(ctx sdk.Context, intermediateAddress string, denom string) ibcexported.Acknowledgement {
	// Parse collection ID from denom (format: badges:COLL_ID:* or badgeslp:COLL_ID:*)
	parts := strings.Split(denom, ":")
	if len(parts) < 3 {
		return types.NewCustomErrorAcknowledgement(fmt.Sprintf("invalid denom format: %s", denom))
	}

	collectionIdStr := parts[1]
	collectionId, err := strconv.ParseUint(collectionIdStr, 10, 64)
	if err != nil {
		return types.NewCustomErrorAcknowledgement(fmt.Sprintf("invalid collection ID in denom: %s", collectionIdStr))
	}

	// Get collection
	collection, found := k.tokenizationKeeper.GetCollectionFromStore(ctx, sdkmath.NewUint(collectionId))
	if !found {
		return types.NewCustomErrorAcknowledgement(fmt.Sprintf("collection not found: %s", collectionIdStr))
	}

	// Get current balance or apply default
	currBalances, _, err := k.tokenizationKeeper.GetBalanceOrApplyDefault(ctx, collection, intermediateAddress)
	if err != nil {
		return types.NewCustomErrorAcknowledgement(fmt.Sprintf("failed to get balance for intermediate address: %v", err))
	}

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
		err = k.tokenizationKeeper.SetBalanceForAddress(ctx, collection, intermediateAddress, currBalances)
		if err != nil {
			return types.NewCustomErrorAcknowledgement("failed to set auto-approve for intermediate address")
		}
	}

	return types.NewSuccessAcknowledgement()
}
