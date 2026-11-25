package tokenization

import (
	"strings"
	"testing"

	"github.com/bitbadges/bitbadgeschain/pkg/evmcompat"
	tokenizationkeeper "github.com/bitbadges/bitbadgeschain/x/tokenization/keeper"
	storetypes "github.com/cosmos/cosmos-sdk/store/v2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestExecuteMultipleRequiredGasIsConstant(t *testing.T) {
	p := NewPrecompile(&tokenizationkeeper.Keeper{}, nil)
	method := p.ABI.Methods[ExecuteMultipleMethod]

	type batchEntry struct {
		MessageType string `abi:"messageType"`
		MsgJson     string `abi:"msgJson"`
	}

	smallMsg, err := p.ABI.Pack(ExecuteMultipleMethod, []batchEntry{
		{MessageType: "/tokenization.MsgTransferTokens", MsgJson: `{}`},
	})
	require.NoError(t, err)

	largeBatch := make([]batchEntry, 50)
	for i := range largeBatch {
		largeBatch[i] = batchEntry{MessageType: "/tokenization.MsgTransferTokens", MsgJson: `{"collectionId":"1"}`}
	}
	largeMsg, err := p.ABI.Pack(ExecuteMultipleMethod, largeBatch)
	require.NoError(t, err)

	// RequiredGas computes base + linear calldata chunk gas + tx buffer without ABI unpacking.
	const txGasBuffer uint64 = 200_000
	expectedSmall := GasExecuteMultipleBase + uint64(len(smallMsg)/32)*GasPerInputChunk + txGasBuffer
	expectedLarge := GasExecuteMultipleBase + uint64(len(largeMsg)/32)*GasPerInputChunk + txGasBuffer
	require.Equal(t, expectedSmall, p.RequiredGas(smallMsg))
	require.Equal(t, expectedLarge, p.RequiredGas(largeMsg))
	require.Equal(t, ExecuteMultipleMethod, method.Name)

	// Malformed ABI tail must not cause panic or expensive unpacking in RequiredGas.
	malformed := append(append([]byte{}, method.ID...), make([]byte, 128)...)
	require.Equal(t, GasExecuteMultipleBase+uint64(len(malformed)/32)*GasPerInputChunk+txGasBuffer, p.RequiredGas(malformed))
}

func TestJSONBudgetRejectsOversizedNumbers(t *testing.T) {
	ctx := sdk.Context{}.WithGasMeter(storetypes.NewInfiniteGasMeter())
	hugeNum := `{"collectionId":` + strings.Repeat("9", evmcompat.DefaultMaxJSONNumberLength+1) + `}`
	require.Error(t, meterJSONInput(ctx, hugeNum))

	okNum := `{"collectionId":` + strings.Repeat("9", 30) + `}`
	require.NoError(t, meterJSONInput(ctx, okNum))
}

func TestParseBoundedUint256String(t *testing.T) {
	v, ok := parseBoundedUint256String("12345")
	require.True(t, ok)
	require.Equal(t, "12345", v.String())

	// Oversized decimal string (> 80 chars) rejected before math.ParseUint
	_, ok = parseBoundedUint256String(strings.Repeat("9", 81))
	require.False(t, ok)

	// Overflow > 2^256-1 rejected
	_, ok = parseBoundedUint256String("115792089237316195423570985008687907853269984665640564039457584007913129639936")
	require.False(t, ok)
}

func TestConvertBech32ToEvmAddressRejectsNon20ByteAddresses(t *testing.T) {
	p := NewPrecompile(&tokenizationkeeper.Keeper{}, nil)
	method := p.ABI.Methods[ConvertBech32ToEvmAddressMethod]

	// Standard 20-byte account address succeeds
	addr20 := sdk.AccAddress(common.HexToAddress("0x1111111111111111111111111111111111111111").Bytes())
	out, err := p.HandleConvertBech32ToEvmAddress(&method, []interface{}{addr20.String()})
	require.NoError(t, err)
	require.NotEmpty(t, out)

	// 32-byte derived account address must return error instead of silently truncating
	raw32 := make([]byte, 32)
	for i := range raw32 {
		raw32[i] = byte(i + 1)
	}
	addr32 := sdk.AccAddress(raw32)
	require.Equal(t, 32, len(addr32))
	_, err = p.HandleConvertBech32ToEvmAddress(&method, []interface{}{addr32.String()})
	require.Error(t, err)
}

func TestConvertEVMAddressToBech32MatchesKeeperFormat(t *testing.T) {
	hexAddr := "0x2222222222222222222222222222222222222222"
	bech32Addr := convertEVMAddressToBech32(hexAddr)
	expected := sdk.AccAddress(common.HexToAddress(hexAddr).Bytes()).String()
	require.Equal(t, expected, bech32Addr)

	// Unrecognized address string passes through unchanged
	require.Equal(t, "not-an-address", convertEVMAddressToBech32("not-an-address"))
}
