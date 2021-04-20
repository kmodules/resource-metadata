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

	networking "k8s.io/api/networking/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func init() {
	Register(IngressClassPrinter{})
}

// ref: https://github.com/kubernetes/kubernetes/blob/v1.21.0/pkg/printers/internalversion/printers.go#L214-L221

type IngressClassPrinter struct{}

var _ ColumnConverter = IngressClassPrinter{}

func (_ IngressClassPrinter) GVK() schema.GroupVersionKind {
	return networking.SchemeGroupVersion.WithKind("IngressClass")
}

func (p IngressClassPrinter) Convert(o runtime.Object) (map[string]interface{}, error) {
	obj, ok := o.(*networking.IngressClass)
	if !ok {
		return nil, fmt.Errorf("expected %v, received %v", p.GVK().Kind, reflect.TypeOf(o))
	}

	row := map[string]interface{}{}

	parameters := "<none>"
	if obj.Spec.Parameters != nil {
		parameters = obj.Spec.Parameters.Kind
		if obj.Spec.Parameters.APIGroup != nil {
			parameters = parameters + "." + *obj.Spec.Parameters.APIGroup
		}
		parameters = parameters + "/" + obj.Spec.Parameters.Name
	}
	createTime := translateTimestampSince(obj.CreationTimestamp)

	row["Name"] = obj.Name
	row["Controller"] = obj.Spec.Controller
	row["Parameters"] = parameters
	row["Age"] = createTime

	return row, nil
}
