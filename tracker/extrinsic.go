package tracker

import (
	"context"
	"errors"
	scalecodec "github.com/itering/scale.go"
	"github.com/itering/scale.go/types"
	"github.com/itering/scale.go/types/scaleBytes"
	"github.com/itering/scale.go/utiles"
	"github.com/itering/substrate-api-rpc/model"
	"github.com/itering/substrate-api-rpc/rpc"
	"github.com/itering/substrate-api-rpc/websocket"
	"math/rand"
)

type BlockData struct {
	Block Block `json:"block"`
}
type Block struct {
	Extrinsics []string `json:"extrinsics"`
}

// get extrinsics raw from block
func getExtrinsics(_ context.Context, p *websocket.PoolConn, blockHash string) ([]string, error) {
	v := &model.JsonRpcResult{}
	if err := websocket.SendWsRequest(p.Conn, v, rpc.ChainGetBlock(rand.Intn(10), blockHash)); err != nil {
		return nil, err
	}
	var block BlockData
	err := v.ToAnyThing(&block)
	if err != nil {
		return nil, err
	}
	return block.Block.Extrinsics, nil
}

// get extrinsic by index
func getExtrinsicByIndex(_ context.Context, extrinsics []string, index int, metadata *types.MetadataStruct) (*scalecodec.GenericExtrinsic, error) {
	if len(extrinsics) <= index {
		return nil, errors.New("index out of range")
	}

	raw := extrinsics[index]
	e := scalecodec.ExtrinsicDecoder{}
	option := types.ScaleDecoderOption{Metadata: metadata}
	e.Init(scaleBytes.ScaleBytes{Data: utiles.HexToBytes(raw)}, &option)
	e.Process()
	return e.Value.(*scalecodec.GenericExtrinsic), nil
}
