package keeper_test

import (
	"math"
	"math/rand"

	sdkmath "cosmossdk.io/math"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
)

func (suite *TestSuite) TestSafeAdd() {
	result, err := types.SafeAdd(sdkmath.NewUint(0), sdkmath.NewUint(1))
	suite.Require().Nil(err, "Error adding: %s")
	AssertUintsEqual(suite, result, sdkmath.NewUint(1))

	result, err = types.SafeAdd(sdkmath.NewUint(math.MaxUint64), sdkmath.NewUint(0))
	suite.Require().Nil(err, "Error adding: %s")
	AssertUintsEqual(suite, result, sdkmath.NewUint(math.MaxUint64))

	_, err = types.SafeAdd(sdkmath.NewUint(math.MaxUint), sdkmath.NewUint(1))
	suite.Require().Nil(err, "Error adding: %s")
	// AssertBalancesEqualEsuite, rror(err, types.ErrOverflow.Error()) With Cosmos SDK Uint now, this error is not returned
}

func (suite *TestSuite) TestSafeSubtract() {
	result, err := types.SafeSubtract(sdkmath.NewUint(1), sdkmath.NewUint(0))
	suite.Require().Nil(err, "Error adding: %s")
	AssertUintsEqual(suite, result, sdkmath.NewUint(1))

	result, err = types.SafeSubtract(sdkmath.NewUint(math.MaxUint64), sdkmath.NewUint(0))
	suite.Require().Nil(err, "Error adding: %s")
	AssertUintsEqual(suite, result, sdkmath.NewUint(math.MaxUint64))

	_, err = types.SafeSubtract(sdkmath.NewUint(1), sdkmath.NewUint(math.MaxUint64))
	suite.Require().Error(err, types.ErrUnderflow.Error())
}

