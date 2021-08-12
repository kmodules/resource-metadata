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
	"bytes"
	"context"
	"fmt"
	"strconv"
	"strings"
	"text/template"

	"kmodules.xyz/resource-metadata/apis/meta/v1alpha1"
	"kmodules.xyz/resource-metadata/hub"
	"kmodules.xyz/resource-metadata/pkg/tableconvertor/printers"

	"gomodules.xyz/encoding/json"
	crd_cs "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metatable "k8s.io/apimachinery/pkg/api/meta/table"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/klog/v2"
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
	err := c.init(filterColumns(columns, v1alpha1.List))
	return c, err
}

func NewForGVR(r *hub.Registry, client crd_cs.CustomResourceDefinitionInterface, gvr schema.GroupVersionResource, priority v1alpha1.Priority) (TableConvertor, error) {
	rd, err := r.LoadByGVR(gvr)
	if err != nil {
		return nil, err
	}

	c := &convertor{
		buf: &bytes.Buffer{},
	}
	err = c.init(filterColumnsWithDefaults(client, gvr, rd.Spec.Columns, priority))
	return c, err
}

type convertor struct {
	buf       *bytes.Buffer
	fieldPath string
	headers   []v1alpha1.ResourceColumnDefinition
}

func filterColumns(columns []v1alpha1.ResourceColumnDefinition, priority v1alpha1.Priority) []v1alpha1.ResourceColumnDefinition {
	out := make([]v1alpha1.ResourceColumnDefinition, 0, len(columns))
	for _, col := range columns {
		if (col.Priority&int32(priority)) == int32(priority) ||
			(priority == v1alpha1.List && col.Priority == 0) {
			out = append(out, col)
		}
	}
	return out
}

func filterColumnsWithDefaults(client crd_cs.CustomResourceDefinitionInterface, gvr schema.GroupVersionResource, columns []v1alpha1.ResourceColumnDefinition, priority v1alpha1.Priority) []v1alpha1.ResourceColumnDefinition {
	// columns are specified in resource description, so use those.
	out := filterColumns(columns, priority)
	if len(out) > 0 {
		return out
	}

	// generate column list by merging default columns + crd additional columns
	var defaultColumns []v1alpha1.ResourceColumnDefinition
	if priority == v1alpha1.List {
		defaultColumns = defaultListColumns()
	} else {
		defaultColumns = defaultDetailsColumns()
	}
	defaultJsonPaths := sets.NewString()
	for _, col := range defaultColumns {
		defaultJsonPaths.Insert(col.Name)
	}

	var additionalColumns []v1alpha1.ResourceColumnDefinition
	if client != nil {
		crd, err := client.Get(context.TODO(), fmt.Sprintf("%s.%s", gvr.Resource, gvr.Group), metav1.GetOptions{})
		if err == nil {
			for _, version := range crd.Spec.Versions {
				if version.Name == gvr.Version && len(version.AdditionalPrinterColumns) > 0 {
					additionalColumns = make([]v1alpha1.ResourceColumnDefinition, 0, len(version.AdditionalPrinterColumns))
					for _, col := range version.AdditionalPrinterColumns {
						if !defaultJsonPaths.Has(col.Name) {
							def := v1alpha1.ResourceColumnDefinition{
								Name:        col.Name,
								Type:        col.Type,
								Format:      col.Format,
								Description: col.Description,
								Priority:    col.Priority,
							}
							col.JSONPath = strings.TrimSpace(col.JSONPath)
							if col.JSONPath != "" {
								def.PathTemplate = fmt.Sprintf(`{{ jp "{%s}" . }}`, col.JSONPath)
							}
							additionalColumns = append(additionalColumns, def)
						}
					}
				}
			}
		}
	}

	return append(defaultColumns, additionalColumns...)
}

func (c *convertor) init(columns []v1alpha1.ResourceColumnDefinition) error {
	for _, col := range columns {

		//desc := fmt.Sprintf("Custom resource definition column (in JSONPath format): %s", col.JSONPath)
		//if len(col.Description) > 0 {
		//	desc = col.Description
		//}

		c.headers = append(c.headers, v1alpha1.ResourceColumnDefinition{
			Name:         col.Name,
			Type:         col.Type,
			Format:       col.Format,
			Description:  col.Description,
			Priority:     col.Priority,
			PathTemplate: col.PathTemplate,
		})
	}
	return nil
}

func (c *convertor) rowFn(data interface{}) ([]interface{}, error) {
	knownCells := map[string]interface{}{}

	if obj, ok := data.(runtime.Unstructured); ok {
		var err error
		knownCells, err = printers.Convert(obj)
		if err != nil {
			return nil, err
		}
		data = obj.UnstructuredContent()
	}

	cells := make([]interface{}, 0, len(c.headers))
	for _, col := range c.headers {
		if v, ok := knownCells[col.Name]; ok {
			cells = append(cells, v)
			continue
		}

		if col.PathTemplate == "" {
			cells = append(cells, nil)
			continue
		}

		tpl, err := template.New("").Funcs(templateFns).Parse(col.PathTemplate)
		if err != nil {
			klog.Infof("Failed to parse. Reason: %v", err)
			return nil, fmt.Errorf("invalid column definition %q", col.PathTemplate)
		}
		err = tpl.Execute(c.buf, data)
		if err != nil {
			klog.Infof("Failed to resolve template. Reason: %v", err)
			return nil, fmt.Errorf("invalid column definition %q", col.PathTemplate)
		}

		v, err := cellForJSONValue(col.Name, col.Type, c.buf.String())
		if err != nil {
			return nil, err
		}
		cells = append(cells, v)
		c.buf.Reset()
	}
	return cells, nil
}

