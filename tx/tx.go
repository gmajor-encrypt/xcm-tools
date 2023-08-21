package tx

import (
	"encoding/json"
	"github.com/itering/substrate-api-rpc/keyring"
	"github.com/shopspring/decimal"
	"strings"

	"github.com/gmajor-encrypt/xcm-tools/util"
	"github.com/itering/scale.go/types"
	"github.com/itering/substrate-api-rpc/metadata"
	"github.com/itering/substrate-api-rpc/rpc"
	"github.com/itering/substrate-api-rpc/websocket"
)

type Client struct {
	Conn       *rpc.Client
	Ump        IXCMP
	Dmp        IXCMP
	Hrmp       IXCMP
	XcmVersion uint
	m          *metadata.Instant
}

// NewClient create a new XCM client
func NewClient(endpoint string) *Client {
	websocket.SetEndpoint(endpoint)
	// websocket init and set metadata
	_, err := websocket.Init()
	if err != nil {
		panic(err)
	}

	raw, err := rpc.GetMetadataByHash(nil)
	if err != nil {
		panic(err)
	}
	c := &rpc.Client{}
	metadataInstant := metadata.RegNewMetadataType(0, raw)
	c.SetMetadata(metadataInstant)
	client := &Client{
		Conn: c,
		Ump:  NewMessage(UMP),
		Dmp:  NewMessage(DMP),
		Hrmp: NewMessage(HRMP),
		m:    metadataInstant,
	}

	client.getXcmLatestVersion()
	return client
}

func (c *Client) SetKeyRing(sk string) {
	c.Conn.SetKeyRing(keyring.New(keyring.Sr25519Type, sk))
}

func (c *Client) Close() {
	websocket.Close()
}

// getCallByName get call by module name and call name
func (c *Client) getCallByName(moduleName, callName string) *types.MetadataCalls {
	module := c.getModule(moduleName)
	if module == nil {
		return nil
	}

	for i, call := range module.Calls {
		if strings.EqualFold(call.Name, callName) {
			return &module.Calls[i]
		}
	}
	return nil
}

// getModule get module by name
func (c *Client) getModule(moduleName string) *types.MetadataModules {
	for i, v := range c.m.Metadata.Modules {
		if strings.EqualFold(v.Name, moduleName) {
			return &c.m.Metadata.Modules[i]
		}
	}
	return nil
}

func (c *Client) getXcmLatestVersion() {
	moduleName := "XcmPallet"
	if c.getModule(moduleName) == nil {
		moduleName = "PolkadotXcm"
	}

	call := c.getCallByName(moduleName, "send")

	if call != nil {
		versionedXcmType := call.Args[1].Type

		r := types.RuntimeType{}
		_, value, _ := r.GetCodecClass(versionedXcmType, 0)

		var mappingTypes types.TypeMapping
		b, _ := json.Marshal(value.Elem().FieldByName("TypeMapping").Interface())
		_ = json.Unmarshal(b, &mappingTypes)

		for _, name := range mappingTypes.Names {
			if strings.HasPrefix(name, "V") {
				c.XcmVersion = uint(util.ToInt(strings.ReplaceAll(name, "V", "")))
			}
		}
		return
	}
	panic("The current chain does not support xcm")
}

// SendUmpTransfer send XCM UMP message
// accountId: the account id of the beneficiary
// amount: the amount of the asset to be transferred
func (c *Client) SendUmpTransfer(accountId string, amount decimal.Decimal) (string, error) {
	callName, args := c.Ump.LimitedReserveTransferAssets(
		SimplifyMultiLocationRelayChain(),
		SimplifyMultiLocationAccountId32(accountId),
		SimplifyMultiAssets(amount),
		0,
		SimplifyUnlimitedWeight(),
	)
	signed, err := c.Conn.SignTransaction(c.Ump.GetModuleName(), callName, args...)
	if err != nil {
		return "", err
	}
	tx, err := c.Conn.SendAuthorSubmitExtrinsic(signed)
	return tx, err
}

// SendDmpTransfer send XCM DMP message
// paraId: the para id of the parachain
// accountId: the account id of the beneficiary
// amount: the amount of the asset to be transferred
func (c *Client) SendDmpTransfer(paraId uint32, accountId string, amount decimal.Decimal) (string, error) {
	callName, args := c.Dmp.LimitedReserveTransferAssets(
		SimplifyMultiLocationParaId(paraId),
		SimplifyMultiLocationAccountId32(accountId),
		SimplifyMultiAssets(amount),
		0,
		SimplifyUnlimitedWeight(),
	)
	signed, err := c.Conn.SignTransaction(c.Dmp.GetModuleName(), callName, args...)
	if err != nil {
		return "", err
	}
	tx, err := c.Conn.SendAuthorSubmitExtrinsic(signed)
	return tx, err
}

// SendHrmpTransfer send XCM HRMP message
// paraId: the para id of the parachain
// accountId: the account id of the beneficiary
// amount: the amount of the asset to be transferred
func (c *Client) SendHrmpTransfer(paraId uint32, accountId string, amount decimal.Decimal) (string, error) {
	callName, args := c.Hrmp.LimitedReserveTransferAssets(
		SimplifyMultiLocationParaId(paraId),
		SimplifyMultiLocationAccountId32(accountId),
		SimplifyMultiAssets(amount),
		0,
		SimplifyUnlimitedWeight(),
	)
	signed, err := c.Conn.SignTransaction(c.Ump.GetModuleName(), callName, args...)
	if err != nil {
		return "", err
	}
	tx, err := c.Conn.SendAuthorSubmitExtrinsic(signed)
	return tx, err
}
