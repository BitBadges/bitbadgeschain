package keeper_test

import (
	"github.com/bitbadges/bitbadgeschain/x/managersplitter/types"
	tokenizationtypes "github.com/bitbadges/bitbadgeschain/x/tokenization/types"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// adminOnlyFieldCases lists the UniversalUpdateCollection fields that only the
// splitter admin may use, each applied on top of a metadata update that the
// executor is otherwise approved for.
func adminOnlyFieldCases() map[string]func(msg *tokenizationtypes.MsgUniversalUpdateCollection) {
	return map[string]func(msg *tokenizationtypes.MsgUniversalUpdateCollection){
		"create collection (zero id)": func(msg *tokenizationtypes.MsgUniversalUpdateCollection) {
			msg.CollectionId = sdkmath.NewUint(0)
		},
		"mint escrow coins to transfer": func(msg *tokenizationtypes.MsgUniversalUpdateCollection) {
			msg.MintEscrowCoinsToTransfer = []*sdk.Coin{{Denom: "ubadge", Amount: sdkmath.NewInt(1)}}
		},
		"invariants": func(msg *tokenizationtypes.MsgUniversalUpdateCollection) {
			msg.Invariants = &tokenizationtypes.InvariantsAddObject{NoCustomOwnershipTimes: true}
		},
		// DefaultBalances is only read when a collection is created, so the
		// create path is where it must be admin-only.
		"default balances on create": func(msg *tokenizationtypes.MsgUniversalUpdateCollection) {
			msg.CollectionId = sdkmath.NewUint(0)
			msg.DefaultBalances = &tokenizationtypes.UserBalanceStore{AutoApproveAllIncomingTransfers: true}
		},
	}
}

func (suite *TestSuite) createSplitterWithMetadataApprover(approver string) string {
	wctx := sdk.WrapSDKContext(suite.ctx)
	perms := GetDefaultPermissions()
	perms.CanUpdateCollectionMetadata.ApprovedAddresses = []string{approver}
	res, err := CreateManagerSplitter(suite, wctx, &types.MsgCreateManagerSplitter{Admin: bob, Permissions: perms})
	suite.Require().NoError(err)
	return res.Address
}

func metadataUpdateMsg(executor, splitter string) *types.MsgExecuteUniversalUpdateCollection {
	return &types.MsgExecuteUniversalUpdateCollection{
		Executor:               executor,
		ManagerSplitterAddress: splitter,
		UniversalUpdateCollectionMsg: &tokenizationtypes.MsgUniversalUpdateCollection{
			Creator:                  executor,
			CollectionId:             sdkmath.NewUint(1),
			UpdateCollectionMetadata: true,
			CollectionMetadata:       &tokenizationtypes.CollectionMetadata{Uri: "https://example.com"},
		},
	}
}

// TestAdminOnlyFields_NonAdminDenied checks that an executor approved only
// for metadata updates cannot use the admin-only fields.
func (suite *TestSuite) TestAdminOnlyFields_NonAdminDenied() {
	splitter := suite.createSplitterWithMetadataApprover(alice)

	for name, apply := range adminOnlyFieldCases() {
		msg := metadataUpdateMsg(alice, splitter)
		apply(msg.UniversalUpdateCollectionMsg)

		_, err := ExecuteUniversalUpdateCollection(suite, sdk.WrapSDKContext(suite.ctx), msg)
		suite.Require().Error(err, "%s: non-admin must be denied", name)
		suite.Require().Contains(err.Error(), "only admin", "%s: error must name the admin restriction", name)
	}
}

// TestAdminOnlyFields_AdminPassesPermissionCheck checks that the admin is not
// stopped by the splitter's own permission gate for the same fields.
func (suite *TestSuite) TestAdminOnlyFields_AdminPassesPermissionCheck() {
	splitter := suite.createSplitterWithMetadataApprover(alice)

	for name, apply := range adminOnlyFieldCases() {
		msg := metadataUpdateMsg(bob, splitter)
		apply(msg.UniversalUpdateCollectionMsg)

		_, err := ExecuteUniversalUpdateCollection(suite, sdk.WrapSDKContext(suite.ctx), msg)
		if err != nil {
			suite.Require().NotContains(err.Error(), "only admin", "%s: admin must pass the splitter permission check", name)
		}
	}
}