func (c *convertor) ConvertToTable(ctx context.Context, obj runtime.Object, tableOptions runtime.Object) (*v1alpha1.Table, error) {
	table := &v1alpha1.Table{
		ColumnDefinitions: c.headers,
	}
	if m, err := meta.ListAccessor(obj); err == nil {
		table.ResourceVersion = m.GetResourceVersion()
		table.Continue = m.GetContinue()
	} else {
		if m, err := meta.CommonAccessor(obj); err == nil {
			table.ResourceVersion = m.GetResourceVersion()
		}
	}

	var err error

	if c.fieldPath == "" {
		table.Rows, err = metaToTableRow(obj, c.rowFn)
	} else {
		arr, ok, err := unstructured.NestedSlice(obj.(runtime.Unstructured).UnstructuredContent(), fields(c.fieldPath)...)
		if err != nil {
			return nil, err
		}
		if !ok {
			return table, nil
		}

		rows := make([]v1alpha1.TableRow, 0, len(arr))
		for _, item := range arr {
			var row v1alpha1.TableRow
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

func fields(path string) []string {
	return strings.Split(strings.Trim(path, "."), ".")
}

func cellForJSONValue(colName, headerType string, value string) (interface{}, error) {
	switch headerType {
	case "integer":
		if value == "" {
			return "<unknown>", nil
		}
		i64, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, err
		}
		return i64, nil
	case "number":
		if value == "" {
			return "<unknown>", nil
		}
		f64, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, err
		}
		return f64, nil
	case "boolean":
		if value == "" {
			return "<unset>", nil
		}
		b, err := strconv.ParseBool(value)
		if err != nil {
			return nil, err
		}
		return b, nil
	case "string":
		return value, nil
	case "date":
		var timestamp metav1.Time
		err := timestamp.UnmarshalQueryParameter(value)
		if err != nil {
			return "<invalid>", nil
		}
		return metatable.ConvertToHumanReadableDateType(timestamp), nil
	case "object":
		if strings.TrimSpace(value) == "" {
			return map[string]interface{}{}, nil
		}
		var obj interface{}
		err := json.Unmarshal([]byte(value), &obj)
		if err != nil {
			return nil, fmt.Errorf("col %s, type %s, err %v, value %s", colName, headerType, err.Error(), value)
		}
		return obj, nil
	}
	return nil, fmt.Errorf("unknown type %s in header %s with value %s", headerType, colName, value)
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
	var row v1alpha1.TableRow
	var err error
	row.Cells, err = rowFn(obj)
	if err != nil {
		return nil, err
	}
	rows = append(rows, row)
	return rows, nil
}

func defaultListColumns() []v1alpha1.ResourceColumnDefinition {
	return []v1alpha1.ResourceColumnDefinition{
		{
			Name:         "Name",
			Type:         "string",
			Format:       "",
			Priority:     int32(v1alpha1.List),
			PathTemplate: `{{ jp "{.metadata.name}" . }}`,
		},
		{
			Name:         "Namespace",
			Type:         "string",
			Format:       "",
			Priority:     int32(v1alpha1.List),
			PathTemplate: `{{ jp "{.metadata.namespace}" . }}`,
		},
		{
			Name:         "Labels",
			Type:         "object",
			Format:       "",
			Priority:     int32(v1alpha1.List),
			PathTemplate: `{{ jp "{.metadata.labels}" . }}`,
		},
		{
			Name:         "Annotations",
			Type:         "object",
			Format:       "",
			Priority:     int32(v1alpha1.List),
			PathTemplate: `{{ jp "{.metadata.annotations}" . }}`,
		},
		{
			Name:         "Age",
			Type:         "date",
			Format:       "",
			Priority:     int32(v1alpha1.List),
			PathTemplate: `{{ jp "{.metadata.creationTimestamp}" . }}`,
		},
	}
}

func defaultDetailsColumns() []v1alpha1.ResourceColumnDefinition {
	return []v1alpha1.ResourceColumnDefinition{
		{
			Name:         "Name",
			Type:         "string",
			Format:       "",
			Priority:     int32(v1alpha1.Field),
			PathTemplate: `{{ jp "{.metadata.name}" . }}`,
		},
		{
			Name:         "Namespace",
			Type:         "string",
			Format:       "",
			Priority:     int32(v1alpha1.Field),
			PathTemplate: `{{ jp "{.metadata.namespace}" . }}`,
		},
		{
			Name:         "Labels",
			Type:         "object",
			Format:       "",
			Priority:     int32(v1alpha1.Field),
			PathTemplate: `{{ jp "{.metadata.labels}" . }}`,
		},
		{
			Name:         "Annotations",
			Type:         "object",
			Format:       "",
			Priority:     int32(v1alpha1.Field),
			PathTemplate: `{{ jp "{.metadata.annotations}" . }}`,
		},
		{
			Name:         "Age",
			Type:         "date",
			Format:       "",
			Priority:     int32(v1alpha1.Field),
			PathTemplate: `{{ jp "{.metadata.creationTimestamp}" . }}`,
		},
		/*
			{
				Name:     "Selector",
				Type:     "object",
				Format:   "selector",
				Priority: int32(v1alpha1.Field),
				JSONPath: ".spec.selector",
			},
			{
				Name:     "Desired Replicas",
				Type:     "integer",
				Format:   "",
				Priority: int32(v1alpha1.Field),
				JSONPath: ".spec.replicas",
			},
		*/
	}
}
