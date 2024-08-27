package tracker

import (
	"github.com/itering/scale.go/types"
	"github.com/itering/substrate-api-rpc/metadata"
	"github.com/itering/substrate-api-rpc/rpc"
	"strings"
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

// PendingAvailability Candidates pending availability by `ParaId
// after runtime 1010000,pendingAvailability rename to v1
func PendingAvailability(paraId uint, blockHash string, m *metadata.Instant) (*Inclusion, error) {
	// check extrinsic call
	callName := "pendingAvailability"
	for _, module := range m.Metadata.Modules {
		if module.Name == "ParaInclusion" {
			for _, st := range module.Storage {
				if strings.EqualFold(st.Name, "v1") {
					callName = "v1"
					break
				}
			}
		}
	}
	raw, err := rpc.ReadStorage(nil, "paraInclusion", callName, blockHash, types.Encode("U32", uint32(paraId)))
	if err != nil {
		return nil, err
	}
	// v1 call
	if callName == "v1" {
		var inclusions []Inclusion
		raw.ToAny(&inclusions)
		if len(inclusions) == 0 {
			return nil, nil
		}
		var inclusion = inclusions[0]
		return &inclusion, nil
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

func SystemNumber(blockHash string) (int64, error) {
	raw, err := rpc.ReadStorage(nil, "System", "Number", blockHash)
	if err != nil {
		return 0, err
	}
	return raw.ToInt64(), nil
}

func timestampNow(blockHash string) (int64, error) {
	raw, err := rpc.ReadStorage(nil, "Timestamp", "Now", blockHash)
	if err != nil {
		return 0, err
	}
	return raw.ToInt64(), nil
}

type BpParachainsParaInfo struct {
	BestHeadHash struct {
		AtRelayBlockNumber uint64 `json:"at_relay_block_number"`
		HeadHash           string `json:"head_hash"`
	} `json:"best_head_hash"`
}

func ParasInfo(module string, paraId uint, blockHash string) (*BpParachainsParaInfo, error) {
	raw, err := rpc.ReadStorage(nil, module, "parasInfo", blockHash, types.Encode("U32", uint32(paraId)))
	if err != nil {
		return nil, err
	}
	if raw == "" {
		return nil, nil
	}
	var info BpParachainsParaInfo
	raw.ToAny(&info)
	return &info, nil
}
