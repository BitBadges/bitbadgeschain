package hooks

import (
	"testing"

	transfertypes "github.com/cosmos/ibc-go/v11/modules/apps/transfer/types"
	"github.com/stretchr/testify/require"
)

func TestExtractDenomFromPacketOnSend(t *testing.T) {
	tests := []struct {
		name        string
		packetDenom string
		want        string
	}{
		{name: "native denom", packetDenom: "ubb", want: "ubb"},
		{
			name:        "local voucher returning over its source channel",
			packetDenom: "transfer/channel-0/uatom",
			want:        transfertypes.ParseDenomTrace("transfer/channel-0/uatom").IBCDenom(),
		},
		{
			name:        "multi-hop local voucher",
			packetDenom: "transfer/channel-1/transfer/channel-2/uatom",
			want:        transfertypes.ParseDenomTrace("transfer/channel-1/transfer/channel-2/uatom").IBCDenom(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, extractDenomFromPacketOnSend("transfer", "channel-0", tt.packetDenom))
		})
	}
}
