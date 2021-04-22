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

	"gomodules.xyz/pointer"
	apps "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func init() {
	Register(ReplicaSetPrinter{})
}

// ref: https://github.com/kubernetes/kubernetes/blob/v1.21.0/pkg/printers/internalversion/printers.go#L135-L146

type ReplicaSetPrinter struct{}

var _ ColumnConverter = ReplicaSetPrinter{}

func (_ ReplicaSetPrinter) GVK() schema.GroupVersionKind {
	return apps.SchemeGroupVersion.WithKind("ReplicaSet")
}

func (p ReplicaSetPrinter) Convert(o runtime.Object) (map[string]interface{}, error) {
	obj := new(apps.ReplicaSet)
	switch to := o.(type) {
	case *unstructured.Unstructured:
		if err := runtime.DefaultUnstructuredConverter.FromUnstructured(to.UnstructuredContent(), obj); err != nil {
			return nil, err
		}
	case *apps.ReplicaSet:
		obj = to
	default:
		return nil, fmt.Errorf("expected %v, received %v", p.GVK().Kind, reflect.TypeOf(o))
	}

	row := map[string]interface{}{}

	desiredReplicas := obj.Spec.Replicas
	currentReplicas := obj.Status.Replicas
	readyReplicas := obj.Status.ReadyReplicas

	row["_Name"] = obj.Name
	row["_Desired"] = int64(pointer.Int32(desiredReplicas))
	row["_Current"] = int64(currentReplicas)
	row["_Ready"] = int64(readyReplicas)
	row["_Age"] = translateTimestampSince(obj.CreationTimestamp)

	names, images := layoutContainerCells(obj.Spec.Template.Spec.Containers)
	row["_Containers"] = names
	row["_Images"] = images
	row["_Selector"] = metav1.FormatLabelSelector(obj.Spec.Selector)

	return row, nil
}
