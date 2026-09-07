package evmcompat

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	cmn "github.com/cosmos/evm/precompiles/common"
	"github.com/ethereum/go-ethereum/core/tracing"
	"github.com/ethereum/go-ethereum/core/vm"
)

// RunNativeAction settles native work on ordinary errors before the upstream
// wrapper converts them to EVM reverts, which preserve remaining gas.
func RunNativeAction(p cmn.Precompile, evm *vm.EVM, contract *vm.Contract, action cmn.NativeAction) ([]byte, error) {
	return p.RunNativeAction(evm, contract, func(ctx sdk.Context) ([]byte, error) {
		initialGas := ctx.GasMeter().GasConsumed()
		bz, err := action(ctx)
		if err != nil {
			cost := ctx.GasMeter().GasConsumed() - initialGas
			if !contract.UseGas(cost, nil, tracing.GasChangeCallPrecompiledContract) {
				return nil, vm.ErrOutOfGas
			}
		}
		return bz, err
	})
}
