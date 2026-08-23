package app

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/module"

	erc20 "github.com/cosmos/evm/x/erc20"
	erc20types "github.com/cosmos/evm/x/erc20/types"
	feemarket "github.com/cosmos/evm/x/feemarket"
	feemarkettypes "github.com/cosmos/evm/x/feemarket/types"
	evmmodule "github.com/cosmos/evm/x/vm"
	evmtypes "github.com/cosmos/evm/x/vm/types"

	appparams "github.com/bitbadges/bitbadgeschain/app/params"
)

// The EVM modules (x/vm, x/feemarket, x/erc20) are registered at runtime in
// registerEVMModules via app.RegisterModules, not through the depinject app
// config. The module.BasicManager that depinject builds therefore does not know
// about them, so `bitbadgeschaind init` emitted a genesis file with no "vm",
// "feemarket" or "erc20" keys at all.
//
// That mattered because the SDK's module manager skips InitGenesis for any
// module with no genesis data. x/vm's InitGenesis is what seeds the process
// global EVM config (SetGlobalConfigVariables, behind the module's own
// sync.Once). With it skipped, the globals stayed unset through genesis and the
// very first gentx failed in the EVM ante handler with
// "global evmCoinInfo is not set yet!".
//
// RegisterEVMModuleBasics adds the three module basics so `init` writes their
// genesis. It is applied to a *copy* of the basic manager used only by InitCmd,
// so the interface/amino registrations these modules perform elsewhere are not
// duplicated into the shared manager.
func RegisterEVMModuleBasics(bm module.BasicManager) {
	bm[evmtypes.ModuleName] = bitbadgesEVMModuleBasic{}
	bm[feemarkettypes.ModuleName] = feemarket.AppModuleBasic{}
	bm[erc20types.ModuleName] = erc20.AppModuleBasic{}
}

// bitbadgesEVMModuleBasic overrides x/vm's default genesis so a freshly
// initialized chain carries BitBadges' denominations instead of upstream's
// "aatom" placeholder.
//
// This is load-bearing, not cosmetic. x/vm InitGenesis calls SetParams with the
// genesis params and then LoadEvmCoinInfo, which resolves the coin info from
// params.EvmDenom plus the bank denom metadata for that denom. With upstream's
// default the lookup is for "aatom", which has no bank metadata on this chain,
// so InitGenesis panics with "denom metadata aatom could not be found".
//
// ExtendedDenomOptions must also be set: LoadEvmCoinInfo rejects a nil value for
// any chain whose display unit is not 18 decimals. BADGE is 9-decimal (ubadge),
// with abadge as the 18-decimal denom the EVM operates in — the mechanism that
// replaced x/precisebank in cosmos/evm v0.7.
type bitbadgesEVMModuleBasic struct {
	evmmodule.AppModuleBasic
}

func (bitbadgesEVMModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	gs := evmtypes.DefaultGenesisState()
	gs.Params.EvmDenom = appparams.BaseCoinUnit
	gs.Params.ExtendedDenomOptions = &evmtypes.ExtendedDenomOptions{
		ExtendedDenom: appparams.ExtendedCoinUnit,
	}
	return cdc.MustMarshalJSON(gs)
}
