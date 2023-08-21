package tx

type Hrmp struct {
	ModuleName string
}

func (h Hrmp) GetModuleName() string {
	return h.ModuleName
}

func (h Hrmp) Send(location *VersionedMultiLocation, i *VersionedXcm) (string, []interface{}) {
	return "send", []interface{}{location.ToScale(), i.ToScale()}
}

func (h Hrmp) LimitedReserveTransferAssets(dest *VersionedMultiLocation, beneficiary *VersionedMultiLocation, assets *MultiAssets, feeAssetItem uint, weightLimit *Weight) (string, []interface{}) {
	return "limited_reserve_transfer_assets", []interface{}{dest.ToScale(), beneficiary.ToScale(), assets.ToScale(), feeAssetItem, weightLimit.ToScale()}
}

func (h Hrmp) ReserveTransferAssets(dest *VersionedMultiLocation, beneficiary *VersionedMultiLocation, assets *MultiAssets, feeAssetItem uint) (string, []interface{}) {
	return "reserve_transfer_assets", []interface{}{dest.ToScale(), beneficiary.ToScale(), assets.ToScale(), feeAssetItem}
}

func (h Hrmp) LimitedTeleportAssets(dest *VersionedMultiLocation, beneficiary *VersionedMultiLocation, assets *MultiAssets, feeAssetItem uint, weightLimit *Weight) (string, []interface{}) {
	return "limited_teleport_assets", []interface{}{dest.ToScale(), beneficiary.ToScale(), assets.ToScale(), feeAssetItem, weightLimit.ToScale()}
}

func (h Hrmp) TeleportAssets(dest *VersionedMultiLocation, beneficiary *VersionedMultiLocation, assets *MultiAssets, feeAssetItem uint) (string, []interface{}) {
	return "teleport_assets", []interface{}{dest.ToScale(), beneficiary.ToScale(), assets.ToScale(), feeAssetItem}
}

func NewHrmp() *Hrmp {
	return &Hrmp{ModuleName: "PolkadotXcm"}
}
