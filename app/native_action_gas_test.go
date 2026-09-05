package app

import (
	"errors"
	"math/big"
	"testing"

	"github.com/bitbadges/bitbadgeschain/pkg/evmcompat"
	gammprecompile "github.com/bitbadges/bitbadgeschain/x/gamm/precompile"
	sendmanagerprecompile "github.com/bitbadges/bitbadgeschain/x/sendmanager/precompile"
	tokenizationprecompile "github.com/bitbadges/bitbadgeschain/x/tokenization/precompile"
	tokenizationtypes "github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/v2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	cmn "github.com/cosmos/evm/precompiles/common"
	"github.com/cosmos/evm/x/vm/statedb"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/tracing"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/require"
)

func TestNativeActionRevertBillsConsumedGas(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContext(false).WithGasMeter(storetypes.NewInfiniteGasMeter())
	ctx.GasMeter().ConsumeGas(123_456, "work before native action")
	db := statedb.New(ctx, app.EVMKeeper, statedb.TxConfig{})
	evm := vm.NewEVM(vm.BlockContext{}, db, params.AllEthashProtocolChanges, vm.Config{})
	contract := vm.NewContract(common.HexToAddress("0x1234"), common.HexToAddress("0x1001"), uint256.NewInt(0), 2_000_000, nil)
	precompile := cmn.Precompile{KvGasConfig: storetypes.KVGasConfig(), TransientKVGasConfig: storetypes.TransientGasConfig()}
	_, err := evmcompat.RunNativeAction(precompile, evm, contract, func(ctx sdk.Context) ([]byte, error) {
		ctx.GasMeter().ConsumeGas(1_000_000, "native workload before ordinary keeper error")
		return nil, errors.New("last batch item fails")
	})
	require.ErrorIs(t, err, vm.ErrExecutionReverted)
	t.Logf("native gas consumed=1000000, EVM gas charged=%d", 2_000_000-contract.Gas)
	require.Equal(t, uint64(1_000_000), contract.Gas, "ordinary reverts must charge only this action's native work")
}

type nativeGasProbe struct {
	cmn.Precompile
	action cmn.NativeAction
}

func (p nativeGasProbe) Name() string              { return "native-gas-probe" }
func (p nativeGasProbe) RequiredGas([]byte) uint64 { return 10_000 }
func (p nativeGasProbe) Run(evm *vm.EVM, contract *vm.Contract, _ bool) ([]byte, error) {
	return evmcompat.RunNativeAction(p.Precompile, evm, contract, p.action)
}

func nativeGasEVM(app *App, ctx sdk.Context) (*vm.EVM, *statedb.StateDB) {
	db := statedb.New(ctx, app.EVMKeeper, statedb.TxConfig{})
	evm := vm.NewEVM(vm.BlockContext{
		BlockNumber: big.NewInt(1), GasLimit: 100_000_000,
		CanTransfer: core.CanTransfer,
		Transfer:    core.Transfer,
	}, db, params.AllEthashProtocolChanges, vm.Config{})
	return evm, db
}

func TestNativeActionGasSettlementAndRollback(t *testing.T) {
	for _, tc := range []struct {
		name          string
		cost          uint64
		actionErr     error
		wantErr       error
		wantRemaining uint64
	}{
		{"success", 1_000_000, nil, nil, 990_000},
		{"ordinary revert", 1_000_000, errors.New("keeper rejected final item"), vm.ErrExecutionReverted, 990_000},
		{"SDK out of gas", 3_000_000, nil, vm.ErrOutOfGas, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			app := Setup(false)
			ctx := app.NewContext(false).WithGasMeter(storetypes.NewInfiniteGasMeter())
			evm, db := nativeGasEVM(app, ctx)
			addr := common.HexToAddress("0x1004")
			probe := nativeGasProbe{Precompile: cmn.Precompile{ContractAddress: addr}, action: func(ctx sdk.Context) ([]byte, error) {
				ctx.KVStore(app.GetKey(tokenizationtypes.StoreKey)).Set([]byte("native-gas-probe"), []byte("written"))
				ctx.GasMeter().ConsumeGas(tc.cost, "native work")
				return []byte("result"), tc.actionErr
			}}
			evm.SetPrecompiles(vm.PrecompiledContracts{addr: probe})
			ret, remaining, err := evm.Call(common.HexToAddress("0x1234"), addr, nil, 2_000_000, uint256.NewInt(0))
			require.ErrorIs(t, err, tc.wantErr)
			require.Equal(t, tc.wantRemaining, remaining, "static and dynamic gas must each be charged exactly once")
			cacheCtx, cacheErr := db.GetCacheContext()
			require.NoError(t, cacheErr)
			got := cacheCtx.KVStore(app.GetKey(tokenizationtypes.StoreKey)).Get([]byte("native-gas-probe"))
			if tc.wantErr == nil {
				require.Equal(t, []byte("written"), got)
				require.Equal(t, []byte("result"), ret)
			} else {
				require.Nil(t, got, "reverted native writes must not survive")
			}
		})
	}
}

