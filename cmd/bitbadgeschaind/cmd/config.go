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
	// In this application, we set the min gas prices to 0.
	srvCfg.MinGasPrices = "0" + appparams.BaseCoinUnit

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