func (suite *TestSuite) TestUpdateAndGetBalancesForIds() {
	err := *new(error)
	balances := []*types.Balance{
		{
			Amount:         sdkmath.NewUint(1),
			OwnershipTimes: GetFullUintRanges(),
			BadgeIds: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(1),
				},
			},
		},
	}

	balances, err = types.UpdateBalance(&types.Balance{
		BadgeIds: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(1),
			},
		},
		OwnershipTimes: GetFullUintRanges(),
		Amount:         sdkmath.NewUint(10),
	}, balances)
	suite.Require().Nil(err, "Error updating balances: %s")

	fetchedBalances, err := types.GetBalancesForIds([]*types.UintRange{
		{
			Start: sdkmath.NewUint(1),
			End:   sdkmath.NewUint(1),
		},
	}, GetFullUintRanges(), balances)
	suite.Require().Nil(err, "Error fetching balances: %s")

	AssertBalancesEqual(suite, balances, []*types.Balance{
		{
			Amount:         sdkmath.NewUint(10),
			OwnershipTimes: GetFullUintRanges(),
			BadgeIds: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(1),
				},
			},
		},
	})
	AssertBalancesEqual(suite, balances, fetchedBalances)

	fetchedBalances, err = types.GetBalancesForIds([]*types.UintRange{
		{
			Start: sdkmath.NewUint(1),
			End:   sdkmath.NewUint(1),
		},
	}, GetFullUintRanges(), balances)
	suite.Require().Nil(err, "Error fetching balances: %s")

	AssertBalancesEqual(suite, fetchedBalances, []*types.Balance{
		{
			Amount:         sdkmath.NewUint(10),
			OwnershipTimes: GetFullUintRanges(),
			BadgeIds: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(1),
				},
			},
		},
	})

	fetchedBalances, err = types.GetBalancesForIds([]*types.UintRange{
		{
			Start: sdkmath.NewUint(1),
			End:   sdkmath.NewUint(2),
		},
	}, GetFullUintRanges(), balances)
	suite.Require().Nil(err, "Error fetching balances: %s")

	AssertBalancesEqual(suite, fetchedBalances, []*types.Balance{
		{
			Amount:         sdkmath.NewUint(10),
			OwnershipTimes: GetFullUintRanges(),
			BadgeIds: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(1),
				},
			},
		},
		{
			Amount:         sdkmath.NewUint(0),
			OwnershipTimes: GetFullUintRanges(),
			BadgeIds: []*types.UintRange{
				{
					Start: sdkmath.NewUint(2),
					End:   sdkmath.NewUint(2),
				},
			},
		},
	})

	fetchedBalances, err = types.GetBalancesForIds([]*types.UintRange{
		{
			Start: sdkmath.NewUint(1),
			End:   sdkmath.NewUint(math.MaxUint64),
		},
	}, GetFullUintRanges(), balances)
	suite.Require().Nil(err, "Error fetching balances: %s")

	AssertBalancesEqual(suite, fetchedBalances, []*types.Balance{
		{
			Amount:         sdkmath.NewUint(0),
			OwnershipTimes: GetFullUintRanges(),
			BadgeIds: []*types.UintRange{
				{
					Start: sdkmath.NewUint(2),
					End:   sdkmath.NewUint(math.MaxUint64),
				},
			},
		},
		{
			Amount:         sdkmath.NewUint(10),
			OwnershipTimes: GetFullUintRanges(),
			BadgeIds: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(1),
				},
			},
		},
	})

	fetchedBalances, err = types.GetBalancesForIds([]*types.UintRange{
		{
			Start: sdkmath.NewUint(3),
			End:   sdkmath.NewUint(math.MaxUint64),
		},
		{
			Start: sdkmath.NewUint(1),
			End:   sdkmath.NewUint(2),
		},
		// {
		// 	Start: sdkmath.NewUint(1),
		// 	End:   sdkmath.NewUint(1),
		// },
	}, GetFullUintRanges(), balances)
	suite.Require().Nil(err, "Error fetching balances: %s")

	AssertBalancesEqual(suite, fetchedBalances, []*types.Balance{
		{
			Amount:         sdkmath.NewUint(0),
			OwnershipTimes: GetFullUintRanges(),
			BadgeIds: []*types.UintRange{
				{
					Start: sdkmath.NewUint(2),
					End:   sdkmath.NewUint(math.MaxUint64),
				},
			},
		},
		{
			Amount:         sdkmath.NewUint(10),
			OwnershipTimes: GetFullUintRanges(),
			BadgeIds: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(1),
				},
			},
		},
	})

	balances, err = types.UpdateBalance(&types.Balance{
		BadgeIds: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(1),
			},
		}, OwnershipTimes: GetFullUintRanges(), Amount: sdkmath.NewUint(5)}, balances)
	suite.Require().Nil(err, "Error fetching balances: %s")

	AssertBalancesEqual(suite, balances, []*types.Balance{
		{
			Amount:         sdkmath.NewUint(5),
			OwnershipTimes: GetFullUintRanges(),
			BadgeIds: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(1),
				},
			},
		},
		// {
		// 	Amount: sdkmath.NewUint(10),
		// 	OwnershipTimes: GetFullUintRanges(),
		// 	BadgeIds: []*types.UintRange{
		// 		{
		// 			Start: sdkmath.NewUint(1),
		// 			End:   sdkmath.NewUint(1),
		// 		},
		// 	},
		// },
	})

	balances, err = types.UpdateBalance(&types.Balance{
		BadgeIds: []*types.UintRange{
			{
				Start: sdkmath.NewUint(2),
				End:   sdkmath.NewUint(math.MaxUint64),
			},
		}, OwnershipTimes: GetFullUintRanges(), Amount: sdkmath.NewUint(5)}, balances)
	suite.Require().Nil(err, "Error fetching balances: %s")

	AssertBalancesEqual(suite, balances, []*types.Balance{
		{
			Amount:         sdkmath.NewUint(5),
			OwnershipTimes: GetFullUintRanges(),
			BadgeIds: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(math.MaxUint64),
				},
			},
		},
		// {
		// 	Amount: sdkmath.NewUint(10),
		// 	OwnershipTimes: GetFullUintRanges(),
		// 	BadgeIds: []*types.UintRange{
		// 		{
		// 			Start: sdkmath.NewUint(1),
		// 			End:   sdkmath.NewUint(1),
		// 		},
		// 	},
		// },
	})

	balances, err = types.UpdateBalance(&types.Balance{
		BadgeIds: []*types.UintRange{
			{
				Start: sdkmath.NewUint(2),
				End:   sdkmath.NewUint(2),
			},
		}, OwnershipTimes: GetFullUintRanges(), Amount: sdkmath.NewUint(10)}, balances)
	suite.Require().Nil(err, "Error fetching balances: %s")

	AssertBalancesEqual(suite, balances, []*types.Balance{
		{
			Amount:         sdkmath.NewUint(5),
			OwnershipTimes: GetFullUintRanges(),
			BadgeIds: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(1),
				},
				{
					Start: sdkmath.NewUint(3),
					End:   sdkmath.NewUint(math.MaxUint64),
				},
			},
		},
		{
			Amount:         sdkmath.NewUint(10),
			OwnershipTimes: GetFullUintRanges(),
			BadgeIds: []*types.UintRange{
				// {
				// 	Start: sdkmath.NewUint(1),
				// 	End:   sdkmath.NewUint(1),
				// },
				{
					Start: sdkmath.NewUint(2),
					End:   sdkmath.NewUint(2),
				},
			},
		},
	})
}

