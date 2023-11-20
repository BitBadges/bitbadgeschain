package eip712

// These are used as fully populated examples to generate EIP712 types
// We want to make sure the type generation includes all default values and empty values
// because that is what the SDK does
// So, we use these filled out schemas and then populate all empty types
func GetSchemas() []string {
	schemas := make([]string, 0)

	schemas = append(schemas, `{
		"type": "badges/CreateAddressMappings",
		"value": {
			"creator": "",
			"addressMappings": [
				{
					"mappingId": "",
					"addresses": [],
					"includeAddresses": false,
					"uri": "",
					"customData": "",
					"createdBy": ""
				}
			]
		}
	}`)

	schemas = append(schemas, `{
		"type": "badges/DeleteCollection",
		"value": {
			"creator": "",
			"collectionId": ""
		}
	}`)

	schemas = append(schemas, `{
		"type": "badges/TransferBadges",
		"value": {
			"creator": "",
			"collectionId": "",
			"transfers": [
				{
					"from": "",
					"toAddresses": [],
					"balances": [
						{
							"amount": "",
							"ownershipTimes": [
								{
									"start": "",
									"end": ""
								}
							],
							"badgeIds": [
								{
									"start": "",
									"end": ""
								}
							]
						}
					],
					"precalculateBalancesFromApproval": {
						"approvalId": "",
						"approvalLevel": "",
						"approverAddress": ""
					},
					"merkleProofs": [
						{
							"leaf": "",
							"aunts": [
								{
									"aunt": "",
									"onRight": false
								}
							]
						}
					],
					"memo": "",
					"prioritizedApprovals": [
						{
							"approvalId": "",
							"approvalLevel": "",
							"approverAddress": ""
						}
					],
					"onlyCheckPrioritizedApprovals": false
				}
			]
		}
	}`)

	schemas = append(schemas, `{
		"type": "badges/UpdateCollection",
		"value": {
			"creator": "",
			"collectionId": "",
			"balancesType": "",
			"defaultOutgoingApprovals": [
				{
					"toMappingId": "",
					"initiatedByMappingId": "",
					"transferTimes": [
						{
							"start": "",
							"end": ""
						}
					],
					"badgeIds": [
						{
							"start": "",
							"end": ""
						}
					],
					"ownershipTimes": [
						{
							"start": "",
							"end": ""
						}
					],
					"amountTrackerId": "",
					"challengeTrackerId": "",
					"uri": "",
					"customData": "",
					"approvalId": "",
					"approvalCriteria": {
						"mustOwnBadges": [
							{
								"collectionId": "",
								"amountRange": {
									"start": "",
									"end": ""
								},
								"ownershipTimes": [
									{
										"start": "",
										"end": ""
									}
								],
								"badgeIds": [
									{
										"start": "",
										"end": ""
									}
								],
								"overrideWithCurrentTime": false,
								"mustOwnAll": false
							}
						],
						"merkleChallenge": {
							"root": "",
							"expectedProofLength": "",
							"useCreatorAddressAsLeaf": false,
							"maxUsesPerLeaf": "",
							"uri": "",
							"customData": ""
						},
						"predeterminedBalances": {
							"manualBalances": [
								{
									"balances": [
										{
											"amount": "",
											"ownershipTimes": [
												{
													"start": "",
													"end": ""
												}
											],
											"badgeIds": [
												{
													"start": "",
													"end": ""
												}
											]
										}
									]
								}
							],
							"incrementedBalances": {
								"startBalances": [
									{
										"amount": "",
										"ownershipTimes": [
											{
												"start": "",
												"end": ""
											}
										],
										"badgeIds": [
											{
												"start": "",
												"end": ""
											}
										]
									}
								],
								"incrementBadgeIdsBy": "",
								"incrementOwnershipTimesBy": ""
							},
							"orderCalculationMethod": {
								"useOverallNumTransfers": false,
								"usePerToAddressNumTransfers": false,
								"usePerFromAddressNumTransfers": false,
								"usePerInitiatedByAddressNumTransfers": false,
								"useMerkleChallengeLeafIndex": false
							}
						},
						"approvalAmounts": {
							"overallApprovalAmount": "",
							"perToAddressApprovalAmount": "",
							"perFromAddressApprovalAmount": "",
							"perInitiatedByAddressApprovalAmount": ""
						},
						"maxNumTransfers": {
							"overallMaxNumTransfers": "",
							"perToAddressMaxNumTransfers": "",
							"perFromAddressMaxNumTransfers": "",
							"perInitiatedByAddressMaxNumTransfers": ""
						},
						"requireToEqualsInitiatedBy": false,
						"requireToDoesNotEqualInitiatedBy": false
					}
				}
			],
			"defaultIncomingApprovals": [
				{
					"fromMappingId": "",
					"initiatedByMappingId": "",
					"transferTimes": [
						{
							"start": "",
							"end": ""
						}
					],
					"badgeIds": [
						{
							"start": "",
							"end": ""
						}
					],
					"ownershipTimes": [
						{
							"start": "",
							"end": ""
						}
					],
					"amountTrackerId": "",
					"challengeTrackerId": "",
					"uri": "",
					"customData": "",
					"approvalId": "",
					"approvalCriteria": {
						"mustOwnBadges": [
							{
								"collectionId": "",
								"amountRange": {
									"start": "",
									"end": ""
								},
								"ownershipTimes": [
									{
										"start": "",
										"end": ""
									}
								],
								"badgeIds": [
									{
										"start": "",
										"end": ""
									}
								],
								"overrideWithCurrentTime": false,
								"mustOwnAll": false
							}
						],
						"merkleChallenge": {
							"root": "",
							"expectedProofLength": "",
							"useCreatorAddressAsLeaf": false,
							"maxUsesPerLeaf": "",
							"uri": "",
							"customData": ""
						},
						"predeterminedBalances": {
							"manualBalances": [
								{
									"balances": [
										{
											"amount": "",
											"ownershipTimes": [
												{
													"start": "",
													"end": ""
												}
											],
											"badgeIds": [
												{
													"start": "",
													"end": ""
												}
											]
										}
									]
								}
							],
							"incrementedBalances": {
								"startBalances": [
									{
										"amount": "",
										"ownershipTimes": [
											{
												"start": "",
												"end": ""
											}
										],
										"badgeIds": [
											{
												"start": "",
												"end": ""
											}
										]
									}
								],
								"incrementBadgeIdsBy": "",
								"incrementOwnershipTimesBy": ""
							},
							"orderCalculationMethod": {
								"useOverallNumTransfers": false,
								"usePerToAddressNumTransfers": false,
								"usePerFromAddressNumTransfers": false,
								"usePerInitiatedByAddressNumTransfers": false,
								"useMerkleChallengeLeafIndex": false
							}
						},
						"approvalAmounts": {
							"overallApprovalAmount": "",
							"perToAddressApprovalAmount": "",
							"perFromAddressApprovalAmount": "",
							"perInitiatedByAddressApprovalAmount": ""
						},
						"maxNumTransfers": {
							"overallMaxNumTransfers": "",
							"perToAddressMaxNumTransfers": "",
							"perFromAddressMaxNumTransfers": "",
							"perInitiatedByAddressMaxNumTransfers": ""
						},
						"requireFromEqualsInitiatedBy": false,
						"requireFromDoesNotEqualInitiatedBy": false
					}
				}
			],
			"defaultAutoApproveSelfInitiatedOutgoingTransfers": false,
			"defaultAutoApproveSelfInitiatedIncomingTransfers": false,
			"defaultUserPermissions": {
				"canUpdateOutgoingApprovals": [
					{
						"toMappingId": "",
						"initiatedByMappingId": "",
						"transferTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"badgeIds": [
							{
								"start": "",
								"end": ""
							}
						],
						"ownershipTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"amountTrackerId": "",
						"challengeTrackerId": "",
						"permittedTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"forbiddenTimes": [
							{
								"start": "",
								"end": ""
							}
						]
					}
				],
				"canUpdateIncomingApprovals": [
					{
						"fromMappingId": "",
						"initiatedByMappingId": "",
						"transferTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"badgeIds": [
							{
								"start": "",
								"end": ""
							}
						],
						"ownershipTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"amountTrackerId": "",
						"challengeTrackerId": "",
						"permittedTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"forbiddenTimes": [
							{
								"start": "",
								"end": ""
							}
						]
					}
				],
				"canUpdateAutoApproveSelfInitiatedOutgoingTransfers": [
					{
						"permittedTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"forbiddenTimes": [
							{
								"start": "",
								"end": ""
							}
						]
					}
				],
				"canUpdateAutoApproveSelfInitiatedIncomingTransfers": [
					{
						"permittedTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"forbiddenTimes": [
							{
								"start": "",
								"end": ""
							}
						]
					}
				]
			},
			"badgesToCreate": [
				{
					"amount": "",
					"ownershipTimes": [
						{
							"start": "",
							"end": ""
						}
					],
					"badgeIds": [
						{
							"start": "",
							"end": ""
						}
					]
				}
			],
			"updateCollectionPermissions": false,
			"collectionPermissions": {
				"canDeleteCollection": [
					{
						"permittedTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"forbiddenTimes": [
							{
								"start": "",
								"end": ""
							}
						]
					}
				],
				"canArchiveCollection": [
					{
						"timelineTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"permittedTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"forbiddenTimes": [
							{
								"start": "",
								"end": ""
							}
						]
					}
				],
				"canUpdateOffChainBalancesMetadata": [
					{
						"timelineTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"permittedTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"forbiddenTimes": [
							{
								"start": "",
								"end": ""
							}
						]
					}
				],
				"canUpdateStandards": [
					{
						"timelineTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"permittedTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"forbiddenTimes": [
							{
								"start": "",
								"end": ""
							}
						]
					}
				],
				"canUpdateCustomData": [
					{
						"timelineTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"permittedTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"forbiddenTimes": [
							{
								"start": "",
								"end": ""
							}
						]
					}
				],
				"canUpdateManager": [
					{
						"timelineTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"permittedTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"forbiddenTimes": [
							{
								"start": "",
								"end": ""
							}
						]
					}
				],
				"canUpdateCollectionMetadata": [
					{
						"timelineTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"permittedTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"forbiddenTimes": [
							{
								"start": "",
								"end": ""
							}
						]
					}
				],
				"canCreateMoreBadges": [
					{
						"badgeIds": [
							{
								"start": "",
								"end": ""
							}
						],
						"ownershipTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"permittedTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"forbiddenTimes": [
							{
								"start": "",
								"end": ""
							}
						]
					}
				],
				"canUpdateBadgeMetadata": [
					{
						"badgeIds": [
							{
								"start": "",
								"end": ""
							}
						],
						"timelineTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"permittedTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"forbiddenTimes": [
							{
								"start": "",
								"end": ""
							}
						]
					}
				],
				"canUpdateCollectionApprovals": [
					{
						"fromMappingId": "",
						"toMappingId": "",
						"initiatedByMappingId": "",
						"transferTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"badgeIds": [
							{
								"start": "",
								"end": ""
							}
						],
						"ownershipTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"amountTrackerId": "",
						"challengeTrackerId": "",
						"permittedTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"forbiddenTimes": [
							{
								"start": "",
								"end": ""
							}
						]
					}
				]
			},
			"updateManagerTimeline": false,
			"managerTimeline": [
				{
					"manager": "",
					"timelineTimes": [
						{
							"start": "",
							"end": ""
						}
					]
				}
			],
			"updateCollectionMetadataTimeline": false,
			"collectionMetadataTimeline": [
				{
					"collectionMetadata": {
						"uri": "",
						"customData": ""
					},
					"timelineTimes": [
						{
							"start": "",
							"end": ""
						}
					]
				}
			],
			"updateBadgeMetadataTimeline": false,
			"badgeMetadataTimeline": [
				{
					"badgeMetadata": [
						{
							"uri": "",
							"customData": "",
							"badgeIds": [
								{
									"start": "",
									"end": ""
								}
							]
						}
					],
					"timelineTimes": [
						{
							"start": "",
							"end": ""
						}
					]
				}
			],
			"updateOffChainBalancesMetadataTimeline": false,
			"offChainBalancesMetadataTimeline": [
				{
					"offChainBalancesMetadata": {
						"uri": "",
						"customData": ""
					},
					"timelineTimes": [
						{
							"start": "",
							"end": ""
						}
					]
				}
			],
			"updateCustomDataTimeline": false,
			"customDataTimeline": [
				{
					"customData": "",
					"timelineTimes": [
						{
							"start": "",
							"end": ""
						}
					]
				}
			],
			"updateCollectionApprovals": false,
			"collectionApprovals": [
				{
					"fromMappingId": "",
					"toMappingId": "",
					"initiatedByMappingId": "",
					"transferTimes": [
						{
							"start": "",
							"end": ""
						}
					],
					"badgeIds": [
						{
							"start": "",
							"end": ""
						}
					],
					"ownershipTimes": [
						{
							"start": "",
							"end": ""
						}
					],
					"amountTrackerId": "",
					"challengeTrackerId": "",
					"uri": "",
					"customData": "",
					"approvalId": "",
					"approvalCriteria": {
						"mustOwnBadges": [
							{
								"collectionId": "",
								"amountRange": {
									"start": "",
									"end": ""
								},
								"ownershipTimes": [
									{
										"start": "",
										"end": ""
									}
								],
								"badgeIds": [
									{
										"start": "",
										"end": ""
									}
								],
								"overrideWithCurrentTime": false,
								"mustOwnAll": false
							}
						],
						"merkleChallenge": {
							"root": "",
							"expectedProofLength": "",
							"useCreatorAddressAsLeaf": false,
							"maxUsesPerLeaf": "",
							"uri": "",
							"customData": ""
						},
						"predeterminedBalances": {
							"manualBalances": [
								{
									"balances": [
										{
											"amount": "",
											"ownershipTimes": [
												{
													"start": "",
													"end": ""
												}
											],
											"badgeIds": [
												{
													"start": "",
													"end": ""
												}
											]
										}
									]
								}
							],
							"incrementedBalances": {
								"startBalances": [
									{
										"amount": "",
										"ownershipTimes": [
											{
												"start": "",
												"end": ""
											}
										],
										"badgeIds": [
											{
												"start": "",
												"end": ""
											}
										]
									}
								],
								"incrementBadgeIdsBy": "",
								"incrementOwnershipTimesBy": ""
							},
							"orderCalculationMethod": {
								"useOverallNumTransfers": false,
								"usePerToAddressNumTransfers": false,
								"usePerFromAddressNumTransfers": false,
								"usePerInitiatedByAddressNumTransfers": false,
								"useMerkleChallengeLeafIndex": false
							}
						},
						"approvalAmounts": {
							"overallApprovalAmount": "",
							"perToAddressApprovalAmount": "",
							"perFromAddressApprovalAmount": "",
							"perInitiatedByAddressApprovalAmount": ""
						},
						"maxNumTransfers": {
							"overallMaxNumTransfers": "",
							"perToAddressMaxNumTransfers": "",
							"perFromAddressMaxNumTransfers": "",
							"perInitiatedByAddressMaxNumTransfers": ""
						},
						"requireToEqualsInitiatedBy": false,
						"requireFromEqualsInitiatedBy": false,
						"requireToDoesNotEqualInitiatedBy": false,
						"requireFromDoesNotEqualInitiatedBy": false,
						"overridesFromOutgoingApprovals": false,
						"overridesToIncomingApprovals": false
					}
				}
			],
			"updateStandardsTimeline": false,
			"standardsTimeline": [
				{
					"standards": [],
					"timelineTimes": [
						{
							"start": "",
							"end": ""
						}
					]
				}
			],
			"updateIsArchivedTimeline": false,
			"isArchivedTimeline": [
				{
					"isArchived": false,
					"timelineTimes": [
						{
							"start": "",
							"end": ""
						}
					]
				}
			]
		}
	}`)

	schemas = append(schemas, `{
		"type": "badges/UpdateUserApprovals",
		"value": {
			"creator": "",
			"collectionId": "",
			"updateOutgoingApprovals": false,
			"outgoingApprovals": [
				{
					"toMappingId": "",
					"initiatedByMappingId": "",
					"transferTimes": [
						{
							"start": "",
							"end": ""
						}
					],
					"badgeIds": [
						{
							"start": "",
							"end": ""
						}
					],
					"ownershipTimes": [
						{
							"start": "",
							"end": ""
						}
					],
					"amountTrackerId": "",
					"challengeTrackerId": "",
					"uri": "",
					"customData": "",
					"approvalId": "",
					"approvalCriteria": {
						"mustOwnBadges": [
							{
								"collectionId": "",
								"amountRange": {
									"start": "",
									"end": ""
								},
								"ownershipTimes": [
									{
										"start": "",
										"end": ""
									}
								],
								"badgeIds": [
									{
										"start": "",
										"end": ""
									}
								],
								"overrideWithCurrentTime": false,
								"mustOwnAll": false
							}
						],
						"merkleChallenge": {
							"root": "",
							"expectedProofLength": "",
							"useCreatorAddressAsLeaf": false,
							"maxUsesPerLeaf": "",
							"uri": "",
							"customData": ""
						},
						"predeterminedBalances": {
							"manualBalances": [
								{
									"balances": [
										{
											"amount": "",
											"ownershipTimes": [
												{
													"start": "",
													"end": ""
												}
											],
											"badgeIds": [
												{
													"start": "",
													"end": ""
												}
											]
										}
									]
								}
							],
							"incrementedBalances": {
								"startBalances": [
									{
										"amount": "",
										"ownershipTimes": [
											{
												"start": "",
												"end": ""
											}
										],
										"badgeIds": [
											{
												"start": "",
												"end": ""
											}
										]
									}
								],
								"incrementBadgeIdsBy": "",
								"incrementOwnershipTimesBy": ""
							},
							"orderCalculationMethod": {
								"useOverallNumTransfers": false,
								"usePerToAddressNumTransfers": false,
								"usePerFromAddressNumTransfers": false,
								"usePerInitiatedByAddressNumTransfers": false,
								"useMerkleChallengeLeafIndex": false
							}
						},
						"approvalAmounts": {
							"overallApprovalAmount": "",
							"perToAddressApprovalAmount": "",
							"perFromAddressApprovalAmount": "",
							"perInitiatedByAddressApprovalAmount": ""
						},
						"maxNumTransfers": {
							"overallMaxNumTransfers": "",
							"perToAddressMaxNumTransfers": "",
							"perFromAddressMaxNumTransfers": "",
							"perInitiatedByAddressMaxNumTransfers": ""
						},
						"requireToEqualsInitiatedBy": false,
						"requireToDoesNotEqualInitiatedBy": false
					}
				}
			],
			"updateIncomingApprovals": false,
			"incomingApprovals": [
				{
					"fromMappingId": "",
					"initiatedByMappingId": "",
					"transferTimes": [
						{
							"start": "",
							"end": ""
						}
					],
					"badgeIds": [
						{
							"start": "",
							"end": ""
						}
					],
					"ownershipTimes": [
						{
							"start": "",
							"end": ""
						}
					],
					"amountTrackerId": "",
					"challengeTrackerId": "",
					"uri": "",
					"customData": "",
					"approvalId": "",
					"approvalCriteria": {
						"mustOwnBadges": [
							{
								"collectionId": "",
								"amountRange": {
									"start": "",
									"end": ""
								},
								"ownershipTimes": [
									{
										"start": "",
										"end": ""
									}
								],
								"badgeIds": [
									{
										"start": "",
										"end": ""
									}
								],
								"overrideWithCurrentTime": false,
								"mustOwnAll": false
							}
						],
						"merkleChallenge": {
							"root": "",
							"expectedProofLength": "",
							"useCreatorAddressAsLeaf": false,
							"maxUsesPerLeaf": "",
							"uri": "",
							"customData": ""
						},
						"predeterminedBalances": {
							"manualBalances": [
								{
									"balances": [
										{
											"amount": "",
											"ownershipTimes": [
												{
													"start": "",
													"end": ""
												}
											],
											"badgeIds": [
												{
													"start": "",
													"end": ""
												}
											]
										}
									]
								}
							],
							"incrementedBalances": {
								"startBalances": [
									{
										"amount": "",
										"ownershipTimes": [
											{
												"start": "",
												"end": ""
											}
										],
										"badgeIds": [
											{
												"start": "",
												"end": ""
											}
										]
									}
								],
								"incrementBadgeIdsBy": "",
								"incrementOwnershipTimesBy": ""
							},
							"orderCalculationMethod": {
								"useOverallNumTransfers": false,
								"usePerToAddressNumTransfers": false,
								"usePerFromAddressNumTransfers": false,
								"usePerInitiatedByAddressNumTransfers": false,
								"useMerkleChallengeLeafIndex": false
							}
						},
						"approvalAmounts": {
							"overallApprovalAmount": "",
							"perToAddressApprovalAmount": "",
							"perFromAddressApprovalAmount": "",
							"perInitiatedByAddressApprovalAmount": ""
						},
						"maxNumTransfers": {
							"overallMaxNumTransfers": "",
							"perToAddressMaxNumTransfers": "",
							"perFromAddressMaxNumTransfers": "",
							"perInitiatedByAddressMaxNumTransfers": ""
						},
						"requireFromEqualsInitiatedBy": false,
						"requireFromDoesNotEqualInitiatedBy": false
					}
				}
			],
			"updateAutoApproveSelfInitiatedOutgoingTransfers": false,
			"autoApproveSelfInitiatedOutgoingTransfers": false,
			"updateAutoApproveSelfInitiatedIncomingTransfers": false,
			"autoApproveSelfInitiatedIncomingTransfers": false,
			"updateUserPermissions": false,
			"userPermissions": {
				"canUpdateOutgoingApprovals": [
					{
						"toMappingId": "",
						"initiatedByMappingId": "",
						"transferTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"badgeIds": [
							{
								"start": "",
								"end": ""
							}
						],
						"ownershipTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"amountTrackerId": "",
						"challengeTrackerId": "",
						"permittedTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"forbiddenTimes": [
							{
								"start": "",
								"end": ""
							}
						]
					}
				],
				"canUpdateIncomingApprovals": [
					{
						"fromMappingId": "",
						"initiatedByMappingId": "",
						"transferTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"badgeIds": [
							{
								"start": "",
								"end": ""
							}
						],
						"ownershipTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"amountTrackerId": "",
						"challengeTrackerId": "",
						"permittedTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"forbiddenTimes": [
							{
								"start": "",
								"end": ""
							}
						]
					}
				],
				"canUpdateAutoApproveSelfInitiatedOutgoingTransfers": [
					{
						"permittedTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"forbiddenTimes": [
							{
								"start": "",
								"end": ""
							}
						]
					}
				],
				"canUpdateAutoApproveSelfInitiatedIncomingTransfers": [
					{
						"permittedTimes": [
							{
								"start": "",
								"end": ""
							}
						],
						"forbiddenTimes": [
							{
								"start": "",
								"end": ""
							}
						]
					}
				]
			}
		}
	}`)

	return schemas
}
