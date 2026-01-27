package keeper

import (
	"fmt"

	"github.com/bitbadges/bitbadgeschain/x/gamm/types"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	transfertypes "github.com/cosmos/ibc-go/v10/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v10/modules/core/02-client/types"

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
	sendManagerKeeper   types.SendManagerKeeper
	transferKeeper      types.TransferKeeper

	// IBC keepers (optional, for IBC transfer functionality)
	ics4Wrapper   types.ICS4Wrapper
	channelKeeper types.ChannelKeeper
}

func NewKeeper(
	cdc codec.BinaryCodec, storeKey storetypes.StoreKey, paramSpace paramtypes.Subspace, accountKeeper types.AccountKeeper, bankKeeper types.BankKeeper, communityPoolKeeper types.CommunityPoolKeeper, badgesKeeper badgeskeeper.Keeper,
	sendManagerKeeper types.SendManagerKeeper,
	transferKeeper types.TransferKeeper,
	ics4Wrapper types.ICS4Wrapper,
	channelKeeper types.ChannelKeeper,
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
		sendManagerKeeper:   sendManagerKeeper,
		transferKeeper:      transferKeeper,
		ics4Wrapper:   ics4Wrapper,
		channelKeeper: channelKeeper,
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

	// IBC v10: scoped keeper removed - capabilities no longer used

	// Validate channel exists
	_, found := k.channelKeeper.GetChannel(ctx, transfertypes.PortID, ibcTransferInfo.SourceChannel)
	if !found {
		return fmt.Errorf("IBC channel %s does not exist", ibcTransferInfo.SourceChannel)
	}

	// IBC v10: Capabilities removed - channel validation is handled by IBC core

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

	// Send IBC packet (IBC v10: capabilities removed)
	_, err = k.ics4Wrapper.SendPacket(
		ctx,
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

// SendCoinsToPoolWithAliasRouting sends coins to a pool, wrapping badges denoms if needed.
// IMPORTANT: Should ONLY be called when to address is a pool address
func (k Keeper) SendCoinsToPoolWithAliasRouting(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, coins sdk.Coins) error {
	return k.badgesKeeper.SendCoinsToPoolWithAliasRouting(ctx, from, to, coins)
}

// SendCoinsFromPoolWithAliasRouting sends coins from a pool, unwrapping badges denoms if needed.
// IMPORTANT: Should ONLY be called when from address is a pool address
func (k Keeper) SendCoinsFromPoolWithAliasRouting(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, coins sdk.Coins) error {
	return k.badgesKeeper.SendCoinsFromPoolWithAliasRouting(ctx, from, to, coins)
}

// FundCommunityPoolWithAliasRouting funds the community pool, wrapping badges denoms if needed.
// Used for taker fees
func (k Keeper) FundCommunityPoolWithAliasRouting(ctx sdk.Context, from sdk.AccAddress, coins sdk.Coins) error {
	return k.sendManagerKeeper.FundCommunityPoolWithAliasRouting(ctx, from, coins)
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
		if k.badgesKeeper.CheckIsAliasDenom(ctx, coin.Denom) {
			collection, err := k.badgesKeeper.ParseCollectionFromDenom(ctx, coin.Denom)
			if err != nil {
				return fmt.Errorf("failed to parse collection from denom: %s: %w", coin.Denom, err)
			}

			// Get the balances that would be needed for the recorded amount
			balancesNeeded, err := badgeskeeper.GetBalancesToTransferWithAlias(collection, coin.Denom, sdkmath.NewUintFromBigInt(coin.Amount.BigInt()))
			if err != nil {
				return fmt.Errorf("failed to get balances to transfer for denom: %s: %w", coin.Denom, err)
			}

			// Get the pool's current balance
			poolBalances, _, err := k.badgesKeeper.GetBalanceOrApplyDefault(ctx, collection, poolAddress.String())
			if err != nil {
				return err
			}

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
		if !badgeskeeper.CheckStartsWithWrappedOrAliasDenom(coin.Denom) {
			continue
		}

		// Parse collection from denom
		collection, err := k.badgesKeeper.ParseCollectionFromDenom(ctx, coin.Denom)
		if err != nil {
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
