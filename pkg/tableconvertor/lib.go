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

package tableconvertor

import (
	"context"

	kmapi "kmodules.xyz/client-go/api/v1"
	"kmodules.xyz/resource-metadata/apis/meta/v1alpha1"
	uiapi "kmodules.xyz/resource-metadata/apis/ui/v1alpha1"
	tabledefs "kmodules.xyz/resource-metadata/hub/resourcetabledefinitions"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type DashboardRendererFunc func(name string) (*uiapi.ResourceDashboard, string, error)

func NewForGVR(kc client.Client, gvr schema.GroupVersionResource, priority v1alpha1.Priority, fn DashboardRendererFunc) (TableConvertor, error) {
	var columns []v1alpha1.ResourceColumnDefinition
	if def, ok := tabledefs.LoadDefaultByGVR(gvr); ok {
		columns = def.Spec.Columns
	}

	var err error
	columns, err = tabledefs.FlattenColumns(columns)
	if err != nil {
		return nil, err
	}
	columns = FilterColumnsWithDefaults(kc, gvr, columns, priority)

	c := &convertor{}
	err = c.init(columns, fn)
	return c, err
}

func TableForAnyList(kc client.Client, items []unstructured.Unstructured, fn DashboardRendererFunc) (*v1alpha1.Table, error) {
	if len(items) == 0 {
		return &v1alpha1.Table{
			Rows: make([]v1alpha1.TableRow, 0),
		}, nil
	}

	gvk := items[0].GetObjectKind().GroupVersionKind()
	rid, err := kmapi.ExtractResourceID(kc.RESTMapper(), kmapi.ResourceID{
		Group:   gvk.Group,
		Version: gvk.Version,
		Name:    "",
		Kind:    gvk.Kind,
		Scope:   "",
	})
	if err != nil {
		return nil, err
	}

	return TableForList(kc, rid.GroupVersionResource(), items, fn)
}

func TableForList(kc client.Client, gvr schema.GroupVersionResource, items []unstructured.Unstructured, fn DashboardRendererFunc) (*v1alpha1.Table, error) {
	c, err := NewForGVR(kc, gvr, v1alpha1.List, fn)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	obj := &unstructured.UnstructuredList{
		Items: items,
	}
	return c.ConvertToTable(ctx, obj)
}

func TableForObject(kc client.Client, obj runtime.Object, fn DashboardRendererFunc) (*v1alpha1.Table, error) {
	gvk := obj.GetObjectKind().GroupVersionKind()
	rid, err := kmapi.ExtractResourceID(kc.RESTMapper(), kmapi.ResourceID{
		Group:   gvk.Group,
		Version: gvk.Version,
		Name:    "",
		Kind:    gvk.Kind,
		Scope:   "",
	})
	if err != nil {
		return nil, err
	}

	c, err := NewForGVR(kc, rid.GroupVersionResource(), v1alpha1.Field, fn)
	if err != nil {
		return nil, err
	}

	return c.ConvertToTable(context.TODO(), obj)
}
