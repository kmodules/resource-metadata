package hub

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegistry_LoadDefaultResourceClass(t *testing.T) {
	reg := NewRegistry("some-uid", NewKVLocal())
	rcs, err := reg.LoadDefaultResourceClassList()
	assert.NoError(t, err)

	for _, rc := range rcs.Items {
		for _, entry := range rc.Spec.Entries {
			println(rc.Name)
			fmt.Println(entry)
		}
	}
}
