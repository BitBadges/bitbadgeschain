package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	"encoding/hex"
	"github.com/bitbadges/bitbadgeschain/app"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/keeper"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/v2/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"testing"
)

// Valid curve signatures force full recovery; the wrong authorized signer prevents early success.
func BenchmarkV35ETHProofBatch(b *testing.B) {
	chain := app.Setup(false)
	key, err := ethcrypto.GenerateKey()
	if err != nil {
		b.Fatal(err)
	}
	signature, err := ethcrypto.Sign(ethcrypto.Keccak256([]byte("benchmark")), key)
	if err != nil {
		b.Fatal(err)
	}
	approval := &types.CollectionApproval{ApprovalId: "bench", ApprovalCriteria: &types.ApprovalCriteria{EthSignatureChallenges: []*types.ETHSignatureChallenge{{Signer: "0x1234567890123456789012345678901234567890", ChallengeTrackerId: "bench"}}}}
	transfer := &types.Transfer{}
	for i := 0; i < 100; i++ {
		transfer.EthSignatureProofs = append(transfer.EthSignatureProofs, &types.ETHSignatureProof{Nonce: "benchmark", Signature: hex.EncodeToString(signature)})
	}
	b.ReportAllocs()
	b.ResetTimer()
	var gas uint64
	for i := 0; i < b.N; i++ {
		ctx := chain.NewContext(false).WithGasMeter(storetypes.NewInfiniteGasMeter())
		_, err := chain.TokenizationKeeper.HandleETHSignatureChallenges(ctx, sdkmath.OneUint(), transfer, approval, keeper.TransferMetadata{InitiatedBy: "benchmark", ApprovalLevel: "collection"})
		if err == nil {
			b.Fatal("wrong signer unexpectedly accepted")
		}
		gas = ctx.GasMeter().GasConsumed()
		if gas < 100*types.ETHSignatureRecoveryGas {
			b.Fatal("proofs were not all metered")
		}
	}
	b.ReportMetric(float64(gas), "gas/op")
}
