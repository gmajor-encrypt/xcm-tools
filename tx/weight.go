package tx

import (
	"encoding/json"
)

type Weight struct {
	Limited   *WeightLimited `json:"Limited,omitempty"`
	Unlimited *string        `json:"Unlimited,omitempty"`
}

type WeightLimited struct {
	ProofSize uint   `json:"proof_size"`
	RefTime   uint64 `json:"ref_time"`
}

func (w *Weight) ToScale() interface{} {
	r := map[string]interface{}{}
	b, _ := json.Marshal(w)
	_ = json.Unmarshal(b, &r)
	return r
}
