/*
Copyright AppsCode Inc. and Contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package printers

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	// on the node it is (was) running.
	NodeUnreachablePodReason = "NodeLost"
)

type ColumnConverter interface {
	GVK() schema.GroupVersionKind
	Convert(obj runtime.Object) (map[string]interface{}, error)
}

var printers = map[schema.GroupVersionKind]ColumnConverter{}

func Register(c ColumnConverter) {
	printers[c.GVK()] = c
}

func Convert(o runtime.Object) (map[string]interface{}, error) {
	gvk := o.GetObjectKind().GroupVersionKind()
	c, ok := printers[gvk]
	if !ok {
		return map[string]interface{}{}, nil
	}
	return c.Convert(o)
}
