package tokenization_test

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/precompile/test/helpers"
	storetypes "github.com/cosmos/cosmos-sdk/store/v2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

func (suite *ERC3643TokenizationTestSuite) TestV35ExampleInitializationAndIssuance() {
	raw, err := os.ReadFile("testdata/v35_examples.json")
	suite.Require().NoError(err)
	var compiled map[string]struct {
		ABI json.RawMessage `json:"abi"`
		Bin string          `json:"bin"`
	}
	suite.Require().NoError(json.Unmarshal(raw, &compiled))
	for _, tc := range []struct {
		name, init, collection string
		constructor, initArgs  []interface{}
	}{
		{"RealEstateSecurityToken", "initializeCollection", "collectionId", []interface{}{"Real Estate", "RE", "Property", big.NewInt(100000), suite.AliceEVM}, []interface{}{big.NewInt(100)}},
		{"PrivateEquityToken", "initializeCollection", "collectionId", []interface{}{"Private Equity", "PE", "Fund", big.NewInt(100000), big.NewInt(1000), big.NewInt(0), suite.AliceEVM, suite.AliceEVM}, nil},
		{"TwoFactorSecurityToken", "initialize", "securityTokenCollectionId", []interface{}{"Two Factor", "TF", suite.AliceEVM}, []interface{}{big.NewInt(100)}},
		{"CarbonCreditToken", "initializeCollection", "collectionId", []interface{}{"Carbon", "CO2", "Verified", "Project", suite.AliceEVM}, nil},
	} {
		suite.Run(tc.name, func() {
			fixture := compiled[tc.name]
			contractABI, err := abi.JSON(bytes.NewReader(fixture.ABI))
			suite.Require().NoError(err)
			bytecode, err := hex.DecodeString(fixture.Bin)
			suite.Require().NoError(err)
			args, err := contractABI.Pack("", tc.constructor...)
			suite.Require().NoError(err)
			suite.Ctx = suite.Ctx.WithGasMeter(storetypes.NewInfiniteGasMeter())
			address, result, err := helpers.DeployContract(suite.Ctx, suite.EVMKeeper, suite.AliceKey, append(bytecode, args...), suite.getChainID())
			suite.Require().NoError(err)
			suite.Require().Empty(result.VmError)
			call := func(method string, view bool, args ...interface{}) []interface{} {
				suite.Ctx = suite.Ctx.WithGasMeter(storetypes.NewInfiniteGasMeter())
				ret, result, err := helpers.CallContractMethod(suite.Ctx, suite.EVMKeeper, suite.AliceKey, address, contractABI, method, args, suite.getChainID(), view)
				suite.Require().NoError(err)
				reason, _ := abi.UnpackRevert(result.Ret)
				suite.Require().Empty(result.VmError, "%s: %s", method, reason)
				values, err := contractABI.Methods[method].Outputs.Unpack(ret)
				suite.Require().NoError(err)
				return values
			}
			call(tc.init, false, tc.initArgs...)
			id := sdkmath.NewUintFromBigInt(call(tc.collection, true)[0].(*big.Int))
			collection, found := suite.TokenizationKeeper.GetCollectionFromStore(suite.Ctx, id)
			suite.Require().True(found)
			suite.Require().Len(collection.CollectionApprovals, 1)
			suite.Require().Equal("!Mint", collection.CollectionApprovals[0].FromListId)
			suite.Require().Equal(sdk.AccAddress(address.Bytes()).String(), collection.CollectionApprovals[0].InitiatedByListId)
			if tc.name == "RealEstateSecurityToken" {
				call("registerIdentity", false, suite.BobEVM)
				call("setAccreditedStatus", false, suite.BobEVM, true)
				call("transfer", false, suite.BobEVM, big.NewInt(1))
				suite.Require().Equal("99", call("balanceOf", true, suite.AliceEVM)[0].(*big.Int).String())
				suite.Require().Equal("1", call("balanceOf", true, suite.BobEVM)[0].(*big.Int).String())
			}
			if tc.name == "PrivateEquityToken" {
				call("setAccreditedInvestor", false, suite.BobEVM, true)
				call("onboardInvestor", false, suite.BobEVM, big.NewInt(1000), "A")
				call("processCapitalCall", false, suite.BobEVM, big.NewInt(1000))
				call("setAccreditedInvestor", false, suite.AliceEVM, true)
				call("onboardInvestor", false, suite.AliceEVM, big.NewInt(1000), "GP")
				call("releaseLockUp", false, suite.AliceEVM)
				call("transferWithApproval", false, suite.BobEVM, big.NewInt(1))
				suite.Require().Equal("98", call("balanceOf", true, suite.AliceEVM)[0].(*big.Int).String())
				suite.Require().Equal("2", call("balanceOf", true, suite.BobEVM)[0].(*big.Int).String())
			}
			if tc.name == "TwoFactorSecurityToken" {
				call("issueTwoFactor", false, suite.AliceEVM)
				suite.Require().Equal(true, call("hasValidTwoFactor", true, suite.AliceEVM)[0])
				call("setKYCStatus", false, suite.BobEVM, true)
				call("transfer", false, suite.BobEVM, big.NewInt(1))
			}
			if tc.name == "CarbonCreditToken" {
				call("createVintage", false, big.NewInt(2026), big.NewInt(suite.Ctx.BlockTime().Unix()+86400), "ipfs://vintage")
				call("verifySeller", false, suite.AliceEVM)
				call("verifyBuyer", false, suite.BobEVM)
				call("issueCredits", false, big.NewInt(2026), suite.AliceEVM, big.NewInt(10))
				call("transfer", false, suite.BobEVM, big.NewInt(2026), big.NewInt(1))
				suite.Require().Equal("9", call("balanceOf", true, suite.AliceEVM, big.NewInt(2026))[0].(*big.Int).String())
				suite.Require().Equal("1", call("balanceOf", true, suite.BobEVM, big.NewInt(2026))[0].(*big.Int).String())
			}
		})
	}
}

func TestV35ExampleFixtureSources(t *testing.T) {
	raw, err := os.ReadFile("testdata/v35_examples_sources.json")
	if err != nil {
		t.Fatal(err)
	}
	var sources map[string]string
	if err := json.Unmarshal(raw, &sources); err != nil {
		t.Fatal(err)
	}
	if len(sources) < 4 {
		t.Fatal("missing Solidity source provenance")
	}
	for path, expected := range sources {
		data, err := os.ReadFile(filepath.Join("../../../../..", path))
		if err != nil {
			t.Fatal(err)
		}
		if actual := fmt.Sprintf("%x", sha256.Sum256(data)); actual != expected {
			t.Errorf("stale Solidity fixture for %s: run python3 scripts/upgrade/test/compile_examples.py", path)
		}
	}
}
