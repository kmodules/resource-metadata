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
	"sync"
	"text/template"

	"kmodules.xyz/resource-metadata/apis/meta/v1alpha1"
	"kmodules.xyz/resource-metadata/hub"

	"github.com/pkg/errors"
	"gomodules.xyz/encoding/json"
	crd_api "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metatable "k8s.io/apimachinery/pkg/api/meta/table"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var pool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

type TableConvertor interface {
	ConvertToTable(ctx context.Context, object runtime.Object, tableOptions runtime.Object) (*v1alpha1.Table, error)
}

// New creates a new table convertor for the provided CRD column definition. If the printer definition cannot be parsed,
// error will be returned along with a default table convertor.
func New(fieldPath string, columns []v1alpha1.ResourceColumnDefinition) (TableConvertor, error) {
	c := &convertor{
		fieldPath: fieldPath,
	}
	err := c.init(columns)
	return c, err
}

func NewForList(fieldPath string, columns []v1alpha1.ResourceColumnDefinition) (TableConvertor, error) {
	c := &convertor{
		fieldPath: fieldPath,
	}
	err := c.init(filterColumns(columns, v1alpha1.List))
	return c, err
}

func NewForGVR(r *hub.Registry, kc client.Client, gvr schema.GroupVersionResource, priority v1alpha1.Priority) (TableConvertor, error) {
	rd, err := r.LoadByGVR(gvr)
	if err != nil {
		return nil, err
	}

	c := &convertor{}
	err = c.init(FilterColumnsWithDefaults(kc, gvr, rd.Spec.Columns, priority))
	return c, err
}

type convertor struct {
	fieldPath string
	headers   []v1alpha1.ResourceColumnDefinition
}

func filterColumns(columns []v1alpha1.ResourceColumnDefinition, priority v1alpha1.Priority) []v1alpha1.ResourceColumnDefinition {
	out := make([]v1alpha1.ResourceColumnDefinition, 0, len(columns))
	for _, col := range columns {
		if (col.Priority&int32(v1alpha1.Metadata)) == int32(v1alpha1.Metadata) ||
			(col.Priority&int32(priority)) == int32(priority) ||
			(priority == v1alpha1.List && col.Priority == 0) {
			out = append(out, col)
		}
	}
	return out
}

