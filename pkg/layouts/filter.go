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
	"context"
	"fmt"

	kmapi "kmodules.xyz/client-go/api/v1"
	rsapi "kmodules.xyz/resource-metadata/apis/meta/v1alpha1"
	uiapi "kmodules.xyz/resource-metadata/apis/ui/v1alpha1"
	blockdefs "kmodules.xyz/resource-metadata/hub/resourceblockdefinitions"
	"kmodules.xyz/resource-metadata/hub/resourceeditors"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GetDefaultResourceOutlineFilter(kc client.Client, outline *rsapi.ResourceOutline) (*uiapi.ResourceOutlineFilter, error) {
	src := outline.Spec.Resource

	var result uiapi.ResourceOutlineFilter
	result.TypeMeta = metav1.TypeMeta{
		Kind:       uiapi.ResourceKindResourceOutlineFilter,
		APIVersion: uiapi.SchemeGroupVersion.String(),
	}
	result.ObjectMeta = outline.ObjectMeta
	result.Spec.Resource = outline.Spec.Resource
	result.Spec.Header = outline.Spec.Header != nil
	result.Spec.TabBar = outline.Spec.TabBar != nil
	if ed, err := resourceeditors.LoadByGVR(kc, outline.Spec.Resource.GroupVersionResource()); err == nil {
		if ed.Spec.UI != nil {
			{
				result.Spec.Actions = make([]uiapi.ActionTemplateGroupFilter, 0, len(ed.Spec.UI.Actions))
				for _, ag := range ed.Spec.UI.Actions {
					ag2 := uiapi.ActionTemplateGroupFilter{
						Name:  ag.Name,
						Items: make([]uiapi.ItemFilter, 0, len(ag.Items)),
					}
					for _, a := range ag.Items {
						a2 := uiapi.ItemFilter{
							Name: a.Name,
							Show: true,
						}
						ag2.Items = append(ag2.Items, a2)
					}
					result.Spec.Actions = append(result.Spec.Actions, ag2)
				}
			}
		}
	}

	result.Spec.Pages = make([]uiapi.ResourcePageOutlineFilter, 0, len(outline.Spec.Pages))

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
		page := uiapi.ResourcePageOutlineFilter{
			Name:     pageOutline.Name,
			Sections: make([]uiapi.SectionOutlineFilter, 0, len(pageOutline.Sections)),
		}
		for _, sectionOutline := range pageOutline.Sections {

			section := uiapi.SectionOutlineFilter{
				Name:    sectionOutline.Name,
				Info:    sectionOutline.Info != nil,
				Insight: sectionOutline.Insight != nil,
			}

			var tables []uiapi.ItemFilter
			for _, block := range sectionOutline.Blocks {
				blocks, err := FlattenPageBlockOutlineFilter(kc, src, block, rsapi.List)
				if err != nil {
					return nil, err
				}
				tables = append(tables, blocks...)
			}
			section.Blocks = tables

			page.Sections = append(page.Sections, section)
		}
		result.Spec.Pages = append(result.Spec.Pages, page)
	}

	return &result, nil
}

func FlattenPageBlockOutlineFilter(
	kc client.Client,
	src kmapi.ResourceID,
	in rsapi.PageBlockOutline,
	priority rsapi.Priority,
) ([]uiapi.ItemFilter, error) {
	if in.Kind == rsapi.TableKindSubTable ||
		in.Kind == rsapi.TableKindConnection ||
		in.Kind == rsapi.TableKindCustom ||
		in.Kind == rsapi.TableKindSelf {
		return []uiapi.ItemFilter{
			{
				Name: in.Name,
				Show: true,
			},
		}, nil
	} else if in.Kind != rsapi.TableKindBlock {
		return nil, fmt.Errorf("unknown block kind %+v", in)
	}

	obj, err := blockdefs.LoadByName(in.Name)
	if err != nil {
		return nil, err
	}
	var result []uiapi.ItemFilter
	for _, block := range obj.Spec.Blocks {
		out, err := FlattenPageBlockOutlineFilter(kc, src, block, priority)
		if err != nil {
			return nil, err
		}
		result = append(result, out...)
	}
	return result, nil
}

func GetResourceOutlineFilter(kc client.Client, outline *rsapi.ResourceOutline) (*uiapi.ResourceOutlineFilter, error) {
	var result uiapi.ResourceOutlineFilter
	err := kc.Get(context.TODO(), client.ObjectKey{Name: outline.Name}, &result)
	if err == nil {
		return &result, nil
	} else if client.IgnoreNotFound(err) != nil {
		return nil, err
	}
	return GetDefaultResourceOutlineFilter(kc, outline)
}
