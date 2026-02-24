package types

import (
	"context"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	"github.com/ethereum/go-ethereum/common"

	poolmanagertypes "github.com/bitbadges/bitbadgeschain/x/poolmanager/types"
)

// AccountKeeper defines the expected interface for the Account module.
type AccountKeeper interface {
	GetAccount(context.Context, sdk.AccAddress) sdk.AccountI
	NewAccountWithAddress(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
	SetAccount(ctx context.Context, acc sdk.AccountI)
	// Methods imported from account should be defined here
}

// BankKeeper defines the expected interface for the Bank module.
type BankKeeper interface {
	SpendableCoins(context.Context, sdk.AccAddress) sdk.Coins
	SendCoins(context.Context, sdk.AccAddress, sdk.AccAddress, sdk.Coins) error
	MintCoins(context.Context, string, sdk.Coins) error
	BurnCoins(context.Context, string, sdk.Coins) error
	SendCoinsFromModuleToAccount(context.Context, string, sdk.AccAddress, sdk.Coins) error
	SendCoinsFromAccountToModule(context.Context, sdk.AccAddress, string, sdk.Coins) error
	GetAllBalances(context.Context, sdk.AccAddress) sdk.Coins
	// Methods imported from bank should be defined here
}

// DistributionKeeper defines the expected interface for the Distribution module.
type DistributionKeeper interface {
	FundCommunityPool(context.Context, sdk.Coins, sdk.AccAddress) error
}

// ParamSubspace defines the expected Subspace interface for parameters.
type ParamSubspace interface {
	Get(context.Context, []byte, interface{})
	Set(context.Context, []byte, interface{})
}

// GammKeeper defines the expected interface for checking liquidity pools.
type GammKeeper interface {
	GetPool(ctx sdk.Context, poolId uint64) (poolmanagertypes.PoolI, error)
}

// EVMKeeper defines the expected interface for checking EVM contracts.
type EVMKeeper interface {
	IsContract(ctx sdk.Context, addr common.Address) bool
	// CallEVMWithData executes a contract call with raw calldata.
	// Set commit=false for read-only calls (staticcall).
	// gasCap limits the gas for the call (nil for no limit).
	CallEVMWithData(ctx sdk.Context, from common.Address, contract *common.Address, data []byte, commit bool, gasCap *big.Int) (*evmtypes.MsgEthereumTxResponse, error)
}

// SendManagerKeeper defines the expected interface for the SendManager module.
type SendManagerKeeper interface {
	SendCoinWithAliasRouting(ctx sdk.Context, fromAddressAcc sdk.AccAddress, toAddressAcc sdk.AccAddress, coin *sdk.Coin) error
	SendCoinsWithAliasRouting(ctx sdk.Context, fromAddressAcc sdk.AccAddress, toAddressAcc sdk.AccAddress, coins sdk.Coins) error
	FundCommunityPoolWithAliasRouting(ctx sdk.Context, fromAddressAcc sdk.AccAddress, coins sdk.Coins) error
	SpendFromCommunityPoolWithAliasRouting(ctx sdk.Context, toAddressAcc sdk.AccAddress, coins sdk.Coins) error
	SendCoinsFromModuleToAccountWithAliasRouting(ctx sdk.Context, moduleName string, toAddressAcc sdk.AccAddress, coins sdk.Coins) error
	SendCoinsFromAccountToModuleWithAliasRouting(ctx sdk.Context, fromAddressAcc sdk.AccAddress, moduleName string, coins sdk.Coins) error
	GetBalanceWithAliasRouting(ctx sdk.Context, address sdk.AccAddress, denom string) (sdk.Coin, error)
}
