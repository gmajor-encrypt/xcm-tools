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
	ExtrinsicIdx int          `json:"extrinsic_idx" `
	ModuleId     string       `json:"module_id"  `
	EventId      string       `json:"event_id" `
	Params       []EventParam `json:"params"`
}

type EventParam struct {
	Type     string      `json:"type"`
	TypeName string      `json:"type_name,omitempty"`
	Value    interface{} `json:"value"`
	Name     string      `json:"name,omitempty"`
}

func getEventsFromChain(_ context.Context, p *recws.RecConn, blockHash string) (string, error) {
	key := storageKey.EncodeStorageKey("system", "events")
	v := &model.JsonRpcResult{}
	if err := websocket.SendWsRequest(p, v, rpc.StateGetStorage(rand.Intn(10000), util.AddHex(key.EncodeKey), blockHash)); err != nil {
		return "", err
	}
	return v.ToString()
}

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

func findEventByEventId(events []Event, index uint, eventId []string) *Event {
	for _, event := range events {
		if util.StringInSlice(event.EventId, eventId) && event.ExtrinsicIdx == int(index) {
			return &event
		}
	}
	return nil
}
