package tx

import "encoding/json"

type VersionedXcm struct {
	V0 *V0 `json:"V0,omitempty"`
	V1 *V1 `json:"V1,omitempty"`
	V2 *V2 `json:"V2,omitempty"`
	V3 *V3 `json:"V3,omitempty"`
	V4 *V4 `json:"V4,omitempty"`
}

type V0 struct {
	WithdrawAsset             *interface{} `json:"WithdrawAsset,omitempty"`
	ReserveAssetDeposit       *interface{} `json:"ReserveAssetDeposit,omitempty"`
	TeleportAsset             *interface{} `json:"TeleportAsset,omitempty"`
	QueryResponse             *interface{} `json:"QueryResponse,omitempty"`
	TransferAsset             *interface{} `json:"TransferAsset,omitempty"`
	TransferReserveAsset      *interface{} `json:"TransferReserveAsset,omitempty"`
	Transact                  *interface{} `json:"Transact,omitempty"`
	HrmpNewChannelOpenRequest *interface{} `json:"HrmpNewChannelOpenRequest,omitempty"`
	HrmpChannelAccepted       *interface{} `json:"HrmpChannelAccepted,omitempty"`
	HrmpChannelClosing        *interface{} `json:"HrmpChannelClosing,omitempty"`
	RelayedFrom               *interface{} `json:"RelayedFrom,omitempty"`
}

type V1 struct {
	WithdrawAsset             *interface{} `json:"WithdrawAsset,omitempty"`
	ReserveAssetDeposit       *interface{} `json:"ReserveAssetDeposit,omitempty"`
	ReceiveTeleportedAsset    *interface{} `json:"ReceiveTeleportedAsset,omitempty"`
	QueryResponse             *interface{} `json:"QueryResponse,omitempty"`
	TransferAsset             *interface{} `json:"TransferAsset,omitempty"`
	TransferReserveAsset      *interface{} `json:"TransferReserveAsset,omitempty"`
	Transact                  *interface{} `json:"Transact,omitempty"`
	HrmpNewChannelOpenRequest *interface{} `json:"HrmpNewChannelOpenRequest,omitempty"`
	HrmpChannelAccepted       *interface{} `json:"HrmpChannelAccepted,omitempty"`
	HrmpChannelClosing        *interface{} `json:"HrmpChannelClosing,omitempty"`
	RelayedFrom               *interface{} `json:"RelayedFrom,omitempty"`
	SubscribeVersion          *interface{} `json:"SubscribeVersion,omitempty"`
	UnsubscribeVersion        *interface{} `json:"UnsubscribeVersion,omitempty"`
}

type V2 []V2XcmInstruction

type V2XcmInstruction struct {
	WithdrawAsset             *interface{}          `json:"WithdrawAsset,omitempty"`
	ReserveAssetDeposited     *interface{}          `json:"ReserveAssetDeposited,omitempty"`
	ReceiveTeleportedAsset    *interface{}          `json:"ReceiveTeleportedAsset,omitempty"`
	QueryResponse             *interface{}          `json:"QueryResponse,omitempty"`
	TransferAsset             *TransferAsset        `json:"TransferAsset,omitempty"`
	TransferReserveAsset      *TransferReserveAsset `json:"TransferReserveAsset,omitempty"`
	Transact                  *Transact             `json:"Transact,omitempty"`
	HrmpNewChannelOpenRequest *interface{}          `json:"HrmpNewChannelOpenRequest,omitempty"`
	HrmpChannelAccepted       *interface{}          `json:"HrmpChannelAccepted,omitempty"`
	HrmpChannelClosing        *interface{}          `json:"HrmpChannelClosing,omitempty"`
	ClearOrigin               *interface{}          `json:"ClearOrigin,omitempty"`
	DescendOrigin             *DescendOrigin        `json:"DescendOrigin,omitempty"`
	ReportError               *interface{}          `json:"ReportError,omitempty"`
	DepositAsset              *DepositAsset         `json:"DepositAsset,omitempty"`
	DepositReserveAsset       *DepositReserveAsset  `json:"DepositReserveAsset,omitempty"`
	ExchangeAsset             *interface{}          `json:"ExchangeAsset,omitempty"`
	InitiateReserveWithdraw   *interface{}          `json:"InitiateReserveWithdraw,omitempty"`
	InitiateTeleport          *InitiateTeleport     `json:"InitiateTeleport,omitempty"`
	QueryHolding              *interface{}          `json:"QueryHolding,omitempty"`
	BuyExecution              *BuyExecution         `json:"BuyExecution,omitempty"`
	RefundSurplus             *interface{}          `json:"RefundSurplus,omitempty"`
	SetErrorHandler           *interface{}          `json:"SetErrorHandler,omitempty"`
	SetAppendix               *interface{}          `json:"SetAppendix,omitempty"`
	ClearError                *interface{}          `json:"ClearError,omitempty"`
	ClaimAsset                *interface{}          `json:"ClaimAsset,omitempty"`
	Trap                      *interface{}          `json:"Trap,omitempty"`
	SubscribeVersion          *interface{}          `json:"SubscribeVersion,omitempty"`
	UnsubscribeVersion        *interface{}          `json:"UnsubscribeVersion,omitempty"`
}

