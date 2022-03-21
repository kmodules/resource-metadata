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

package resourcedashboards

import (
	"context"
	"embed"
	"fmt"
	iofs "io/fs"
	"path/filepath"
	"reflect"
	"sort"
	"sync"

	kmapi "kmodules.xyz/client-go/api/v1"
	"kmodules.xyz/resource-metadata/apis/meta/v1alpha1"

	"github.com/pkg/errors"
	ioutilx "gomodules.xyz/x/ioutil"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/yaml"
)

var (
	//go:embed **/**/*.yaml **/**/**/*.yaml trigger
	fs embed.FS

	m       sync.Mutex
	rdMap   map[string]*v1alpha1.ResourceDashboard
	rdPerGK map[schema.GroupVersionKind]*v1alpha1.ResourceDashboard
	rdPerGR map[schema.GroupVersionResource]*v1alpha1.ResourceDashboard

	loader = ioutilx.NewReloader(
		filepath.Join("/tmp", "hub", "resourcedashboards"),
		fs,
		func(fsys iofs.FS) {
			rdMap = map[string]*v1alpha1.ResourceDashboard{}
			rdPerGK = map[schema.GroupVersionKind]*v1alpha1.ResourceDashboard{}
			rdPerGR = map[schema.GroupVersionResource]*v1alpha1.ResourceDashboard{}

			if err := iofs.WalkDir(fsys, ".", func(path string, d iofs.DirEntry, err error) error {
				if d.IsDir() || err != nil {
					return errors.Wrap(err, path)
				}
				ext := filepath.Ext(d.Name())
				if ext != ".yaml" && ext != ".yml" && ext != ".json" {
					return nil
				}

				data, err := iofs.ReadFile(fsys, path)
				if err != nil {
					return errors.Wrap(err, path)
				}
				var obj v1alpha1.ResourceDashboard
				err = yaml.Unmarshal(data, &obj)
				if err != nil {
					return errors.Wrap(err, path)
				}
				rdMap[obj.Name] = &obj

				if obj.Spec.DefaultDashboards {
					gvr := obj.Spec.Resource.GroupVersionResource()
					expectedName := DefaultDashboardName(gvr)
					if obj.Name != expectedName {
						return fmt.Errorf("expected default %s name to be %s, found %s", reflect.TypeOf(v1alpha1.ResourceDashboard{}), expectedName, obj.Name)
					}

					gvk := obj.Spec.Resource.GroupVersionKind()
					if rv, ok := rdPerGK[gvk]; !ok {
						rdPerGK[gvk] = &obj
					} else {
						return fmt.Errorf("multiple %s found for %+v: %s and %s", reflect.TypeOf(v1alpha1.ResourceDashboard{}), gvk, rv.Name, obj.Name)
					}
					if rv, ok := rdPerGR[gvr]; !ok {
						rdPerGR[gvr] = &obj
					} else {
						return fmt.Errorf("multiple %s found for %+v: %s and %s", reflect.TypeOf(v1alpha1.ResourceDashboard{}), gvk, rv.Name, obj.Name)
					}
				}
				return nil
			}); err != nil {
				panic(errors.Wrapf(err, "failed to load %s", reflect.TypeOf(v1alpha1.ResourceDashboard{})))
			}
		},
	)
)

func init() {
	loader.ReloadIfTriggered()
}

func EmbeddedFS() iofs.FS {
	return fs
}

func DefaultDashboardName(gvr schema.GroupVersionResource) string {
	if gvr.Group == "" && gvr.Version == "v1" {
		return fmt.Sprintf("core-v1-%s", gvr.Resource)
	}
	return fmt.Sprintf("%s-%s-%s", gvr.Group, gvr.Version, gvr.Resource)
}

func LoadByName(name string) (*v1alpha1.ResourceDashboard, error) {
	m.Lock()
	defer m.Unlock()
	loader.ReloadIfTriggered()

	if obj, ok := rdMap[name]; ok {
		return obj, nil
	}
	return nil, apierrors.NewNotFound(v1alpha1.Resource(v1alpha1.ResourceKindResourceDashboard), name)
}

func LoadDefaultByGVK(gvk schema.GroupVersionKind) (*v1alpha1.ResourceDashboard, bool) {
	m.Lock()
	defer m.Unlock()
	loader.ReloadIfTriggered()

	rv, found := rdPerGK[gvk]
	return rv, found
}

func LoadDefaultByGVR(gvr schema.GroupVersionResource) (*v1alpha1.ResourceDashboard, bool) {
	m.Lock()
	defer m.Unlock()
	loader.ReloadIfTriggered()

	rv, found := rdPerGR[gvr]
	return rv, found
}

func LoadByGVR(kc client.Client, gvr schema.GroupVersionResource) (*v1alpha1.ResourceDashboard, bool) {
	var ed v1alpha1.ResourceDashboard
	err := kc.Get(context.TODO(), client.ObjectKey{Name: DefaultDashboardName(gvr)}, &ed)
	if err == nil {
		return &ed, true
	} else if client.IgnoreNotFound(err) != nil {
		klog.V(3).InfoS(fmt.Sprintf("failed to load resource dashboard for %+v", gvr))
	}
	return LoadDefaultByGVR(gvr)
}

func LoadByResourceID(kc client.Client, rid *kmapi.ResourceID) (*v1alpha1.ResourceDashboard, bool) {
	if rid == nil {
		return nil, false
	}

	gvr := rid.GroupVersionResource()
	if gvr.Version == "" || gvr.Resource == "" {
		id, err := kmapi.ExtractResourceID(kc.RESTMapper(), *rid)
		if client.IgnoreNotFound(err) != nil {
			klog.V(3).InfoS(fmt.Sprintf("failed to extract resource id for %+v", *rid))
		}
		gvr = id.GroupVersionResource()
	}
	return LoadByGVR(kc, gvr)
}

func List() []v1alpha1.ResourceDashboard {
	m.Lock()
	defer m.Unlock()
	loader.ReloadIfTriggered()

	out := make([]v1alpha1.ResourceDashboard, 0, len(rdMap))
	for _, rl := range rdMap {
		out = append(out, *rl)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Name < out[j].Name
	})
	return out
}

func Names() []string {
	m.Lock()
	defer m.Unlock()
	loader.ReloadIfTriggered()

	out := make([]string, 0, len(rdMap))
	for name := range rdMap {
		out = append(out, name)
	}
	sort.Strings(out)
	return out
}
