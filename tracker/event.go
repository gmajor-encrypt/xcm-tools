package tracker

import (
	"context"
	"encoding/json"
	scalecodec "github.com/itering/scale.go"
	"github.com/itering/scale.go/types"
	"github.com/itering/scale.go/types/scaleBytes"
	"github.com/itering/scale.go/utiles"
	"github.com/itering/substrate-api-rpc/model"
	"github.com/itering/substrate-api-rpc/pkg/recws"
	"github.com/itering/substrate-api-rpc/rpc"
	"github.com/itering/substrate-api-rpc/storageKey"
	"github.com/itering/substrate-api-rpc/util"
	"github.com/itering/substrate-api-rpc/websocket"
	"math/rand"
)

type Event struct {
	// extrinsic index in block
	ExtrinsicIdx int          `json:"extrinsic_idx" `
	BlockNum     int          `json:"block_num" `
	ModuleId     string       `json:"module_id"  `
	EventId      string       `json:"event_id" `
	Params       []EventParam `json:"params"`
	// block timestamp
	BlockTime int64 `json:"block_time"`

	TxHash       string      `json:"tx_hash"`
	ResourceData interface{} `json:"resource_data"`
}

type EventParam struct {
	Type     string      `json:"type"`
	TypeName string      `json:"type_name,omitempty"`
	Value    interface{} `json:"value"`
	Name     string      `json:"name,omitempty"`
}

// get events raw from chain
func getEventsFromChain(_ context.Context, p *recws.RecConn, blockHash string) (string, error) {
	key := storageKey.EncodeStorageKey("system", "events")
	v := &model.JsonRpcResult{}
	if err := websocket.SendWsRequest(p, v, rpc.StateGetStorage(rand.Intn(10000), util.AddHex(key.EncodeKey), blockHash)); err != nil {
		return "", err
	}
	return v.ToString()
}

// get events from chain, and decode it
func getEvents(ctx context.Context, p *websocket.PoolConn, blockHash string, metadata *types.MetadataStruct) ([]Event, error) {
	raw, err := getEventsFromChain(ctx, p.Conn, blockHash)
	if err != nil {
		return nil, err
	}

	e := scalecodec.EventsDecoder{}
	option := types.ScaleDecoderOption{Metadata: metadata}
	e.Init(scaleBytes.ScaleBytes{Data: utiles.HexToBytes(raw)}, &option)
	e.Process()

	var events []Event
	bytes, _ := json.Marshal(e.Value)
	err = json.Unmarshal(bytes, &events)
	return events, err
}

// find event by event id
func findEventByEventId(events []Event, index int, eventId []string) *Event {
	for _, event := range events {
		if util.StringInSlice(event.EventId, eventId) {
			if index >= 0 {
				if event.ExtrinsicIdx != index {
					continue
				}
			}
			return &event
		}
	}
	return nil
}
