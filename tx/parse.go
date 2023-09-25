package tx

import (
	"errors"
	"github.com/itering/scale.go/utiles"
	"github.com/itering/substrate-api-rpc/storage"
	"github.com/itering/substrate-api-rpc/util"
)

var MessageRawIsEmptyErr = errors.New("messageRaw is empty")

func (c *Client) ParseXcmMessageInstruction(messageRaw string) (*VersionedXcm, error) {
	typeName, err := GetMessageVersionType(messageRaw, c.XcmVersionTypeName)
	if err != nil {
		return nil, err
	}
	raw, err := storage.Decode(messageRaw, typeName, nil)
	if err != nil {
		return nil, err
	}
	var instruction VersionedXcm
	raw.ToAny(&instruction)
	return &instruction, nil
}

func GetMessageVersionType(messageRaw, defaultType string) (string, error) {
	typeName := "VersionedXcm"
	bytes := util.HexToBytes(messageRaw)
	if len(bytes) == 0 {
		return "", MessageRawIsEmptyErr
	}
	XcmVersion := int(utiles.U256(util.BytesToHex(bytes[0:1])).Uint64())
	if XcmVersion > 1 {
		typeName = defaultType
	}
	return typeName, nil
}