type V3 []XcmInstructionV3

type XcmInstructionV3 struct {
	WithdrawAsset             *interface{}            `json:"WithdrawAsset,omitempty"`
	ReserveAssetDeposited     *interface{}            `json:"ReserveAssetDeposited,omitempty"`
	ReceiveTeleportedAsset    *interface{}            `json:"ReceiveTeleportedAsset,omitempty"`
	QueryResponse             *interface{}            `json:"QueryResponse,omitempty"`
	TransferAsset             *TransferAsset          `json:"TransferAsset,omitempty"`
	TransferReserveAsset      *TransferReserveAssetV3 `json:"TransferReserveAsset,omitempty"`
	Transact                  *TransactV3             `json:"Transact,omitempty"`
	HrmpNewChannelOpenRequest *interface{}            `json:"HrmpNewChannelOpenRequest,omitempty"`
	HrmpChannelAccepted       *interface{}            `json:"HrmpChannelAccepted,omitempty"`
	HrmpChannelClosing        *interface{}            `json:"HrmpChannelClosing,omitempty"`
	ClearOrigin               *interface{}            `json:"ClearOrigin,omitempty"`
	DescendOrigin             *DescendOrigin          `json:"DescendOrigin,omitempty"`
	ReportError               *interface{}            `json:"ReportError,omitempty"`
	DepositAsset              *DepositAsset           `json:"DepositAsset,omitempty"`
	DepositReserveAsset       *DepositReserveAssetV3  `json:"DepositReserveAsset,omitempty"`
	ExchangeAsset             *interface{}            `json:"ExchangeAsset,omitempty"`
	InitiateReserveWithdraw   *interface{}            `json:"InitiateReserveWithdraw,omitempty"`
	InitiateTeleport          *InitiateTeleportV3     `json:"InitiateTeleport,omitempty"`
	QueryHolding              *interface{}            `json:"QueryHolding,omitempty"`
	BuyExecution              *BuyExecution           `json:"BuyExecution,omitempty"`
	RefundSurplus             *interface{}            `json:"RefundSurplus,omitempty"`
	SetErrorHandler           *interface{}            `json:"SetErrorHandler,omitempty"`
	SetAppendix               *interface{}            `json:"SetAppendix,omitempty"`
	ClearError                *interface{}            `json:"ClearError,omitempty"`
	ClaimAsset                *interface{}            `json:"ClaimAsset,omitempty"`
	Trap                      *interface{}            `json:"Trap,omitempty"`
	SubscribeVersion          *interface{}            `json:"SubscribeVersion,omitempty"`
	UnsubscribeVersion        *interface{}            `json:"UnsubscribeVersion,omitempty"`
	BurnAsset                 *interface{}            `json:"BurnAsset,omitempty"`
	ExpectAsset               *interface{}            `json:"ExpectAsset,omitempty"`
	ExpectOrigin              *interface{}            `json:"ExpectOrigin,omitempty"`
	ExpectError               *interface{}            `json:"ExpectError,omitempty"`
	ExpectTransactStatus      *interface{}            `json:"ExpectTransactStatus,omitempty"`
	QueryPallet               *interface{}            `json:"QueryPallet,omitempty"`
	ExpectPallet              *interface{}            `json:"ExpectPallet,omitempty"`
	ReportTransactStatus      *interface{}            `json:"ReportTransactStatus,omitempty"`
	ClearTransactStatus       *interface{}            `json:"ClearTransactStatus,omitempty"`
	UniversalOrigin           *interface{}            `json:"UniversalOrigin,omitempty"`
	ExportMessage             *interface{}            `json:"ExportMessage,omitempty"`
	LockAsset                 *interface{}            `json:"LockAsset,omitempty"`
	UnlockAsset               *interface{}            `json:"UnlockAsset,omitempty"`
	NoteUnlockable            *interface{}            `json:"NoteUnlockable,omitempty"`
	RequestUnlock             *interface{}            `json:"RequestUnlock,omitempty"`
	SetFeesMode               *interface{}            `json:"SetFeesMode,omitempty"`
	SetTopic                  *string                 `json:"SetTopic,omitempty"`
	ClearTopic                *interface{}            `json:"ClearTopic,omitempty"`
	AliasOrigin               *interface{}            `json:"AliasOrigin,omitempty"`
	UnpaidExecution           *interface{}            `json:"UnpaidExecution,omitempty"`
}

