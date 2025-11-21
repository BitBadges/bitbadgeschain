package keeper

import (
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

	badgeskeeper "github.com/bitbadges/bitbadgeschain/x/badges/keeper"
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

	// IBC keepers (optional, for IBC transfer functionality)
	ics4Wrapper   types.ICS4Wrapper
	channelKeeper types.ChannelKeeper
	scopedKeeper  types.ScopedKeeper
}

func NewKeeper(
	cdc codec.BinaryCodec, storeKey storetypes.StoreKey, paramSpace paramtypes.Subspace, accountKeeper types.AccountKeeper, bankKeeper types.BankKeeper, communityPoolKeeper types.CommunityPoolKeeper, badgesKeeper badgeskeeper.Keeper,
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
	packetData := transfertypes.FungibleTokenPacketData{
		Denom:    token.Denom,
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
