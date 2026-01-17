package json

import (
	"testing"

	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
	"matheusd.com/depvendoredtestify/require"
	"matheusd.com/mdcapnp/capnpser"
)

// TestCapWithComplexMethodJSON is a basic test for encoding types as JSON.
func TestCapWithComplexMethodJSON(t *testing.T) {
	mb, err := capnpser.NewMessageBuilder(capnpser.DefaultSimpleSingleAllocator)
	require.NoError(t, err)

	in := `{   "strValue": "i am the walrus", "int64Value": 555999  }`

	b, err := mb.NewRootStruct(testCapWithComplexMethod_ComplexMethod_ParamsSize)
	require.NoError(t, err)
	pb := testCapWithComplexMethod_ComplexMethod_ParamsBuilder(b)

	err = json.Unmarshal([]byte(in), &pb)
	require.NoError(t, err)

	// pb.SetInt64Value(555999)
	// pb.SetStrValue("i am the walrus")

	reader := mb.MessageReader()
	readerRoot, err := reader.GetRoot()
	require.NoError(t, err)
	pr := testCapWithComplexMethod_ComplexMethod_Params(readerRoot)
	out, err := json.Marshal(pr, jsontext.Multiline(true))
	require.NoError(t, err)

	t.Log(string(out))
}
