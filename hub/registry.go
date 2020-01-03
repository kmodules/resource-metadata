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

package hub

import (
	"fmt"
	"strings"
	"sync"

	"kmodules.xyz/resource-metadata/apis/meta/v1alpha1"
	"kmodules.xyz/resource-metadata/hub/resourceclasses"
	"kmodules.xyz/resource-metadata/hub/resourcedescriptors"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
)

type Registry struct {
	uid    string
	cache  KV
	m      sync.RWMutex
	regGVK map[schema.GroupVersionKind]*v1alpha1.ResourceID
	regGVR map[schema.GroupVersionResource]*v1alpha1.ResourceID
}

func NewRegistry(uid string, cache KV) *Registry {
	r := &Registry{
		uid:    uid,
		cache:  cache,
		regGVK: map[schema.GroupVersionKind]*v1alpha1.ResourceID{},
		regGVR: map[schema.GroupVersionResource]*v1alpha1.ResourceID{},
	}

	r.cache.Visit(func(key string, val *v1alpha1.ResourceDescriptor) {
		v := val.Spec.Resource // copy
		r.regGVK[v.GroupVersionKind()] = &v
		r.regGVR[v.GroupVersionResource()] = &v
	})
	return r
}

func NewRegistryOfKnownResources() *Registry {
	return NewRegistry(KnownUID, KnownResources)
}

func (r *Registry) Register(gvr schema.GroupVersionResource, dc discovery.ServerResourcesInterface) error {
	r.m.RLock()
	if _, found := r.regGVR[gvr]; found {
		r.m.RUnlock()
		return nil
	}
	r.m.RUnlock()

	reg, err := r.createRegistry(dc)
	if err != nil {
		return err
	}

	r.m.Lock()
	for filename, rd := range reg {
		if _, found := r.cache.Get(filename); !found {
			r.regGVK[rd.Spec.Resource.GroupVersionKind()] = &rd.Spec.Resource
			r.regGVR[rd.Spec.Resource.GroupVersionResource()] = &rd.Spec.Resource
			r.cache.Set(filename, rd)
		}
	}
	r.m.Unlock()

	return nil
}

func (r *Registry) createRegistry(dc discovery.ServerResourcesInterface) (map[string]*v1alpha1.ResourceDescriptor, error) {
	rsLists, err := dc.ServerPreferredResources()
	if err != nil && !discovery.IsGroupDiscoveryFailedError(err) {
		return nil, err
	}

	reg := make(map[string]*v1alpha1.ResourceDescriptor)
	for _, rsList := range rsLists {
		for i := range rsList.APIResources {
			rs := rsList.APIResources[i]

			gv, err := schema.ParseGroupVersion(rsList.GroupVersion)
			if err != nil {
				return nil, err
			}
			rs.Group = gv.Group
			rs.Version = gv.Version

			scope := v1alpha1.ClusterScoped
			if rs.Namespaced {
				scope = v1alpha1.NamespaceScoped
			}

			filename := fmt.Sprintf("%s/%s/%s.yaml", rs.Group, rs.Version, rs.Name)
			reg[filename] = &v1alpha1.ResourceDescriptor{
				TypeMeta: metav1.TypeMeta{
					APIVersion: v1alpha1.SchemeGroupVersion.String(),
					Kind:       v1alpha1.ResourceKindResourceDescriptor,
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: fmt.Sprintf("%s-%s-%s", rs.Group, rs.Version, rs.Name),
					Labels: map[string]string{
						"k8s.io/group":    rs.Group,
						"k8s.io/version":  rs.Version,
						"k8s.io/resource": rs.Name,
						"k8s.io/kind":     rs.Kind,
					},
				},
				Spec: v1alpha1.ResourceDescriptorSpec{
					Resource: v1alpha1.ResourceID{
						Group:   rs.Group,
						Version: rs.Version,
						Name:    rs.Name,
						Kind:    rs.Kind,
						Scope:   scope,
					},
				},
			}
		}
	}

	for _, name := range resourcedescriptors.AssetNames() {
		delete(reg, name)
	}
	return reg, nil
}

func (r *Registry) Visit(f func(key string, val *v1alpha1.ResourceDescriptor)) {
	r.cache.Visit(f)
}

func (r *Registry) GVR(gvk schema.GroupVersionKind) (schema.GroupVersionResource, error) {
	r.m.RLock()
	defer r.m.RUnlock()
	rid, exist := r.regGVK[gvk]
	if !exist {
		return schema.GroupVersionResource{}, UnregisteredErr{gvk.String()}
	}
	return rid.GroupVersionResource(), nil
}

func (r *Registry) TypeMeta(gvr schema.GroupVersionResource) (metav1.TypeMeta, error) {
	r.m.RLock()
	defer r.m.RUnlock()
	rid, exist := r.regGVR[gvr]
	if !exist {
		return metav1.TypeMeta{}, UnregisteredErr{gvr.String()}
	}
	return rid.TypeMeta(), nil
}

func (r *Registry) GVK(gvr schema.GroupVersionResource) (schema.GroupVersionKind, error) {
	r.m.RLock()
	defer r.m.RUnlock()
	rid, exist := r.regGVR[gvr]
	if !exist {
		return schema.GroupVersionKind{}, UnregisteredErr{gvr.String()}
	}
	return rid.GroupVersionKind(), nil
}

