package api

import (
	"fmt"
	"sync"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	plugins = map[schema.GroupVersionKind]ResourceCalculator{}
	lock    sync.RWMutex
)

func Register(gvk schema.GroupVersionKind, c ResourceCalculator) {
	lock.Lock()
	plugins[gvk] = c
	lock.Unlock()
}

func Load(obj map[string]interface{}) (ResourceCalculator, error) {
	u := unstructured.Unstructured{Object: obj}
	gvk := u.GroupVersionKind()

	lock.RLock()
	c, ok := plugins[gvk]
	lock.RUnlock()
	if !ok {
		return nil, NotRegistered{gvk}
	}
	return c, nil
}

type NotRegistered struct {
	gvk schema.GroupVersionKind
}

var _ error = NotRegistered{}

func (e NotRegistered) Error() string {
	return fmt.Sprintf("no calculator registered for %v", e.gvk)
}
