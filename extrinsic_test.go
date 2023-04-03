package scalecodec_test

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/itering/scale.go"
	"github.com/itering/scale.go/source"
	"github.com/itering/scale.go/types"
	"github.com/itering/scale.go/types/scaleBytes"
	"github.com/itering/scale.go/utiles"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestV14ExtrinsicDecoder(t *testing.T) {
	m := scalecodec.MetadataDecoder{}
	m.Init(utiles.HexToBytes(kusamaV14))
	_ = m.Process()
	c, err := ioutil.ReadFile(fmt.Sprintf("%s.json", "network/kusama"))
	if err != nil {
		panic(err)
	}
	types.RegCustomTypes(source.LoadTypeRegistry(c))
	e := scalecodec.ExtrinsicDecoder{}
	option := types.ScaleDecoderOption{Metadata: &m.Metadata, Spec: 9111}
	extrinsicRaw := "0x1d03840018c7717a3c5d2930f10eb5b0f67c191210e018bc55481bc44c1c1c7e810e945c01922c584c1c205b9747e76aad28cf2e73f729a9b3180072c47cd3abd205bb4b54f78a2627fa62a799f363fde25b5db74e5f8d8f3bde7a9a7f2ea8c95033d84e8585030800630301000400000000070088526a74080700000000070088526a74005ed0b200000000005ed0b20000000000000101000100000000010100707fd754e80e531ad411eb8b433548acbe669bf46a7e896e440feadc5ef3530800bca06501000000"
	e.Init(scaleBytes.ScaleBytes{Data: utiles.HexToBytes(extrinsicRaw)}, &option)
	e.Process()
}

func TestCallEncode(t *testing.T) {
	m := scalecodec.MetadataDecoder{}
	m.Init(utiles.HexToBytes(kusamaV14))
	_ = m.Process()

	callMapInterface := map[string]interface{}{
		"call_index": "0x0001",
		"params":     []types.ExtrinsicParam{{Type: "Bytes", Value: "0x11111"}},
	}
	assert.Equal(t, types.EncodeWithOpt("Call", callMapInterface, &types.ScaleDecoderOption{Metadata: &m.Metadata}), "00010c111110")

	call := types.BoxCall{
		CallIndex: "0400",
		Params: []types.ExtrinsicParam{
			{Type: "sp_runtime:multiaddress:MultiAddress", Value: map[string]interface{}{"Id": "9094c424429709a324e65c64f151630e6c3700192bba8abd3c8e2218b61c0a7a"}},
			{Type: "Compact<u128>", Value: decimal.RequireFromString("1000000000000000000")},
		}}

	assert.Equal(t, "0400009094c424429709a324e65c64f151630e6c3700192bba8abd3c8e2218b61c0a7a13000064a7b3b6e00d", types.EncodeWithOpt("Call", call, &types.ScaleDecoderOption{Metadata: &m.Metadata}), "")
}

func TestExtrinsicEncode(t *testing.T) {
	m := scalecodec.MetadataDecoder{}
	m.Init(utiles.HexToBytes(Kusama9370))
	_ = m.Process()
	option := types.ScaleDecoderOption{Metadata: &m.Metadata}
	extrinsicRaw := []string{"280402000b91bdddd38601", "3d028400e673cb35ffaaf7ab98c4e9268bfa9b4a74e49d41c8225121c346db7a7dd06d8800fce9453b1442bba86c2781e755a29c8a215ccf4b65ce81eeaa5b5a04dcdb79a54525cc86969f910c71c05f84aeab9c205022ecd4aa2abb4a3c3667f09dd16e0b0000000400000770e0831a275b534f7507c8ebd9f5f982a55053c9dc672da886ef41a6b5c62807003c7bb7fe"}
	for _, raw := range extrinsicRaw {
		e := scalecodec.ExtrinsicDecoder{}
		e.Init(scaleBytes.ScaleBytes{Data: utiles.HexToBytes(raw)}, &option)
		e.Process()
		assert.Equal(t, e.Value.(*scalecodec.GenericExtrinsic).EncodeWithOpt(&option), raw)
	}
}
