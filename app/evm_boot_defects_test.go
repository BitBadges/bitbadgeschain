package app

import (
	"math/big"
	"testing"

	clienthelpers "cosmossdk.io/client/v2/helpers"
	log "cosmossdk.io/log/v2"
	sdkmath "cosmossdk.io/math"
	dbm "github.com/cosmos/cosmos-db"
	sdkserver "github.com/cosmos/cosmos-sdk/server"
	simapp "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	precisebanktypes "github.com/cosmos/evm/contrib/x/precisebank/types"
	erc20types "github.com/cosmos/evm/x/erc20/types"
	feemarkettypes "github.com/cosmos/evm/x/feemarket/types"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	appparams "github.com/bitbadges/bitbadgeschain/app/params"
)

// Regression tests for the boot defects found while bringing this chain onto
// cosmos/evm v0.7. Each one nearly shipped, and each was invisible to the suite
// as it stood — see the individual comments for why.

// ubadgeToWeiFactor is 10^(18-9): the chain is 9-decimal (ubadge) and the EVM
// operates in the 18-decimal extended denom (abadge).
var ubadgeToWeiFactor = new(big.Int).Exp(big.NewInt(10), big.NewInt(18-appparams.BaseCoinDecimals), nil)

// TestEVMBalanceEqualsBankBalanceTimes1e9 is the regression guard for the
// precisebank removal.
//
// Dropping x/precisebank made every EVM balance read as zero: eth_getBalance
// returned 0 for accounts holding real ubadge, so no EVM transaction could ever
// pay for gas and the chain was EVM-dead while looking healthy from Cosmos.
//
// Nothing in the suite caught it. app/precisebank_test.go "tests" the
// conversion by recomputing 1 * 10^9 in the test body and asserting it equals
// 1e9 — arithmetic on local variables that never touches a keeper, so it passes
// with the module removed, present, or replaced by a stub.
//
// This asserts the real thing: fund an account through x/bank, then read it back
// through exactly the code path eth_getBalance uses. rpc/backend.GetBalance
// calls QueryClient.Balance, which is x/vm's Balance gRPC handler, which is what
// is called here. If the extended-denom bridge is broken, this returns 0 and the
// test fails.
func TestEVMBalanceEqualsBankBalanceTimes1e9(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContext(false)

	addr := newEVMAddress(t)
	accAddr := sdk.AccAddress(addr.Bytes())

	// A deliberately awkward amount: not a power of ten, so a conversion that
	// merely returns something nonzero cannot pass by accident.
	const ubadgeAmount = int64(123_456_789)
	funds := sdk.NewCoins(sdk.NewCoin(appparams.BaseCoinUnit, sdkmath.NewInt(ubadgeAmount)))
	require.NoError(t, app.BankKeeper.MintCoins(ctx, "mint", funds))
	require.NoError(t, app.BankKeeper.SendCoinsFromModuleToAccount(ctx, "mint", accAddr, funds))

	// Sanity: the bank side really does hold what we think it does.
	bankBalance := app.BankKeeper.GetBalance(ctx, accAddr, appparams.BaseCoinUnit)
	require.Equal(t, sdkmath.NewInt(ubadgeAmount), bankBalance.Amount,
		"precondition: bank must hold the funded ubadge")

	want := new(big.Int).Mul(big.NewInt(ubadgeAmount), ubadgeToWeiFactor)

	// The eth_getBalance path.
	res, err := app.EVMKeeper.Balance(ctx, &evmtypes.QueryBalanceRequest{Address: addr.Hex()})
	require.NoError(t, err)
	got, ok := new(big.Int).SetString(res.Balance, 10)
	require.True(t, ok, "unparseable balance %q", res.Balance)
	require.Zero(t, want.Cmp(got),
		"eth_getBalance must report bank ubadge x 10^%d: want %s, got %s (a zero here means the "+
			"extended-denom bridge is gone and the chain is EVM-dead)",
		18-appparams.BaseCoinDecimals, want, got)

	// The in-EVM path (StateDB / BALANCE opcode) must agree with the RPC path.
	// They are different functions in x/vm and only one of them was ever
	// exercised before.
	stateDBBalance := app.EVMKeeper.GetBalance(ctx, addr)
	require.NotNil(t, stateDBBalance, "GetBalance returned nil")
	require.Zero(t, want.Cmp(stateDBBalance.ToBig()),
		"the StateDB balance must match eth_getBalance: want %s, got %s", want, stateDBBalance)
}