func (r *Registry) IsNamespaced(gvr schema.GroupVersionResource) (bool, error) {
	r.m.RLock()
	defer r.m.RUnlock()
	rid, exist := r.regGVR[gvr]
	if !exist {
		return false, UnregisteredErr{gvr.String()}
	}
	return rid.Scope == v1alpha1.NamespaceScoped, nil
}

func (r *Registry) Types() []metav1.TypeMeta {
	r.m.RLock()
	defer r.m.RUnlock()

	types := make([]metav1.TypeMeta, 0, len(r.regGVK))
	for _, v := range r.regGVK {
		types = append(types, v.TypeMeta())
	}
	return types
}

func (r *Registry) Resources() []schema.GroupVersionResource {
	r.m.RLock()
	defer r.m.RUnlock()

	resources := make([]schema.GroupVersionResource, 0, len(r.regGVR))
	for k := range r.regGVR {
		resources = append(resources, k)
	}
	return resources
}

func (r *Registry) LoadByGVR(gvr schema.GroupVersionResource) (*v1alpha1.ResourceDescriptor, error) {
	var filename string
	if gvr.Group == "" && gvr.Version == "v1" {
		filename = fmt.Sprintf("core/v1/%s.yaml", gvr.Resource)
	} else {
		filename = fmt.Sprintf("%s/%s/%s.yaml", gvr.Group, gvr.Version, gvr.Resource)
	}
	return r.LoadByFile(filename)
}

func (r *Registry) LoadByName(name string) (*v1alpha1.ResourceDescriptor, error) {
	filename := strings.Replace(name, "-", "/", 2) + ".yaml"
	return r.LoadByFile(filename)
}

func (r *Registry) LoadByFile(filename string) (*v1alpha1.ResourceDescriptor, error) {
	obj, ok := r.cache.Get(filename)
	if !ok {
		return nil, UnregisteredErr{filename}
	}
	return obj, nil
}

func (r *Registry) DefaultResourcePanel() (*v1alpha1.ResourcePanel, error) {
	return r.newPanel(true, false)
}

func (r *Registry) AvailableResourcePanel() (*v1alpha1.ResourcePanel, error) {
	return r.newPanel(false, true)
}

func (r *Registry) newPanel(skipK8sGroups, mutateRequiredSections bool) (*v1alpha1.ResourcePanel, error) {
	sections := make(map[string]*v1alpha1.PanelSection)

	// first add the known required sections
	for group, rc := range KnownClasses {
		if !rc.IsRequired() {
			continue
		}

		section := &v1alpha1.PanelSection{
			Name:              rc.Name,
			ResourceClassInfo: rc.Spec.ResourceClassInfo,
		}
		for _, entry := range rc.Spec.Entries {
			pe := v1alpha1.PanelEntry{
				Entry:      entry,
				Namespaced: false,
				Icons:      nil,
			}
			if entry.Type != nil {
				if rd, err := r.LoadByGVR(entry.Type.GVR()); err != nil {
					pe.Namespaced = rd.Spec.Resource.Scope == v1alpha1.NamespaceScoped
					pe.Icons = rd.Spec.Icons
				}
			}
			section.Entries = append(section.Entries, pe)
		}
		sections[group] = section
	}

	// now, auto discover sections from registry
	r.Visit(func(_ string, rd *v1alpha1.ResourceDescriptor) {
		if skipK8sGroups && (rd.Spec.Resource.Group == "" ||
			strings.ContainsRune(rd.Spec.Resource.Group, '.') ||
			strings.HasSuffix(rd.Spec.Resource.Group, ".k8s.io")) {
			return // skip k8s.io api groups
		}

		name := resourceclasses.ResourceClassName(rd.Spec.Resource.Group)

		section, found := sections[rd.Spec.Resource.Group]
		if found {
			if !mutateRequiredSections {
				return // this api group was manually configured with required entries
			}
		} else {
			if rc, found := KnownClasses[rd.Spec.Resource.Group]; found {
				section = &v1alpha1.PanelSection{
					Name:              rc.Name,
					ResourceClassInfo: rc.Spec.ResourceClassInfo,
				}
			} else {
				// unknown api group, so use CRD icon
				section = &v1alpha1.PanelSection{
					Name: name,
					ResourceClassInfo: v1alpha1.ResourceClassInfo{
						APIGroup: rd.Spec.Resource.Group,
						Icons: []v1alpha1.ImageSpec{
							{
								Source: "https://cdn.appscode.com/k8s/icons/apiextensions.k8s.io/crd.svg",
								Type:   "image/svg+xml",
							},
						},
					},
				}
			}
			sections[rd.Spec.Resource.Group] = section
		}

		if !section.Contains(rd) {
			section.Entries = append(section.Entries, v1alpha1.PanelEntry{
				Entry: v1alpha1.Entry{
					Name: rd.Spec.Resource.Kind,
					Type: &v1alpha1.GroupVersionResource{
						Group:    rd.Spec.Resource.Group,
						Version:  rd.Spec.Resource.Version,
						Resource: rd.Spec.Resource.Name,
					},
				},
				Namespaced: rd.Spec.Resource.Scope == v1alpha1.NamespaceScoped,
				Icons:      rd.Spec.Icons,
			})
		}
	})

	out := &v1alpha1.ResourcePanel{
		Sections: make([]v1alpha1.PanelSection, 0, len(sections)),
	}
	for key := range sections {
		out.Sections = append(out.Sections, *sections[key])
	}
	return out, nil
}

type UnregisteredErr struct {
	t string
}

var _ error = UnregisteredErr{}

func (e UnregisteredErr) Error() string {
	return e.t + " isn't registered"
}
