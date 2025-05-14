package tx

import "github.com/shopspring/decimal"

// SimplifyMultiLocationParaId Simplify paraId to VersionedMultiLocation
func SimplifyMultiLocationParaId(paraId uint32) *VersionedMultiLocation {
	return &VersionedMultiLocation{
		V3: &V3MultiLocation{
			Interior: V3MultiLocationJunctions{
				X1: &XCMJunctionV3{
					Parachain: &paraId,
				},
			},
			Parents: 0,
		},
	}
}

// SimplifyMultiLocationRelayChain Simplify parent relay chain to VersionedMultiLocation
func SimplifyMultiLocationRelayChain() *VersionedMultiLocation {
	return &VersionedMultiLocation{
		V3: &V3MultiLocation{
			Parents:  1,
			Interior: V3MultiLocationJunctions{Here: "NULL"},
		},
	}
}

// SimplifyMultiLocationAccountId32 Simplify accountId to VersionedMultiLocation
func SimplifyMultiLocationAccountId32(accountId string) *VersionedMultiLocation {
	return &VersionedMultiLocation{
		V4: &V3MultiLocation{
			Interior: V3MultiLocationJunctions{
				X1: &XCMJunctionV3{
					AccountId32: &XCMJunctionV3AccountId32{
						Id:      accountId,
						Network: nil,
					},
				},
			},
			Parents: 0,
		},
	}
}

// SimplifyMultiAssets Simplify sovereignty token to MultiAssets
func SimplifyMultiAssets(amount decimal.Decimal) *MultiAssets {
	return &MultiAssets{
		V3: []V3MultiAssets{
			{
				Id: AssetsIdV3{
					Concrete: &V3MultiLocation{
						Interior: V3MultiLocationJunctions{Here: "NULL"},
						Parents:  1,
					},
				},
				Fun: AssetsFun{
					Fungible: &amount,
				},
			},
		},
	}
}

// SimplifyV3MultiAssets Simplify sovereignty token to V3MultiAssets
func SimplifyV3MultiAssets(amount decimal.Decimal) *MultiAssets {
	return &MultiAssets{
		V3: []V3MultiAssets{
			{
				Id: AssetsIdV3{
					Concrete: &V3MultiLocation{
						Interior: V3MultiLocationJunctions{Here: "NULL"},
						Parents:  1,
					},
				},
				Fun: AssetsFun{
					Fungible: &amount,
				},
			},
		},
	}
}

// SimplifyEthereumAssets Simplify Ethereum assets to V3MultiAssets
// chainId: Ethereum chain id
// tokenContract: Ethereum token contract address
// amount: Ethereum token amount
func SimplifyEthereumAssets(chainId uint64, tokenContract string, amount decimal.Decimal) V3MultiAssets {
	return V3MultiAssets{
		Id: AssetsIdV3{
			Concrete: &V3MultiLocation{
				Interior: V3MultiLocationJunctions{
					X2: map[string]XCMJunctionV3{
						"col0": {GlobalConsensus: &GlobalConsensusNetworkId{Ethereum: &chainId}},
						"col1": {AccountKey20: &XCMJunctionV3AccountKey20{Key: tokenContract}},
					},
				},
				Parents: 2,
			},
		},
		Fun: AssetsFun{
			Fungible: &amount,
		},
	}
}

// SimplifyUnlimitedWeight Simplify unlimited weight
func SimplifyUnlimitedWeight() *Weight {
	var unlimited string
	return &Weight{Unlimited: &unlimited}
}
