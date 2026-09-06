package tokenization_test

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm/runtime"
	"github.com/stretchr/testify/require"
)

func TestTransferJSONSenderForms(t *testing.T) {
	const artifact = "../../../../../contracts/test/test_TransferJSONTestContract_sol_TransferJSONTestContract"
	abiBytes, err := os.ReadFile(artifact + ".abi")
	require.NoError(t, err)
	contractABI, err := abi.JSON(strings.NewReader(string(abiBytes)))
	require.NoError(t, err)
	binBytes, err := os.ReadFile(artifact + ".bin")
	require.NoError(t, err)
	bytecode, err := hex.DecodeString(strings.TrimSpace(string(binBytes)))
	require.NoError(t, err)
	code, _, err := runtime.Execute(bytecode, nil, nil)
	require.NoError(t, err)

	from := common.HexToAddress("0xabcdef1234567890123456789012345678901234")
	recipients := []common.Address{common.HexToAddress("0x1234"), common.HexToAddress("0xabcd")}
	amount := new(big.Int).Lsh(big.NewInt(1), 128)
	const tokenIDs = `[{"start":"2","end":"9"}]`
	const times = `[{"start":"1","end":"18446744073709551615"}]`
	for _, tc := range []struct {
		name string
		from *common.Address
	}{{"default", nil}, {"explicit", &from}, {"explicit zero", new(common.Address)}} {
		t.Run(tc.name, func(t *testing.T) {
			method := "defaultSender"
			args := []interface{}{big.NewInt(7)}
			fromField := ""
			if tc.from != nil {
				method = "explicitSender"
				args = append(args, *tc.from)
				fromField = fmt.Sprintf(`"from":"%s",`, strings.ToLower(tc.from.Hex()))
			}
			args = append(args, recipients, amount, tokenIDs, times)
			input, err := contractABI.Pack(method, args...)
			require.NoError(t, err)
			output, _, err := runtime.Execute(code, input, nil)
			require.NoError(t, err)
			values, err := contractABI.Methods[method].Outputs.Unpack(output)
			require.NoError(t, err)
			require.Len(t, values, 1)
			expected := fmt.Sprintf(`{"collectionId":"7","transfers":[{%s"toAddresses":["%s","%s"],"balances":[{"amount":"%s","tokenIds":%s,"ownershipTimes":%s}]}]}`,
				fromField, strings.ToLower(recipients[0].Hex()), strings.ToLower(recipients[1].Hex()), amount, tokenIDs, times)
			require.JSONEq(t, expected, values[0].(string))
		})
	}
}
