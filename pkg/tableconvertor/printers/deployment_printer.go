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

	apps "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func init() {
	Register(DeploymentPrinter{})
}

// ref: https://github.com/kubernetes/kubernetes/blob/v1.21.0/pkg/printers/internalversion/printers.go#L135-L146

type DeploymentPrinter struct{}

var _ ColumnConverter = DeploymentPrinter{}

func (_ DeploymentPrinter) GVK() schema.GroupVersionKind {
	return apps.SchemeGroupVersion.WithKind("Deployment")
}

func (p DeploymentPrinter) Convert(o runtime.Object) (map[string]interface{}, error) {
	obj := new(apps.Deployment)
	switch to := o.(type) {
	case *unstructured.Unstructured:
		if err := runtime.DefaultUnstructuredConverter.FromUnstructured(to.UnstructuredContent(), obj); err != nil {
			return nil, err
		}
	case *apps.Deployment:
		obj = to
	default:
		return nil, fmt.Errorf("expected %v, received %v", p.GVK().Kind, reflect.TypeOf(o))
	}

	row := map[string]interface{}{}

	names, images := layoutContainerCells(obj.Spec.Template.Spec.Containers)
	row["_Containers"] = names
	row["Images"] = images
	row["_Selector"] = metav1.FormatLabelSelector(obj.Spec.Selector)

	return row, nil
}
