package cmd

import (
	"strconv"

	cmtcfg "github.com/cometbft/cometbft/config"
	serverconfig "github.com/cosmos/cosmos-sdk/server/config"
	cosmosevmserverconfig "github.com/cosmos/evm/server/config"

	appparams "github.com/bitbadges/bitbadgeschain/app/params"
	evmtypes "github.com/cosmos/evm/x/vm/types"
)

// initCometBFTConfig helps to override default CometBFT Config values.
// return cmtcfg.DefaultConfig if no custom configuration is required for the application.
func initCometBFTConfig() *cmtcfg.Config {
	cfg := cmtcfg.DefaultConfig()

	// The EVM mempool is enabled on this chain, and cosmos/evm cross-validates
	// the two config files at startup: with it enabled, config.toml's
	// mempool.type must be "app" or the node refuses to boot with
	//
	//	EVM mempool enabled, but comet-bft has invalid config.toml:mempool.type
	//	(want 'app', got 'flood')
	//
	// CometBFT's own default is "flood", so leaving this alone made every
	// freshly initialised v34 node unstartable. The v33 -> v34 upgrade path did
	// not expose it, because there the pre-upgrade hook rewrites an existing
	// config.toml; nothing rewrites a config that `init` has just created. That
	// is the path taken by anyone adding a sentry, rebuilding a node for state
	// sync, joining the network after v34, or running a local dev chain.
	//
	// Enabling the mempool and defaulting its transport are two halves of one
	// decision, so they belong together. Operators who disable the EVM mempool
	// must set this back to "flood": the cross-check is symmetric and rejects
	// app-mempool-without-EVM-mempool too. Note that disabling it requires
	// passing --mempool.max-txs=-1 to `start` - editing app.toml alone is not
	// enough, because cosmos/evm's start command force-sets that flag to 0 and
	// a changed flag shadows the config file in viper (server/start.go:188-190).
	cfg.Mempool.Type = cmtcfg.MempoolTypeApp

	// these values put a higher strain on node memory
	// cfg.P2P.MaxNumInboundPeers = 100
	// cfg.P2P.MaxNumOutboundPeers = 40

	return cfg
}

// initAppConfig helps to override default appConfig template and configs.
// return "", nil if no custom configuration is required for the application.
func initAppConfig() (string, interface{}) {
	// Optionally allow the chain developer to overwrite the SDK's default
	// server config.
	srvCfg := serverconfig.DefaultConfig()
	// The SDK's default minimum gas price is set to "" (empty value) inside
	// app.toml. If left empty by validators, the node will halt on startup.
	// However, the chain developer can set a default app.toml value for their
	// validators here.
	//
	// In summary:
	// - if you leave srvCfg.MinGasPrices = "", all validators MUST tweak their
	//   own app.toml config,
	// - if you set srvCfg.MinGasPrices non-empty, validators CAN tweak their
	//   own app.toml to override, or use this default value.
	//
	// v35 aligns the node default with the fee-market floor the v35 upgrade
	// sets on chain (see app/upgrades/v35), so a fresh app.toml and the
	// chain agree. Validators may still override it.
	srvCfg.MinGasPrices = "10" + appparams.BaseCoinUnit

	// 0 ("enabled, unbounded") rather than the SDK default of -1 ("disabled").
	// The runtime enables the EVM mempool regardless of this file (see
	// initCometBFTConfig above); writing -1 here produced a generated
	// config.toml/app.toml pair that contradicted each other as written and
	// only booted because of upstream's force-set. The two files must agree on
	// their own.
	srvCfg.Mempool.MaxTxs = 0

	// Default to local dev EVM chain ID (90123) for the app.toml configuration.
	// The actual EVM chain ID used at runtime is set via build flags (ldflags):
	//   - Mainnet: 50024 (set via LDFLAGS_MAINNET)
	//   - Testnet: 50025 (set via LDFLAGS_TESTNET)
	//   - Local dev: 90123 (default if no build flag is set)
	// See app/evm.go (keeper initialization) and app/params/constants.go for details.
	evmChainID, err := strconv.ParseUint(appparams.EVMChainIDLocalDev, 10, 64)
	if err != nil {
		// Fallback to default if parsing fails
		evmChainID = evmtypes.DefaultEVMChainID
	}

	evmCfg := cosmosevmserverconfig.DefaultEVMConfig()
	evmCfg.EVMChainID = evmChainID

	customAppConfig := EVMAppConfig{
		Config:  *srvCfg,
		EVM:     *evmCfg,
		JSONRPC: *cosmosevmserverconfig.DefaultJSONRPCConfig(),
		TLS:     *cosmosevmserverconfig.DefaultTLSConfig(),
	}

	return EVMAppTemplate, customAppConfig
}

type EVMAppConfig struct {
	serverconfig.Config

	EVM     cosmosevmserverconfig.EVMConfig
	JSONRPC cosmosevmserverconfig.JSONRPCConfig
	TLS     cosmosevmserverconfig.TLSConfig
}

const EVMAppTemplate = serverconfig.DefaultConfigTemplate + cosmosevmserverconfig.DefaultEVMConfigTemplate
