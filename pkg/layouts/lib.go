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

package layouts

import (
	"fmt"
	"strings"

	kmapi "kmodules.xyz/client-go/api/v1"
	"kmodules.xyz/resource-metadata/apis/meta/v1alpha1"
	"kmodules.xyz/resource-metadata/hub"
	blockdefs "kmodules.xyz/resource-metadata/hub/resourceblockdefinitions"
	"kmodules.xyz/resource-metadata/hub/resourceoutlines"
	tabledefs "kmodules.xyz/resource-metadata/hub/resourcetabledefinitions"
	"kmodules.xyz/resource-metadata/pkg/tableconvertor"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const BasicPage = "Basic"

var reg = hub.NewRegistryOfKnownResources()

func LoadResourceLayoutForGVR(kc client.Client, gvr schema.GroupVersionResource) (*v1alpha1.ResourceLayout, error) {
	outline, found := resourceoutlines.LoadDefaultByGVR(gvr)
	if found {
		return GetResourceLayout(kc, outline)
	}

	mapper := kc.RESTMapper()
	rid, err := resourceIDForGVR(mapper, gvr)
	if err != nil {
		return nil, err
	}
	return generateDefaultLayout(kc, *rid)
}

func resourceIDForGVR(mapper meta.RESTMapper, gvr schema.GroupVersionResource) (*kmapi.ResourceID, error) {
	rid, err := kmapi.ExtractResourceID(mapper, kmapi.ResourceID{
		Group:   gvr.Group,
		Version: gvr.Version,
		Name:    gvr.Resource,
		Kind:    "",
		Scope:   "",
	})
	if err != nil {
		rid, err = reg.ResourceIDForGVR(gvr)
		if err != nil {
			return nil, err
		}
		if rid == nil {
			return nil, apierrors.NewNotFound(v1alpha1.Resource(v1alpha1.ResourceKindResourceOutline), gvr.String())
		}
	}
	return rid, nil
}

func LoadResourceLayoutForGVK(kc client.Client, gvk schema.GroupVersionKind) (*v1alpha1.ResourceLayout, error) {
	outline, found := resourceoutlines.LoadDefaultByGVK(gvk)
	if found {
		return GetResourceLayout(kc, outline)
	}

	rid, err := resourceIDForGVK(kc, gvk)
	if err != nil {
		return nil, err
	}
	return generateDefaultLayout(kc, *rid)
}

func resourceIDForGVK(kc client.Client, gvk schema.GroupVersionKind) (*kmapi.ResourceID, error) {
	mapper := kc.RESTMapper()
	rid, err := kmapi.ExtractResourceID(mapper, kmapi.ResourceID{
		Group:   gvk.Group,
		Version: gvk.Version,
		Name:    "",
		Kind:    gvk.Kind,
		Scope:   "",
	})
	if err != nil {
		rid, err = reg.ResourceIDForGVK(gvk)
		if err != nil {
			return nil, err
		}
		if rid == nil {
			return nil, apierrors.NewNotFound(v1alpha1.Resource(v1alpha1.ResourceKindResourceOutline), gvk.String())
		}
	}
	return rid, nil
}

func generateDefaultLayout(kc client.Client, rid kmapi.ResourceID) (*v1alpha1.ResourceLayout, error) {
	outline := &v1alpha1.ResourceOutline{
		TypeMeta: metav1.TypeMeta{
			Kind:       v1alpha1.ResourceKindResourceOutline,
			APIVersion: v1alpha1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: resourceoutlines.DefaultLayoutName(rid.GroupVersionResource()),
			Labels: map[string]string{
				"k8s.io/group":    rid.Group,
				"k8s.io/version":  rid.Version,
				"k8s.io/resource": rid.Name,
				"k8s.io/kind":     rid.Kind,
			},
		},
		Spec: v1alpha1.ResourceOutlineSpec{
			Resource:      rid,
			DefaultLayout: true,
			Header:        nil,
			TabBar:        nil,
			//Pages: []v1alpha1.ResourcePageOutline{
			//	{
			//		Name: "Basic",
			//		Info: &v1alpha1.PageBlockOutline{
			//			Kind:        v1alpha1.TableKindSelf,
			//			DisplayMode: v1alpha1.DisplayModeField,
			//		},
			//		// Insight *PageBlockOutline  `json:"insight,omitempty"`
			//		// Blocks  []PageBlockOutline `json:"blocks" json:"blocks,omitempty"`
			//	},
			//},
			UI: nil,
		},
	}
	return GetResourceLayout(kc, outline)
}

func LoadResourceLayout(kc client.Client, name string) (*v1alpha1.ResourceLayout, error) {
	outline, err := resourceoutlines.LoadByName(name)
	if apierrors.IsNotFound(err) {
		parts := strings.SplitN(name, "-", 3)
		if len(parts) != 3 {
			return nil, err
		}
		var group string
		if parts[0] != "core" {
			group = parts[0]
		}
		return LoadResourceLayoutForGVR(kc, schema.GroupVersionResource{
			Group:    group,
			Version:  parts[1],
			Resource: parts[2],
		})
	} else if err != nil {
		return nil, err
	}

	return GetResourceLayout(kc, outline)
}

func GetResourceLayout(kc client.Client, outline *v1alpha1.ResourceOutline) (*v1alpha1.ResourceLayout, error) {
	src := outline.Spec.Resource

	var result v1alpha1.ResourceLayout
	result.TypeMeta = metav1.TypeMeta{
		Kind:       v1alpha1.ResourceKindResourceLayout,
		APIVersion: v1alpha1.SchemeGroupVersion.String(),
	}
	result.ObjectMeta = outline.ObjectMeta
	result.Spec.DefaultLayout = outline.Spec.DefaultLayout
	result.Spec.Resource = outline.Spec.Resource
	result.Spec.UI = outline.Spec.UI
	if outline.Spec.Header != nil {
		tables, err := FlattenPageBlockOutline(kc, src, *outline.Spec.Header, v1alpha1.Field)
		if err != nil {
			return nil, err
		}
		if len(tables) != 1 {
			return nil, fmt.Errorf("ResourceOutline %s uses multiple headers", outline.Name)
		}
		result.Spec.Header = &tables[0]
	}
	if outline.Spec.TabBar != nil {
		tables, err := FlattenPageBlockOutline(kc, src, *outline.Spec.TabBar, v1alpha1.Field)
		if err != nil {
			return nil, err
		}
		if len(tables) != 1 {
			return nil, fmt.Errorf("ResourceOutline %s uses multiple tab bars", outline.Name)
		}
		result.Spec.TabBar = &tables[0]
	}

	result.Spec.Pages = make([]v1alpha1.ResourcePageLayout, 0, len(outline.Spec.Pages))

	pages := outline.Spec.Pages
	if outline.Spec.DefaultLayout && (len(outline.Spec.Pages) == 0 || outline.Spec.Pages[0].Name != BasicPage) {
		pages = append([]v1alpha1.ResourcePageOutline{
			{
				Name: BasicPage,
				//Info: &v1alpha1.PageBlockOutline{
				//	Kind:        v1alpha1.TableKindSelf,
				//	DisplayMode: v1alpha1.DisplayModeField,
				//},
				// Insight *PageBlockOutline  `json:"insight,omitempty"`
				// Blocks  []PageBlockOutline `json:"blocks" json:"blocks,omitempty"`
			},
		}, outline.Spec.Pages...)
	}
	if pages[0].Name == BasicPage && pages[0].Info == nil {
		pages[0].Info = &v1alpha1.PageBlockOutline{
			Kind:        v1alpha1.TableKindSelf,
			DisplayMode: v1alpha1.DisplayModeField,
		}
	}

	for _, pageOutline := range pages {
		page := v1alpha1.ResourcePageLayout{
			Name:    pageOutline.Name,
			Info:    nil,
			Insight: nil,
			Blocks:  nil,
		}
		if pageOutline.Info != nil {
			tables, err := FlattenPageBlockOutline(kc, src, *pageOutline.Info, v1alpha1.Field)
			if err != nil {
				return nil, err
			}
			if len(tables) != 1 {
				return nil, fmt.Errorf("ResourceOutline %s page %s uses multiple basic blocks", outline.Name, page.Name)
			}
			page.Info = &tables[0]
		}
		if pageOutline.Insight != nil {
			tables, err := FlattenPageBlockOutline(kc, src, *pageOutline.Insight, v1alpha1.Field)
			if err != nil {
				return nil, err
			}
			if len(tables) != 1 {
				return nil, fmt.Errorf("ResourceOutline %s page %s uses multiple insight blocks", outline.Name, page.Name)
			}
			page.Insight = &tables[0]
		}

		var tables []v1alpha1.PageBlockLayout
		for _, block := range pageOutline.Blocks {
			blocks, err := FlattenPageBlockOutline(kc, src, block, v1alpha1.List)
			if err != nil {
				return nil, err
			}
			tables = append(tables, blocks...)
		}
		page.Blocks = tables

		result.Spec.Pages = append(result.Spec.Pages, page)
	}

	return &result, nil
}

func FlattenPageBlockOutline(
	kc client.Client,
	src kmapi.ResourceID,
	in v1alpha1.PageBlockOutline,
	priority v1alpha1.Priority,
) ([]v1alpha1.PageBlockLayout, error) {
	if in.Kind == v1alpha1.TableKindSubTable ||
		in.Kind == v1alpha1.TableKindConnection ||
		in.Kind == v1alpha1.TableKindSelf {
		out, err := Convert_PageBlockOutline_To_PageBlockLayout(kc, src, in, priority)
		if err != nil {
			return nil, err
		}
		return []v1alpha1.PageBlockLayout{out}, nil
	} else if in.Kind != v1alpha1.TableKindBlock {
		return nil, fmt.Errorf("unknown block kind %+v", in)
	}

	obj, err := blockdefs.LoadByName(in.Name)
	if err != nil {
		return nil, err
	}
	var result []v1alpha1.PageBlockLayout
	for _, block := range obj.Spec.Blocks {
		out, err := FlattenPageBlockOutline(kc, src, block, priority)
		if err != nil {
			return nil, err
		}
		result = append(result, out...)
	}
	return result, nil
}

func Convert_PageBlockOutline_To_PageBlockLayout(
	kc client.Client,
	src kmapi.ResourceID,
	in v1alpha1.PageBlockOutline,
	priority v1alpha1.Priority,
) (v1alpha1.PageBlockLayout, error) {
	var columns []v1alpha1.ResourceColumnDefinition
	if in.View != nil {
		if in.View.Name != "" {
			obj, err := tabledefs.LoadByName(in.View.Name)
			if err != nil {
				return v1alpha1.PageBlockLayout{}, err
			}
			columns = obj.Spec.Columns
		} else {
			columns = in.View.Columns
		}
	}

	columns, err := tabledefs.FlattenColumns(columns)
	if err != nil {
		return v1alpha1.PageBlockLayout{}, err
	}

	if in.Kind == v1alpha1.TableKindSubTable && len(columns) == 0 {
		return v1alpha1.PageBlockLayout{}, fmt.Errorf("missing columns for SubTable %s with fieldPath %s", in.Name, in.FieldPath)
	}

	if in.Kind == v1alpha1.TableKindConnection && kc != nil {
		mapping, err := kc.RESTMapper().RESTMapping(schema.GroupKind{Group: in.Ref.Group, Kind: in.Ref.Kind})
		if meta.IsNoMatchError(err) {
			columns = tableconvertor.FilterColumnsWithDefaults(nil, schema.GroupVersionResource{} /*ignore*/, columns, priority)
		} else if err == nil {
			if in.View == nil || (in.View.Name == "" && len(in.View.Columns) == 0) {
				if rv, ok := tabledefs.LoadDefaultByGVK(mapping.GroupVersionKind); ok {
					columns = rv.Spec.Columns
				}
				columns, err = tabledefs.FlattenColumns(columns)
				if err != nil {
					return v1alpha1.PageBlockLayout{}, err
				}
			}
			columns = tableconvertor.FilterColumnsWithDefaults(kc, mapping.Resource, columns, priority)
		} else {
			return v1alpha1.PageBlockLayout{}, err
		}
	} else if in.Kind == v1alpha1.TableKindSelf {
		columns = tableconvertor.FilterColumnsWithDefaults(kc, src.GroupVersionResource(), columns, priority)
	}

	return v1alpha1.PageBlockLayout{
		Kind:            in.Kind,
		Name:            in.Name,
		FieldPath:       in.FieldPath,
		ResourceLocator: in.ResourceLocator,
		DisplayMode:     in.DisplayMode,
		Actions:         in.Actions,
		View: v1alpha1.PageBlockTableDefinition{
			Columns: columns,
		},
	}, nil
}
