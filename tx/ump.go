package tx

type Ump struct {
	ModuleName string
}

func (u Ump) GetModuleName() string {
	return u.ModuleName
}

func NewUmp() *Ump {
	return &Ump{ModuleName: "PolkadotXcm"}
}

func (u Ump) DefaultUmpDest() *VersionedMultiLocation {
	return nil
}

func (u Ump) Send(location *VersionedMultiLocation, i *VersionedXcm) (string, []interface{}) {
	return "send", []interface{}{location.ToScale(), i.ToScale()}
}

// LimitedReserveTransferAssets
// (dest VersionedMultiLocation, beneficiary VersionedMultiLocation, assets VersionedMultiAssets, fee_asset_item u32, weight_limit WeightLimit)
func (u Ump) LimitedReserveTransferAssets(dest *VersionedMultiLocation, beneficiary *VersionedMultiLocation, assets *MultiAssets, feeAssetItem uint, weightLimit *Weight) (string, []interface{}) {
	return "limited_reserve_transfer_assets", []interface{}{dest.ToScale(), beneficiary.ToScale(), assets.ToScale(), feeAssetItem, weightLimit.ToScale()}
}

// ReserveTransferAssets
// (dest VersionedMultiLocation, beneficiary VersionedMultiLocation, assets VersionedMultiAssets, fee_asset_item u32)
func (u Ump) ReserveTransferAssets(dest *VersionedMultiLocation, beneficiary *VersionedMultiLocation, assets *MultiAssets, feeAssetItem uint) (string, []interface{}) {
	return "reserve_transfer_assets", []interface{}{dest.ToScale(), beneficiary.ToScale(), assets.ToScale(), feeAssetItem}
}

// LimitedTeleportAssets
// (dest VersionedMultiLocation, beneficiary VersionedMultiLocation, assets VersionedMultiAssets, fee_asset_item u32, weight_limit WeightLimit)
func (u Ump) LimitedTeleportAssets(dest *VersionedMultiLocation, beneficiary *VersionedMultiLocation, assets *MultiAssets, feeAssetItem uint, weightLimit *Weight) (string, []interface{}) {
	return "limited_teleport_assets", []interface{}{dest.ToScale(), beneficiary.ToScale(), assets.ToScale(), feeAssetItem, weightLimit.ToScale()}
}

// TeleportAssets
// (dest VersionedMultiLocation, beneficiary VersionedMultiLocation, assets VersionedMultiAssets, fee_asset_item u32)
func (u Ump) TeleportAssets(dest *VersionedMultiLocation, beneficiary *VersionedMultiLocation, assets *MultiAssets, feeAssetItem uint) (string, []interface{}) {
	return "teleport_assets", []interface{}{dest.ToScale(), beneficiary.ToScale(), assets.ToScale(), feeAssetItem}
}
