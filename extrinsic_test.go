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

	extrinsicRaw := []string{
		"280402000b91bdddd38601",
		"3d028400e673cb35ffaaf7ab98c4e9268bfa9b4a74e49d41c8225121c346db7a7dd06d8800fce9453b1442bba86c2781e755a29c8a215ccf4b65ce81eeaa5b5a04dcdb79a54525cc86969f910c71c05f84aeab9c205022ecd4aa2abb4a3c3667f09dd16e0b0000000400000770e0831a275b534f7507c8ebd9f5f982a55053c9dc672da886ef41a6b5c62807003c7bb7fe",
		"4d0384000ed3eb9c7be61e0f8e13abb04857d82be899c1d7c3203998e02538720052ad34019edcf93212fc47fd9934311d10363a2cdbaa6efb52660bfa335325a2eb78c52d1052ff7b3d4c926ec2ed02e9ce1ddacf34211a97d467b8e1fe9e8b5f63fa9f8d550005040000009d01524d524b3a3a45515549503a3a322e302e303a3a31333235343234332d3865646566396635326238316339306132612d434c4f544845532d53555350454e444552532d30303030303030333a3a626173652d31333235323835322d444742502e434c4f54484553",
		"bd1384001641233add39a0bfb8941c27f3091e320fab4032db6153c4597408dc2733e52e01ae293222b8a9e7356522429a0bfa448bd5aeedeb8fafd214f5e36507c8e2b8221b19fff42d179072667ecfa8eade81d9be9ae80e7890f991c7ceb6fbfe4d638f350065070018023000007901524d524b3a3a4c4953543a3a322e302e303a3a31353837383636302d3334326631323130366561623664393034632d594f55444c4543484553542d594f55444c4543484553542d30303031313231343a3a3836353030303030303030303000006901524d524b3a3a4c4953543a3a322e302e303a3a31353838303530312d3334326631323130366561623664393034632d594f55444c454143432d594f55444c454143432d30303032333638353a3a3433323530303030303030303000007101524d524b3a3a4c4953543a3a322e302e303a3a31353838303237352d3334326631323130366561623664393034632d594f55444c45484149522d594f55444c45484149522d30303032323339353a3a3137333030303030303030303000007101524d524b3a3a4c4953543a3a322e302e303a3a31353837393232312d3334326631323130366561623664393034632d594f55444c45464143452d594f55444c45464143452d30303031353432383a3a3334363030303030303030303000007101524d524b3a3a4c4953543a3a322e302e303a3a31353837373734332d3334326631323130366561623664393034632d594f55444c45534b494e2d594f55444c45534b494e2d30303030363934363a3a3137333030303030303030303000007101524d524b3a3a4c4953543a3a322e302e303a3a31353837393732372d3334326631323130366561623664393034632d594f55444c45455945532d594f55444c45455945532d30303031383537383a3a3639323030303030303030303000007101524d524b3a3a4c4953543a3a322e302e303a3a31353837383839372d3334326631323130366561623664393034632d594f55444c45454152532d594f55444c45454152532d30303031333131303a3a3235393530303030303030303000007101524d524b3a3a4c4953543a3a322e302e303a3a31353837393835382d3334326631323130366561623664393034632d594f55444c45484149522d594f55444c45484149522d30303031393636373a3a3433323530303030303030303000006d01524d524b3a3a4c4953543a3a322e302e303a3a31353838303835302d3334326631323130366561623664393034632d594f55444c454143432d594f55444c454143432d30303032353530383a3a32313632353030303030303030300000a101524d524b3a3a4c4953543a3a322e302e303a3a31353837363832362d3334326631323130366561623664393034632d594f55444c454241434b47524f554e442d594f55444c454241434b47524f554e442d30303030313231393a3a3433323530303030303030303000007101524d524b3a3a4c4953543a3a322e302e303a3a31353838363037372d3334326631323130366561623664393034632d594f55444c4541524d2d594f55444c4541524d532d30303030303639393a3a313833333530303030303030303000005501524d524b3a3a4c4953543a3a322e302e303a3a31353837353238362d3334326631323130366561623664393034632d594f55444c452d594f55444c452d30303030303737303a3a3839303935303030303030303030",
		"ad0284004a7c11337fdf53dd1ac3cdbfa86f90a059684f0c01a00827ab00259d946d5f5b01aed0bd7afbf2ad7a926bbab1157d3de28aa60f9e268e0400ba1165b1af4f5b56daa7faf62cd8219a1644268fc6c52946c749d840779b0677d2fd220019f71d8fc500d902001e06004e6709193edd9ad1773ef48e549fe565ea7d6e75ee1dfc1f631dbca7106ee26fc17cbec65603a1e6d434f8e3651db8020f453b8fad25d98b7822b9ef55f77665",
	}

	for _, raw := range extrinsicRaw {
		e := scalecodec.ExtrinsicDecoder{}
		e.Init(scaleBytes.ScaleBytes{Data: utiles.HexToBytes(raw)}, &option)
		e.Process()
		encode, err := e.Value.(*scalecodec.GenericExtrinsic).Encode(&option)
		assert.NoError(t, err)
		assert.Equal(t, raw, encode)
	}

	// will raise extrinsic params length not match error
	genericExtrinsic := scalecodec.GenericExtrinsic{VersionInfo: "04", CallCode: "0200", Params: nil}
	_, err := genericExtrinsic.Encode(&option)
	assert.Error(t, err)
}
