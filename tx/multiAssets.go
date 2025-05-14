package tx

import (
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
)

type MultiAssets struct {
	V0 []V0MultiAssets `json:"V0,omitempty"`
	V1 []V1MultiAssets `json:"V1,omitempty"`
	V2 []V2MultiAssets `json:"V2,omitempty"`
	V3 []V3MultiAssets `json:"V3,omitempty"`
}

type V2MultiAssets V1MultiAssets
type V3MultiAssets V3MultiAsset

type V0MultiAssets struct {
	None                Enum `json:"None,omitempty"`
	All                 Enum `json:"All,omitempty"`
	AllFungible         Enum `json:"AllFungible,omitempty"`
	AllNonFungible      Enum `json:"AllNonFungible,omitempty"`
	AllAbstractFungible *struct {
		Id string `json:"id"`
	} `json:"AllAbstractFungible,omitempty"`
	AllAbstractNonFungible *struct {
		Class string `json:"class"`
	} `json:"AllAbstractNonFungible,omitempty"`
	AllConcreteFungible    *V0MultiLocation `json:"AllConcreteFungible,omitempty"`
	AllConcreteNonFungible *V0MultiLocation `json:"AllConcreteNonFungible,omitempty"`
	AbstractFungible       *struct {
		Id     string          `json:"id"`
		Amount decimal.Decimal `json:"amount"`
	} `json:"AbstractFungible,omitempty"`
	AbstractNonFungible *struct {
		Class    string        `json:"id"`
		Instance AssetInstance `json:"instance"`
	} `json:"AbstractNonFungible,omitempty"`
	ConcreteFungible *struct {
		Id     *V0MultiLocation `json:"id"`
		Amount decimal.Decimal  `json:"amount"`
	} `json:"ConcreteFungible,omitempty"`
	ConcreteNonFungible *struct {
		Id       *V0MultiLocation `json:"id"`
		Instance AssetInstance    `json:"instance"`
	} `json:"ConcreteNonFungible,omitempty"`
}

type V1MultiAssets struct {
	Id  AssetsId  `json:"id"`  // AssetId
	Fun AssetsFun `json:"fun"` // Fungibility
}

type V3MultiAsset struct {
	Id  AssetsIdV3 `json:"id"`  // AssetId
	Fun AssetsFun  `json:"fun"` // Fungibility
}

type AssetsId struct {
	Concrete *V1MultiLocation `json:"Concrete,omitempty"`
	Abstract *string          `json:"Abstract,omitempty"`
}

type AssetsIdV3 struct {
	Concrete *V3MultiLocation `json:"Concrete,omitempty"`
	Abstract *string          `json:"Abstract,omitempty"`
}
type AssetsFun struct {
	Fungible    *decimal.Decimal `json:"Fungible,omitempty"`
	NonFungible *AssetInstance   `json:"NonFungible,omitempty"`
}

type AssetInstance struct {
	Undefined *Enum            `json:"Undefined,omitempty"`
	Index     *decimal.Decimal `json:"Index,omitempty"`
	Array4    *string          `json:"Array4,omitempty"`
	Array8    *string          `json:"Array8,omitempty"`
	Array16   *string          `json:"Array16,omitempty"`
	Array32   *string          `json:"Array32,omitempty"`
	Blob      *string          `json:"Blob,omitempty"`
}

func (m *MultiAssets) ToScale() interface{} {
	r := map[string]interface{}{}
	b, _ := json.Marshal(m)
	_ = json.Unmarshal(b, &r)
	return r
}
