package keeper

import (
	"context"
	"fmt"

	"github.com/bitbadges/bitbadgeschain/x/gamm/types"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	host "github.com/cosmos/ibc-go/v8/modules/core/24-host"

	sdkmath "cosmossdk.io/math"
	badgeskeeper "github.com/bitbadges/bitbadgeschain/x/badges/keeper"
	badgestypes "github.com/bitbadges/bitbadgeschain/x/badges/types"
	poolmanagertypes "github.com/bitbadges/bitbadgeschain/x/poolmanager/types"
)

func permContains(perms []string, perm string) bool {
	for _, v := range perms {
		if v == perm {
			return true
		}
	}

	return false
}

type Keeper struct {
	storeKey storetypes.StoreKey
	cdc      codec.BinaryCodec

	paramSpace paramtypes.Subspace

	// keepers
	accountKeeper       types.AccountKeeper
	bankKeeper          types.BankKeeper
	communityPoolKeeper types.CommunityPoolKeeper
	poolManager         types.PoolManager
	badgesKeeper        badgeskeeper.Keeper
	transferKeeper      types.TransferKeeper

	// IBC keepers (optional, for IBC transfer functionality)
	ics4Wrapper   types.ICS4Wrapper
	channelKeeper types.ChannelKeeper
	scopedKeeper  types.ScopedKeeper
}

func NewKeeper(
	cdc codec.BinaryCodec, storeKey storetypes.StoreKey, paramSpace paramtypes.Subspace, accountKeeper types.AccountKeeper, bankKeeper types.BankKeeper, communityPoolKeeper types.CommunityPoolKeeper, badgesKeeper badgeskeeper.Keeper,
	transferKeeper types.TransferKeeper,
	ics4Wrapper types.ICS4Wrapper,
	channelKeeper types.ChannelKeeper,
	scopedKeeper types.ScopedKeeper,
) Keeper {
	// Ensure that the module account are set.
	moduleAddr, perms := accountKeeper.GetModuleAddressAndPermissions(types.ModuleName)
	if moduleAddr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}
	if !permContains(perms, authtypes.Minter) {
		panic(fmt.Sprintf("%s module account should have the minter permission", types.ModuleName))
	}
	if !permContains(perms, authtypes.Burner) {
		panic(fmt.Sprintf("%s module account should have the burner permission", types.ModuleName))
	}
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}
	return Keeper{
		storeKey:   storeKey,
		cdc:        cdc,
		paramSpace: paramSpace,
		// keepers
		accountKeeper:       accountKeeper,
		bankKeeper:          bankKeeper,
		communityPoolKeeper: communityPoolKeeper,
		badgesKeeper:        badgesKeeper,
		transferKeeper:      transferKeeper,
		ics4Wrapper:         ics4Wrapper,
		channelKeeper:       channelKeeper,
		scopedKeeper:        scopedKeeper,
	}
}

func (k *Keeper) SetPoolManager(poolManager types.PoolManager) {
	k.poolManager = poolManager
}

// ExecuteIBCTransfer executes an IBC transfer
// This method is similar to the custom hooks keeper's ExecuteIBCTransfer
func (k Keeper) ExecuteIBCTransfer(ctx sdk.Context, sender sdk.AccAddress, ibcTransferInfo *types.IBCTransferInfo, token sdk.Coin) error {
	if ibcTransferInfo == nil {
		return fmt.Errorf("ibc_transfer_info cannot be nil")
	}

	if k.ics4Wrapper == nil {
		return fmt.Errorf("ICS4 wrapper not set")
	}

	if k.channelKeeper == nil {
		return fmt.Errorf("channel keeper not set")
	}

	if k.scopedKeeper == nil {
		return fmt.Errorf("scoped keeper not set")
	}

	// Validate channel exists
	_, found := k.channelKeeper.GetChannel(ctx, transfertypes.PortID, ibcTransferInfo.SourceChannel)
	if !found {
		return fmt.Errorf("IBC channel %s does not exist", ibcTransferInfo.SourceChannel)
	}

	// Get channel capability
	capPath := host.ChannelCapabilityPath(transfertypes.PortID, ibcTransferInfo.SourceChannel)
	channelCap, ok := k.scopedKeeper.GetCapability(ctx, capPath)
	if !ok {
		return fmt.Errorf("channel capability not found for channel %s", ibcTransferInfo.SourceChannel)
	}

	// Use zero height for timeout (no timeout height)
	timeoutHeight := clienttypes.ZeroHeight()

	// Use timeout_timestamp from message, or default to current time + 5 minutes
	timeoutTimestamp := ibcTransferInfo.TimeoutTimestamp
	if timeoutTimestamp == 0 {
		timeoutTimestamp = uint64(ctx.BlockTime().UnixNano()) + uint64(5*60*1e9)
	}

	// Create transfer packet data
	denom, err := types.ExpandIBCDenomToFullPath(ctx, token.Denom, k.transferKeeper)
	if err != nil {
		return fmt.Errorf("failed to expand IBC denom: %w", err)
	}

	packetData := transfertypes.FungibleTokenPacketData{
		Denom:    denom,
		Amount:   token.Amount.String(),
		Sender:   sender.String(),
		Receiver: ibcTransferInfo.Receiver,
		Memo:     ibcTransferInfo.Memo,
	}

	// Marshal packet data
	data, err := transfertypes.ModuleCdc.MarshalJSON(&packetData)
	if err != nil {
		return fmt.Errorf("failed to marshal packet data: %w", err)
	}

	// Send IBC packet
	_, err = k.ics4Wrapper.SendPacket(
		ctx,
		channelCap,
		transfertypes.PortID,
		ibcTransferInfo.SourceChannel,
		timeoutHeight,
		timeoutTimestamp,
		data,
	)
	if err != nil {
		return fmt.Errorf("failed to send IBC packet: %w", err)
	}

	return nil
}

