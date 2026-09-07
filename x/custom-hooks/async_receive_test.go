package customhooks_test

import (
	ibchooks "github.com/bitbadges/bitbadgeschain/x/ibc-hooks"
	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v11/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v11/modules/core/04-channel/types"
	ibcexported "github.com/cosmos/ibc-go/v11/modules/core/exported"
)

type pendingReceiveModule struct{ mockIBCModule }

func (m *pendingReceiveModule) OnRecvPacket(ctx sdk.Context, _ string, _ channeltypes.Packet, _ sdk.AccAddress) ibcexported.Acknowledgement {
	ctx.EventManager().EmitEvent(sdk.NewEvent("pending_receive"))
	return nil
}

func (s *HooksTestSuite) TestPendingReceiveCustomHookAtomicity() {
	for _, memo := range []string{"", `{"swap_and_action":{}}`} {
		data := transfertypes.FungibleTokenPacketData{Denom: "uatom", Amount: "10", Sender: "sender", Receiver: s.TestAccs[0].String(), Memo: memo}
		bz, err := transfertypes.ModuleCdc.MarshalJSON(&data)
		s.Require().NoError(err)
		ctx := s.Ctx.WithEventManager(sdk.NewEventManager())
		im := ibchooks.NewIBCMiddleware(&pendingReceiveModule{}, nil)
		s.Require().NotPanics(func() {
			ack := s.customHooks.OnRecvPacketOverride(im, ctx, "ics20-1", channeltypes.Packet{Data: bz}, nil)
			if memo == "" {
				s.Require().Nil(ack)
				s.Require().Len(ctx.EventManager().Events(), 1)
			} else {
				s.Require().NotNil(ack)
				s.Require().False(ack.Success())
				s.Require().Empty(ctx.EventManager().Events())
			}
		})
	}
}
