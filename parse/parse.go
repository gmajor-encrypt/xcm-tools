package parse

import (
	"encoding/json"
	"errors"
	"github.com/gmajor-encrypt/xcm-tools/tx"
	ut "github.com/gmajor-encrypt/xcm-tools/util"
	"github.com/itering/scale.go/types"
	"github.com/itering/scale.go/types/scaleBytes"
	"github.com/itering/scale.go/utiles"
	"github.com/itering/substrate-api-rpc/metadata"
	"github.com/itering/substrate-api-rpc/storage"
	"github.com/itering/substrate-api-rpc/util"
	"strings"
)

type Instant struct {
	// xcm version scale type name
	xcmVersionTypeName string
	// xcm version number, like 2,3,4,5
	xcmVersion int
	// metadata instant
	m *metadata.Instant
}

func New(m *metadata.Instant) *Instant {
	c := &Instant{
		m: m,
	}
	c.getXcmLatestVersion()
	return c
}

var MessageRawIsEmptyErr = errors.New("messageRaw is empty")

func (c *Instant) ParseXcmMessageInstruction(messageRaw string) (*tx.VersionedXcm, error) {
	typeName, err := c.getMessageVersionType(messageRaw, c.xcmVersionTypeName)
	if err != nil {
		return nil, err
	}

	// scale decode message
	raw, err := storage.Decode(messageRaw, typeName, nil)
	if err != nil {
		return nil, err
	}

	var instruction tx.VersionedXcm
	raw.ToAny(&instruction)
	return &instruction, nil
}

func (c *Instant) getMessageVersionType(messageRaw, defaultType string) (string, error) {
	// default type name is VersionedXcm
	typeName := "VersionedXcm"
	bytes := util.HexToBytes(messageRaw)
	if len(bytes) == 0 {
		return "", MessageRawIsEmptyErr
	}

	XcmVersion := int(utiles.U256(util.BytesToHex(bytes[0:1])).Uint64())
	// because xcm V0,V1 has been removed, so we need to check if the version is greater than 1
	// default type name is metadata scale type name
	if XcmVersion > 1 {
		typeName = defaultType
	}
	return typeName, nil
}

func (c *Instant) getXcmLatestVersion() {
	moduleName := "XcmPallet"

	// parachain xcm module name is PolkadotXcm
	if tx.GetModule(moduleName, c.m) == nil {
		moduleName = "PolkadotXcm"
	}

	call := tx.GetCallByName(moduleName, "send", c.m)

	if call != nil {
		// get xcm version scale type name
		c.xcmVersionTypeName = call.Args[1].Type

		r := types.RuntimeType{}
		_, value, _ := r.GetCodecClass(c.xcmVersionTypeName, 0)

		var mappingTypes types.TypeMapping
		b, _ := json.Marshal(value.Elem().FieldByName("TypeMapping").Interface())
		_ = json.Unmarshal(b, &mappingTypes)

		for _, name := range mappingTypes.Names {
			if strings.HasPrefix(name, "V") {
				c.xcmVersion = ut.ToInt(strings.ReplaceAll(name, "V", ""))
			}
		}
		return
	}
	panic("The current chain does not support xcm")
}

func (c *Instant) DecodeFixedMessage(messageRaw string) []string {
	if messageRaw == "" {
		return []string{messageRaw}
	}
	var messages []string
	var tryCount int

	var retryDecode = func(raw []byte, decodeType string) (last []byte) {
		defer func() {
			if r := recover(); r != nil {
				last = nil
			}
		}()
		scale := types.ScaleDecoder{}
		scale.Init(scaleBytes.ScaleBytes{Data: raw}, nil)
		scale.ProcessAndUpdateData(decodeType)
		// if remaining length is 0, means all messages are decoded
		if scale.Data.GetRemainingLength() == 0 {
			return nil
		}
		return scale.Data.GetNextBytes(scale.Data.GetRemainingLength())
	}
	// hrmp message is a fixed count messages, but raw message not have array length
	// so we need to try decode message until all messages are decoded
	for {
		if tryCount > 200 {
			break
		}
		tryCount++
		remain := retryDecode(util.HexToBytes(messageRaw), c.xcmVersionTypeName)
		if len(remain) == 0 {
			messages = append(messages, messageRaw)
			break
		}
		messages = append(messages, messageRaw[0:len(messageRaw)-len(util.BytesToHex(remain))])
		messageRaw = util.BytesToHex(remain)
	}
	return messages
}