// GetParams returns the total set params.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the total set of params.
func (k Keeper) setParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

// SetParam sets a specific gamm module's parameter with the provided parameter.
func (k Keeper) SetParam(ctx sdk.Context, key []byte, value interface{}) {
	k.paramSpace.Set(ctx, key, value)
}

// Wrapper methods that delegate to badges keeper for pool integration

// SendCoinsToPoolWithWrapping sends coins to a pool, wrapping badges denoms if needed.
// IMPORTANT: Should ONLY be called when to address is a pool address
func (k Keeper) SendCoinsToPoolWithWrapping(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, coins sdk.Coins) error {
	// Create an adapter that implements badges BankKeeper interface
	bankKeeperAdapter := &bankKeeperAdapter{bankKeeper: k.bankKeeper}
	return k.badgesKeeper.SendCoinsToPoolWithWrapping(ctx, bankKeeperAdapter, from, to, coins)
}

// SendCoinsFromPoolWithUnwrapping sends coins from a pool, unwrapping badges denoms if needed.
// IMPORTANT: Should ONLY be called when from address is a pool address
func (k Keeper) SendCoinsFromPoolWithUnwrapping(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, coins sdk.Coins) error {
	// Create an adapter that implements badges BankKeeper interface
	bankKeeperAdapter := &bankKeeperAdapter{bankKeeper: k.bankKeeper}
	return k.badgesKeeper.SendCoinsFromPoolWithUnwrapping(ctx, bankKeeperAdapter, from, to, coins)
}

// FundCommunityPoolWithWrapping funds the community pool, wrapping badges denoms if needed.
// Used for taker fees
func (k Keeper) FundCommunityPoolWithWrapping(ctx sdk.Context, from sdk.AccAddress, coins sdk.Coins) error {
	// Create an adapter that implements badges BankKeeper interface
	bankKeeperAdapter := &bankKeeperAdapter{bankKeeper: k.bankKeeper}
	// Create an adapter for community pool keeper
	communityPoolKeeperAdapter := &communityPoolKeeperAdapter{communityPoolKeeper: k.communityPoolKeeper}
	return k.badgesKeeper.FundCommunityPoolWithWrapping(ctx, bankKeeperAdapter, communityPoolKeeperAdapter, from, coins)
}

// CheckPoolLiquidityInvariant checks that the pool address has enough underlying assets for all recorded pool liquidity.
// This is a global invariant check that compares recorded liquidity with actual balances behind the scenes.
func (k Keeper) CheckPoolLiquidityInvariant(ctx sdk.Context, pool poolmanagertypes.PoolI) error {
	poolAddress := pool.GetAddress()

	// Convert to CFMMPoolI to access GetTotalPoolLiquidity
	cfmmPool, ok := pool.(types.CFMMPoolI)
	if !ok {
		return fmt.Errorf("pool does not implement CFMMPoolI")
	}

	poolLiquidity := cfmmPool.GetTotalPoolLiquidity(ctx)

	// Iterate over all denoms in the pool's liquidity
	for _, coin := range poolLiquidity {
		// Check if this is a wrapped badges denom
		if k.badgesKeeper.CheckIsWrappedDenom(ctx, coin.Denom) {
			collection, err := k.badgesKeeper.ParseCollectionFromDenom(ctx, coin.Denom)
			if err != nil {
				return fmt.Errorf("failed to parse collection from denom: %s: %w", coin.Denom, err)
			}

			// Get the balances that would be needed for the recorded amount
			balancesNeeded, err := badgeskeeper.GetBalancesToTransfer(collection, coin.Denom, sdkmath.NewUintFromBigInt(coin.Amount.BigInt()))
			if err != nil {
				return fmt.Errorf("failed to get balances to transfer for denom: %s: %w", coin.Denom, err)
			}

			// Get the pool's current balance
			poolBalances, _ := k.badgesKeeper.GetBalanceOrApplyDefault(ctx, collection, poolAddress.String())

			// Try to subtract needed balances from pool balances - will error on underflow
			poolBalancesCopy := badgestypes.DeepCopyBalances(poolBalances.Balances)
			_, err = badgestypes.SubtractBalances(ctx, balancesNeeded, poolBalancesCopy)
			if err != nil {
				return fmt.Errorf("pool address %s has insufficient badges liquidity for denom %s",
					poolAddress.String(), coin.Denom)
			}
		} else {
			// For all other denoms (including IBC and native), check bank balance
			allBalances := k.bankKeeper.GetAllBalances(ctx, poolAddress)
			poolBalance := allBalances.AmountOf(coin.Denom)

			if poolBalance.LT(coin.Amount) {
				return fmt.Errorf("pool address %s has insufficient liquidity: required %s, available %s for denom %s",
					poolAddress.String(), coin.Amount.String(), poolBalance.String(), coin.Denom)
			}
		}
	}

	return nil
}

