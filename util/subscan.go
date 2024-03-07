package util

import (
	"context"
	"encoding/json"
	"fmt"
)

var SubscanEndpoint = map[string]string{
	"polkadot":          "https://polkadot.api.subscan.io",
	"kusama":            "https://kusama.api.subscan.io",
	"assethub-rococo":   "https://assethub-rococo.api.subscan.io",
	"bridgehub-rococo":  "https://bridgehub-rococo.api.subscan.io",
	"assethub-kusama":   "https://assethub-kusama.api.subscan.io",
	"assethub-polkadot": "https://assethub-polkadot.api.subscan.io",
}

const SubscanAPIKey = "98802de223864e7987d3266b3af6e521"

var subscanApiHeaders = map[string]string{
	"X-API-Key": SubscanAPIKey,
}

type SubscanRes[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data,omitempty"`
}
type SubscanBlock struct {
	BlockNum       int    `json:"block_num"`
	BlockTimestamp int    `json:"block_timestamp"`
	Hash           string `json:"hash"`
	ParentHash     string `json:"parent_hash"`
}

func SubscanGetBlockByTime(ctx context.Context, network string, blockTime uint) (*SubscanBlock, error) {
	path := "/api/scan/block"
	url := SubscanEndpoint[network] + path
	if url == "" {
		return nil, fmt.Errorf("network not supported")
	}

	params := map[string]interface{}{"block_timestamp": blockTime, "only_head": true}
	paramsBytes, _ := json.Marshal(params)
	data, err := HttpPost(ctx, paramsBytes, url, subscanApiHeaders)
	if err != nil {
		return nil, err
	}
	var r SubscanRes[SubscanBlock]
	if err = json.Unmarshal(data, &r); err != nil {
		return nil, err
	}
	return &r.Data, nil
}

type SubscanEventRequestParams struct {
	Row        int    `json:"row"`
	Page       int    `json:"page"`
	Module     string `json:"module"`
	EventId    string `json:"event_id"`
	BlockRange string `json:"block_range"`
	Order      string `json:"order"`
}

type SubscanEvent struct {
	Count  int                      `json:"count"`
	Events []SubscanEventWithParams `json:"events"`
}

type SubscanEventWithParams struct {
	Id             int64                `json:"id"`
	BlockTimestamp int                  `json:"block_timestamp"`
	EventIndex     string               `json:"event_index"`
	ExtrinsicIndex string               `json:"extrinsic_index"`
	Phase          int                  `json:"phase"`
	ModuleId       string               `json:"module_id"`
	EventId        string               `json:"event_id"`
	ExtrinsicHash  string               `json:"extrinsic_hash"`
	Finalized      bool                 `json:"finalized"`
	Params         []SubscanEventParams `json:"params"`
}

type SubscanEventParamsRes struct {
	EventIndex string               `json:"event_index"`
	Params     []SubscanEventParams `json:"params"`
}

type SubscanEventParams struct {
	Type     string      `json:"type"`
	TypeName string      `json:"type_name"`
	Value    interface{} `json:"value"`
	Name     string      `json:"name"`
}

func SubscanGetEvents(ctx context.Context, network string, params *SubscanEventRequestParams) ([]SubscanEventWithParams, error) {
	path := "/api/v2/scan/events"
	url := SubscanEndpoint[network] + path
	if url == "" {
		return nil, fmt.Errorf("network not supported")
	}
	paramsBytes, _ := json.Marshal(params)
	data, err := HttpPost(ctx, paramsBytes, url, subscanApiHeaders)
	if err != nil {
		return nil, err
	}
	var r SubscanRes[SubscanEvent]
	if err = json.Unmarshal(data, &r); err != nil {
		return nil, err
	}
	var eventIndex []string
	for _, event := range r.Data.Events {
		eventIndex = append(eventIndex, event.EventIndex)
	}

	paramsBytes, _ = json.Marshal(map[string]interface{}{"event_index": eventIndex})
	url = SubscanEndpoint[network] + "/api/scan/event/params"
	data, err = HttpPost(ctx, paramsBytes, url, subscanApiHeaders)
	if err != nil {
		return nil, err
	}
	var eventParams SubscanRes[[]SubscanEventParamsRes]
	if err = json.Unmarshal(data, &eventParams); err != nil {
		return nil, err
	}
	var eventParamsMap = make(map[string][]SubscanEventParams)
	for _, param := range eventParams.Data {
		eventParamsMap[param.EventIndex] = param.Params
	}
	for i, event := range r.Data.Events {
		r.Data.Events[i].Params = eventParamsMap[event.EventIndex]
	}
	return r.Data.Events, nil
}
