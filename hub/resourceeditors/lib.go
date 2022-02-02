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

package resourceeditors

import (
	"embed"
	"fmt"
	iofs "io/fs"
	"reflect"
	"sort"

	"kmodules.xyz/resource-metadata/apis/meta/v1alpha1"

	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/yaml"
)

//go:embed **/**/*.yaml
var fs embed.FS

func FS() embed.FS {
	return fs
}

var (
	reMap = map[string]*v1alpha1.ResourceEditor{}
)

func init() {
	if err := iofs.WalkDir(fs, ".", func(path string, d iofs.DirEntry, err error) error {
		if d.IsDir() || err != nil {
			return err
		}
		data, err := fs.ReadFile(path)
		if err != nil {
			return errors.Wrap(err, path)
		}
		var obj v1alpha1.ResourceEditor
		err = yaml.Unmarshal(data, &obj)
		if err != nil {
			return errors.Wrap(err, path)
		}
		reMap[obj.Name] = &obj

		return nil
	}); err != nil {
		panic(errors.Wrapf(err, "failed to load %s", reflect.TypeOf(v1alpha1.ResourceEditor{})))
	}
}

func DefaultEditorName(gvr schema.GroupVersionResource) string {
	if gvr.Group == "" && gvr.Version == "v1" {
		return fmt.Sprintf("core-v1-%s", gvr.Resource)
	}
	return fmt.Sprintf("%s-%s-%s", gvr.Group, gvr.Version, gvr.Resource)
}

func LoadByName(name string) (*v1alpha1.ResourceEditor, error) {
	if obj, ok := reMap[name]; ok {
		return obj, nil
	}
	return nil, apierrors.NewNotFound(v1alpha1.Resource(v1alpha1.ResourceKindResourceEditor), name)
}

func LoadForGVR(gvr schema.GroupVersionResource) (*v1alpha1.ResourceEditor, bool) {
	name := DefaultEditorName(gvr)
	obj, ok := reMap[name]
	return obj, ok
}

func List() []v1alpha1.ResourceEditor {
	out := make([]v1alpha1.ResourceEditor, 0, len(reMap))
	for _, rl := range reMap {
		out = append(out, *rl)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Name < out[j].Name
	})
	return out
}

func Names() []string {
	out := make([]string, 0, len(reMap))
	for name := range reMap {
		out = append(out, name)
	}
	sort.Strings(out)
	return out
}
