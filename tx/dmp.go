package tx

type Dmp struct {
	ModuleName string
}

func (d Dmp) GetModuleName() string {
	return d.ModuleName
}

func (d Dmp) LimitedReserveTransferAssets(dest *VersionedMultiLocation, beneficiary *VersionedMultiLocation, assets *MultiAssets, feeAssetItem uint, weightLimit *Weight) (string, []interface{}) {
	return "limited_reserve_transfer_assets", []interface{}{dest.ToScale(), beneficiary.ToScale(), assets.ToScale(), feeAssetItem, weightLimit.ToScale()}
}

func (d Dmp) ReserveTransferAssets(dest *VersionedMultiLocation, beneficiary *VersionedMultiLocation, assets *MultiAssets, feeAssetItem uint) (string, []interface{}) {
	return "reserve_transfer_assets", []interface{}{dest.ToScale(), beneficiary.ToScale(), assets.ToScale(), feeAssetItem}
}

func (d Dmp) LimitedTeleportAssets(dest *VersionedMultiLocation, beneficiary *VersionedMultiLocation, assets *MultiAssets, feeAssetItem uint, weightLimit *Weight) (string, []interface{}) {
	return "limited_teleport_assets", []interface{}{dest.ToScale(), beneficiary.ToScale(), assets.ToScale(), feeAssetItem, weightLimit.ToScale()}
}

func (d Dmp) TeleportAssets(dest *VersionedMultiLocation, beneficiary *VersionedMultiLocation, assets *MultiAssets, feeAssetItem uint) (string, []interface{}) {
	return "teleport_assets", []interface{}{dest.ToScale(), beneficiary.ToScale(), assets.ToScale(), feeAssetItem}
}

func (d Dmp) Send(location *VersionedMultiLocation, i *VersionedXcm) (string, []interface{}) {
	return "send", []interface{}{location.ToScale(), i.ToScale()}
}

// TransferAssets
// (dest VersionedMultiLocation, beneficiary VersionedMultiLocation, assets VersionedMultiAssets, feeAssetItem u32, weightLimit Weight)
func (d Dmp) TransferAssets(*VersionedMultiLocation, *VersionedMultiLocation, *MultiAssets, uint, *Weight) (string, []interface{}) {
	panic("dmp does not support TransferAssets")
}

func NewDmp() *Dmp {
	return &Dmp{ModuleName: "XcmPallet"}
}
