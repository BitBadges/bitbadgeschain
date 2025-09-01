package balancer

import (
	"github.com/bitbadges/bitbadgeschain/third_party/osmomath"
	"github.com/bitbadges/bitbadgeschain/x/gamm/types"
)

func NewPoolParams(spreadFactor, exitFee osmomath.Dec) PoolParams {
	return PoolParams{
		SwapFee: spreadFactor,
		ExitFee: exitFee,
	}
}

func (params PoolParams) Validate(poolWeights []PoolAsset) error {
	if params.ExitFee.IsNegative() {
		return types.ErrNegativeExitFee
	}

	if params.ExitFee.GTE(osmomath.OneDec()) {
		return types.ErrTooMuchExitFee
	}

	if params.SwapFee.IsNegative() {
		return types.ErrNegativeSpreadFactor
	}

	if params.SwapFee.GTE(osmomath.OneDec()) {
		return types.ErrTooMuchSpreadFactor
	}

	return nil
}

func (params PoolParams) GetPoolSpreadFactor() osmomath.Dec {
	return params.SwapFee
}

func (params PoolParams) GetPoolExitFee() osmomath.Dec {
	return params.ExitFee
}