type V4 []XcmInstructionV4

type XcmInstructionV4 struct {
	WithdrawAsset             *interface{} `json:"WithdrawAsset,omitempty"`
	ReserveAssetDeposited     *interface{} `json:"ReserveAssetDeposited,omitempty"`
	ReceiveTeleportedAsset    *interface{} `json:"ReceiveTeleportedAsset,omitempty"`
	QueryResponse             *interface{} `json:"QueryResponse,omitempty"`
	TransferAsset             *interface{} `json:"TransferAsset,omitempty"`
	TransferReserveAsset      *interface{} `json:"TransferReserveAsset,omitempty"`
	Transact                  *interface{} `json:"Transact,omitempty"`
	HrmpNewChannelOpenRequest *interface{} `json:"HrmpNewChannelOpenRequest,omitempty"`
	HrmpChannelAccepted       *interface{} `json:"HrmpChannelAccepted,omitempty"`
	HrmpChannelClosing        *interface{} `json:"HrmpChannelClosing,omitempty"`
	ClearOrigin               *interface{} `json:"ClearOrigin,omitempty"`
	DescendOrigin             *interface{} `json:"DescendOrigin,omitempty"`
	ReportError               *interface{} `json:"ReportError,omitempty"`
	DepositAsset              *interface{} `json:"DepositAsset,omitempty"`
	DepositReserveAsset       *interface{} `json:"DepositReserveAsset,omitempty"`
	ExchangeAsset             *interface{} `json:"ExchangeAsset,omitempty"`
	InitiateReserveWithdraw   *interface{} `json:"InitiateReserveWithdraw,omitempty"`
	InitiateTeleport          *interface{} `json:"InitiateTeleport,omitempty"`
	ReportHolding             *interface{} `json:"ReportHolding,omitempty"`
	BuyExecution              *interface{} `json:"BuyExecution,omitempty"`
	RefundSurplus             *interface{} `json:"RefundSurplus,omitempty"`
	SetErrorHandler           *interface{} `json:"SetErrorHandler,omitempty"`
	SetAppendix               *interface{} `json:"SetAppendix,omitempty"`
	ClearError                *interface{} `json:"ClearError,omitempty"`
	ClaimAsset                *interface{} `json:"ClaimAsset,omitempty"`
	Trap                      *interface{} `json:"Trap,omitempty"`
	SubscribeVersion          *interface{} `json:"SubscribeVersion,omitempty"`
	UnsubscribeVersion        *interface{} `json:"UnsubscribeVersion,omitempty"`
	BurnAsset                 *interface{} `json:"BurnAsset,omitempty"`
	ExpectAsset               *interface{} `json:"ExpectAsset,omitempty"`
	ExpectOrigin              *interface{} `json:"ExpectOrigin,omitempty"`
	ExpectError               *interface{} `json:"ExpectError,omitempty"`
	ExpectTransactStatus      *interface{} `json:"ExpectTransactStatus,omitempty"`
	QueryPallet               *interface{} `json:"QueryPallet,omitempty"`
	ExpectPallet              *interface{} `json:"ExpectPallet,omitempty"`
	ReportTransactStatus      *interface{} `json:"ReportTransactStatus,omitempty"`
	ClearTransactStatus       *interface{} `json:"ClearTransactStatus,omitempty"`
	UniversalOrigin           *interface{} `json:"UniversalOrigin,omitempty"`
	ExportMessage             *interface{} `json:"ExportMessage,omitempty"`
	LockAsset                 *interface{} `json:"LockAsset,omitempty"`
	UnlockAsset               *interface{} `json:"UnlockAsset,omitempty"`
	NoteUnlockable            *interface{} `json:"NoteUnlockable,omitempty"`
	RequestUnlock             *interface{} `json:"RequestUnlock,omitempty"`
	SetFeesMode               *interface{} `json:"SetFeesMode,omitempty"`
	SetTopic                  *string      `json:"SetTopic,omitempty"`
	ClearTopic                *interface{} `json:"ClearTopic,omitempty"`
	AliasOrigin               *interface{} `json:"AliasOrigin,omitempty"`
	UnpaidExecution           *interface{} `json:"UnpaidExecution,omitempty"`
}

