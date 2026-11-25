package precompile

import (
	"strings"
	"testing"

	"github.com/bitbadges/bitbadgeschain/pkg/evmcompat"
	sendmanagerkeeper "github.com/bitbadges/bitbadgeschain/x/sendmanager/keeper"
	storetypes "github.com/cosmos/cosmos-sdk/store/v2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/stretchr/testify/require"
)

func TestSendManagerJSONBudgetBounds(t *testing.T) {
	p := NewPrecompile(sendmanagerkeeper.Keeper{}, nil)
	ctx := sdk.Context{}.WithGasMeter(storetypes.NewInfiniteGasMeter())
	caller := common.HexToAddress("0x1111111111111111111111111111111111111111")
	contract := vm.NewContract(caller, p.ContractAddress, nil, 1_000_000, nil)

	// Deep nesting rejected
	deep := strings.Repeat("[", evmcompat.DefaultMaxJSONDepth+2) + "0" + strings.Repeat("]", evmcompat.DefaultMaxJSONDepth+2)
	_, err := p.unmarshalMsgFromJSON(ctx, SendMethod, deep, contract)
	require.Error(t, err)

	// Oversized JSON number rejected
	hugeNum := `{"amount":[{"denom":"ubadge","amount":` + strings.Repeat("9", evmcompat.DefaultMaxJSONNumberLength+1) + `}]}`
	_, err = p.unmarshalMsgFromJSON(ctx, SendMethod, hugeNum, contract)
	require.Error(t, err)

	// Trailing garbage rejected
	_, err = p.unmarshalMsgFromJSON(ctx, SendMethod, `{} {}`, contract)
	require.Error(t, err)
}
