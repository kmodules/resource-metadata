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

package clusterprofiles

import (
	"context"
	"embed"
	iofs "io/fs"
	"path/filepath"
	"reflect"
	"sort"
	"sync"

	"kmodules.xyz/resource-metadata/apis/ui/v1alpha1"

	"github.com/pkg/errors"
	ioutilx "gomodules.xyz/x/ioutil"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/yaml"
)

var (
	//go:embed **/*.yaml trigger
	fs embed.FS

	m     sync.Mutex
	cpMap map[string]*v1alpha1.ClusterProfile

	loader = ioutilx.NewReloader(
		filepath.Join("/tmp", "hub", "clusterprofiles"),
		fs,
		func(fsys iofs.FS) {
			cpMap = map[string]*v1alpha1.ClusterProfile{}

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
				var obj v1alpha1.ClusterProfile
				err = yaml.Unmarshal(data, &obj)
				if err != nil {
					return errors.Wrap(err, path)
				}
				cpMap[obj.Name] = &obj
				return nil
			}); err != nil {
				panic(errors.Wrapf(err, "failed to load %s", reflect.TypeOf(v1alpha1.ClusterProfile{})))
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

func LoadInternalByName(name string) (*v1alpha1.ClusterProfile, error) {
	m.Lock()
	defer m.Unlock()
	loader.ReloadIfTriggered()

	if obj, ok := cpMap[name]; ok {
		return obj, nil
	}
	return nil, apierrors.NewNotFound(v1alpha1.Resource(v1alpha1.ResourceKindClusterProfile), name)
}

func LoadByName(kc client.Reader, name string) (*v1alpha1.ClusterProfile, error) {
	var ed v1alpha1.ClusterProfile
	err := kc.Get(context.TODO(), client.ObjectKey{Name: name}, &ed)
	if meta.IsNoMatchError(err) || apierrors.IsNotFound(err) {
		return LoadInternalByName(name)
	} else if err != nil {
		return nil, err
	}
	return &ed, nil
}

func List(kc client.Reader) ([]v1alpha1.ClusterProfile, error) {
	m.Lock()
	defer m.Unlock()
	loader.ReloadIfTriggered()

	var list v1alpha1.ClusterProfileList
	err := kc.List(context.TODO(), &list)
	if err != nil && !(meta.IsNoMatchError(err) || apierrors.IsNotFound(err)) {
		return nil, err
	}

	profiles := map[string]v1alpha1.ClusterProfile{}
	for name, cp := range cpMap {
		profiles[name] = *cp
	}
	for _, obj := range list.Items {
		profiles[obj.Name] = obj
	}

	out := make([]v1alpha1.ClusterProfile, 0, len(profiles))
	for _, cp := range profiles {
		out = append(out, cp)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Name < out[j].Name
	})
	return out, nil
}

func Names(kc client.Reader) ([]string, error) {
	list, err := List(kc)
	if err != nil {
		return nil, err
	}

	out := make([]string, 0, len(list))
	for _, obj := range list {
		out = append(out, obj.Name)
	}
	return out, nil
}
