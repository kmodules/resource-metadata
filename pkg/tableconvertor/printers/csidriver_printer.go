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
	"fmt"
	"reflect"
	"strings"

	storage "k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func init() {
	Register(CSIDriverPrinter{})
}

// ref: https://github.com/kubernetes/kubernetes/blob/v1.21.0/pkg/printers/internalversion/printers.go#L513-L534

type CSIDriverPrinter struct{}

var _ ColumnConverter = CSIDriverPrinter{}

func (_ CSIDriverPrinter) GVK() schema.GroupVersionKind {
	return storage.SchemeGroupVersion.WithKind("CSIDriver")
}

func (p CSIDriverPrinter) Convert(o runtime.Object) (map[string]interface{}, error) {
	obj, ok := o.(*storage.CSIDriver)
	if !ok {
		return nil, fmt.Errorf("expected %v, received %v", p.GVK().Kind, reflect.TypeOf(o))
	}

	row := map[string]interface{}{}

	attachRequired := true
	if obj.Spec.AttachRequired != nil {
		attachRequired = *obj.Spec.AttachRequired
	}
	podInfoOnMount := false
	if obj.Spec.PodInfoOnMount != nil {
		podInfoOnMount = *obj.Spec.PodInfoOnMount
	}
	allModes := []string{}
	for _, mode := range obj.Spec.VolumeLifecycleModes {
		allModes = append(allModes, string(mode))
	}
	modes := strings.Join(allModes, ",")
	if len(modes) == 0 {
		modes = None
	}

	row["Name"] = obj.Name
	row["AttachRequired"] = attachRequired
	row["PodInfoOnMount"] = podInfoOnMount

	/*
		storageCapacity := false
		if obj.Spec.StorageCapacity != nil {
			storageCapacity = *obj.Spec.StorageCapacity
		}
		row["StorageCapacity"] = storageCapacity

		tokenRequests := "<unset>"
		if obj.Spec.TokenRequests != nil {
			audiences := []string{}
			for _, t := range obj.Spec.TokenRequests {
				audiences = append(audiences, t.Audience)
			}
			tokenRequests = strings.Join(audiences, ",")
		}
		requiresRepublish := false
		if obj.Spec.RequiresRepublish != nil {
			requiresRepublish = *obj.Spec.RequiresRepublish
		}

		row["TokenRequests"] = tokenRequests
		row["RequiresRepublish"] = requiresRepublish
	*/
	row["Modes"] = modes
	row["Age"] = translateTimestampSince(obj.CreationTimestamp)

	return row, nil
}
