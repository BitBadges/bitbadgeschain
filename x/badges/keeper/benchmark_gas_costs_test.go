package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/rodaine/table"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

var tbl = table.New("Function Name", "Gas Consumed")

type GasFunction struct {
	Name      string
	IgnoreGas bool
	F         func()
}

func RunFunctionsAndPrintGasCosts(suite *TestSuite, tbl table.Table, functions []GasFunction) {
	for _, f := range functions {
		if f.IgnoreGas {
			f.F()
		} else {
			startGas := suite.ctx.GasMeter().GasConsumed()
			f.F()
			endGas := suite.ctx.GasMeter().GasConsumed()
			tbl.AddRow(f.Name, endGas-startGas)
		}
	}

	firstColumnFormatter := func(format string, vals ...interface{}) string {
		return ""
	}

	tbl.WithFirstColumnFormatter(firstColumnFormatter).Print()
}

func (suite *TestSuite) TestGasCosts() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	badgesToCreate := []BadgesToCreate{
		{
			Badge: types.MsgNewBadge{
				Uri:          validUri,
				Permissions:  46,
				SubassetUris: validUri,
			},
			Amount:  1,
			Creator: bob,
		},
	}

	badgesToCreate2 := []BadgesToCreate{
		{
			Badge: types.MsgNewBadge{
				Uri:          validUri,
				Permissions:  62,
				SubassetUris: validUri,
			},
			Amount:  1,
			Creator: bob,
		},
	}

	addressesThroughTenThousand := make([]uint64, 10000)
	for j := 0; j < 10000; j++ {
		addressesThroughTenThousand[j] = uint64(j)
	}

	addresses := []uint64{}
	for i := 0; i < 1000; i++ {
		addresses = append(addresses, firstAccountNumCreated+1+uint64(i))
	}

	badgesToCreateAllInOne := []BadgesToCreate{
		{
			Badge: types.MsgNewBadge{
				Uri:                     validUri,
				Permissions:             62,
				SubassetUris:            validUri,
				SubassetSupplys:         []uint64{1000000, 1, 10000},
				SubassetAmountsToCreate: []uint64{1, 10000, 10000},
				FreezeAddressRanges: []*types.IdRange{
					{
						Start: 1000,
						End:   1000,
					},
				},
			},
			Amount:  1,
			Creator: bob,
		},
	}

	RunFunctionsAndPrintGasCosts(suite, tbl, []GasFunction{
		{F: func() { CreateBadges(suite, wctx, badgesToCreate) }},
		{F: func() { GetBadge(suite, wctx, 0) }},
		{F: func() {
			err := CreateSubBadges(suite, wctx, bob, 0, []uint64{1000000}, []uint64{1})
			suite.Require().Nil(err, "Error creating subbadge")
		}},
		{F: func() {
			err := CreateSubBadges(suite, wctx, bob, 0, []uint64{1}, []uint64{10000})
			suite.Require().Nil(err, "Error creating subbadge")
		}},
		{F: func() {
			err := CreateSubBadges(suite, wctx, bob, 0, []uint64{10000}, []uint64{10000})
			suite.Require().Nil(err, "Error creating subbadge")
		}},
		{F: func() {
			err := FreezeAddresses(suite, wctx, bob, []*types.IdRange{{Start: 1000, End: 1000}}, 0, 0, true)
			suite.Require().Nil(err, "Error creating subbadge")
		}},
		{F: func() { GetBadge(suite, wctx, 0) }},
		{F: func() {
			err := FreezeAddresses(suite, wctx, bob, []*types.IdRange{{Start: 0, End: 9999}}, 0, 0, true)
			suite.Require().Nil(err, "Error creating subbadge")
		}},
		{F: func() { GetBadge(suite, wctx, 0) }},
		{F: func() {
			err := FreezeAddresses(suite, wctx, bob, []*types.IdRange{{Start: 1000, End: 1000}}, 0, 0, false)
			suite.Require().Nil(err, "Error creating subbadge")
		}},
		{F: func() { GetBadge(suite, wctx, 0) }},
		{F: func() {
			err := FreezeAddresses(suite, wctx, bob, []*types.IdRange{{Start: 0, End: 9999}}, 0, 0, false)
			suite.Require().Nil(err, "Error creating subbadge")
		}},
		{F: func() {
			GetBadge(suite, wctx, 0)
		}},
		{F: func() {
			err := TransferBadge(suite, wctx, bob, firstAccountNumCreated, []uint64{firstAccountNumCreated + 1}, []uint64{1}, 0, []*types.IdRange{{Start: 0, End: 0}}, 0)
			suite.Require().Nil(err, "Error transferring badge")
		}},
		{F: func() {
			GetUserBalance(suite, wctx, 0, 0, firstAccountNumCreated+1)
		}},
		{F: func() {
			err := HandlePendingTransfers(suite, wctx, alice, true, 0, []*types.IdRange{{Start: 0, End: 0}}, false)
			suite.Require().Nil(err, "Error accepting badge")
		}},
		{F: func() {
			GetUserBalance(suite, wctx, 0, 0, firstAccountNumCreated+1)
		}},
		{F: func() {
			err := TransferBadge(suite, wctx, bob, firstAccountNumCreated, []uint64{firstAccountNumCreated + 1}, []uint64{1}, 0, []*types.IdRange{{Start: 0, End: 999}}, 0)
			suite.Require().Nil(err, "Error transferring badge")
		}},
		{F: func() {
			GetUserBalance(suite, wctx, 0, 0, firstAccountNumCreated+1)
		}},
		{F: func() {
			err := HandlePendingTransfers(suite, wctx, alice, true, 0, []*types.IdRange{{Start: 0, End: 999}}, false)
			suite.Require().Nil(err, "Error transferring badge")
		}},
		{F: func() {
			GetUserBalance(suite, wctx, 0, 0, firstAccountNumCreated+1)
		}},
		{F: func() {
			err := TransferBadge(suite, wctx, bob, firstAccountNumCreated, []uint64{firstAccountNumCreated + 1}, []uint64{1}, 0, []*types.IdRange{{Start: 0, End: 0}}, 0)
			suite.Require().Nil(err, "Error transferring badge")
		}},
		{F: func() {
			GetUserBalance(suite, wctx, 0, 0, firstAccountNumCreated+1)
		}},
		{F: func() {
			err := HandlePendingTransfers(suite, wctx, alice, false, 0, []*types.IdRange{{Start: 2, End: 2}}, false)
			suite.Require().Nil(err, "Error accepting badge")
		}},
		{
			IgnoreGas: true,
			F: func() {
				err := HandlePendingTransfers(suite, wctx, bob, false, 0, []*types.IdRange{{Start: 2, End: 2}}, false)
				suite.Require().Nil(err, "Error accepting badge")
			},
		},
		{F: func() {
			GetUserBalance(suite, wctx, 0, 0, firstAccountNumCreated+1)
		}},
		{F: func() {
			err := TransferBadge(suite, wctx, bob, firstAccountNumCreated, []uint64{firstAccountNumCreated + 1}, []uint64{1}, 0, []*types.IdRange{{Start: 1000, End: 1999}}, 0)
			suite.Require().Nil(err, "Error transferring badge")
		}},
		{F: func() {
			GetUserBalance(suite, wctx, 0, 0, firstAccountNumCreated+1)
		}},
		{F: func() {
			err := HandlePendingTransfers(suite, wctx, alice, false, 0, []*types.IdRange{{Start: 0, End: 999}}, false)
			suite.Require().Nil(err, "Error transferring badge")
		}},
		{F: func() {
			err := HandlePendingTransfers(suite, wctx, bob, false, 0, []*types.IdRange{{Start: 0, End: 999}}, true)
			suite.Require().Nil(err, "Error transferring badge")
		},
			IgnoreGas: true,
		},
		{F: func() {
			GetUserBalance(suite, wctx, 0, 0, firstAccountNumCreated+1)
		}},
		{F: func() {
			err := RequestTransferBadge(suite, wctx, alice, firstAccountNumCreated, 1, 0, []*types.IdRange{{Start: 0, End: 0}}, 0)
			suite.Require().Nil(err, "Error transferring badge")
		}},
		{F: func() {
			GetUserBalance(suite, wctx, 0, 0, firstAccountNumCreated+1)
		}},
		{F: func() {
			err := HandlePendingTransfers(suite, wctx, bob, true, 0, []*types.IdRange{{Start: 4, End: 4}}, true)
			suite.Require().Nil(err, "Error accepting badge")
		}},
		{F: func() {
			GetUserBalance(suite, wctx, 0, 0, firstAccountNumCreated+1)
		}},
		{F: func() {
			err := RequestTransferBadge(suite, wctx, alice, firstAccountNumCreated, 1, 0, []*types.IdRange{{Start: 1000, End: 1999}}, 0)
			suite.Require().Nil(err, "Error transferring badge")
		}},
		{F: func() {
			GetUserBalance(suite, wctx, 0, 0, firstAccountNumCreated+1)
		}},
		{F: func() {
			err := HandlePendingTransfers(suite, wctx, bob, true, 0, []*types.IdRange{{Start: 0, End: 999}}, true)
			suite.Require().Nil(err, "Error transferring badge")
		}},
		{F: func() {
			GetUserBalance(suite, wctx, 0, 0, firstAccountNumCreated+1)
		}},
		{F: func() {
			err := RequestTransferBadge(suite, wctx, alice, firstAccountNumCreated, 1, 0, []*types.IdRange{{Start: 0, End: 0}}, 0)
			suite.Require().Nil(err, "Error transferring badge")
		}},
		{F: func() {
			GetUserBalance(suite, wctx, 0, 0, firstAccountNumCreated+1)
		}},
		{F: func() {
			err := HandlePendingTransfers(suite, wctx, bob, false, 0, []*types.IdRange{{Start: 0, End: 999}}, false)
			suite.Require().Nil(err, "Error accepting badge")
		}},
		{F: func() {
			err := HandlePendingTransfers(suite, wctx, alice, false, 0, []*types.IdRange{{Start: 0, End: 999}}, false)
			suite.Require().Nil(err, "Error accepting badge")
		},
			IgnoreGas: true,
		},
		{F: func() {
			GetUserBalance(suite, wctx, 0, 0, firstAccountNumCreated+1)
		}},
		{F: func() {
			err := RequestTransferBadge(suite, wctx, alice, firstAccountNumCreated, 1, 0, []*types.IdRange{{Start: 0, End: 999}}, 0)
			suite.Require().Nil(err, "Error transferring badge")
		}},
		{F: func() {
			GetUserBalance(suite, wctx, 0, 0, firstAccountNumCreated+1)
		}},
		{F: func() {
			err := HandlePendingTransfers(suite, wctx, bob, false, 0, []*types.IdRange{{Start: 0, End: 999}}, false)
			suite.Require().Nil(err, "Error transferring badge")
		}},
		{F: func() {
			err := HandlePendingTransfers(suite, wctx, alice, false, 0, []*types.IdRange{{Start: 0, End: 999}}, false)
			suite.Require().Nil(err, "Error accepting badge")
		},
			IgnoreGas: true,
		},
		{F: func() {
			GetUserBalance(suite, wctx, 0, 0, firstAccountNumCreated+1)
		}},
		{F: func() {
			GetUserBalance(suite, wctx, 0, 0, firstAccountNumCreated+1)
		}},
		{F: func() {
			err := RevokeBadges(suite, wctx, bob, []uint64{firstAccountNumCreated + 1}, []uint64{1}, 0, []*types.IdRange{{Start: 0, End: 0}})
			suite.Require().Nil(err, "Error transferring badge")
		}},
		{F: func() {
			GetUserBalance(suite, wctx, 0, 0, firstAccountNumCreated+1)
		}},
		{F: func() {
			err := RevokeBadges(suite, wctx, bob, []uint64{firstAccountNumCreated + 1}, []uint64{1}, 0, []*types.IdRange{{Start: 1000, End: 1999}})
			suite.Require().Nil(err, "Error transferring badge")
		}},
		{F: func() {
			GetUserBalance(suite, wctx, 0, 0, firstAccountNumCreated+1)
		}},
		{F: func() {
			GetUserBalance(suite, wctx, 0, 0, firstAccountNumCreated)
		}},
		{F: func() {
			err := SetApproval(suite, wctx, bob, 10000, firstAccountNumCreated+1, 0, []*types.IdRange{{Start: 0, End: 0}}, 0)
			suite.Require().Nil(err, "Error transferring badge")
		}},
		{F: func() {
			GetUserBalance(suite, wctx, 0, 0, firstAccountNumCreated)
		}},
		{F: func() {
			err := SetApproval(suite, wctx, bob, 10, firstAccountNumCreated+1, 0, []*types.IdRange{{Start: 0, End: 0}}, 0)
			suite.Require().Nil(err, "Error transferring badge")
		}},
		{F: func() {
			GetUserBalance(suite, wctx, 0, 0, firstAccountNumCreated)
		}},
		{F: func() {
			err := SetApproval(suite, wctx, bob, 10, firstAccountNumCreated+1, 0, []*types.IdRange{{Start: 0, End: 999}}, 0)
			suite.Require().Nil(err, "Error transferring badge")
		}},
		{F: func() {
			GetUserBalance(suite, wctx, 0, 0, firstAccountNumCreated)
		}},
		{F: func() {
			CreateBadges(suite, wctx, badgesToCreate2)
		}},
		{F: func() {
			GetBadge(suite, wctx, 1)
		}},
		{F: func() {
			err := CreateSubBadges(suite, wctx, bob, 1, []uint64{10000}, []uint64{10000})
			suite.Require().Nil(err, "Error creating subbadge")
		}},
		{F: func() {
			err := TransferBadge(suite, wctx, bob, firstAccountNumCreated, []uint64{firstAccountNumCreated + 1}, []uint64{1}, 1, []*types.IdRange{{Start: 0, End: 0}}, 0)
			suite.Require().Nil(err, "Error transferring badge")
		}},
		{F: func() {
			GetUserBalance(suite, wctx, 0, 0, firstAccountNumCreated+1)
		}},
		{F: func() {
			err := TransferBadge(suite, wctx, bob, firstAccountNumCreated, []uint64{firstAccountNumCreated + 1}, []uint64{1}, 1, []*types.IdRange{{Start: 0, End: 999}}, 0)
			suite.Require().Nil(err, "Error transferring badge")
		}},
		{F: func() {
			GetUserBalance(suite, wctx, 0, 0, firstAccountNumCreated+1)
		}},
		{F: func() {

			err := TransferBadge(suite, wctx, bob, firstAccountNumCreated, addresses, []uint64{1}, 1, []*types.IdRange{{Start: 0, End: 0}}, 0)
			suite.Require().Nil(err, "Error transferring badge")
			suite.Require().Nil(err, "Error transferring badge")
		}},
		{F: func() {
			GetUserBalance(suite, wctx, 0, 0, firstAccountNumCreated+1)
		}},
		{F: func() {
			err := TransferBadge(suite, wctx, bob, firstAccountNumCreated, addresses, []uint64{1}, 1, []*types.IdRange{{Start: 0, End: 999}}, 0)
			suite.Require().Nil(err, "Error transferring badge")
			suite.Require().Nil(err, "Error transferring badge")
		}},
		{F: func() {
			GetUserBalance(suite, wctx, 0, 0, firstAccountNumCreated+1)
		}},
		{F: func() {
			for i := uint64(0); i < 1000; i++ {
				err := TransferBadge(suite, wctx, bob, firstAccountNumCreated, []uint64{firstAccountNumCreated + 1}, []uint64{1}, 1, []*types.IdRange{{Start: i * 2, End: i * 2}}, 0)
				suite.Require().Nil(err, "Error transferring badge")
				suite.Require().Nil(err, "Error transferring badge")
			}
		}},
		{F: func() {
			GetUserBalance(suite, wctx, 1, 0, firstAccountNumCreated)
		}},
		{F: func() { CreateBadges(suite, wctx, badgesToCreateAllInOne) }},
	})
}