func FilterColumnsWithDefaults(
	kc client.Client,
	gvr schema.GroupVersionResource,
	columns []v1alpha1.ResourceColumnDefinition,
	priority v1alpha1.Priority,
) []v1alpha1.ResourceColumnDefinition {

	// columns are specified in resource description, so use those.
	out := filterColumns(columns, priority)
	if len(out) > 0 {
		return out
	}

	// generate column list by merging default columns + crd additional columns
	var defaultColumns []v1alpha1.ResourceColumnDefinition
	if priority == v1alpha1.List {
		defaultColumns = DefaultListColumns()
	} else {
		defaultColumns = DefaultDetailsColumns()
	}
	defaultJsonPaths := sets.NewString()
	for _, col := range defaultColumns {
		defaultJsonPaths.Insert(col.Name)
	}

	var additionalColumns []v1alpha1.ResourceColumnDefinition
	if kc != nil {
		var crd crd_api.CustomResourceDefinition
		err := kc.Get(context.TODO(), client.ObjectKey{Name: fmt.Sprintf("%s.%s", gvr.Resource, gvr.Group)}, &crd)
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
	c.headers = append(c.headers, columns...)
	return nil
}

func (c *convertor) rowFn(obj interface{}) ([]v1alpha1.TableCell, error) {
	data := obj
	if o, ok := obj.(runtime.Unstructured); ok {
		data = o.UnstructuredContent()
	}

	buf := pool.Get().(*bytes.Buffer)
	defer pool.Put(buf)

	cells := make([]v1alpha1.TableCell, 0, len(c.headers))
	for _, col := range c.headers {

		var cell v1alpha1.TableCell
		{
			if v, err := renderTemplate(data, columnOptions{
				Name:     col.Name,
				Type:     col.Type,
				Template: col.PathTemplate,
			}, buf); err != nil {
				return nil, err
			} else {
				cell.Data = v
			}
		}
		if col.Sort != nil && col.Sort.Enable && col.Sort.Template != "" {
			if v, err := renderTemplate(data, columnOptions{
				Name:     col.Name,
				Type:     col.Sort.Type,
				Template: col.Sort.Template,
			}, buf); err != nil {
				return nil, err
			} else {
				cell.Sort = v
			}
		}
		if col.Link != nil && col.Link.Enable && col.Link.Template != "" {
			if v, err := renderTemplate(data, columnOptions{
				Name:     col.Name,
				Type:     "string",
				Template: col.Link.Template,
			}, buf); err != nil {
				return nil, err
			} else {
				cell.Link = v.(string)
			}
		}
		if col.Icon != nil && col.Icon.Enable && col.Icon.Template != "" {
			if v, err := renderTemplate(data, columnOptions{
				Name:     col.Name,
				Type:     "string",
				Template: col.Icon.Template,
			}, buf); err != nil {
				return nil, err
			} else {
				cell.Icon = v.(string)
			}
		}
		if col.Color != nil && col.Color.Template != "" {
			if v, err := renderTemplate(data, columnOptions{
				Name:     col.Name,
				Type:     "string",
				Template: col.Color.Template,
			}, buf); err != nil {
				return nil, err
			} else {
				cell.Color = v.(string)
			}
		}

		cells = append(cells, cell)
	}
	return cells, nil
}

type columnOptions struct {
	Name     string
	Type     string
	Template string
}

func renderTemplate(data interface{}, col columnOptions, buf *bytes.Buffer) (interface{}, error) {
	if col.Template == "" {
		return nil, nil
	}

	tpl, err := template.New("").Funcs(templateFns).Parse(col.Template)
	if err != nil {
		klog.ErrorS(err, "failed to parse column template", "name", col.Name, "template", col.Template)
		return nil, errors.Wrapf(err, "falied to parse column %+v", col)
	}
	// Do nothing and continue execution.
	// If printed, the result of the index operation is the string "<no value>".
	// We mitigate that later.
	tpl.Option("missingkey=default")
	buf.Reset()
	err = tpl.Execute(buf, data)
	if err != nil {
		klog.ErrorS(err, "failed to render column template", "name", col.Name, "template", col.Template)
		return nil, errors.Wrapf(err, "falied to render column %+v", col)
	}
	return cellForJSONValue(col, strings.ReplaceAll(buf.String(), "<no value>", ""))
}

func (c *convertor) ConvertToTable(_ context.Context, obj runtime.Object, _ runtime.Object) (*v1alpha1.Table, error) {
	table := &v1alpha1.Table{
		Columns: make([]v1alpha1.ResourceColumn, 0, len(c.headers)),
		Rows:    make([]v1alpha1.TableRow, 0),
	}

	for _, def := range c.headers {
		table.Columns = append(table.Columns, v1alpha1.Convert_ResourceColumnDefinition_To_ResourceColumn(def))
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
	table.Rows, err = metaToTableRow(obj, c.fieldPath, c.rowFn)
	return table, err
}

func fields(path string) []string {
	return strings.Split(strings.Trim(path, "."), ".")
}

func cellForJSONValue(col columnOptions, value string) (interface{}, error) {
	value = strings.TrimSpace(value)
	switch col.Type {
	case "integer":
		if value == "" {
			return UnknownValue, nil
		}
		i64, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, err
		}
		return i64, nil
	case "number":
		if value == "" {
			return UnknownValue, nil
		}
		f64, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, err
		}
		return f64, nil
	case "boolean":
		if value == "" {
			return UnknownValue, nil
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
		if value == "" || value == "null" {
			return map[string]interface{}{}, nil
		}
		var obj interface{}
		err := json.Unmarshal([]byte(value), &obj)
		if err != nil {
			return nil, fmt.Errorf("col %s, type %s, err %v, value %s", col.Name, col.Type, err.Error(), value)
		}
		return obj, nil
	}
	return nil, fmt.Errorf("unknown type %s in header %s with value %s", col.Type, col.Name, value)
}

// metaToTableRow converts a list or object into one or more table rows. The provided rowFn is invoked for
// each accessed item, with name and age being passed to each.
func metaToTableRow(obj runtime.Object, fieldPath string, rowFn func(obj interface{}) ([]v1alpha1.TableCell, error)) ([]v1alpha1.TableRow, error) {
	if meta.IsListType(obj) {
		rows := make([]v1alpha1.TableRow, 0, 16)
		err := meta.EachListItem(obj, func(obj runtime.Object) error {
			nestedRows, err := metaToTableRow(obj, fieldPath, rowFn)
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

	if fieldPath == "" {
		// obj to row
		cells, err := rowFn(obj)
		if err != nil {
			return nil, err
		}
		return []v1alpha1.TableRow{
			{
				Cells: cells,
			},
		}, nil
	}

	// subtable
	arr, ok, err := unstructured.NestedSlice(obj.(runtime.Unstructured).UnstructuredContent(), fields(fieldPath)...)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}
	rows := make([]v1alpha1.TableRow, 0, len(arr))
	for _, item := range arr {
		var row v1alpha1.TableRow
		row.Cells, err = rowFn(item)
		if err != nil {
			return nil, err
		}
		rows = append(rows, row)
	}
	return rows, nil
}

func DefaultListColumns() []v1alpha1.ResourceColumnDefinition {
	return []v1alpha1.ResourceColumnDefinition{
		{
			Name:         "Name",
			Type:         "string",
			Format:       "",
			Priority:     int32(v1alpha1.List),
			PathTemplate: `{{ .metadata.name }}`,
			Sort: &v1alpha1.SortDefinition{
				Enable: true,
				// Template: "",
			},
			Link: &v1alpha1.AttributeDefinition{
				Enable: true,
				// Template: "",
			},
			//Shape ShapeProperty `json:"shape,omitempty"`
			//Icon  bool          `json:"icon,omitempty"`
			//Color ColorProperty `json:"color,omitempty"`
		},
		{
			Name:         "Namespace",
			Type:         "string",
			Format:       "",
			Priority:     int32(v1alpha1.List),
			PathTemplate: `{{ .metadata.namespace }}`,
		},
		{
			Name:         "Labels",
			Type:         "object",
			Format:       "",
			Priority:     int32(v1alpha1.List),
			PathTemplate: `{{ .metadata.labels | toRawJson }}`,
		},
		{
			Name:         "Annotations",
			Type:         "object",
			Format:       "",
			Priority:     int32(v1alpha1.List),
			PathTemplate: `{{ .metadata.annotations | toRawJson }}`,
		},
		{
			Name:         "Age",
			Type:         "date",
			Format:       "",
			Priority:     int32(v1alpha1.List),
			PathTemplate: `{{ .metadata.creationTimestamp }}`,
			Sort: &v1alpha1.SortDefinition{
				Enable:   true,
				Template: `{{ .metadata.creationTimestamp | toDate "2006-01-02T15:04:05Z07:00" | unixEpoch }}`,
				Type:     "integer",
			},
		},
	}
}

func DefaultDetailsColumns() []v1alpha1.ResourceColumnDefinition {
	return []v1alpha1.ResourceColumnDefinition{
		{
			Name:         "Name",
			Type:         "string",
			Format:       "",
			Priority:     int32(v1alpha1.Field | v1alpha1.List),
			PathTemplate: `{{ .metadata.name }}`,
			Sort: &v1alpha1.SortDefinition{
				Enable: true,
				// Template: "",
			},
			Link: &v1alpha1.AttributeDefinition{
				Enable: true,
				// Template: "",
			},
			//Shape ShapeProperty `json:"shape,omitempty"`
			//Icon  bool          `json:"icon,omitempty"`
			//Color ColorProperty `json:"color,omitempty"`
		},
		{
			Name:         "Namespace",
			Type:         "string",
			Format:       "",
			Priority:     int32(v1alpha1.Field | v1alpha1.List),
			PathTemplate: `{{ .metadata.namespace }}`,
		},
		{
			Name:         "Labels",
			Type:         "object",
			Format:       "",
			Priority:     int32(v1alpha1.Field | v1alpha1.List),
			PathTemplate: `{{ .metadata.labels | toRawJson }}`,
		},
		{
			Name:         "Annotations",
			Type:         "object",
			Format:       "",
			Priority:     int32(v1alpha1.Field | v1alpha1.List),
			PathTemplate: `{{ .metadata.annotations | toRawJson }}`,
		},
		{
			Name:         "Age",
			Type:         "date",
			Format:       "",
			Priority:     int32(v1alpha1.Field | v1alpha1.List),
			PathTemplate: `{{ .metadata.creationTimestamp }}`,
			Sort: &v1alpha1.SortDefinition{
				Enable:   true,
				Template: `{{ .metadata.creationTimestamp | toDate "2006-01-02T15:04:05Z07:00" | unixEpoch }}`,
				Type:     "integer",
			},
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
