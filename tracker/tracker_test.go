package tracker

import (
	"github.com/gmajor-encrypt/xcm-tools/tx"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_TrackTx(t *testing.T) {
	// ump
	// https://polkadot.subscan.io/xcm_message/polkadot-bce0edc64c3c2af9f903a55e537037b88f35503f
	event, err := TrackXcmMessage("4310901-13", tx.UMP, "wss://moonbeam.api.onfinality.io/public-ws", "wss://polkadot.api.onfinality.io/public-ws", "")
	assert.NoError(t, err)
	assert.NotNil(t, event)
	assert.Equal(t, "0xf053b9f150fbff79347bb6ed438e9cdf20083b7dbfa203078e9648b5bbfaa902", event.Params[0].Value.(string))

	// dmp
	// https://polkadot.subscan.io/xcm_message/polkadot-670f66a3cf0f8f0523fe7d490ee435cd474fa672
	event, err = TrackXcmMessage("17053966-2", tx.DMP, "wss://polkadot.api.onfinality.io/public-ws", "wss://moonbeam.api.onfinality.io/public-ws", "")
	assert.NoError(t, err)
	assert.Equal(t, "0x99f8179f1a3ca331998e7369ced93ac187036f287dcd3d015f3bcc585df92fa4", event.Params[0].Value.(string))

	// hrmp
	// https://polkadot.subscan.io/xcm_message/polkadot-13586a835ebe97b4e2d046233aac26657f64da04
	event, err = TrackXcmMessage("4325642-7", tx.HRMP, "wss://astar.api.onfinality.io/public-ws", "wss://rpc.hydradx.cloud", "wss://polkadot.api.onfinality.io/public-ws")
	assert.NoError(t, err)
	assert.Equal(t, "0x5d81466ae4b2d9fb1fd140cd690bb25276b0bfafabecd62840c67e0b062c8181", event.Params[0].Value.(string))
}
