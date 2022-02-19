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

package resourcetabledefinitions

import (
	"embed"
	"fmt"
	iofs "io/fs"
	"path/filepath"
	"reflect"
	"sort"
	"sync"

	"kmodules.xyz/resource-metadata/apis/meta/v1alpha1"

	"github.com/pkg/errors"
	ioutilx "gomodules.xyz/x/ioutil"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/yaml"
)

var (
	//go:embed **/**/*.yaml **/**/**/*.yaml trigger
	fs embed.FS

	m        sync.Mutex
	rtdMap   map[string]*v1alpha1.ResourceTableDefinition
	rtdPerGK map[schema.GroupVersionKind]*v1alpha1.ResourceTableDefinition
	rtdPerGR map[schema.GroupVersionResource]*v1alpha1.ResourceTableDefinition

	loader = ioutilx.NewReloader(
		filepath.Join("/tmp", "hub", "resourcetabledefinitions"),
		fs,
		func(fsys iofs.FS) {
			rtdMap = map[string]*v1alpha1.ResourceTableDefinition{}
			rtdPerGK = map[schema.GroupVersionKind]*v1alpha1.ResourceTableDefinition{}
			rtdPerGR = map[schema.GroupVersionResource]*v1alpha1.ResourceTableDefinition{}

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
				var obj v1alpha1.ResourceTableDefinition
				err = yaml.Unmarshal(data, &obj)
				if err != nil {
					return errors.Wrap(err, path)
				}
				rtdMap[obj.Name] = &obj

				if obj.Spec.Resource != nil && obj.Spec.DefaultView {
					gvk := obj.Spec.Resource.GroupVersionKind()
					if rv, ok := rtdPerGK[gvk]; !ok {
						rtdPerGK[gvk] = &obj
					} else {
						return fmt.Errorf("multiple %s found for %+v: %s and %s", reflect.TypeOf(v1alpha1.ResourceTableDefinition{}), gvk, rv.Name, obj.Name)
					}
					gvr := obj.Spec.Resource.GroupVersionResource()
					if rv, ok := rtdPerGR[gvr]; !ok {
						rtdPerGR[gvr] = &obj
					} else {
						return fmt.Errorf("multiple %s found for %+v: %s and %s", reflect.TypeOf(v1alpha1.ResourceTableDefinition{}), gvk, rv.Name, obj.Name)
					}
				}
				return nil
			}); err != nil {
				panic(errors.Wrapf(err, "failed to load %s", reflect.TypeOf(v1alpha1.ResourceTableDefinition{})))
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

func GetName(gvr schema.GroupVersionResource) string {
	if gvr.Group == "" && gvr.Version == "v1" {
		return fmt.Sprintf("core-v1-%s", gvr.Resource)
	}
	return fmt.Sprintf("%s-%s-%s", gvr.Group, gvr.Version, gvr.Resource)
}

func LoadByName(name string) (*v1alpha1.ResourceTableDefinition, error) {
	m.Lock()
	defer m.Unlock()
	loader.ReloadIfTriggered()

	return loadByName(name)
}

func loadByName(name string) (*v1alpha1.ResourceTableDefinition, error) {
	if obj, ok := rtdMap[name]; ok {
		return obj, nil
	}
	return nil, apierrors.NewNotFound(v1alpha1.Resource(v1alpha1.ResourceKindResourceTableDefinition), name)
}

func LoadDefaultByGVK(gvk schema.GroupVersionKind) (*v1alpha1.ResourceTableDefinition, bool) {
	m.Lock()
	defer m.Unlock()
	loader.ReloadIfTriggered()

	rv, found := rtdPerGK[gvk]
	return rv, found
}

func LoadDefaultByGVR(gvr schema.GroupVersionResource) (*v1alpha1.ResourceTableDefinition, bool) {
	m.Lock()
	defer m.Unlock()
	loader.ReloadIfTriggered()

	rv, found := rtdPerGR[gvr]
	return rv, found
}

func List() []v1alpha1.ResourceTableDefinition {
	m.Lock()
	defer m.Unlock()
	loader.ReloadIfTriggered()

	out := make([]v1alpha1.ResourceTableDefinition, 0, len(rtdMap))
	for _, rl := range rtdMap {
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

	out := make([]string, 0, len(rtdMap))
	for name := range rtdMap {
		out = append(out, name)
	}
	sort.Strings(out)
	return out
}

func FlattenColumns(in []v1alpha1.ResourceColumnDefinition) ([]v1alpha1.ResourceColumnDefinition, error) {
	m.Lock()
	defer m.Unlock()
	loader.ReloadIfTriggered()

	return flattenColumns(in)
}

func flattenColumns(in []v1alpha1.ResourceColumnDefinition) ([]v1alpha1.ResourceColumnDefinition, error) {
	var foundRef bool
	for _, c := range in {
		if c.Type == v1alpha1.ColumnTypeRef {
			foundRef = true
			break
		}
	}
	if !foundRef {
		return in, nil
	}

	var out []v1alpha1.ResourceColumnDefinition
	for _, c := range in {
		if c.Type == v1alpha1.ColumnTypeRef {
			def, err := loadByName(c.Name)
			if err != nil {
				return nil, err
			}
			cols, err := FlattenColumns(def.Spec.Columns)
			if err != nil {
				return nil, err
			}
			out = append(out, cols...)
		} else {
			out = append(out, c)
		}
	}
	return out, nil
}
