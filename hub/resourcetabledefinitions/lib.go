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
	rtdMap   = map[string]*v1alpha1.ResourceTableDefinition{}
	rtdPerGK = map[schema.GroupVersionKind]*v1alpha1.ResourceTableDefinition{}
	rtdPerGR = map[schema.GroupVersionResource]*v1alpha1.ResourceTableDefinition{}
)

func init() {
	if err := iofs.WalkDir(fs, ".", func(path string, d iofs.DirEntry, err error) error {
		if d.IsDir() || err != nil {
			return errors.Wrap(err, path)
		}
		data, err := fs.ReadFile(path)
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
}

func LoadByName(name string) (*v1alpha1.ResourceTableDefinition, error) {
	if obj, ok := rtdMap[name]; ok {
		return obj, nil
	}
	return nil, apierrors.NewNotFound(v1alpha1.Resource(v1alpha1.ResourceKindResourceTableDefinition), name)
}

func DefaultTableDefinitionForGVK(gvk schema.GroupVersionKind) (*v1alpha1.ResourceTableDefinition, bool) {
	rv, found := rtdPerGK[gvk]
	return rv, found
}

func DefaultTableDefinitionForGVR(gvr schema.GroupVersionResource) (*v1alpha1.ResourceTableDefinition, bool) {
	rv, found := rtdPerGR[gvr]
	return rv, found
}

func List() []v1alpha1.ResourceTableDefinition {
	out := make([]v1alpha1.ResourceTableDefinition, 0, len(rtdMap))
	for _, rl := range rtdMap {
		out = append(out, *rl)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Name < out[j].Name
	})
	return out
}
