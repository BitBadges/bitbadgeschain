package keeper_test

import (
	"fmt"
	"math/big"
	"testing"

	storetypes "github.com/cosmos/cosmos-sdk/store/v2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/evm/x/vm/statedb"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	testutilkeeper "github.com/bitbadges/bitbadgeschain/x/tokenization/testutil/keeper"
)

type contextQueryKeeper struct {
	call func(sdk.Context, *big.Int) (*evmtypes.MsgEthereumTxResponse, error)
}

func (*contextQueryKeeper) IsContract(sdk.Context, common.Address) bool { return true }

func (m *contextQueryKeeper) CallEVMWithData(ctx sdk.Context, _ *statedb.StateDB, _ common.Address, _ *common.Address, _ []byte, _ bool, _ bool, gas *big.Int) (*evmtypes.MsgEthereumTxResponse, error) {
	return m.call(ctx, gas)
}

func TestEVMQueryRejectsNestedQuery(t *testing.T) {
	k, ctx := testutilkeeper.TokenizationKeeper(t)
	calls := 0
	k.SetEVMKeeper(&contextQueryKeeper{call: func(child sdk.Context, _ *big.Int) (*evmtypes.MsgEthereumTxResponse, error) {
		calls++
		if calls == 4 {
			return &evmtypes.MsgEthereumTxResponse{GasUsed: 100}, fmt.Errorf("bounded test stop")
		}
		_, err := k.ExecuteEVMQuery(child, gasMeteringContract, nil, 500000)
		return &evmtypes.MsgEthereumTxResponse{GasUsed: 100, VmError: "nested query rejected"}, err
	}})
	_, err := k.ExecuteEVMQuery(ctx, gasMeteringContract, nil, 500000)
	require.Error(t, err)
	require.Equal(t, 1, calls)
	_, err = k.ExecuteEVMQuery(ctx, gasMeteringContract, nil, 500000)
	require.Error(t, err)
	require.Equal(t, 2, calls, "independent queries must remain available")
}

func TestEVMQueryUsesRemainingParentGas(t *testing.T) {
	k, ctx := testutilkeeper.TokenizationKeeper(t)
	ctx = ctx.WithGasMeter(storetypes.NewGasMeter(100000))
	ctx.GasMeter().ConsumeGas(60000, "prior work")
	k.SetEVMKeeper(&contextQueryKeeper{call: func(child sdk.Context, cap *big.Int) (*evmtypes.MsgEthereumTxResponse, error) {
		require.LessOrEqual(t, cap.Uint64(), uint64(40000))
		require.Equal(t, cap.Uint64(), child.GasMeter().Limit())
		return &evmtypes.MsgEthereumTxResponse{GasUsed: 100}, nil
	}})
	_, err := k.ExecuteEVMQuery(ctx, gasMeteringContract, nil, 500000)
	require.NoError(t, err)
}
