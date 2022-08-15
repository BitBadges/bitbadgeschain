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

	//Initial variable definitions
	badgesToCreate := []BadgesToCreate{
		{
			Badge: types.MsgNewBadge{
				Uri:         &types.UriObject{
					Uri: 	[]byte("example.com/"),
					Scheme: 1,
					IdxRangeToRemove: &types.IdRange{},
					InsertSubassetBytesIdx: 0,
					
					InsertIdIdx: 10,
				},
				Permissions:  46,
				
			},
			Amount:  1,
			Creator: bob,
		},
	}

	badgesToCreate2 := []BadgesToCreate{
		{
			Badge: types.MsgNewBadge{
				Uri:         &types.UriObject{
					Uri: 	[]byte("example.com/"),
					Scheme: 1,
					IdxRangeToRemove: &types.IdRange{},
					InsertSubassetBytesIdx: 0,
					
					InsertIdIdx: 10,
				},
				Permissions:  62,
				
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
		addresses = append(addresses, aliceAccountNum+uint64(i))
	}

	badgesToCreateAllInOne := []BadgesToCreate{
		{
			Badge: types.MsgNewBadge{
				Uri:         &types.UriObject{
					Uri: 	[]byte("example.com/"),
					Scheme: 1,
					IdxRangeToRemove: &types.IdRange{},
					InsertSubassetBytesIdx: 0,
					
					InsertIdIdx: 10,
				},
				Permissions:             62,
				
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
		{F: func() { CreateSubBadges(suite, wctx, bob, 0, []uint64{1000000}, []uint64{1}) }},
		{F: func() { CreateSubBadges(suite, wctx, bob, 0, []uint64{1}, []uint64{10000})}},
		{F: func() { CreateSubBadges(suite, wctx, bob, 0, []uint64{10000}, []uint64{10000})}},
		{F: func() { FreezeAddresses(suite, wctx, bob, []*types.IdRange{{Start: 1000, End: 1000}}, 0, true)}},
		{F: func() { GetBadge(suite, wctx, 0) }},
		{F: func() { FreezeAddresses(suite, wctx, bob, []*types.IdRange{{Start: 0, End: 9999}}, 0, true)}},
		{F: func() { GetBadge(suite, wctx, 0) }},
		{F: func() { FreezeAddresses(suite, wctx, bob, []*types.IdRange{{Start: 1000, End: 1000}}, 0, false)}},
		{F: func() { GetBadge(suite, wctx, 0) }},
		{F: func() { FreezeAddresses(suite, wctx, bob, []*types.IdRange{{Start: 0, End: 9999}}, 0, false)}},
		{F: func() { GetBadge(suite, wctx, 0) }},
		{F: func() { TransferBadge(suite, wctx, bob, bobAccountNum, []uint64{aliceAccountNum}, []uint64{1}, 0, []*types.IdRange{{Start: 0, End: 0}}, 0, 0)}},
		{F: func() { GetUserBalance(suite, wctx, 0, aliceAccountNum) }},
		{F: func() { HandlePendingTransfers(suite, wctx, alice, true, 0, []*types.IdRange{{Start: 0, End: 0}}, false) }},
		{F: func() { GetUserBalance(suite, wctx, 0, aliceAccountNum) }},
		{F: func() { TransferBadge(suite, wctx, bob, bobAccountNum, []uint64{aliceAccountNum}, []uint64{1}, 0, []*types.IdRange{{Start: 0, End: 999}}, 0, 0)}},
		{F: func() { GetUserBalance(suite, wctx, 0, aliceAccountNum) }},
		{F: func() { HandlePendingTransfers(suite, wctx, alice, true, 0, []*types.IdRange{{Start: 0, End: 999}}, false)}},
		{F: func() { GetUserBalance(suite, wctx, 0, aliceAccountNum) }},
		{F: func() { TransferBadge(suite, wctx, bob, bobAccountNum, []uint64{aliceAccountNum}, []uint64{1}, 0, []*types.IdRange{{Start: 0, End: 0}}, 0, 0)}},
		{F: func() { GetUserBalance(suite, wctx, 0, aliceAccountNum) }},
		{F: func() { HandlePendingTransfers(suite, wctx, alice, false, 0, []*types.IdRange{{Start: 2, End: 2}}, false) }},
		{ IgnoreGas: true, F: func() { HandlePendingTransfers(suite, wctx, bob, false, 0, []*types.IdRange{{Start: 2, End: 2}}, false)},},
		{F: func() { GetUserBalance(suite, wctx, 0, aliceAccountNum) }},
		{F: func() { TransferBadge(suite, wctx, bob, bobAccountNum, []uint64{aliceAccountNum}, []uint64{1}, 0, []*types.IdRange{{Start: 1000, End: 1999}}, 0, 0)}},
		{F: func() { GetUserBalance(suite, wctx, 0, aliceAccountNum) }},
		{F: func() { HandlePendingTransfers(suite, wctx, alice, false, 0, []*types.IdRange{{Start: 0, End: 999}}, false)}},
		{F: func() { HandlePendingTransfers(suite, wctx, bob, false, 0, []*types.IdRange{{Start: 0, End: 999}}, true)}, IgnoreGas: true,},
		{F: func() { GetUserBalance(suite, wctx, 0, aliceAccountNum) }},
		{F: func() { RequestTransferBadge(suite, wctx, alice, bobAccountNum, 1, 0, []*types.IdRange{{Start: 0, End: 0}}, 0, 0)}},
		{F: func() { GetUserBalance(suite, wctx, 0, aliceAccountNum) }},
		{F: func() { HandlePendingTransfers(suite, wctx, bob, true, 0, []*types.IdRange{{Start: 4, End: 4}}, true) }},
		{F: func() { GetUserBalance(suite, wctx, 0, aliceAccountNum) }},
		{F: func() { RequestTransferBadge(suite, wctx, alice, bobAccountNum, 1, 0, []*types.IdRange{{Start: 1000, End: 1999}}, 0, 0)}},
		{F: func() { GetUserBalance(suite, wctx, 0, aliceAccountNum) }},
		{F: func() { HandlePendingTransfers(suite, wctx, bob, true, 0, []*types.IdRange{{Start: 0, End: 999}}, true)}},
		{F: func() { GetUserBalance(suite, wctx, 0, aliceAccountNum) }},
		{F: func() { RequestTransferBadge(suite, wctx, alice, bobAccountNum, 1, 0, []*types.IdRange{{Start: 0, End: 0}}, 0, 0)}},
		{F: func() { GetUserBalance(suite, wctx, 0, aliceAccountNum) }},
		{F: func() { HandlePendingTransfers(suite, wctx, bob, false, 0, []*types.IdRange{{Start: 0, End: 999}}, false) }},
		{F: func() { HandlePendingTransfers(suite, wctx, alice, false, 0, []*types.IdRange{{Start: 0, End: 999}}, false) }, IgnoreGas: true,},
		{F: func() { GetUserBalance(suite, wctx, 0, aliceAccountNum) }},
		{F: func() { RequestTransferBadge(suite, wctx, alice, bobAccountNum, 1, 0, []*types.IdRange{{Start: 0, End: 999}}, 0, 0)}},
		{F: func() { GetUserBalance(suite, wctx, 0, aliceAccountNum) }},
		{F: func() { HandlePendingTransfers(suite, wctx, bob, false, 0, []*types.IdRange{{Start: 0, End: 999}}, false)}},
		{F: func() { HandlePendingTransfers(suite, wctx, alice, false, 0, []*types.IdRange{{Start: 0, End: 999}}, false) }, IgnoreGas: true,},
		{F: func() { GetUserBalance(suite, wctx, 0, aliceAccountNum) }},
		{F: func() { GetUserBalance(suite, wctx, 0, aliceAccountNum) }},
		{F: func() { RevokeBadges(suite, wctx, bob, []uint64{aliceAccountNum}, []uint64{1}, 0, []*types.IdRange{{Start: 0, End: 0}})}},
		{F: func() { GetUserBalance(suite, wctx, 0, aliceAccountNum) }},
		{F: func() { RevokeBadges(suite, wctx, bob, []uint64{aliceAccountNum}, []uint64{1}, 0, []*types.IdRange{{Start: 1000, End: 1999}})}},
		{F: func() { GetUserBalance(suite, wctx, 0, aliceAccountNum) }},
		{F: func() { GetUserBalance(suite, wctx, 0, bobAccountNum) }},
		{F: func() { SetApproval(suite, wctx, bob, 10000, aliceAccountNum, 0, []*types.IdRange{{Start: 0, End: 0}})}},
		{F: func() { GetUserBalance(suite, wctx, 0, bobAccountNum) }},
		{F: func() { SetApproval(suite, wctx, bob, 10, aliceAccountNum, 0, []*types.IdRange{{Start: 0, End: 0}})}},
		{F: func() { GetUserBalance(suite, wctx, 0, bobAccountNum) }},
		{F: func() { SetApproval(suite, wctx, bob, 10, aliceAccountNum, 0, []*types.IdRange{{Start: 0, End: 999}})}},
		{F: func() { GetUserBalance(suite, wctx, 0, bobAccountNum) }},
		{F: func() { CreateBadges(suite, wctx, badgesToCreate2) }},
		{F: func() { GetBadge(suite, wctx, 1) }},
		{F: func() { CreateSubBadges(suite, wctx, bob, 1, []uint64{10000}, []uint64{10000})}},
		{F: func() { TransferBadge(suite, wctx, bob, bobAccountNum, []uint64{aliceAccountNum}, []uint64{1}, 1, []*types.IdRange{{Start: 0, End: 0}}, 0, 0)}},
		{F: func() { GetUserBalance(suite, wctx, 0, aliceAccountNum) }},
		{F: func() { TransferBadge(suite, wctx, bob, bobAccountNum, []uint64{aliceAccountNum}, []uint64{1}, 1, []*types.IdRange{{Start: 0, End: 999}}, 0, 0)}},
		{F: func() { GetUserBalance(suite, wctx, 0, aliceAccountNum) }},
		{F: func() { TransferBadge(suite, wctx, bob, bobAccountNum, addresses, []uint64{1}, 1, []*types.IdRange{{Start: 0, End: 0}}, 0, 0)}},
		{F: func() { GetUserBalance(suite, wctx, 0, aliceAccountNum) }},
		{F: func() { TransferBadge(suite, wctx, bob, bobAccountNum, addresses, []uint64{1}, 1, []*types.IdRange{{Start: 0, End: 999}}, 0, 0)}},
		{F: func() { GetUserBalance(suite, wctx, 0, aliceAccountNum) }},
		{F: func() { 
			for i := uint64(0); i < 1000; i++ {
 				TransferBadge(suite, wctx, bob, bobAccountNum, []uint64{aliceAccountNum}, []uint64{1}, 1, []*types.IdRange{{Start: i * 2, End: i * 2}}, 0, 0)
			}
		}},
		{F: func() { GetUserBalance(suite, wctx, 1, bobAccountNum) }},
		{F: func() { CreateBadges(suite, wctx, badgesToCreateAllInOne) }},
	})
}
