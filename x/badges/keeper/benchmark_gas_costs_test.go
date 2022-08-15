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

const PRINT_MODE = false

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

	if PRINT_MODE {
		firstColumnFormatter := func(format string, vals ...interface{}) string {
			return ""
		}

		tbl.WithFirstColumnFormatter(firstColumnFormatter).Print()
	}
	
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


func (suite *TestSuite) TestGasCostsOldVersion() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	for i := 0; i < 1; i++ {
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
		tbl := table.New("Function Name", "Gas Consumed")

		startGas := suite.ctx.GasMeter().GasConsumed()
		CreateBadges(suite, wctx, badgesToCreate)
		endGas := suite.ctx.GasMeter().GasConsumed()
		tbl.AddRow("CreateBadge", endGas-startGas)

		startGas = suite.ctx.GasMeter().GasConsumed()
		_, _ = GetBadge(suite, wctx, 0)
		endGas = suite.ctx.GasMeter().GasConsumed()
		tbl.AddRow("GetBadge", endGas-startGas)

		startGas = suite.ctx.GasMeter().GasConsumed()
		err := CreateSubBadges(suite, wctx, bob, 0, []uint64{1000000}, []uint64{1})
		suite.Require().Nil(err, "Error creating subbadge")
		endGas = suite.ctx.GasMeter().GasConsumed()
		tbl.AddRow("CreateSubBadge 1 (Supply 10000)", endGas-startGas)

		startGas = suite.ctx.GasMeter().GasConsumed()
		err = CreateSubBadges(suite, wctx, bob, 0, []uint64{1}, []uint64{10000})
		suite.Require().Nil(err, "Error creating subbadge")
		endGas = suite.ctx.GasMeter().GasConsumed()
		tbl.AddRow("CreateSubBadge 10000 (Supply 1)", endGas-startGas)

		startGas = suite.ctx.GasMeter().GasConsumed()
		err = CreateSubBadges(suite, wctx, bob, 0, []uint64{10000}, []uint64{10000})
		suite.Require().Nil(err, "Error creating subbadge")
		endGas = suite.ctx.GasMeter().GasConsumed()
		tbl.AddRow("CreateSubBadge 10000 (Supply 10000)", endGas-startGas)

		startGas = suite.ctx.GasMeter().GasConsumed()
		err = FreezeAddresses(suite, wctx, bob, []*types.IdRange{{Start: 1000, End: 1000}}, 0, true)
		suite.Require().Nil(err, "Error creating subbadge")
		endGas = suite.ctx.GasMeter().GasConsumed()
		tbl.AddRow("Freeze 1 Address", endGas-startGas)

		startGas = suite.ctx.GasMeter().GasConsumed()
		_, _ = GetBadge(suite, wctx, 0)
		endGas = suite.ctx.GasMeter().GasConsumed()
		tbl.AddRow("GetBadge", endGas-startGas)

		addressesThroughTenThousand := make([]uint64, 10000)
		for j := 0; j < 10000; j++ {
			addressesThroughTenThousand[j] = uint64(j)
		}
		startGas = suite.ctx.GasMeter().GasConsumed()
		err = FreezeAddresses(suite, wctx, bob, []*types.IdRange{{Start: 0, End: 9999}}, 0, true)
		suite.Require().Nil(err, "Error creating subbadge")
		endGas = suite.ctx.GasMeter().GasConsumed()
		tbl.AddRow("Freeze 10000 Addresses", endGas-startGas)

		startGas = suite.ctx.GasMeter().GasConsumed()
		_, _ = GetBadge(suite, wctx, 0)
		endGas = suite.ctx.GasMeter().GasConsumed()
		tbl.AddRow("GetBadge", endGas-startGas)

		startGas = suite.ctx.GasMeter().GasConsumed()
		err = FreezeAddresses(suite, wctx, bob, []*types.IdRange{{Start: 1000, End: 1000}}, 0, false)
		suite.Require().Nil(err, "Error creating subbadge")
		endGas = suite.ctx.GasMeter().GasConsumed()
		tbl.AddRow("Unfreeze 1 Address", endGas-startGas)

		startGas = suite.ctx.GasMeter().GasConsumed()
		_, _ = GetBadge(suite, wctx, 0)
		endGas = suite.ctx.GasMeter().GasConsumed()
		tbl.AddRow("GetBadge", endGas-startGas)

		startGas = suite.ctx.GasMeter().GasConsumed()
		err = FreezeAddresses(suite, wctx, bob, []*types.IdRange{{Start: 0, End: 9999}}, 0, false)
		suite.Require().Nil(err, "Error creating subbadge")
		endGas = suite.ctx.GasMeter().GasConsumed()
		tbl.AddRow("Unfreeze 10000 Address", endGas-startGas)

		startGas = suite.ctx.GasMeter().GasConsumed()
		_, _ = GetBadge(suite, wctx, 0)
		endGas = suite.ctx.GasMeter().GasConsumed()
		tbl.AddRow("GetBadge", endGas-startGas)

		startGas = suite.ctx.GasMeter().GasConsumed()
		err = TransferBadge(suite, wctx, bob, bobAccountNum, []uint64{aliceAccountNum}, []uint64{1}, 0, []*types.IdRange{{Start: 0, End: 0}}, 0, 0)
		suite.Require().Nil(err, "Error transferring badge")
		endGas = suite.ctx.GasMeter().GasConsumed()
		tbl.AddRow("TransferBadge - Pending 1", endGas-startGas)

		startGas = suite.ctx.GasMeter().GasConsumed()
		_, _ = GetUserBalance(suite, wctx, 0, aliceAccountNum)
		endGas = suite.ctx.GasMeter().GasConsumed()
		tbl.AddRow("GetUserBalance", endGas-startGas)

		startGas = suite.ctx.GasMeter().GasConsumed()
		err = HandlePendingTransfers(suite, wctx, alice, true, 0, []*types.IdRange{{Start: 0, End: 0}}, false)
		suite.Require().Nil(err, "Error accepting badge")
		endGas = suite.ctx.GasMeter().GasConsumed()
		tbl.AddRow("HandlePendingTransfer - 1", endGas-startGas)

		startGas = suite.ctx.GasMeter().GasConsumed()
		_, _ = GetUserBalance(suite, wctx, 0, aliceAccountNum)
		endGas = suite.ctx.GasMeter().GasConsumed()
		tbl.AddRow("GetUserBalance", endGas-startGas)

		startGas = suite.ctx.GasMeter().GasConsumed()
		err = TransferBadge(suite, wctx, bob, bobAccountNum, []uint64{aliceAccountNum}, []uint64{1}, 0, []*types.IdRange{{Start: 0, End: 999}}, 0, 0)
		suite.Require().Nil(err, "Error transferring badge")
		endGas = suite.ctx.GasMeter().GasConsumed()
		tbl.AddRow("TransferBadge - Pending 1000 Diff IDs", endGas-startGas)

		startGas = suite.ctx.GasMeter().GasConsumed()
		_, _ = GetUserBalance(suite, wctx, 0, aliceAccountNum)
		endGas = suite.ctx.GasMeter().GasConsumed()
		tbl.AddRow("GetUserBalance", endGas-startGas)

		startGas = suite.ctx.GasMeter().GasConsumed()
		err = HandlePendingTransfers(suite, wctx, alice, true, 0, []*types.IdRange{{Start: 0, End: 999}}, false)
		suite.Require().Nil(err, "Error transferring badge")
		endGas = suite.ctx.GasMeter().GasConsumed()
		tbl.AddRow("HandlePendingTransfer - 1000", endGas-startGas)

		startGas = suite.ctx.GasMeter().GasConsumed()
		_, _ = GetUserBalance(suite, wctx, 0, aliceAccountNum)
		endGas = suite.ctx.GasMeter().GasConsumed()
		tbl.AddRow("GetUserBalance", endGas-startGas)

		startGas = suite.ctx.GasMeter().GasConsumed()
		err = TransferBadge(suite, wctx, bob, bobAccountNum, []uint64{aliceAccountNum}, []uint64{1}, 0, []*types.IdRange{{Start: 0, End: 0}}, 0, 0)
		suite.Require().Nil(err, "Error transferring badge")
		endGas = suite.ctx.GasMeter().GasConsumed()
		tbl.AddRow("TransferBadge - Pending 1", endGas-startGas)

		startGas = suite.ctx.GasMeter().GasConsumed()
		_, _ = GetUserBalance(suite, wctx, 0, aliceAccountNum)
		endGas = suite.ctx.GasMeter().GasConsumed()
		tbl.AddRow("GetUserBalance", endGas-startGas)

		startGas = suite.ctx.GasMeter().GasConsumed()
		err = HandlePendingTransfers(suite, wctx, alice, false, 0, []*types.IdRange{{Start: 2, End: 2}}, false)
		suite.Require().Nil(err, "Error accepting badge")
		endGas = suite.ctx.GasMeter().GasConsumed()
		tbl.AddRow("HandlePendingTransfer (Reject) - 1", endGas-startGas)
		err = HandlePendingTransfers(suite, wctx, bob, false, 0, []*types.IdRange{{Start: 2, End: 2}}, false)
		suite.Require().Nil(err, "Error accepting badge")

		startGas = suite.ctx.GasMeter().GasConsumed()
		_, _ = GetUserBalance(suite, wctx, 0, aliceAccountNum)
		endGas = suite.ctx.GasMeter().GasConsumed()
		tbl.AddRow("GetUserBalance", endGas-startGas)

		startGas = suite.ctx.GasMeter().GasConsumed()
		err = TransferBadge(suite, wctx, bob, bobAccountNum, []uint64{aliceAccountNum}, []uint64{1}, 0, []*types.IdRange{{Start: 1000, End: 1999}}, 0, 0)
		suite.Require().Nil(err, "Error transferring badge")
		endGas = suite.ctx.GasMeter().GasConsumed()
		tbl.AddRow("TransferBadge - Pending 1000", endGas-startGas)

		startGas = suite.ctx.GasMeter().GasConsumed()
		_, _ = GetUserBalance(suite, wctx, 0, aliceAccountNum)
		endGas = suite.ctx.GasMeter().GasConsumed()
		tbl.AddRow("GetUserBalance", endGas-startGas)

		startGas = suite.ctx.GasMeter().GasConsumed()
		err = HandlePendingTransfers(suite, wctx, alice, false, 0, []*types.IdRange{{Start: 0, End: 999}}, false)
		suite.Require().Nil(err, "Error transferring badge")
		endGas = suite.ctx.GasMeter().GasConsumed()
		tbl.AddRow("HandlePendingTransfer (Reject) - 1000", endGas-startGas)
		err = HandlePendingTransfers(suite, wctx, bob, false, 0, []*types.IdRange{{Start: 0, End: 999}}, true)
		suite.Require().Nil(err, "Error transferring badge")

		startGas = suite.ctx.GasMeter().GasConsumed()
		_, _ = GetUserBalance(suite, wctx, 0, aliceAccountNum)
		endGas = suite.ctx.GasMeter().GasConsumed()
		tbl.AddRow("GetUserBalance", endGas-startGas)

		startGas = suite.ctx.GasMeter().GasConsumed()
		err = RequestTransferBadge(suite, wctx, alice, bobAccountNum, 1, 0, []*types.IdRange{{Start: 0, End: 0}}, 0, 0)
		suite.Require().Nil(err, "Error transferring badge")
		endGas = suite.ctx.GasMeter().GasConsumed()
		tbl.AddRow("RequestTransferBadge - 1", endGas-startGas)

		startGas = suite.ctx.GasMeter().GasConsumed()
		_, _ = GetUserBalance(suite, wctx, 0, aliceAccountNum)
		endGas = suite.ctx.GasMeter().GasConsumed()
		tbl.AddRow("GetUserBalance", endGas-startGas)

		startGas = suite.ctx.GasMeter().GasConsumed()
		err = HandlePendingTransfers(suite, wctx, bob, true, 0, []*types.IdRange{{Start: 4, End: 4}}, true)
		suite.Require().Nil(err, "Error accepting badge")
		endGas = suite.ctx.GasMeter().GasConsumed()
		tbl.AddRow("Accept RequestTransfer - 1", endGas-startGas)
		// err = HandlePendingTransfers(suite, wctx, alice, true, 0, []*types.IdRange{{Start: 0, End: 0}}, 0, 0)
		// suite.Require().Nil(err, "Error accepting badge")

		startGas = suite.ctx.GasMeter().GasConsumed()
		_, _ = GetUserBalance(suite, wctx, 0, aliceAccountNum)
		endGas = suite.ctx.GasMeter().GasConsumed()
		tbl.AddRow("GetUserBalance", endGas-startGas)

		startGas = suite.ctx.GasMeter().GasConsumed()
		err = RequestTransferBadge(suite, wctx, alice, bobAccountNum, 1, 0, []*types.IdRange{{Start: 1000, End: 1999}}, 0, 0)
		suite.Require().Nil(err, "Error transferring badge")
		endGas = suite.ctx.GasMeter().GasConsumed()
		tbl.AddRow("Request Transfer- 1000", endGas-startGas)

		startGas = suite.ctx.GasMeter().GasConsumed()
		_, _ = GetUserBalance(suite, wctx, 0, aliceAccountNum)
		endGas = suite.ctx.GasMeter().GasConsumed()
		tbl.AddRow("GetUserBalance", endGas-startGas)

		startGas = suite.ctx.GasMeter().GasConsumed()
		err = HandlePendingTransfers(suite, wctx, bob, true, 0, []*types.IdRange{{Start: 0, End: 999}}, true)
		suite.Require().Nil(err, "Error transferring badge")
		endGas = suite.ctx.GasMeter().GasConsumed()
		tbl.AddRow("Accept RequestTransfer - 1000", endGas-startGas)
		// err = HandlePendingTransfers(suite, wctx, alice, true, 0, []*types.IdRange{{Start: 0, End: 999}}, 0, 0)
		// suite.Require().Nil(err, "Error transferring badge")

		startGas = suite.ctx.GasMeter().GasConsumed()
		_, _ = GetUserBalance(suite, wctx, 0, aliceAccountNum)
		endGas = suite.ctx.GasMeter().GasConsumed()
		tbl.AddRow("GetUserBalance", endGas-startGas)

		startGas = suite.ctx.GasMeter().GasConsumed()
		err = RequestTransferBadge(suite, wctx, alice, bobAccountNum, 1, 0, []*types.IdRange{{Start: 0, End: 0}}, 0, 0)
		suite.Require().Nil(err, "Error transferring badge")
		endGas = suite.ctx.GasMeter().GasConsumed()
		tbl.AddRow("RequestTransferBadge - 1", endGas-startGas)

		startGas = suite.ctx.GasMeter().GasConsumed()
		_, _ = GetUserBalance(suite, wctx, 0, aliceAccountNum)
		endGas = suite.ctx.GasMeter().GasConsumed()
		tbl.AddRow("GetUserBalance", endGas-startGas)

		startGas = suite.ctx.GasMeter().GasConsumed()
		err = HandlePendingTransfers(suite, wctx, bob, false, 0, []*types.IdRange{{Start: 0, End: 999}}, false)
		suite.Require().Nil(err, "Error accepting badge")
		endGas = suite.ctx.GasMeter().GasConsumed()
		tbl.AddRow("Reject RequestTransfer - 1", endGas-startGas)
		err = HandlePendingTransfers(suite, wctx, alice, false, 0, []*types.IdRange{{Start: 0, End: 999}}, false)
		suite.Require().Nil(err, "Error accepting badge")

		startGas = suite.ctx.GasMeter().GasConsumed()
		_, _ = GetUserBalance(suite, wctx, 0, aliceAccountNum)
		endGas = suite.ctx.GasMeter().GasConsumed()
		tbl.AddRow("GetUserBalance", endGas-startGas)

		startGas = suite.ctx.GasMeter().GasConsumed()
		err = RequestTransferBadge(suite, wctx, alice, bobAccountNum, 1, 0, []*types.IdRange{{Start: 0, End: 999}}, 0, 0)
		suite.Require().Nil(err, "Error transferring badge")
		endGas = suite.ctx.GasMeter().GasConsumed()
		tbl.AddRow("Request Transfer- 1000", endGas-startGas)

		startGas = suite.ctx.GasMeter().GasConsumed()
		_, _ = GetUserBalance(suite, wctx, 0, aliceAccountNum)
		endGas = suite.ctx.GasMeter().GasConsumed()
		tbl.AddRow("GetUserBalance", endGas-startGas)

		startGas = suite.ctx.GasMeter().GasConsumed()
		err = HandlePendingTransfers(suite, wctx, bob, false, 0, []*types.IdRange{{Start: 0, End: 999}}, false)
		suite.Require().Nil(err, "Error transferring badge")
		endGas = suite.ctx.GasMeter().GasConsumed()
		tbl.AddRow("Reject RequestTransfer - 1000", endGas-startGas)
		err = HandlePendingTransfers(suite, wctx, alice, false, 0, []*types.IdRange{{Start: 0, End: 999}}, false)
		suite.Require().Nil(err, "Error accepting badge")

		startGas = suite.ctx.GasMeter().GasConsumed()
		_, _ = GetUserBalance(suite, wctx, 0, aliceAccountNum)
		endGas = suite.ctx.GasMeter().GasConsumed()
		tbl.AddRow("GetUserBalance", endGas-startGas)

		startGas = suite.ctx.GasMeter().GasConsumed()
		_, _ = GetUserBalance(suite, wctx, 0, aliceAccountNum)
		endGas = suite.ctx.GasMeter().GasConsumed()
		tbl.AddRow("GetUserBalance", endGas-startGas)

		startGas = suite.ctx.GasMeter().GasConsumed()
		err = RevokeBadges(suite, wctx, bob, []uint64{aliceAccountNum}, []uint64{1}, 0, []*types.IdRange{{Start: 0, End: 0}})
		suite.Require().Nil(err, "Error transferring badge")
		endGas = suite.ctx.GasMeter().GasConsumed()
		tbl.AddRow("Revoke Badge 1", endGas-startGas)

		startGas = suite.ctx.GasMeter().GasConsumed()
		_, _ = GetUserBalance(suite, wctx, 0, aliceAccountNum)
		endGas = suite.ctx.GasMeter().GasConsumed()
		tbl.AddRow("GetUserBalance", endGas-startGas)

		startGas = suite.ctx.GasMeter().GasConsumed()
		err = RevokeBadges(suite, wctx, bob, []uint64{aliceAccountNum}, []uint64{1}, 0, []*types.IdRange{{Start: 1000, End: 1999}})
		suite.Require().Nil(err, "Error transferring badge")
		endGas = suite.ctx.GasMeter().GasConsumed()
		tbl.AddRow("Revoke Badge 1000", endGas-startGas)

		startGas = suite.ctx.GasMeter().GasConsumed()
		_, _ = GetUserBalance(suite, wctx, 0, aliceAccountNum)
		endGas = suite.ctx.GasMeter().GasConsumed()
		tbl.AddRow("GetUserBalance", endGas-startGas)

		startGas = suite.ctx.GasMeter().GasConsumed()
		_, _ = GetUserBalance(suite, wctx, 0, bobAccountNum)
		endGas = suite.ctx.GasMeter().GasConsumed()
		tbl.AddRow("GetUserBalance", endGas-startGas)

		startGas = suite.ctx.GasMeter().GasConsumed()
		err = SetApproval(suite, wctx, bob, 10000, aliceAccountNum, 0, []*types.IdRange{{Start: 0, End: 0}})
		suite.Require().Nil(err, "Error transferring badge")
		endGas = suite.ctx.GasMeter().GasConsumed()
		tbl.AddRow("Set Approval 1", endGas-startGas)

		startGas = suite.ctx.GasMeter().GasConsumed()
		_, _ = GetUserBalance(suite, wctx, 0, bobAccountNum)
		endGas = suite.ctx.GasMeter().GasConsumed()
		tbl.AddRow("GetUserBalance", endGas-startGas)

		startGas = suite.ctx.GasMeter().GasConsumed()
		err = SetApproval(suite, wctx, bob, 10, aliceAccountNum, 0, []*types.IdRange{{Start: 0, End: 0}})
		suite.Require().Nil(err, "Error transferring badge")
		endGas = suite.ctx.GasMeter().GasConsumed()
		tbl.AddRow("Remove Approval 1", endGas-startGas)

		startGas = suite.ctx.GasMeter().GasConsumed()
		_, _ = GetUserBalance(suite, wctx, 0, bobAccountNum)
		endGas = suite.ctx.GasMeter().GasConsumed()
		tbl.AddRow("GetUserBalance", endGas-startGas)

		startGas = suite.ctx.GasMeter().GasConsumed()
		err = SetApproval(suite, wctx, bob, 10, aliceAccountNum, 0, []*types.IdRange{{Start: 0, End: 999}})
		suite.Require().Nil(err, "Error transferring badge")
		endGas = suite.ctx.GasMeter().GasConsumed()
		tbl.AddRow("Set Approval 1000 diff Subbadge IDs", endGas-startGas)

		startGas = suite.ctx.GasMeter().GasConsumed()
		_, _ = GetUserBalance(suite, wctx, 0, bobAccountNum)
		endGas = suite.ctx.GasMeter().GasConsumed()
		tbl.AddRow("GetUserBalance", endGas-startGas)

		badgesToCreate = []BadgesToCreate{
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

		startGas = suite.ctx.GasMeter().GasConsumed()
		CreateBadges(suite, wctx, badgesToCreate)
		endGas = suite.ctx.GasMeter().GasConsumed()
		tbl.AddRow("CreateBadge", endGas-startGas)

		startGas = suite.ctx.GasMeter().GasConsumed()
		_, _ = GetBadge(suite, wctx, 1)
		endGas = suite.ctx.GasMeter().GasConsumed()
		tbl.AddRow("GetBadge", endGas-startGas)

		startGas = suite.ctx.GasMeter().GasConsumed()
		err = CreateSubBadges(suite, wctx, bob, 1, []uint64{10000}, []uint64{10000})
		suite.Require().Nil(err, "Error creating subbadge")
		endGas = suite.ctx.GasMeter().GasConsumed()
		tbl.AddRow("CreateSubBadge 10000 (Supply 10000)", endGas-startGas)

		startGas = suite.ctx.GasMeter().GasConsumed()
		err = TransferBadge(suite, wctx, bob, bobAccountNum, []uint64{aliceAccountNum}, []uint64{1}, 1, []*types.IdRange{{Start: 0, End: 0}}, 0, 0)
		suite.Require().Nil(err, "Error transferring badge")
		endGas = suite.ctx.GasMeter().GasConsumed()
		tbl.AddRow("TransferBadge - Forceful 1", endGas-startGas)

		startGas = suite.ctx.GasMeter().GasConsumed()
		_, _ = GetUserBalance(suite, wctx, 0, aliceAccountNum)
		endGas = suite.ctx.GasMeter().GasConsumed()
		tbl.AddRow("GetUserBalance", endGas-startGas)

		startGas = suite.ctx.GasMeter().GasConsumed()
		err = TransferBadge(suite, wctx, bob, bobAccountNum, []uint64{aliceAccountNum}, []uint64{1}, 1, []*types.IdRange{{Start: 0, End: 999}}, 0, 0)
		suite.Require().Nil(err, "Error transferring badge")
		endGas = suite.ctx.GasMeter().GasConsumed()
		tbl.AddRow("TransferBadge - Forceful 1000 (Same Address, Diff Subbadge ID)", endGas-startGas)

		startGas = suite.ctx.GasMeter().GasConsumed()
		_, _ = GetUserBalance(suite, wctx, 0, aliceAccountNum)
		endGas = suite.ctx.GasMeter().GasConsumed()
		tbl.AddRow("GetUserBalance", endGas-startGas)

		startGas = suite.ctx.GasMeter().GasConsumed()
		addresses := []uint64{}
		for i := 0; i < 1000; i++ {
			addresses = append(addresses, aliceAccountNum+uint64(i))
		}
		err = TransferBadge(suite, wctx, bob, bobAccountNum, addresses, []uint64{1}, 1, []*types.IdRange{{Start: 0, End: 0}}, 0, 0)
		suite.Require().Nil(err, "Error transferring badge")
		suite.Require().Nil(err, "Error transferring badge")
		endGas = suite.ctx.GasMeter().GasConsumed()
		tbl.AddRow("TransferBadge - Forceful 1000 (Different Addresses, Same SubbadgeId)", endGas-startGas)

		startGas = suite.ctx.GasMeter().GasConsumed()
		_, _ = GetUserBalance(suite, wctx, 0, aliceAccountNum)
		endGas = suite.ctx.GasMeter().GasConsumed()
		tbl.AddRow("GetUserBalance", endGas-startGas)

		err = TransferBadge(suite, wctx, bob, bobAccountNum, addresses, []uint64{1}, 1, []*types.IdRange{{Start: 0, End: 999}}, 0, 0)
		suite.Require().Nil(err, "Error transferring badge")
		suite.Require().Nil(err, "Error transferring badge")
		endGas = suite.ctx.GasMeter().GasConsumed()
		tbl.AddRow("TransferBadge - Forceful 1000 (Different Addresses, Diff SubbadgeId)", endGas-startGas)

		startGas = suite.ctx.GasMeter().GasConsumed()
		_, _ = GetUserBalance(suite, wctx, 0, aliceAccountNum)
		endGas = suite.ctx.GasMeter().GasConsumed()
		tbl.AddRow("GetUserBalance", endGas-startGas)

		startGas = suite.ctx.GasMeter().GasConsumed()
		for i := uint64(0); i < 1000; i++ {
			err = TransferBadge(suite, wctx, bob, bobAccountNum, []uint64{aliceAccountNum}, []uint64{1}, 1, []*types.IdRange{{Start: i * 2, End: i * 2}}, 0, 0)
			suite.Require().Nil(err, "Error transferring badge")
			suite.Require().Nil(err, "Error transferring badge")
		}
		endGas = suite.ctx.GasMeter().GasConsumed()
		tbl.AddRow("TransferBadge - 1000 alternating subbadge IDs", endGas-startGas)

		startGas = suite.ctx.GasMeter().GasConsumed()
		userBalanceInfo, _ := GetUserBalance(suite, wctx, 1, bobAccountNum)
		endGas = suite.ctx.GasMeter().GasConsumed()
		tbl.AddRow("GetUserBalance", endGas-startGas)

		_ = userBalanceInfo

		firstColumnFormatter := func(format string, vals ...interface{}) string {
			return ""
		}

		tbl.WithFirstColumnFormatter(firstColumnFormatter).Print()

	}
}