package tracker

import (
	"context"
	"errors"
	"github.com/gmajor-encrypt/xcm-tools/tx"
	"github.com/gmajor-encrypt/xcm-tools/util"
	"github.com/itering/scale.go/types"
	"github.com/itering/scale.go/utiles"
	"github.com/itering/substrate-api-rpc/hasher"
	"github.com/itering/substrate-api-rpc/metadata"
	"github.com/itering/substrate-api-rpc/rpc"
	"github.com/itering/substrate-api-rpc/websocket"
	"strings"
)

type ITracker interface {
	Track(ctx context.Context) (*Event, error)
}

var (
	NotfoundXcmMessageErr = errors.New("not found xcm message")
	InvalidDestParaId     = errors.New("invalid dest para id")
)

func CreateSnapshotClient(endpoint string) (*websocket.PoolConn, *metadata.Instant, func()) {
	websocket.SetEndpoint(endpoint)
	// websocket init and set metadata
	p, err := websocket.Init()
	if err != nil {
		panic(err)
	}
	raw, err := rpc.GetMetadataByHash(nil)
	if err != nil {
		panic(err)
	}
	metadataInstant := metadata.RegNewMetadataType(0, raw)
	return p, metadataInstant, func() {
		websocket.Close()
	}
}

const (
	Parachain  = "parachain"
	Relaychain = "relaychain"
	Solo       = "solo"
)

// checkChain check chain is parachain/relaychain or not
func checkChain(m *metadata.Instant) string {
	switch {
	case tx.GetModule("XcmPallet", m) != nil:
		return Relaychain
	case tx.GetModule("PolkadotXcm", m) != nil:
		return Parachain
	default:
		return Solo
	}
}

type ExtrinsicIndex struct {
	BlockNum uint
	Index    uint
}

func findOutBlockByExtrinsicIndex(extrinsicIndex string) *ExtrinsicIndex {
	if sliceIndex := strings.Split(extrinsicIndex, "-"); len(sliceIndex) == 2 {
		return &ExtrinsicIndex{BlockNum: util.ToUint(sliceIndex[0]), Index: util.ToUint(sliceIndex[1])}
	}
	return nil
}

func hash(hex string) string {
	return utiles.AddHex(utiles.BytesToHex(hasher.HashByCryptoName(utiles.HexToBytes(hex), "Blake2_256")))
}

func TrackXcmMessage(extrinsicIndex string, protocol tx.Protocol, originEndpoint, destEndpoint, relayEndpoint string) (*Event, error) {
	_, metadataInstant, cancel := CreateSnapshotClient(originEndpoint)
	chain := checkChain(metadataInstant)
	if chain == Solo {
		return nil, errors.New("originEndpoint not parachain or relaychain")
	}
	ctx := context.Background()
	cancel()
	defer func() {
		types.Clean()
	}()
	switch protocol {
	case tx.UMP:
		if chain != Parachain {
			return nil, errors.New("originEndpoint not parachain")
		}
		u := Ump{ExtrinsicIndex: extrinsicIndex, OriginEndpoint: originEndpoint, DestEndpoint: destEndpoint}
		return u.Track(ctx)
	case tx.HRMP:
		if chain != Parachain {
			return nil, errors.New("originEndpoint not parachain")
		}
		h := Hrmp{extrinsicIndex: extrinsicIndex, originEndpoint: originEndpoint, destEndpoint: destEndpoint, relayChainEndpoint: relayEndpoint}
		return h.Track(ctx)
	case tx.DMP:
		if chain != Relaychain {
			return nil, errors.New("originEndpoint not relaychain")
		}
		d := Dmp{extrinsicIndex: extrinsicIndex, originEndpoint: originEndpoint, destEndpoint: destEndpoint}
		return d.Track(ctx)
	default:
		return nil, errors.New("not support solo chain")
	}
}
