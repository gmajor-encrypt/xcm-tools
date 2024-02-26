package tx

// IXCMP XCMP interface

// https://substrate.stackexchange.com/questions/37/how-can-i-transfer-assets-using-xcm
// https://wiki.polkadot.network/docs/learn-xcm#reserve-asset-transfer

type IXCMP interface {
	// LimitedReserveTransferAssets
	// Transfer some assets from the local chain to the sovereign account of a destination chain and forward a notification XCM.
	// [dest, beneficiary, assets, fee_asset_item, weight_limit]
	LimitedReserveTransferAssets(*VersionedMultiLocation, *VersionedMultiLocation, *MultiAssets, uint, *Weight) (string, []interface{})

	// ReserveTransferAssets
	// Transfer some assets from the local chain to the sovereign account of a destination chain and forward a notification XCM.
	// [dest, beneficiary, assets, fee_asset_item]
	ReserveTransferAssets(*VersionedMultiLocation, *VersionedMultiLocation, *MultiAssets, uint) (string, []interface{})

	// LimitedTeleportAssets Teleport some assets from the local chain to some destination chain.
	// [dest, beneficiary, assets, fee_asset_item, weight_limit]
	LimitedTeleportAssets(*VersionedMultiLocation, *VersionedMultiLocation, *MultiAssets, uint, *Weight) (string, []interface{})
	// TeleportAssets
	// Teleport some assets from the local chain to some destination chain.
	// [dest, beneficiary, assets, fee_asset_item]
	TeleportAssets(*VersionedMultiLocation, *VersionedMultiLocation, *MultiAssets, uint) (string, []interface{})

	// Send [dest,message]
	Send(*VersionedMultiLocation, *VersionedXcm) (string, []interface{})

	// TransferAssets
	// Transfer some assets from the local chain to the destination chain through their local,
	// destination or remote reserve, or through teleports.
	// [dest, beneficiary, assets, feeAssetItem, weightLimit]
	TransferAssets(*VersionedMultiLocation, *VersionedMultiLocation, *MultiAssets, uint, *Weight) (string, []interface{})

	// GetModuleName
	// Extrinsics module name
	GetModuleName() string
}