func TestNativeActionCaughtRevertsRetainGasCharges(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContext(false).WithGasMeter(storetypes.NewInfiniteGasMeter())
	evm, db := nativeGasEVM(app, ctx)
	addr := common.HexToAddress("0x1004")
	called := 0
	probe := nativeGasProbe{Precompile: cmn.Precompile{ContractAddress: addr}, action: func(ctx sdk.Context) ([]byte, error) {
		called++
		ctx.KVStore(app.GetKey(tokenizationtypes.StoreKey)).Set([]byte("caught-revert-probe"), []byte("written"))
		ctx.GasMeter().ConsumeGas(1_000_000, "native work")
		return nil, errors.New("caught keeper error")
	}}
	evm.SetPrecompiles(vm.PrecompiledContracts{addr: probe})
	// CALL twice, discard each false result, then successfully return from the parent.
	code := []byte{}
	for i := 0; i < 2; i++ {
		code = append(code, 0x60, 0, 0x60, 0, 0x60, 0, 0x60, 0, 0x60, 0, 0x73)
		code = append(code, addr.Bytes()...)
		code = append(code, 0x62, 0x16, 0xe3, 0x60, 0xf1, 0x50) // forward 1,500,000 gas
	}
	code = append(code, 0x60, 0, 0x60, 0, 0xf3)
	parent := common.HexToAddress("0x3333")
	db.SetCode(parent, code, tracing.CodeChangeUnspecified)
	_, remaining, err := evm.Call(common.HexToAddress("0x1234"), parent, nil, 4_000_000, uint256.NewInt(0))
	require.NoError(t, err)
	require.Equal(t, 2, called)
	require.GreaterOrEqual(t, 4_000_000-remaining, uint64(2_020_000))
	require.Less(t, 4_000_000-remaining, uint64(2_050_000), "caught reverts must not burn unused child gas")
	cacheCtx, err := db.GetCacheContext()
	require.NoError(t, err)
	require.Nil(t, cacheCtx.KVStore(app.GetKey(tokenizationtypes.StoreKey)).Get([]byte("caught-revert-probe")))
}

func TestCustomPrecompileFailuresBillNativeWork(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContext(false).WithGasMeter(storetypes.NewInfiniteGasMeter())
	caller := common.HexToAddress("0x1234")
	token := tokenizationprecompile.NewPrecompile(app.TokenizationKeeper, app.BankKeeper)
	gamm := gammprecompile.NewPrecompile(app.GammKeeper, app.BankKeeper)
	send := sendmanagerprecompile.NewPrecompile(app.SendmanagerKeeper, app.BankKeeper)
	tokenInput, err := token.ABI.Pack("setCustomData", `{"collectionId":"999","customData":"test"}`)
	require.NoError(t, err)
	gammInput, err := gamm.ABI.Pack("swapExactAmountIn", `{"routes":[{"pool_id":999,"token_out_denom":"stake"}],"token_in":{"denom":"ubadge","amount":"1"},"token_out_min_amount":"1"}`)
	require.NoError(t, err)
	sendInput, err := send.ABI.Pack("send", `{"toAddress":"`+sdk.AccAddress(common.HexToAddress("0x5678").Bytes()).String()+`","amount":[{"denom":"ubadge","amount":"1"}]}`)
	require.NoError(t, err)
	for _, tc := range []struct {
		p     vm.PrecompiledContract
		input []byte
	}{{token, tokenInput}, {gamm, gammInput}, {send, sendInput}} {
		t.Run(tc.p.Name(), func(t *testing.T) {
			evm, _ := nativeGasEVM(app, ctx)
			contract := vm.NewContract(caller, tc.p.Address(), uint256.NewInt(0), 2_000_000, nil)
			contract.Input = tc.input
			ret, err := tc.p.Run(evm, contract, false)
			require.ErrorIs(t, err, vm.ErrExecutionReverted)
			reason, _ := abi.UnpackRevert(ret)
			require.Less(t, contract.Gas, uint64(2_000_000), "native work must be charged in addition to RequiredGas: %s", reason)
			require.Positive(t, contract.Gas, "ordinary failure must retain unused gas")
		})
	}
}
