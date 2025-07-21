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

	kmapi "kmodules.xyz/client-go/api/v1"
	meta_util "kmodules.xyz/client-go/meta"
	rsapi "kmodules.xyz/resource-metadata/apis/meta/v1alpha1"
	"kmodules.xyz/resource-metadata/apis/shared"
	"kmodules.xyz/resource-metadata/hub"
	blockdefs "kmodules.xyz/resource-metadata/hub/resourceblockdefinitions"
	"kmodules.xyz/resource-metadata/hub/resourceeditors"
	"kmodules.xyz/resource-metadata/hub/resourceoutlines"
	tabledefs "kmodules.xyz/resource-metadata/hub/resourcetabledefinitions"
	"kmodules.xyz/resource-metadata/pkg/tableconvertor"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
	releasesapi "x-helm.dev/apimachinery/apis/releases/v1alpha1"
)

const HomePage = "Overview"

var reg = hub.NewRegistryOfKnownResources()

func LoadResourceLayoutForGVR(kc client.Client, gvr schema.GroupVersionResource) (*rsapi.ResourceLayout, error) {
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
			return nil, apierrors.NewNotFound(rsapi.Resource(rsapi.ResourceKindResourceOutline), gvr.String())
		}
	}
	return rid, nil
}

func LoadResourceLayoutForGVK(kc client.Client, gvk schema.GroupVersionKind) (*rsapi.ResourceLayout, error) {
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
			return nil, apierrors.NewNotFound(rsapi.Resource(rsapi.ResourceKindResourceOutline), gvk.String())
		}
	}
	return rid, nil
}

