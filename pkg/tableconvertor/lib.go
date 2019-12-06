/*
Copyright The Kmodules Authors.

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

	crd_cs "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func TableForList(r *hub.Registry, client crd_cs.CustomResourceDefinitionInterface, gvr schema.GroupVersionResource, items []unstructured.Unstructured) (*v1alpha1.Table, error) {
	c, err := NewForGVR(r, client, gvr, v1alpha1.List)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	obj := &unstructured.UnstructuredList{
		Items: items,
	}
	return c.ConvertToTable(ctx, obj, nil)
}

func TableForObject(r *hub.Registry, client crd_cs.CustomResourceDefinitionInterface, obj runtime.Object) (*v1alpha1.Table, error) {
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

	c, err := NewForGVR(r, client, gvr, v1alpha1.Field)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	table, err := c.ConvertToTable(ctx, obj, nil)
	if err != nil {
		return nil, err
	}

	for _, st := range rd.Spec.SubTables {
		c2, err := New(st.FieldPath, st.Columns)
		if err != nil {
			return nil, err
		}
		t2, err := c2.ConvertToTable(ctx, obj, nil)
		if err != nil {
			return nil, err
		}
		table.SubTables = append(table.SubTables, v1alpha1.SubTable{
			Name:              st.Name,
			ColumnDefinitions: t2.ColumnDefinitions,
			Rows:              t2.Rows,
		})
	}

	return table, nil
}