type TransferAsset struct {
	Assets      interface{}      `json:"assets"`
	Beneficiary *V1MultiLocation `json:"beneficiary"`
}

type TransferReserveAssetV3 struct {
	Assets interface{}        `json:"assets"`
	Dest   *V1MultiLocation   `json:"dest"`
	Xcm    []XcmInstructionV3 `json:"xcm"`
}

type V1XCMRuntimeCall struct {
	Noop                    *interface{}  `json:"Noop,omitempty"`
	DepositAsset            *DepositAsset `json:"DepositAsset,omitempty"`
	DepositReserveAsset     *interface{}  `json:"DepositReserveAsset,omitempty"`
	ExchangeAsset           *interface{}  `json:"ExchangeAsset,omitempty"`
	InitiateReserveWithdraw *interface{}  `json:"InitiateReserveWithdraw,omitempty"`
	InitiateTeleport        *interface{}  `json:"InitiateTeleport,omitempty"`
	QueryHolding            *interface{}  `json:"QueryHolding,omitempty"`
	BuyExecution            *interface{}  `json:"BuyExecution,omitempty"`
}

type DepositAsset struct {
	Assets      interface{}      `json:"assets"`
	MaxAssets   int              `json:"max_assets"`
	Beneficiary *V1MultiLocation `json:"beneficiary"`
}

type ReserveAssetDeposited struct {
	Assets  *interface{}       `json:"assets"`
	Effects []V1XCMRuntimeCall `json:"effects"`
}

type TransferReserveAsset struct {
	Assets interface{}        `json:"assets"`
	Dest   *V1MultiLocation   `json:"dest"`
	Xcm    []V2XcmInstruction `json:"xcm"`
}

type DescendOrigin *V0MultiLocation

type Transact struct {
	Call                string      `json:"call"`
	OriginType          string      `json:"origin_type"`
	RequireWeightAtMost interface{} `json:"require_weight_at_most"`
}

type TransactV3 struct {
	Call                string      `json:"call"`
	OriginKind          string      `json:"origin_kind"`
	RequireWeightAtMost interface{} `json:"require_weight_at_most"`
}

type DepositReserveAsset struct {
	Assets    interface{} `json:"assets"`
	Dest      interface{} `json:"dest"`
	MaxAssets int         `json:"max_assets"`
	Xcm       V2          `json:"xcm"`
}

type DepositReserveAssetV3 struct {
	Assets interface{} `json:"assets"`
	Dest   interface{} `json:"dest"`
	Xcm    V3          `json:"xcm"`
}

type InitiateTeleport struct {
	Assets interface{} `json:"assets"`
	Dest   interface{} `json:"dest"`
	Xcm    V2          `json:"xcm"`
}

type BuyExecution struct {
	Fees          *V1MultiAssets `json:"fees"`
	WeightLimitV2 *WeightLimitV2 `json:"weight_limit"`
}

type WeightLimitV2 struct {
	Unlimited *string      `json:"Unlimited,omitempty"`
	Limited   *interface{} `json:"Limited,omitempty"`
}

type InitiateTeleportV3 struct {
	Assets interface{} `json:"assets"`
	Dest   interface{} `json:"dest"`
	Xcm    V3          `json:"xcm"`
}

func (v *VersionedXcm) ToScale() interface{} {
	r := map[string]interface{}{}
	b, _ := json.Marshal(v)
	_ = json.Unmarshal(b, &r)
	return r
}