// TestEVMBalanceIsZeroForUnfundedAccount is the negative control for the test
// above: it proves the assertion can distinguish funded from unfunded, rather
// than the query returning something nonzero regardless.
func TestEVMBalanceIsZeroForUnfundedAccount(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContext(false)

	addr := newEVMAddress(t)

	res, err := app.EVMKeeper.Balance(ctx, &evmtypes.QueryBalanceRequest{Address: addr.Hex()})
	require.NoError(t, err)
	require.Equal(t, "0", res.Balance, "an unfunded account must report zero")
}

// TestRegisterEVMModuleBasicsDefaultGenesis guards the `init` path.
//
// Without RegisterEVMModuleBasics, `bitbadgeschaind init` writes a genesis with
// no vm/feemarket/erc20 app_state at all. With the stock x/vm basic instead of
// the BitBadges override it writes upstream's "aatom" placeholder, and
// InitGenesis then panics with "denom metadata aatom could not be found",
// because this chain has no bank metadata for aatom.
//
// Both failures only surface when someone runs `init` on a fresh machine, which
// no test did.
func TestRegisterEVMModuleBasicsDefaultGenesis(t *testing.T) {
	app := Setup(false)

	bm := module.BasicManager{}
	RegisterEVMModuleBasics(bm)

	// x/vm's module name is "evm", not "vm" — the genesis key an operator sees
	// in genesis.json. Using the constants rather than literals keeps that
	// straight.
	evmModules := []string{evmtypes.ModuleName, feemarkettypes.ModuleName, erc20types.ModuleName}

	for _, name := range evmModules {
		require.Contains(t, bm, name, "`init` must write app_state for %q", name)
	}

	genesis := bm.DefaultGenesis(app.appCodec)
	for _, name := range evmModules {
		require.Contains(t, genesis, name, "default genesis is missing the %q section", name)
	}

	// Decoded with the app codec, not encoding/json: genesis is proto-JSON and
	// enum fields are name strings that encoding/json cannot handle.
	var vmGenesis evmtypes.GenesisState
	app.appCodec.MustUnmarshalJSON(genesis[evmtypes.ModuleName], &vmGenesis)

	require.Equal(t, appparams.BaseCoinUnit, vmGenesis.Params.EvmDenom,
		"evm_denom must be %q, not upstream's placeholder; a wrong value panics InitGenesis "+
			"with \"denom metadata ... could not be found\"", appparams.BaseCoinUnit)

	require.NotNil(t, vmGenesis.Params.ExtendedDenomOptions,
		"extended_denom_options must be set: LoadEvmCoinInfo rejects nil for any chain whose "+
			"display unit is not 18 decimals, and this chain is 9-decimal")
	require.Equal(t, appparams.ExtendedCoinUnit, vmGenesis.Params.ExtendedDenomOptions.ExtendedDenom,
		"extended_denom must be %q", appparams.ExtendedCoinUnit)

	// Upstream's placeholder must not survive anywhere in the EVM sections.
	for _, name := range evmModules {
		require.NotContains(t, string(genesis[name]), "aatom",
			"the %q section still carries upstream's aatom placeholder", name)
	}
}

// TestGenesisModuleOrderEVMBeforeGenutil pins the ordering constraint behind the
// "global evmCoinInfo is not set yet!" boot failure.
//
// x/vm's InitGenesis is what installs the process-global EVM coin config.
// genutil's InitGenesis delivers the gentxs, which on this chain can reach EVM
// code. If genutil runs first, the chain dies at genesis with an error that
// names a global variable and points at nothing actionable.
//
// The constraint is invisible in the source — it is "index of evm < index of
// genutil" in a 30-element list that people reorder for unrelated reasons — so
// it gets a test rather than a comment.
func TestGenesisModuleOrderEVMBeforeGenutil(t *testing.T) {
	index := map[string]int{}
	for i, name := range genesisModuleOrder {
		require.NotContains(t, index, name, "module %q appears twice in genesisModuleOrder", name)
		index[name] = i
	}

	for _, name := range []string{
		banktypes.ModuleName, feemarkettypes.ModuleName, evmtypes.ModuleName,
		erc20types.ModuleName, precisebanktypes.ModuleName, genutiltypes.ModuleName,
	} {
		require.Contains(t, index, name, "genesisModuleOrder is missing %q", name)
	}

	require.Less(t, index[evmtypes.ModuleName], index[genutiltypes.ModuleName],
		"x/vm InitGenesis must run before genutil: genutil delivers gentxs that can reach EVM "+
			"code, and x/vm is what sets the process-global EVM coin config. The reverse order "+
			"fails genesis with \"global evmCoinInfo is not set yet!\"")

	require.Less(t, index[banktypes.ModuleName], index[evmtypes.ModuleName],
		"x/vm resolves its coin info from the bank denom metadata for the EVM denom, so bank "+
			"must initialise first")

	require.Less(t, index[feemarkettypes.ModuleName], index[evmtypes.ModuleName],
		"feemarket must initialise before vm")

	require.Less(t, index[evmtypes.ModuleName], index[erc20types.ModuleName],
		"erc20 depends on the EVM keeper, so it must initialise after vm")

	require.Less(t, index[evmtypes.ModuleName], index[precisebanktypes.ModuleName],
		"precisebank's InitGenesis needs the EVM keeper")
}