// Adjust these values to test more or less
const NUM_RUNS = 1
const NUM_IDS = 10
const NUM_OPERATIONS = 10

func (suite *TestSuite) TestBalancesFuzz() {
	for a := 0; a < NUM_RUNS; a++ {
		userBalance := &types.UserBalanceStore{}
		balances := make([]sdkmath.Uint, NUM_IDS)
		for i := 0; i < NUM_IDS; i++ {
			balances[i] = sdkmath.NewUint(0)
		}

		// adds := make([]*types.UintRange, NUM_OPERATIONS)
		// subs := make([]*types.UintRange, NUM_OPERATIONS)
		for i := 0; i < NUM_OPERATIONS; i++ { //NUM_OPERATIONS iterations
			//Get random start value
			start := (uint64(rand.Intn(NUM_IDS / 2)))
			//Get random end value
			end := (uint64(NUM_IDS/2 + rand.Intn(NUM_IDS/2)))

			amount := sdkmath.NewUint(uint64(rand.Intn(100)))
			err := *new(error)

			userBalance.Balances, err = types.AddBalance(userBalance.Balances, &types.Balance{
				Amount: amount,
				BadgeIds: []*types.UintRange{
					{
						Start: sdkmath.NewUint(start),
						End:   sdkmath.NewUint(end),
					},
				},
				OwnershipTimes: GetFullUintRanges(),
			})
			suite.Require().Nil(err, "error adding balance to approval")

			// adds = append(adds, &types.UintRange{
			// 	Start: sdkmath.NewUint(start),
			// 	End:   sdkmath.NewUint(end),
			// })
			// println("adding", start, end, amount.String())

			for j := start; j <= end; j++ {
				balances[j] = balances[j].Add(amount)
			}

			start = (uint64(rand.Intn(NUM_IDS / 2)))
			end = (uint64(NUM_IDS/2 + rand.Intn(NUM_IDS/2)))
			amount = sdkmath.NewUint(uint64(rand.Intn(20))) //Make this substantially less than add, so we have less chance of underflow
			// println("removing", start, end, amount.String())

			userBalancesCopy := make([]*types.Balance, len(userBalance.Balances))
			copy(userBalancesCopy, userBalance.Balances)

			userBalance.Balances, err = types.SubtractBalance(userBalance.Balances, &types.Balance{
				Amount: amount,
				BadgeIds: []*types.UintRange{
					{
						Start: sdkmath.NewUint(start),
						End:   sdkmath.NewUint(end),
					},
				},
				OwnershipTimes: GetFullUintRanges(),
			})

			if err != nil {
				suite.Require().EqualError(err, types.ErrUnderflow.Error())
				userBalance.Balances = userBalancesCopy //revert to previous balances
			} else {
				// if sdkmath.NewUint(256) < start && sdkmath.NewUint(256) < end {
				// 	println("removing", start, end, amount)
				// }
				// subs = append(subs, &types.UintRange{
				// 	Start: sdkmath.NewUint(start),
				// 	End:   sdkmath.NewUint(end),
				// })

				for j := start; j <= end; j++ {
					balances[j] = balances[j].Sub(amount)
				}
			}

		}

		for i := 0; i < NUM_IDS; i++ {
			fetchedBalances, _ := types.GetBalancesForIds(
				[]*types.UintRange{
					{
						Start: sdkmath.NewUint(uint64(i)),
						End:   sdkmath.NewUint(uint64(i)),
					},
				},
				GetFullUintRanges(),
				userBalance.Balances,
			)

			// println("fetched", i, fetchedBalances[0].Amount.String(), balances[i].String())
			AssertUintsEqual(suite, fetchedBalances[0].Amount, balances[i])
		}
	}
}

/* --------------------------------------START TESTING WITH TIMES-------------------------------------- */
//Previously, everything was just FullUintRanges() for times

//Adjust these values to test more or less

