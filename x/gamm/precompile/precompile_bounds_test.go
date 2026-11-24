package gamm

import (
	"math/big"
	"strings"
	"testing"

	"github.com/bitbadges/bitbadgeschain/pkg/evmcompat"
	gammkeeper "github.com/bitbadges/bitbadgeschain/x/gamm/keeper"
	gammtypes "github.com/bitbadges/bitbadgeschain/x/gamm/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/v2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/stretchr/testify/require"
)

func TestGammJSONBudgetRejectsDeepAndOversizedNumbers(t *testing.T) {
	p := NewPrecompile(gammkeeper.Keeper{}, nil)
	ctx := sdk.Context{}.WithGasMeter(storetypes.NewInfiniteGasMeter())
	caller := common.HexToAddress("0x1111111111111111111111111111111111111111")
	contract := vm.NewContract(caller, p.ContractAddress, nil, 1_000_000, nil)

	deep := strings.Repeat("[", evmcompat.DefaultMaxJSONDepth+2) + "1" + strings.Repeat("]", evmcompat.DefaultMaxJSONDepth+2)
	_, err := p.unmarshalMsgFromJSON(ctx, CreatePoolMethod, deep, contract)
	require.Error(t, err)

	hugeNum := `{"poolId":` + strings.Repeat("9", evmcompat.DefaultMaxJSONNumberLength+1) + `}`
	_, err = p.unmarshalQueryFromJSON(ctx, GetPoolMethod, hugeNum)
	require.Error(t, err)
}

func TestGammQueryPoolsPaginationCap(t *testing.T) {
	p := NewPrecompile(gammkeeper.Keeper{}, nil)
	ctx := sdk.Context{}.WithGasMeter(storetypes.NewInfiniteGasMeter())
	method := p.ABI.Methods[GetPoolsMethod]

	// Excessive limit is rejected before hitting the store
	_, err := p.HandleQuery(ctx, &method, `{"pagination":{"limit":1000000}}`)
	require.Error(t, err)
	require.Contains(t, err.Error(), "exceeds maximum allowed")
}

func TestGammSwapExactAmountInWithIBCTransferMethodOutputPacking(t *testing.T) {
	p := NewPrecompile(gammkeeper.Keeper{}, nil)
	method, ok := p.ABI.Methods[SwapExactAmountInWithIBCTransferMethod]
	require.True(t, ok)

	resp := &gammtypes.MsgSwapExactAmountInWithIBCTransferResponse{
		TokenOutAmount: gammtypes.OneShare,
	}
	tokenOut := resp.TokenOutAmount.BigInt()
	require.NoError(t, CheckUint256Overflow(tokenOut, "tokenOutAmount"))
	packed, err := method.Outputs.Pack(tokenOut)
	require.NoError(t, err)

	unpacked, err := method.Outputs.Unpack(packed)
	require.NoError(t, err)
	require.Len(t, unpacked, 1)
	require.Equal(t, 0, unpacked[0].(*big.Int).Cmp(tokenOut))
}
