package tracker

import (
	"github.com/itering/scale.go/types"
	"github.com/itering/substrate-api-rpc/rpc"
)

func HRMPWatermark(blockHash string) (int, error) {
	raw, err := rpc.ReadStorage(nil, "parachainSystem", "hrmpWatermark", blockHash)
	if err != nil {
		return 0, err
	}
	return raw.ToInt(), nil
}

type InboundDownwardMessage struct {
	SentAt uint   `json:"sent_at"`
	Msg    string `json:"msg"`
}

func DownwardMessageQueues(paraId uint, blockHash string) ([]InboundDownwardMessage, error) {
	raw, err := rpc.ReadStorage(nil, "DMP", "DownwardMessageQueues", blockHash, types.Encode("U32", uint32(paraId)))
	if err != nil {
		return nil, err
	}
	var list []InboundDownwardMessage
	raw.ToAny(&list)
	return list, nil
}

type Inclusion struct {
	Descriptor struct {
		ParaHead string `json:"para_head"`
	} `json:"descriptor"`
}

func PendingAvailability(paraId uint, blockHash string) (*Inclusion, error) {
	raw, err := rpc.ReadStorage(nil, "paraInclusion", "pendingAvailability", blockHash, types.Encode("U32", uint32(paraId)))
	if err != nil {
		return nil, err
	}
	var inclusion Inclusion
	raw.ToAny(&inclusion)
	return &inclusion, nil
}

func UpwardMessages(blockHash string) ([]string, error) {
	raw, err := rpc.ReadStorage(nil, "parachainSystem", "upwardMessages", blockHash)
	if err != nil {
		return nil, err
	}
	return raw.ToStringSlice(), nil
}

type OutboundMessage struct {
	Recipient uint   `json:"recipient"`
	Data      string `json:"data"`
}

func HrmpOutboundMessages(blockHash string) (list []OutboundMessage, err error) {
	raw, err := rpc.ReadStorage(nil, "parachainSystem", "hrmpOutboundMessages", blockHash)
	if err != nil {
		return
	}
	raw.ToAny(&list)
	return
}
