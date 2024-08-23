package tx

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_HRMP_func(t *testing.T) {
	hrmp := NewHrmp()
	assert.Equal(t, hrmp.ModuleName, "PolkadotXcm")
	t.Run("LimitedReserveTransferAssets", func(t *testing.T) {
		funcName, _ := hrmp.LimitedReserveTransferAssets(nil, nil, nil, 0, nil)
		assert.Equal(t, funcName, "limited_reserve_transfer_assets")
	})
	t.Run("ReserveTransferAssets", func(t *testing.T) {
		funcName, _ := hrmp.ReserveTransferAssets(nil, nil, nil, 0)
		assert.Equal(t, funcName, "reserve_transfer_assets")
	})
	t.Run("LimitedTeleportAssets", func(t *testing.T) {
		funcName, _ := hrmp.LimitedTeleportAssets(nil, nil, nil, 0, nil)
		assert.Equal(t, funcName, "limited_teleport_assets")
	})
	t.Run("TeleportAssets", func(t *testing.T) {
		funcName, _ := hrmp.TeleportAssets(nil, nil, nil, 0)
		assert.Equal(t, funcName, "teleport_assets")
	})
	t.Run("Send", func(t *testing.T) {
		funcName, _ := hrmp.Send(nil, nil)
		assert.Equal(t, funcName, "send")
	})
	t.Run("ReserveTransfer", func(t *testing.T) {
		funcName, _ := hrmp.TransferAssets(nil, nil, nil, 0, nil)
		assert.Equal(t, funcName, "transfer_assets")
	})
}
