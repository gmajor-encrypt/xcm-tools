package tracker

import (
	"context"
	"github.com/gmajor-encrypt/xcm-tools/tx"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_TrackTx(t *testing.T) {
	// ump
	// https://polkadot.subscan.io/xcm_message/polkadot-bce0edc64c3c2af9f903a55e537037b88f35503f
	event, err := TrackXcmMessage("4310901-13", tx.UMP, "wss://moonbeam-rpc.dwellir.com", "wss://polkadot-rpc.dwellir.com", "")
	assert.NoError(t, err)
	assert.NotNil(t, event)
	assert.Equal(t, "0xf053b9f150fbff79347bb6ed438e9cdf20083b7dbfa203078e9648b5bbfaa902", event.Params[0].Value.(string))

	// dmp
	// https://polkadot.subscan.io/xcm_message/polkadot-670f66a3cf0f8f0523fe7d490ee435cd474fa672
	event, err = TrackXcmMessage("17053966-2", tx.DMP, "wss://polkadot-rpc.dwellir.com", "wss://moonbeam-rpc.dwellir.com", "")
	assert.NoError(t, err)
	assert.Equal(t, "0x99f8179f1a3ca331998e7369ced93ac187036f287dcd3d015f3bcc585df92fa4", event.Params[0].Value.(string))

	// hrmp
	// https://polkadot.subscan.io/xcm_message/polkadot-13586a835ebe97b4e2d046233aac26657f64da04
	event, err = TrackXcmMessage("4325642-7", tx.HRMP, "wss://astar-rpc.dwellir.com", "wss://rpc.hydradx.cloud", "wss://polkadot-rpc.dwellir.com")
	assert.NoError(t, err)
	assert.Equal(t, "0x5d81466ae4b2d9fb1fd140cd690bb25276b0bfafabecd62840c67e0b062c8181", event.Params[0].Value.(string))
}

func TestTrackBridgeMessage(t *testing.T) {
	ctx := context.Background()

	var err error
	// ethereum => polkadot
	// https://sepolia.etherscan.io/tx/0x799f01445e2be3103a1a751e33b395c4b894529ce3b320d2fd94c22d4e3d6e01
	event, err := TrackBridgeMessage(ctx, &TrackBridgeMessageOptions{Tx: "0x799f01445e2be3103a1a751e33b395c4b894529ce3b320d2fd94c22d4e3d6e01"})
	assert.NoError(t, err)
	assert.Equal(t, event.BlockNum, 2542078)
	assert.Equal(t, event.ExtrinsicIdx, 2)

	// polkadot => ethereum
	// https://assethub-rococo.subscan.io/extrinsic/3879712-2
	event, err = TrackBridgeMessage(ctx, &TrackBridgeMessageOptions{
		ExtrinsicIndex:    "3879712-2",
		BridgeHubEndpoint: "wss://rococo-bridge-hub-rpc.polkadot.io",
		OriginEndpoint:    "wss://rococo-rockmine-rpc.polkadot.io",
		RelayEndpoint:     "wss://rococo-rpc.polkadot.io",
	})
	assert.NoError(t, err)
	assert.Equal(t, event.TxHash, "0x1435866e5c320adac9fed7827934ce6c34f28bf6cc2b5fae1ab3f5512fd0db76")

	_, err = TrackBridgeMessage(ctx, &TrackBridgeMessageOptions{})
	// will raise extrinsicIndexEmptyError error
	assert.ErrorIs(t, err, extrinsicIndexEmptyError)

	// will raise not found message id error
	_, err = TrackBridgeMessage(ctx, &TrackBridgeMessageOptions{Tx: "0x2eb3249c5d64a617cc31a0c01e29f98a03400d066699d525a1fc17a0d0210660"})
	assert.ErrorIs(t, err, notFindBridgeMessageIdError)

}
