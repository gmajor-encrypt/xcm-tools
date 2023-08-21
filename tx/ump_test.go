package tx

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Ump_func(t *testing.T) {
	dmp := NewUmp()
	assert.Equal(t, dmp.ModuleName, "PolkadotXcm")
	t.Run("LimitedReserveTransferAssets", func(t *testing.T) {
		funcName, _ := dmp.LimitedReserveTransferAssets(nil, nil, nil, 0, nil)
		assert.Equal(t, funcName, "limited_reserve_transfer_assets")
	})
	t.Run("ReserveTransferAssets", func(t *testing.T) {
		funcName, _ := dmp.ReserveTransferAssets(nil, nil, nil, 0)
		assert.Equal(t, funcName, "reserve_transfer_assets")
	})
	t.Run("LimitedTeleportAssets", func(t *testing.T) {
		funcName, _ := dmp.LimitedTeleportAssets(nil, nil, nil, 0, nil)
		assert.Equal(t, funcName, "limited_teleport_assets")
	})
	t.Run("TeleportAssets", func(t *testing.T) {
		funcName, _ := dmp.TeleportAssets(nil, nil, nil, 0)
		assert.Equal(t, funcName, "teleport_assets")
	})
	t.Run("Send", func(t *testing.T) {
		funcName, _ := dmp.Send(nil, nil)
		assert.Equal(t, funcName, "send")
	})
}
