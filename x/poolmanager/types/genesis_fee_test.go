package types

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"
)

func TestGenesisRejectsInvalidPairFees(t *testing.T) {
	for _, fee := range []sdkmath.LegacyDec{sdkmath.LegacyOneDec(), sdkmath.LegacyNewDec(-1), {}} {
		gs := DefaultGenesis()
		gs.DenomPairTakerFeeStore = []DenomPairTakerFee{{TokenInDenom: "ubadge", TokenOutDenom: "uatom", TakerFee: fee}}
		require.Error(t, gs.Validate())
	}
}
