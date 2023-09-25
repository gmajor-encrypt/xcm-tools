package tx

import (
	"encoding/json"
	"fmt"
)

type VersionedMultiLocation struct {
	V1 *V1MultiLocation `json:"V1,omitempty"`
	V2 *V1MultiLocation `json:"V2,omitempty"`
	V3 *V1MultiLocation `json:"V3,omitempty"`
}

func (v *VersionedMultiLocation) GetParaId() uint {
	if v.V1 != nil {
		return v.V1.GetParaId()
	}
	if v.V2 != nil {
		return v.V2.GetParaId()
	}
	if v.V3 != nil {
		return v.V3.GetParaId()
	}
	return 0
}

func (v *V1MultiLocation) GetParaId() uint {
	m := v.Interior
	var junctions map[string]XCMJunction
	switch {
	case m.X1 != nil:
		junctions = map[string]XCMJunction{"col0": *m.X1}
	case m.X2 != nil:
		junctions = m.X2
	case m.X3 != nil:
		junctions = m.X3
	case m.X4 != nil:
		junctions = m.X4
	case m.X5 != nil:
		junctions = m.X5
	case m.X6 != nil:
		junctions = m.X6
	case m.X7 != nil:
		junctions = m.X7
	case m.X8 != nil:
		junctions = m.X8
	}
	for i := 0; i <= 7; i++ {
		if junction, ok := junctions[fmt.Sprintf("col%d", i)]; ok && junction.Parachain != nil {
			return uint(*junction.Parachain)
		}
	}
	return 0
}

func (v *VersionedMultiLocation) ToScale() interface{} {
	r := map[string]interface{}{}
	b, _ := json.Marshal(v)
	_ = json.Unmarshal(b, &r)
	return r
}

type V1MultiLocation struct {
	Interior V0MultiLocation `json:"interior"`
	Parents  uint            `json:"parents"`
}

type Enum map[string]string

type XCMJunction struct {
	Parent         *string                    `json:"Parent,omitempty"`
	Parachain      *uint32                    `json:"Parachain,omitempty"`
	AccountId32    *XCMJunctionAccountId32    `json:"AccountId32,omitempty"`
	AccountIndex64 *XCMJunctionAccountIndex64 `json:"AccountIndex64,omitempty"`
	AccountKey20   *XCMJunctionAccountKey20   `json:"AccountKey20,omitempty"`
	PalletInstance *uint32                    `json:"PalletInstance,omitempty"`
	GeneralIndex   *string                    `json:"GeneralIndex,omitempty"`
	GeneralKey     *interface{}               `json:"GeneralKey,omitempty"`
	OnlyChild      *interface{}               `json:"OnlyChild,omitempty"`
	Plurality      *map[string]interface{}    `json:"Plurality,omitempty"`
}

type XCMJunctionAccountId32 struct {
	Network Enum   `json:"network"`
	Id      string `json:"id"`
}

type XCMJunctionAccountIndex64 struct {
	Network Enum `json:"network"`
	Index   uint `json:"index"`
}

type XCMJunctionAccountKey20 struct {
	Network Enum   `json:"network"`
	Key     string `json:"key"`
}

type V0MultiLocation struct {
	Here string                 `json:"Here,omitempty"`
	NULL string                 `json:"NULL,omitempty"`
	X1   *XCMJunction           `json:"X1,omitempty"`
	X2   map[string]XCMJunction `json:"X2,omitempty"`
	X3   map[string]XCMJunction `json:"X3,omitempty"`
	X4   map[string]XCMJunction `json:"X4,omitempty"`
	X5   map[string]XCMJunction `json:"X5,omitempty"`
	X6   map[string]XCMJunction `json:"X6,omitempty"`
	X7   map[string]XCMJunction `json:"X7,omitempty"`
	X8   map[string]XCMJunction `json:"X8,omitempty"`
}

func (v *VersionedMultiLocation) UnmarshalJSON(data []byte) error {
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	r := make(map[string]interface{})
	if _, ok := m["V0"]; ok {
		r = m
	} else if _, ok = m["V1"]; ok {
		r = m
	} else if _, ok = m["V2"]; ok {
		r = m
	} else if _, ok = m["V3"]; ok {
		r = m
	} else {
		r["V0"] = m
	}
	type T VersionedMultiLocation
	b, _ := json.Marshal(r)
	return json.Unmarshal(b, (*T)(v))
}
