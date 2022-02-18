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

package resourceblockdefinitions

import (
	"embed"
	iofs "io/fs"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"sync"

	"kmodules.xyz/resource-metadata/apis/meta/v1alpha1"

	"github.com/pkg/errors"
	ioutilx "gomodules.xyz/x/ioutil"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/yaml"
)

var (
	//go:embed **/**/*.yaml trigger
	fs embed.FS

	m     sync.Mutex
	rbMap map[string]*v1alpha1.ResourceBlockDefinition

	loader = ioutilx.NewReloader(
		filepath.Join(os.TempDir(), "hub", "resourceblockdefinitions"),
		fs,
		func(fsys iofs.FS) {
			rbMap = map[string]*v1alpha1.ResourceBlockDefinition{}

			if err := iofs.WalkDir(fsys, ".", func(path string, d iofs.DirEntry, err error) error {
				if d.IsDir() || d.Name() == ioutilx.TriggerFile || err != nil {
					return errors.Wrap(err, path)
				}

				data, err := iofs.ReadFile(fsys, path)
				if err != nil {
					return errors.Wrap(err, path)
				}
				var obj v1alpha1.ResourceBlockDefinition
				err = yaml.Unmarshal(data, &obj)
				if err != nil {
					return errors.Wrap(err, path)
				}
				rbMap[obj.Name] = &obj
				return nil
			}); err != nil {
				panic(errors.Wrapf(err, "failed to load %s", reflect.TypeOf(v1alpha1.ResourceBlockDefinition{})))
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

func LoadByName(name string) (*v1alpha1.ResourceBlockDefinition, error) {
	m.Lock()
	defer m.Unlock()
	loader.ReloadIfTriggered()

	if obj, ok := rbMap[name]; ok {
		return obj, nil
	}
	return nil, apierrors.NewNotFound(v1alpha1.Resource(v1alpha1.ResourceKindResourceBlockDefinition), name)
}

func List() []v1alpha1.ResourceBlockDefinition {
	m.Lock()
	defer m.Unlock()
	loader.ReloadIfTriggered()

	out := make([]v1alpha1.ResourceBlockDefinition, 0, len(rbMap))
	for _, rl := range rbMap {
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

	out := make([]string, 0, len(rbMap))
	for name := range rbMap {
		out = append(out, name)
	}
	sort.Strings(out)
	return out
}
