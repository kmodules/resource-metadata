/*
Copyright 2018 The Kubernetes Authors.

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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"k8s.io/apimachinery/pkg/api/meta"
	metatable "k8s.io/apimachinery/pkg/api/meta/table"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/util/jsonpath"
	"kmodules.xyz/resource-metadata/apis/meta/v1alpha1"
	hub "kmodules.xyz/resource-metadata/hub/v1alpha1"
)

type TableConvertor interface {
	ConvertToTable(ctx context.Context, object runtime.Object, tableOptions runtime.Object) (*v1alpha1.Table, error)
}

// New creates a new table convertor for the provided CRD column definition. If the printer definition cannot be parsed,
// error will be returned along with a default table convertor.
func New(fieldPath string, columns []v1alpha1.ResourceColumnDefinition) (TableConvertor, error) {
	c := &convertor{
		fieldPath: fieldPath,
		buf:       &bytes.Buffer{},
	}
	err := c.init(columns, v1alpha1.List)
	return c, err
}

func NewForGVR(gvr schema.GroupVersionResource, priority v1alpha1.Priority) (TableConvertor, error) {
	rd, err := hub.LoadByGVR(gvr)
	if err != nil {
		return nil, err
	}

	c := &convertor{
		buf: &bytes.Buffer{},
	}
	err = c.init(rd.Spec.Columns, priority)
	return c, err
}

type convertor struct {
	buf       *bytes.Buffer
	fieldPath string
	headers   []v1alpha1.ResourceColumnDefinition
	columns   []*jsonpath.JSONPath
}

func (c *convertor) init(columns []v1alpha1.ResourceColumnDefinition, priority v1alpha1.Priority) error {
	for _, col := range columns {
		if (col.Priority&int32(priority)) == int32(priority) ||
			(priority == v1alpha1.List && col.Priority == 0) {
			path := jsonpath.New(col.Name)

			col.JSONPath = strings.TrimSpace(col.JSONPath)
			if !strings.HasPrefix(col.JSONPath, "{") {
				col.JSONPath = fmt.Sprintf("{%s}", col.JSONPath)
			}
			if err := path.Parse(col.JSONPath); err != nil {
				return fmt.Errorf("unrecognized column definition %q", col.JSONPath)
			}
			path.AllowMissingKeys(true)

			desc := fmt.Sprintf("Custom resource definition column (in JSONPath format): %s", col.JSONPath)
			if len(col.Description) > 0 {
				desc = col.Description
			}

			c.columns = append(c.columns, path)
			c.headers = append(c.headers, v1alpha1.ResourceColumnDefinition{
				Name:        col.Name,
				Type:        col.Type,
				Format:      col.Format,
				Description: desc,
				Priority:    col.Priority,
			})
		}
	}
	return nil
}

func (c *convertor) rowFn(data interface{}) ([]interface{}, error) {
	cells := make([]interface{}, 0, len(c.columns))
	for i, column := range c.columns {
		results, err := column.FindResults(data)
		if err != nil || len(results) == 0 || len(results[0]) == 0 {
			cells = append(cells, nil)
			continue
		}

		// as we only support simple JSON path, we can assume to have only one result (or none, filtered out above)
		value := results[0][0].Interface()
		if c.headers[i].Type == "string" {
			if err := column.PrintResults(c.buf, []reflect.Value{reflect.ValueOf(value)}); err == nil {
				cells = append(cells, c.buf.String())
				c.buf.Reset()
			} else {
				cells = append(cells, nil)
			}
		} else {
			cells = append(cells, cellForJSONValue(c.headers[i].Type, value))
		}
	}
	return cells, nil
}

func (c *convertor) ConvertToTable(ctx context.Context, obj runtime.Object, tableOptions runtime.Object) (*v1alpha1.Table, error) {
	table := &v1alpha1.Table{
		ColumnDefinitions: c.headers,
	}
	if m, err := meta.ListAccessor(obj); err == nil {
		table.ResourceVersion = m.GetResourceVersion()
		table.SelfLink = m.GetSelfLink()
		table.Continue = m.GetContinue()
	} else {
		if m, err := meta.CommonAccessor(obj); err == nil {
			table.ResourceVersion = m.GetResourceVersion()
			table.SelfLink = m.GetSelfLink()
		}
	}

	var err error

	if c.fieldPath == "" {
		table.Rows, err = metaToTableRow(obj, c.rowFn)
	} else {
		arr, ok, err := unstructured.NestedSlice(obj.(runtime.Unstructured).UnstructuredContent(), strings.Split(c.fieldPath, ".")...)
		if err != nil {
			return nil, err
		}
		if !ok {
			return table, nil
		}

		rows := make([]v1alpha1.TableRow, 0, len(arr))
		for _, item := range arr {
			row := v1alpha1.TableRow{
				// Object: runtime.RawExtension{Object: obj},
			}
			row.Cells, err = c.rowFn(item)
			if err != nil {
				return nil, err
			}
			rows = append(rows, row)
		}
		table.Rows = rows
	}

	return table, err
}

func cellForJSONValue(headerType string, value interface{}) interface{} {
	if value == nil {
		return nil
	}

	switch headerType {
	case "integer":
		switch typed := value.(type) {
		case int64:
			return typed
		case float64:
			return int64(typed)
		case json.Number:
			if i64, err := typed.Int64(); err == nil {
				return i64
			}
		}
	case "number":
		switch typed := value.(type) {
		case int64:
			return float64(typed)
		case float64:
			return typed
		case json.Number:
			if f, err := typed.Float64(); err == nil {
				return f
			}
		}
	case "boolean":
		if b, ok := value.(bool); ok {
			return b
		}
	case "string":
		if s, ok := value.(string); ok {
			return s
		}
	case "date":
		if typed, ok := value.(string); ok {
			var timestamp metav1.Time
			err := timestamp.UnmarshalQueryParameter(typed)
			if err != nil {
				return "<invalid>"
			}
			return metatable.ConvertToHumanReadableDateType(timestamp)
		}
		// TODO: Fix things
	case "object":
		return value
	}

	return nil
}

// metaToTableRow converts a list or object into one or more table rows. The provided rowFn is invoked for
// each accessed item, with name and age being passed to each.
func metaToTableRow(obj runtime.Object, rowFn func(obj interface{}) ([]interface{}, error)) ([]v1alpha1.TableRow, error) {
	if meta.IsListType(obj) {
		rows := make([]v1alpha1.TableRow, 0, 16)
		err := meta.EachListItem(obj, func(obj runtime.Object) error {
			nestedRows, err := metaToTableRow(obj, rowFn)
			if err != nil {
				return err
			}
			rows = append(rows, nestedRows...)
			return nil
		})
		if err != nil {
			return nil, err
		}
		return rows, nil
	}

	rows := make([]v1alpha1.TableRow, 0, 1)
	row := v1alpha1.TableRow{}
	var err error
	row.Cells, err = rowFn(obj.(runtime.Unstructured).UnstructuredContent())
	if err != nil {
		return nil, err
	}
	rows = append(rows, row)
	return rows, nil
}