// ValidatePoolCreationAllowed checks if pool creation is allowed for all badges assets in the given coins.
// Returns an error if any badges asset has disablePoolCreation set to true.
func (k Keeper) ValidatePoolCreationAllowed(ctx sdk.Context, coins sdk.Coins) error {
	for _, coin := range coins {
		// Check if this is a badges denom
		if !badgeskeeper.CheckStartsWithBadges(coin.Denom) {
			continue
		}

		// Parse collection from denom
		collection, err := k.badgesKeeper.ParseCollectionFromDenom(ctx, coin.Denom)
		if err != nil {
			// If we can't parse the collection, skip it (might be a malformed denom)
			continue
		}

		// Check if the collection disables pool creation
		// If invariants is nil or disablePoolCreation is not explicitly set to true, allow it (default behavior)
		// Only fail if disablePoolCreation is explicitly set to true
		if collection.Invariants != nil && collection.Invariants.DisablePoolCreation {
			return fmt.Errorf("pool creation is not allowed for collection %s (denom: %s). The collection's disablePoolCreation invariant is set to true",
				collection.CollectionId.String(),
				coin.Denom)
		}
	}

	return nil
}

// CheckIsWrappedDenom checks if a denom is a wrapped badges denom.
// This method is required by the custom-hooks GammKeeper interface.
func (k Keeper) CheckIsWrappedDenom(ctx sdk.Context, denom string) bool {
	return k.badgesKeeper.CheckIsWrappedDenom(ctx, denom)
}

// SendNativeTokensFromPool sends native badges tokens from a pool address.
// This method is required by the custom-hooks GammKeeper interface.
func (k Keeper) SendNativeTokensFromPool(ctx sdk.Context, poolAddress string, recipientAddress string, denom string, amount sdkmath.Uint) error {
	return k.badgesKeeper.SendNativeTokensFromPool(ctx, poolAddress, recipientAddress, denom, amount)
}

// SendNativeTokensToPool sends native badges tokens to a pool address.
// This method is used by tests and may be used by other modules.
func (k Keeper) SendNativeTokensToPool(ctx sdk.Context, recipientAddress string, poolAddress string, denom string, amount sdkmath.Uint) error {
	return k.badgesKeeper.SendNativeTokensToPool(ctx, recipientAddress, poolAddress, denom, amount)
}

// ParseCollectionFromDenom parses a collection from a badges denom.
// This method is used by tests.
func (k Keeper) ParseCollectionFromDenom(ctx sdk.Context, denom string) (*badgestypes.TokenCollection, error) {
	return k.badgesKeeper.ParseCollectionFromDenom(ctx, denom)
}

// bankKeeperAdapter adapts gamm BankKeeper to badges BankKeeper interface
type bankKeeperAdapter struct {
	bankKeeper types.BankKeeper
}

func (a *bankKeeperAdapter) SpendableCoins(ctx context.Context, addr sdk.AccAddress) sdk.Coins {
	// Not used in pool integration, return empty coins
	return sdk.Coins{}
}

func (a *bankKeeperAdapter) SendCoins(ctx context.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error {
	return a.bankKeeper.SendCoins(ctx, fromAddr, toAddr, amt)
}

func (a *bankKeeperAdapter) MintCoins(ctx context.Context, moduleName string, amt sdk.Coins) error {
	return a.bankKeeper.MintCoins(ctx, moduleName, amt)
}

func (a *bankKeeperAdapter) BurnCoins(ctx context.Context, name string, amt sdk.Coins) error {
	return a.bankKeeper.BurnCoins(ctx, name, amt)
}

func (a *bankKeeperAdapter) SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error {
	return a.bankKeeper.SendCoinsFromModuleToAccount(ctx, senderModule, recipientAddr, amt)
}

func (a *bankKeeperAdapter) SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error {
	return a.bankKeeper.SendCoinsFromAccountToModule(ctx, senderAddr, recipientModule, amt)
}

func (a *bankKeeperAdapter) GetAllBalances(ctx context.Context, addr sdk.AccAddress) sdk.Coins {
	return a.bankKeeper.GetAllBalances(ctx, addr)
}

// communityPoolKeeperAdapter adapts gamm CommunityPoolKeeper to badges DistributionKeeper interface
type communityPoolKeeperAdapter struct {
	communityPoolKeeper types.CommunityPoolKeeper
}

func (a *communityPoolKeeperAdapter) FundCommunityPool(ctx context.Context, amount sdk.Coins, sender sdk.AccAddress) error {
	return a.communityPoolKeeper.FundCommunityPool(ctx, amount, sender)
}
