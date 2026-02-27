//go:build test
// +build test

package ibc

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	sdkmath "cosmossdk.io/math"

	clienthelpers "cosmossdk.io/client/v2/helpers"
	"cosmossdk.io/log"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/crypto/ed25519"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmtypes "github.com/cometbft/cometbft/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	simapp "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	ibctesting "github.com/cosmos/ibc-go/v10/testing"

	evmtypes "github.com/cosmos/evm/x/vm/types"

	"github.com/bitbadges/bitbadgeschain/app"
)

// DefaultTestingAppInit is the default app initializer for ibctesting
var DefaultTestingAppInit = SetupBitBadgesTestingApp

// DefaultConsensusParams defines the consensus params used for IBC testing
var DefaultConsensusParams = &tmproto.ConsensusParams{
	Block: &tmproto.BlockParams{
		MaxBytes: 200000,
		MaxGas:   2000000,
	},
	Evidence: &tmproto.EvidenceParams{
		MaxAgeNumBlocks: 302400,
		MaxAgeDuration:  504 * time.Hour, // 3 weeks is the max duration
		MaxBytes:        10000,
	},
	Validator: &tmproto.ValidatorParams{
		PubKeyTypes: []string{
			tmtypes.ABCIPubKeyTypeEd25519,
		},
	},
}

func init() {
	// Set the default app init function for ibctesting
	// Note: ibc-go's setupWithGenesisValSet overwrites the bank genesis with empty metadata.
	// We handle this in app.go by setting bank metadata in the InitChainer before
	// module InitGenesis runs.
	ibctesting.DefaultTestingAppInit = func() (ibctesting.TestingApp, map[string]json.RawMessage) {
		return SetupBitBadgesTestingApp()
	}
}

// SetupBitBadgesTestingApp returns a BitBadges app and genesis state for use with ibctesting
func SetupBitBadgesTestingApp() (ibctesting.TestingApp, map[string]json.RawMessage) {
	// Reset EVM config for testing - required when running parallel tests
	// because the EVM module uses global state that persists between test runs
	evmtypes.NewEVMConfigurator().ResetTestConfig()

	// Create a mock validator using ed25519 (IBC v10 removed testing/mock package)
	privKey := ed25519.GenPrivKey()
	pubKey := privKey.PubKey()

	// create validator set with single validator
	validator := tmtypes.NewValidator(pubKey, 824639203360100)
	valSet := tmtypes.NewValidatorSet([]*tmtypes.Validator{validator})

	// generate genesis account with substantial funds for testing
	senderPrivKey := secp256k1.GenPrivKey()
	acc := authtypes.NewBaseAccount(senderPrivKey.PubKey().Address().Bytes(), senderPrivKey.PubKey(), 0, 0)
	balance := banktypes.Balance{
		Address: acc.GetAddress().String(),
		Coins: sdk.NewCoins(
			sdk.NewCoin("ubadge", sdkmath.NewInt(100000000000000)),
			sdk.NewCoin("ustake", sdkmath.NewInt(100000000000000)),
		),
	}

	origDefault, err := clienthelpers.GetNodeHomeDirectory("." + app.Name)
	if err != nil {
		panic(err)
	}

	randomHomeDir := origDefault + "/ibc_test_" + fmt.Sprint(rand.Int63n(1000000))

	db := dbm.NewMemDB()
	bitbadgesApp, err := app.New(log.NewNopLogger(), db, nil, true, simapp.NewAppOptionsWithFlagHome(randomHomeDir))
	if err != nil {
		panic(err)
	}

	// Get default genesis state and merge with validator set
	genesisState := bitbadgesApp.DefaultGenesis()
	genesisState = app.GenesisStateWithValSet(bitbadgesApp, genesisState, valSet, []authtypes.GenesisAccount{acc}, balance)

	// Ensure bank module has ubadge metadata for EVM module
	genesisState = ensureBankMetadata(bitbadgesApp, genesisState)

	return bitbadgesApp, genesisState
}

// ensureBankMetadata ensures the bank module genesis has ubadge denom metadata
func ensureBankMetadata(bitbadgesApp *app.App, genesisState map[string]json.RawMessage) map[string]json.RawMessage {
	var bankGenesis banktypes.GenesisState
	if bankGenesisBytes, ok := genesisState[banktypes.ModuleName]; ok {
		bitbadgesApp.AppCodec().MustUnmarshalJSON(bankGenesisBytes, &bankGenesis)
	}

	// Check if ubadge metadata already exists
	hasUbadgeMetadata := false
	for _, metadata := range bankGenesis.DenomMetadata {
		if metadata.Base == "ubadge" {
			hasUbadgeMetadata = true
			break
		}
	}

	// Add ubadge metadata if it doesn't exist
	if !hasUbadgeMetadata {
		ubadgeMetadata := banktypes.Metadata{
			Description: "The native token of BitBadges Chain",
			DenomUnits: []*banktypes.DenomUnit{
				{Denom: "ubadge", Exponent: 0},
				{Denom: "badge", Exponent: 9},
			},
			Base:    "ubadge",
			Display: "badge",
			Name:    "Badge",
			Symbol:  "BADGE",
		}
		bankGenesis.DenomMetadata = append(bankGenesis.DenomMetadata, ubadgeMetadata)
	}

	genesisState[banktypes.ModuleName] = bitbadgesApp.AppCodec().MustMarshalJSON(&bankGenesis)

	// Also ensure EVM genesis params use ubadge as the EVM denom
	if evmGenesisBytes, ok := genesisState["evm"]; ok {
		var evmGenesis evmtypes.GenesisState
		bitbadgesApp.AppCodec().MustUnmarshalJSON(evmGenesisBytes, &evmGenesis)
		evmGenesis.Params.EvmDenom = "ubadge"
		genesisState["evm"] = bitbadgesApp.AppCodec().MustMarshalJSON(&evmGenesis)
	}

	return genesisState
}

// SetupBitBadgesTestingAppWithGenesis initializes a testing app with custom genesis state
func SetupBitBadgesTestingAppWithGenesis(genState map[string]json.RawMessage) ibctesting.TestingApp {
	// Reset EVM config for testing
	evmtypes.NewEVMConfigurator().ResetTestConfig()

	origDefault, err := clienthelpers.GetNodeHomeDirectory("." + app.Name)
	if err != nil {
		panic(err)
	}

	randomHomeDir := origDefault + "/ibc_test_" + fmt.Sprint(rand.Int63n(1000000))

	db := dbm.NewMemDB()
	bitbadgesApp, err := app.New(log.NewNopLogger(), db, nil, true, simapp.NewAppOptionsWithFlagHome(randomHomeDir))
	if err != nil {
		panic(err)
	}

	// Marshal genesis state
	stateBytes, err := json.MarshalIndent(genState, "", " ")
	if err != nil {
		panic(err)
	}

	// Initialize the chain
	bitbadgesApp.InitChain(
		&abci.RequestInitChain{
			Validators:      []abci.ValidatorUpdate{},
			ConsensusParams: DefaultConsensusParams,
			AppStateBytes:   stateBytes,
		},
	)

	return bitbadgesApp
}
