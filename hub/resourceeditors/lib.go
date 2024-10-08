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
	"context"
	"embed"
	"fmt"
	iofs "io/fs"
	"path/filepath"
	"reflect"
	"sort"
	"sync"

	kmapi "kmodules.xyz/client-go/api/v1"
	"kmodules.xyz/resource-metadata/apis/ui/v1alpha1"

	"github.com/pkg/errors"
	ioutilx "gomodules.xyz/x/ioutil"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/yaml"
)

var (
	//go:embed **/**/*.yaml trigger
	fs embed.FS

	m     sync.Mutex
	reMap map[string]*v1alpha1.ResourceEditor

	loader = ioutilx.NewReloader(
		filepath.Join("/tmp", "hub", "resourceeditors"),
		fs,
		func(fsys iofs.FS) {
			reMap = map[string]*v1alpha1.ResourceEditor{}

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
		},
	)
)

func init() {
	loader.ReloadIfTriggered()
}

func EmbeddedFS() iofs.FS {
	return fs
}

func DefaultEditorName(gvr schema.GroupVersionResource) string {
	if gvr.Group == "" && gvr.Version == "v1" {
		return fmt.Sprintf("core-v1-%s", gvr.Resource)
	}
	return fmt.Sprintf("%s-%s-%s", gvr.Group, gvr.Version, gvr.Resource)
}

func LoadInternalByName(name string) (*v1alpha1.ResourceEditor, error) {
	m.Lock()
	defer m.Unlock()
	loader.ReloadIfTriggered()

	if obj, ok := reMap[name]; ok {
		return obj, nil
	}
	return nil, apierrors.NewNotFound(v1alpha1.Resource(v1alpha1.ResourceKindResourceEditor), name)
}

func LoadInternalByGVR(gvr schema.GroupVersionResource) (*v1alpha1.ResourceEditor, error) {
	return LoadInternalByName(DefaultEditorName(gvr))
}

func LoadByName(kc client.Client, name string) (*v1alpha1.ResourceEditor, error) {
	var ed v1alpha1.ResourceEditor
	err := kc.Get(context.TODO(), client.ObjectKey{Name: name}, &ed)
	if meta.IsNoMatchError(err) || apierrors.IsNotFound(err) {
		return LoadInternalByName(name)
	} else if err != nil {
		return nil, err
	}
	d, _ := LoadInternalByName(name)
	return merge(&ed, d), nil
}

func LoadByGVR(kc client.Client, gvr schema.GroupVersionResource) (*v1alpha1.ResourceEditor, error) {
	return LoadByName(kc, DefaultEditorName(gvr))
}

func merge(in, d *v1alpha1.ResourceEditor) *v1alpha1.ResourceEditor {
	if d == nil {
		return in
	}

	in.Labels = d.Labels
	in.Spec.Resource = d.Spec.Resource

	if d.Spec.UI != nil {
		if in.Spec.UI == nil {
			in.Spec.UI = &v1alpha1.UIParameters{
				Actions:            d.Spec.UI.Actions,
				InstanceLabelPaths: d.Spec.UI.InstanceLabelPaths,
			}
		}

		if d.Spec.UI.Options != nil && in.Spec.UI.Options == nil {
			in.Spec.UI.Options = d.Spec.UI.Options
		}
		if d.Spec.UI.Editor != nil && in.Spec.UI.Editor == nil {
			in.Spec.UI.Editor = d.Spec.UI.Editor
		}
	}

	if len(in.Spec.Icons) == 0 {
		in.Spec.Icons = d.Spec.Icons
	}
	if in.Spec.Installer == nil {
		in.Spec.Installer = d.Spec.Installer
	}
	return in
}

func LoadByResourceID(kc client.Client, rid *kmapi.ResourceID) (*v1alpha1.ResourceEditor, error) {
	if rid == nil {
		return nil, errors.New("missing ResourceID")
	}

	gvr := rid.GroupVersionResource()
	if gvr.Version == "" || gvr.Resource == "" {
		id, err := kmapi.ExtractResourceID(kc.RESTMapper(), *rid)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to extract resource id for %+v", *rid)
		}
		gvr = id.GroupVersionResource()
	}
	return LoadByGVR(kc, gvr)
}

func List() []v1alpha1.ResourceEditor {
	m.Lock()
	defer m.Unlock()
	loader.ReloadIfTriggered()

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
	m.Lock()
	defer m.Unlock()
	loader.ReloadIfTriggered()

	out := make([]string, 0, len(reMap))
	for name := range reMap {
		out = append(out, name)
	}
	sort.Strings(out)
	return out
}