func (suite *TestSuite) TestBalancesWithTimesFuzz() {
	for a := 0; a < NUM_RUNS; a++ {
		userBalance := &types.UserBalanceStore{}
		balances := make([][]sdkmath.Uint, NUM_IDS)
		for i := 0; i < NUM_IDS; i++ {
			balances[i] = make([]sdkmath.Uint, NUM_IDS)
			for j := 0; j < NUM_IDS; j++ {
				balances[i][j] = sdkmath.NewUint(0)
			}
		}

		// adds := make([]*types.UintRange, NUM_OPERATIONS)
		// subs := make([]*types.UintRange, NUM_OPERATIONS)
		for i := 0; i < NUM_OPERATIONS; i++ { //NUM_OPERATIONS iterations
			start := (uint64(rand.Intn(NUM_IDS / 2)))
			end := (uint64(NUM_IDS/2 + rand.Intn(NUM_IDS/2)))
			startTime := (uint64(rand.Intn(NUM_IDS / 2)))
			endTime := (uint64(NUM_IDS/2 + rand.Intn(NUM_IDS/2)))

			amount := sdkmath.NewUint(uint64(rand.Intn(100)))
			err := *new(error)

			userBalance.Balances, err = types.AddBalance(userBalance.Balances, &types.Balance{
				Amount: amount,
				BadgeIds: []*types.UintRange{
					{
						Start: sdkmath.NewUint(start),
						End:   sdkmath.NewUint(end),
					},
				},
				OwnershipTimes: []*types.UintRange{
					{
						Start: sdkmath.NewUint(startTime),
						End:   sdkmath.NewUint(endTime),
					},
				},
			})
			suite.Require().Nil(err, "error adding balance to approval")

			// adds = append(adds, &types.UintRange{
			// 	Start: sdkmath.NewUint(start),
			// 	End:   sdkmath.NewUint(end),
			// })
			// println("adding", start, end, startTime, endTime, amount.String())

			for j := start; j <= end; j++ {
				for k := startTime; k <= endTime; k++ {
					balances[j][k] = balances[j][k].Add(amount)
				}
			}

			start = (uint64(rand.Intn(NUM_IDS / 2)))
			end = (uint64(NUM_IDS/2 + rand.Intn(NUM_IDS/2)))
			amount = sdkmath.NewUint(uint64(rand.Intn(20))) //Make this substantially less than add, so we have less chance of underflow
			startTime = (uint64(rand.Intn(NUM_IDS / 2)))
			endTime = (uint64(NUM_IDS/2 + rand.Intn(NUM_IDS/2)))
			// println("removing", start, end, startTime, endTime, amount.String())

			userBalancesCopy := make([]*types.Balance, len(userBalance.Balances))
			copy(userBalancesCopy, userBalance.Balances)

			//removing 18 from IDs 1-7 Times 1-6

			userBalance.Balances, err = types.SubtractBalance(userBalance.Balances, &types.Balance{
				Amount: amount,
				BadgeIds: []*types.UintRange{
					{
						Start: sdkmath.NewUint(start),
						End:   sdkmath.NewUint(end),
					},
				},
				OwnershipTimes: []*types.UintRange{
					{
						Start: sdkmath.NewUint(startTime),
						End:   sdkmath.NewUint(endTime),
					},
				},
			})

			if err != nil {
				suite.Require().EqualError(err, types.ErrUnderflow.Error())
				userBalance.Balances = userBalancesCopy //revert to previous balances
				// println("reverted")
			} else {
				// if sdkmath.NewUint(256) < start && sdkmath.NewUint(256) < end {
				// 	println("removing", start, end, amount)
				// }
				// subs = append(subs, &types.UintRange{
				// 	Start: sdkmath.NewUint(start),
				// 	End:   sdkmath.NewUint(end),
				// })

				for j := start; j <= end; j++ {
					for k := startTime; k <= endTime; k++ {
						balances[j][k] = balances[j][k].Sub(amount)
					}
				}
			}
		}

		for i := 0; i < NUM_IDS; i++ {
			for j := 0; j < NUM_IDS; j++ {
				fetchedBalances, _ := types.GetBalancesForIds(
					[]*types.UintRange{
						{
							Start: sdkmath.NewUint(uint64(i)),
							End:   sdkmath.NewUint(uint64(i)),
						},
					},
					[]*types.UintRange{
						{
							Start: sdkmath.NewUint(uint64(j)),
							End:   sdkmath.NewUint(uint64(j)),
						},
					},
					userBalance.Balances,
				)

				//1 1 12 13
				// add 31
				// remove 13
				// remove 18
				// remove 5
				// println("fetched", i, j, fetchedBalances[0].Amount.String(), balances[i][j].String())
				AssertUintsEqual(suite, fetchedBalances[0].Amount, balances[i][j])
			}
		}
	}
}
