package tracker

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHRMPWatermark(t *testing.T) {
	_, closeF := CreateSnapshotClient("wss://moonbeam-rpc.dwellir.com")
	_ = RegDefaultMetadata()
	blockNum, err := HRMPWatermark("0xebaa16b5cc4c53595acb541ea46fdf9d6625f00ea335e91c53391d13cb599f36")
	assert.NoError(t, err)
	assert.Equal(t, 17040495, blockNum)
	closeF()
}

func TestParasInfo(t *testing.T) {
	_, closeF := CreateSnapshotClient("wss://polkadot-bridge-hub-rpc.polkadot.io")
	_ = RegDefaultMetadata()
	paraInfo, err := ParasInfo("bridgeKusamaParachains", 1002, "0x18c4e703423958f0f5571362db411b1bf386822059ae0cd948163f16aa9c1975")
	assert.NoError(t, err)
	assert.Equal(t, paraInfo.BestHeadHash.HeadHash, "0xd8e88896f7dd8cc280e4cc337f49c3048ed3bbca6cae9fb6afaeef905aa6577c")
	closeF()
}
