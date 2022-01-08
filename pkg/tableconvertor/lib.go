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

	"kmodules.xyz/resource-metadata/apis/meta/v1alpha1"
	"kmodules.xyz/resource-metadata/hub"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Deprecated
func TableForList(r *hub.Registry, kc client.Client, gvr schema.GroupVersionResource, items []unstructured.Unstructured) (*v1alpha1.Table, error) {
	c, err := NewForGVR(r, kc, gvr, v1alpha1.List)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	obj := &unstructured.UnstructuredList{
		Items: items,
	}
	return c.ConvertToTable(ctx, obj, nil)
}

// Deprecated
func TableForObject(r *hub.Registry, kc client.Client, obj runtime.Object) (*v1alpha1.Table, error) {
	gvk := obj.GetObjectKind().GroupVersionKind()
	t := metav1.TypeMeta{APIVersion: gvk.GroupVersion().String(), Kind: gvk.Kind}
	gvr, err := r.GVR(t.GroupVersionKind())
	if err != nil {
		return nil, err
	}

	rd, err := r.LoadByGVR(gvr)
	if err != nil {
		return nil, err
	}

	c, err := NewForGVR(r, kc, gvr, v1alpha1.Field)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	table, err := c.ConvertToTable(ctx, obj, nil)
	if err != nil {
		return nil, err
	}

	for _, st := range rd.Spec.SubTables {
		c2, err := New(st.FieldPath, st.ColumnDefinitions)
		if err != nil {
			return nil, err
		}
		t2, err := c2.ConvertToTable(ctx, obj, nil)
		if err != nil {
			return nil, err
		}
		table.SubTables = append(table.SubTables, v1alpha1.SubTable{
			Name:    st.Name,
			Columns: t2.Columns,
			Rows:    t2.Rows,
		})
	}

	return table, nil
}
