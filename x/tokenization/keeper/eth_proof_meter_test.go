package keeper_test

import (
	"encoding/hex"
	"fmt"
	"strings"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/keeper"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/v2/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
)

func (suite *TestSuite) TestETHProofAttemptsConsumeGas() {
	ctx := suite.ctx.WithGasMeter(storetypes.NewInfiniteGasMeter())
	approval := &types.CollectionApproval{ApprovalId: "proof-meter", ApprovalCriteria: &types.ApprovalCriteria{EthSignatureChallenges: []*types.ETHSignatureChallenge{{Signer: "0x1234567890123456789012345678901234567890", ChallengeTrackerId: "challenge"}}}}
	transfer := &types.Transfer{EthSignatureProofs: []*types.ETHSignatureProof{{Nonce: "ticket", Signature: strings.Repeat("00", 65)}, {Nonce: "ticket2", Signature: strings.Repeat("00", 65)}}}
	_, err := suite.app.TokenizationKeeper.HandleETHSignatureChallenges(ctx, sdkmath.OneUint(), transfer, approval, keeper.TransferMetadata{InitiatedBy: bob, ApprovalLevel: "collection"})
	suite.Require().Error(err)
	suite.Require().GreaterOrEqual(ctx.GasMeter().GasConsumed(), uint64(6000))
}

func (suite *TestSuite) TestETHSignatureRejectsLegacyFormat() {
	secret, signer, err := generateTestETHPrivateKey()
	suite.Require().NoError(err)
	key, err := ethcrypto.HexToECDSA(secret)
	suite.Require().NoError(err)
	message := "nonce-" + bob + "-1--collection-approval-challenge"
	digest := ethcrypto.Keccak256([]byte(fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message)))
	signature, err := ethcrypto.Sign(digest, key)
	suite.Require().NoError(err)
	approval := &types.CollectionApproval{ApprovalId: "approval", ApprovalCriteria: &types.ApprovalCriteria{EthSignatureChallenges: []*types.ETHSignatureChallenge{{Signer: signer, ChallengeTrackerId: "challenge"}}}}
	transfer := &types.Transfer{EthSignatureProofs: []*types.ETHSignatureProof{{Nonce: "nonce", Signature: hex.EncodeToString(signature)}}}
	_, err = suite.app.TokenizationKeeper.HandleETHSignatureChallenges(suite.ctx, sdkmath.OneUint(), transfer, approval, keeper.TransferMetadata{InitiatedBy: bob, ApprovalLevel: "collection"})
	suite.Require().Error(err, "legacy ambiguous messages must not remain an accepted fallback")
}