// TestEVMMempoolHandlersInstalled proves the app-side mempool is actually wired
// under the configuration a node runs by default.
//
// The failure this guards against produced no error at all: with the handlers
// missing, cosmos/evm's cross-config check still saw the EVM mempool as enabled
// and required mempool.type = "app", so CometBFT routed to an app mempool that
// had no ReapTxs handler. The node produced empty blocks forever and logged
// "ReapTxs handler not set" once a block, while every health check said it was
// fine.
//
// The app is built through the same New() the daemon uses, so if
// configureEVMMempool silently bailed out, EVMMempool would be nil here.
//
// mempool.max-txs is set explicitly rather than left to Setup's default. Setup
// disables the mempool in test apps because it leaks goroutines (see the note
// in test_helpers.go), so this test has to ask for the node configuration it is
// making a claim about. 0 is the value cosmos/evm's own start command sets for
// that flag - see server/start.go, which both registers the flag with a default
// of 0 and then force-sets it - so this is what a default node runs with.
func TestEVMMempoolHandlersInstalled(t *testing.T) {
	app := SetupWithAppOptions(false, map[string]interface{}{
		sdkserver.FlagMempoolMaxTxs: 0,
	})

	require.NotNil(t, app.EVMMempool,
		"the EVM mempool must be constructed under the default configuration; a nil value means "+
			"configureEVMMempool returned early and no Insert/Reap/CheckTx handler is installed, "+
			"which produces a node that mints empty blocks forever")

	require.NotNil(t, app.Mempool(), "BaseApp must have the mempool set")
	require.Same(t, any(app.EVMMempool), any(app.Mempool()),
		"BaseApp's mempool must be the EVM mempool that was just built, not some other instance")
}

// TestConfigureEVMMempoolErrorIsFatal pins that a mempool configuration failure
// stops the node instead of being logged and stepped over.
//
// The previous code logged "failed to configure EVM mempool" and continued, on
// the stated grounds that the mempool is optional. It is not optional once the
// app-side mempool is expected: continuing produces precisely the silent
// dead-node symptom described on TestEVMMempoolHandlersInstalled. A node that
// refuses to start is recoverable; a node that mints empty blocks is not
// noticed.
//
// ValidateReapBounds is driven by operator-supplied app.toml values, so this is
// reachable from a plain misconfiguration, not just from a bug.
func TestConfigureEVMMempoolErrorIsFatal(t *testing.T) {
	// ValidateReapBounds rejects an admission cap larger than the matching reap
	// cap. CometBFT's default mempool.max_tx_bytes is 1 MiB, so a reap_max_bytes
	// of 1 is a plain operator misconfiguration that trips it.
	//
	// New() is called directly rather than through SetupWithAppOptions because
	// that helper panics on error, and the whole point here is that an error is
	// returned rather than swallowed.
	evmtypes.NewEVMConfigurator().ResetTestConfig()
	home, dirErr := clienthelpers.GetNodeHomeDirectory("." + Name)
	require.NoError(t, dirErr)

	appOpts := overlaidAppOptions{
		base: simapp.NewAppOptionsWithFlagHome(home + "/test_mempool_fatal"),
		overlay: map[string]interface{}{
			"mempool.reap_max_bytes": int64(1),
		},
	}

	_, err := New(log.NewNopLogger(), dbm.NewMemDB(), true, appOpts)
	require.Error(t, err,
		"a mempool configuration error must abort app startup; swallowing it leaves BaseApp with "+
			"no Reap handler and the node mints empty blocks forever")
	require.Contains(t, err.Error(), "EVM mempool",
		"the error must name the mempool so an operator can act on it, got: %v", err)
}

// newEVMAddress generates a fresh EVM address.
func newEVMAddress(t *testing.T) ethcommon.Address {
	t.Helper()
	key, err := ethcrypto.GenerateKey()
	require.NoError(t, err)
	return ethcrypto.PubkeyToAddress(key.PublicKey)
}