func generateDefaultLayout(kc client.Client, rid kmapi.ResourceID) (*rsapi.ResourceLayout, error) {
	outline := &rsapi.ResourceOutline{
		TypeMeta: metav1.TypeMeta{
			Kind:       rsapi.ResourceKindResourceOutline,
			APIVersion: rsapi.SchemeGroupVersion.String(),
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
		Spec: rsapi.ResourceOutlineSpec{
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
		},
	}
	return GetResourceLayout(kc, outline)
}

func LoadResourceLayout(kc client.Client, name string) (*rsapi.ResourceLayout, error) {
	outline, err := resourceoutlines.LoadByName(name)
	if apierrors.IsNotFound(err) {
		gvr, e2 := hub.ParseGVR(name)
		if e2 != nil {
			return nil, err
		}
		return LoadResourceLayoutForGVR(kc, *gvr)
	} else if err != nil {
		return nil, err
	}

	return GetResourceLayout(kc, outline)
}

func GetResourceLayout(kc client.Client, outline *rsapi.ResourceOutline) (*rsapi.ResourceLayout, error) {
	filter, err := GetResourceOutlineFilter(kc, outline)
	if err != nil {
		return nil, err
	}

	src := outline.Spec.Resource

	var result rsapi.ResourceLayout
	result.TypeMeta = metav1.TypeMeta{
		Kind:       rsapi.ResourceKindResourceLayout,
		APIVersion: rsapi.SchemeGroupVersion.String(),
	}
	result.ObjectMeta = outline.ObjectMeta
	result.Spec.DefaultLayout = outline.Spec.DefaultLayout
	result.Spec.Resource = outline.Spec.Resource
	if ed, err := resourceeditors.LoadByGVR(kc, outline.Spec.Resource.GroupVersionResource()); err == nil {
		if ed.Spec.UI != nil {
			result.Spec.UI = &shared.UIParameterTemplate{
				InstanceLabelPaths: ed.Spec.UI.InstanceLabelPaths,
			}

			expand := func(ref *releasesapi.ChartSourceRef) *releasesapi.ChartSourceRef {
				if ref == nil {
					return nil
				}
				if ref.SourceRef.Namespace == "" {
					ref.SourceRef.Namespace = meta_util.PodNamespace()
				}
				return ref
			}
			{
				result.Spec.UI.Editor = expand(ed.Spec.UI.Editor)
				result.Spec.UI.Options = expand(ed.Spec.UI.Options)
				result.Spec.UI.EnforceQuota = ed.Spec.UI.EnforceQuota
			}
			{
				result.Spec.UI.Actions = make([]*shared.ActionTemplateGroup, 0, len(ed.Spec.UI.Actions))
				for _, ag := range ed.Spec.UI.Actions {
					agFilter := filter.Spec.GetAction(ag.Name)
					if !agFilter.Show {
						continue
					}
					ag2 := shared.ActionTemplateGroup{
						ActionInfo: ag.ActionInfo,
						Items:      make([]shared.ActionTemplate, 0, len(ag.Items)),
					}
					for _, a := range ag.Items {
						if agFilter.Items[a.Name] {
							a2 := shared.ActionTemplate{
								ActionInfo:       a.ActionInfo,
								Icons:            a.Icons,
								OperationID:      a.OperationID,
								Flow:             a.Flow,
								DisabledTemplate: a.DisabledTemplate,
								EnforceQuota:     a.EnforceQuota,
							}
							a2.Editor = expand(a.Editor)

							ag2.Items = append(ag2.Items, a2)
						}
					}
					result.Spec.UI.Actions = append(result.Spec.UI.Actions, &ag2)
				}
			}
		}
	}
	if outline.Spec.Header != nil && filter.Spec.Header {
		tables, err := FlattenPageBlockOutline(kc, src, *outline.Spec.Header, rsapi.Field)
		if err != nil {
			return nil, err
		}
		if len(tables) != 1 {
			return nil, fmt.Errorf("ResourceOutline %s uses multiple headers", outline.Name)
		}
		result.Spec.Header = &tables[0]
	}
	if outline.Spec.TabBar != nil && filter.Spec.TabBar {
		tables, err := FlattenPageBlockOutline(kc, src, *outline.Spec.TabBar, rsapi.Field)
		if err != nil {
			return nil, err
		}
		if len(tables) != 1 {
			return nil, fmt.Errorf("ResourceOutline %s uses multiple tab bars", outline.Name)
		}
		result.Spec.TabBar = &tables[0]
	}

	result.Spec.Pages = make([]rsapi.ResourcePageLayout, 0, len(outline.Spec.Pages))

	pages := outline.Spec.Pages
	if outline.Spec.DefaultLayout && (len(pages) == 0 || pages[0].Name != HomePage) {
		pages = append([]rsapi.ResourcePageOutline{
			{
				Name:     HomePage,
				Sections: nil,
			},
		}, outline.Spec.Pages...)
	}

	if len(pages) > 0 && pages[0].Name == HomePage {
		if len(pages[0].Sections) == 0 {
			pages[0].Sections = []rsapi.SectionOutline{
				{
					Info: &rsapi.PageBlockOutline{
						Kind:        rsapi.TableKindSelf,
						DisplayMode: rsapi.DisplayModeField,
					},
				},
			}
		} else {
			if pages[0].Sections[0].Info == nil {
				pages[0].Sections[0].Info = &rsapi.PageBlockOutline{
					Kind:        rsapi.TableKindSelf,
					DisplayMode: rsapi.DisplayModeField,
				}
			}
		}
	}

	for _, pageOutline := range pages {
		pageFilter := filter.Spec.GetPage(pageOutline.Name)
		if !pageFilter.Show {
			continue
		}
		page := rsapi.ResourcePageLayout{
			Name:                pageOutline.Name,
			RequiredFeatureSets: pageOutline.RequiredFeatureSets,
			Sections:            make([]rsapi.SectionLayout, 0, len(pageOutline.Sections)),
		}
		for _, sectionOutline := range pageOutline.Sections {
			sectionFilter := pageFilter.GetSection(sectionOutline.Name)
			if !sectionFilter.Show {
				continue
			}

			section := rsapi.SectionLayout{
				Name:                sectionOutline.Name,
				Icons:               sectionOutline.Icons,
				Info:                nil,
				Insight:             nil,
				RequiredFeatureSets: sectionOutline.RequiredFeatureSets,
			}
			if sectionOutline.Info != nil && sectionFilter.Info["basic"] {
				tables, err := FlattenPageBlockOutline(kc, src, *sectionOutline.Info, rsapi.Field)
				if err != nil {
					return nil, err
				}
				if len(tables) != 1 {
					return nil, fmt.Errorf("ResourceOutline %s page %s uses multiple basic blocks", outline.Name, section.Name)
				}
				section.Info = &tables[0]

				section.Info.Filters = make(map[string]bool)
				for typ, show := range sectionFilter.Info {
					if typ == "basic" {
						continue
					}
					section.Info.Filters[typ] = show
				}
			}
			if sectionOutline.Insight != nil && sectionFilter.Insight {
				tables, err := FlattenPageBlockOutline(kc, src, *sectionOutline.Insight, rsapi.Field)
				if err != nil {
					return nil, err
				}
				if len(tables) != 1 {
					return nil, fmt.Errorf("ResourceOutline %s page %s uses multiple insight blocks", outline.Name, section.Name)
				}
				section.Insight = &tables[0]
			}

			var tables []rsapi.PageBlockLayout
			for _, block := range sectionOutline.Blocks {
				blocks, err := FlattenPageBlockOutline(kc, src, block, rsapi.List)
				if err != nil {
					return nil, err
				}
				tables = append(tables, blocks...)
			}

			// https://go.dev/wiki/SliceTricks#filtering-without-allocating
			b := tables[:0]
			for _, x := range tables {
				if sectionFilter.Blocks[x.Name] {
					b = append(b, x)
				}
			}
			section.Blocks = b

			page.Sections = append(page.Sections, section)
		}
		result.Spec.Pages = append(result.Spec.Pages, page)
	}

	return &result, nil
}

func FlattenPageBlockOutline(
	kc client.Client,
	src kmapi.ResourceID,
	in rsapi.PageBlockOutline,
	priority rsapi.Priority,
) ([]rsapi.PageBlockLayout, error) {
	if in.Kind == rsapi.TableKindSubTable ||
		in.Kind == rsapi.TableKindConnection ||
		in.Kind == rsapi.TableKindCustom ||
		in.Kind == rsapi.TableKindSelf {
		out, err := Convert_PageBlockOutline_To_PageBlockLayout(kc, src, in, priority)
		if err != nil {
			return nil, err
		}
		return []rsapi.PageBlockLayout{out}, nil
	} else if in.Kind != rsapi.TableKindBlock {
		return nil, fmt.Errorf("unknown block kind %+v", in)
	}

	obj, err := blockdefs.LoadByName(in.Name)
	if err != nil {
		return nil, err
	}
	var result []rsapi.PageBlockLayout
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
	in rsapi.PageBlockOutline,
	priority rsapi.Priority,
) (rsapi.PageBlockLayout, error) {
	var columns []rsapi.ResourceColumnDefinition
	if in.View != nil {
		if in.View.Name != "" {
			obj, err := tabledefs.LoadByName(in.View.Name)
			if err != nil {
				return rsapi.PageBlockLayout{}, err
			}
			columns = obj.Spec.Columns
		} else {
			columns = in.View.Columns
		}
	}

	columns, err := tabledefs.FlattenColumns(columns)
	if err != nil {
		return rsapi.PageBlockLayout{}, err
	}

	if in.Kind == rsapi.TableKindSubTable && len(columns) == 0 {
		return rsapi.PageBlockLayout{}, fmt.Errorf("missing columns for SubTable %s with fieldPath %s", in.Name, in.FieldPath)
	}

	if in.Kind == rsapi.TableKindConnection && kc != nil {
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
					return rsapi.PageBlockLayout{}, err
				}
			}
			columns = tableconvertor.FilterColumnsWithDefaults(kc, mapping.Resource, columns, priority)
		} else {
			return rsapi.PageBlockLayout{}, err
		}
	} else if in.Kind == rsapi.TableKindSelf {
		columns = tableconvertor.FilterColumnsWithDefaults(kc, src.GroupVersionResource(), columns, priority)
	}

	result := rsapi.PageBlockLayout{
		Kind:                in.Kind,
		Name:                in.Name,
		Width:               in.Width,
		Icons:               in.Icons,
		FieldPath:           in.FieldPath,
		ResourceLocator:     in.ResourceLocator,
		DisplayMode:         in.DisplayMode,
		Actions:             in.Actions,
		RequiredFeatureSets: in.RequiredFeatureSets,
	}
	if len(columns) > 0 {
		result.View = &rsapi.PageBlockTableDefinition{
			Columns: columns,
		}
		if in.View != nil && in.View.Sort != nil {
			result.View.Sort = in.View.Sort
		}
	}
	return result, nil
}
