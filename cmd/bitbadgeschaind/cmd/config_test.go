package cmd

import (
	"testing"

	cmtcfg "github.com/cometbft/cometbft/config"
	cosmosevmserverconfig "github.com/cosmos/evm/server/config"
)

// TestInitCometBFTConfigAgreesWithTheEVMMempoolDefault pins that the two
// config files this binary generates can actually run together.
//
// `bitbadgeschaind init` writes both config.toml (from initCometBFTConfig) and
// app.toml (from initAppConfig). cosmos/evm v0.7 cross-validates them at
// startup: if the EVM mempool is enabled it demands config.toml's
// mempool.type = "app", and refuses to start otherwise.
//
// v34 enables the EVM mempool by default — cosmos/evm's own start command
// registers mempool.max-txs with a default of 0 and then force-sets it, so
// app.toml cannot turn it off (see server/start.go:188-190). But
// initCometBFTConfig returned CometBFT's stock default, whose mempool type is
// "flood". The result was that a freshly initialised v34 node could not start
// at all:
//
//	EVM mempool enabled, but comet-bft has invalid config.toml:mempool.type
//	(want 'app', got 'flood'): error in app.toml
//
// The upgrade path hid this, because there the pre-upgrade hook rewrites an
// existing config.toml. What it did NOT cover is every path that creates a
// *new* config: a validator adding a sentry, an operator rebuilding a node for
// state sync, anyone joining the network after v34, or a local dev chain.
// Those all run `init` and would have hit a node that refuses to boot with no
// documented remedy.
//
// This asserts the invariant rather than the literal string, by running the
// same cross-check the node runs at startup.
func TestInitCometBFTConfigAgreesWithTheEVMMempoolDefault(t *testing.T) {
	cmtConfig := initCometBFTConfig()

	_, rawAppConfig := initAppConfig()
	appConfig, ok := rawAppConfig.(EVMAppConfig)
	if !ok {
		t.Fatalf("initAppConfig returned %T, want EVMAppConfig", rawAppConfig)
	}

	evmConfig := &cosmosevmserverconfig.Config{
		Config:  appConfig.Config,
		EVM:     appConfig.EVM,
		JSONRPC: appConfig.JSONRPC,
		TLS:     appConfig.TLS,
	}

	// Model the value the node actually runs with, not the one app.toml holds.
	// cosmos/evm's start command force-sets this flag to 0, which takes
	// precedence over the config file, so the EVM mempool is enabled on a
	// default node whatever app.toml says.
	evmConfig.Mempool.MaxTxs = 0

	if err := cosmosevmserverconfig.ValidateCrossConfig(cmtConfig, evmConfig); err != nil {
		t.Fatalf("the config.toml and app.toml that `init` generates cannot run together: %v\n"+
			"config.toml mempool.type = %q. A node created with `bitbadgeschaind init` will "+
			"refuse to start, and unlike the upgrade path there is no pre-upgrade hook to fix it.",
			err, cmtConfig.Mempool.Type)
	}
}

// TestInitCometBFTConfigUsesTheAppMempool is the narrow, readable companion to
// the cross-check above: it names the value rather than deriving it, so a
// reader who breaks the test can see immediately what is expected.
func TestInitCometBFTConfigUsesTheAppMempool(t *testing.T) {
	if got := initCometBFTConfig().Mempool.Type; got != cmtcfg.MempoolTypeApp {
		t.Fatalf("initCometBFTConfig().Mempool.Type = %q, want %q", got, cmtcfg.MempoolTypeApp)
	}
}
