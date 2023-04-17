package types

import (
	"fmt"
	"sync"
	"testing"

	"github.com/itering/scale.go/source"
)

func TestRegCustomTypesConcurrency(t *testing.T) {
	wg := sync.WaitGroup{}
	count := 0
	for {
		count++
		go func() {
			RegCustomTypes(map[string]source.TypeStruct{fmt.Sprintf("%d", count): {Type: "string", TypeString: "u32"}})
			wg.Add(1)
			wg.Done()
		}()
		if count > 100 {
			break
		}
	}

	wg.Wait()
}
