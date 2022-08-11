package types

import (
	"fmt"
	"strings"

	"github.com/itering/scale.go/utiles"
)

type BTreeMap struct{ ScaleDecoder }

func (b *BTreeMap) Process() {
	elementCount := b.ProcessAndUpdateData("Compact<u32>").(int)
	var result []interface{}
	if elementCount > 1000 {
		panic(fmt.Sprintf("BTreeMap length %d exceeds %d", elementCount, 1000))
	}
	for i := 0; i < elementCount; i++ {
		subType := strings.Split(b.SubType, ",")
		key := utiles.ToString(b.ProcessAndUpdateData(subType[0]))
		result = append(result, map[string]interface{}{
			key: b.ProcessAndUpdateData(subType[1]),
		})
	}
	b.Value = result
}
