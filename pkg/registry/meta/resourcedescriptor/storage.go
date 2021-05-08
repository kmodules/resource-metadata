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

package resourcedescriptor

import (
	"context"
	"fmt"
	"io/fs"
	"path"
	"strconv"

	"kmodules.xyz/resource-metadata/apis/meta"
	"kmodules.xyz/resource-metadata/apis/meta/v1alpha1"
	"kmodules.xyz/resource-metadata/hub/resourcedescriptors"

	kerr "k8s.io/apimachinery/pkg/api/errors"
	metainternalversion "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/registry/rest"
)

type Storage struct {
	convertor rest.TableConvertor
}

var _ rest.GroupVersionKindProvider = &Storage{}
var _ rest.Scoper = &Storage{}
var _ rest.Getter = &Storage{}
var _ rest.Lister = &Storage{}

func NewStorage() *Storage {
	return &Storage{
		convertor: rest.NewDefaultTableConvertor(schema.GroupResource{
			Group:    v1alpha1.SchemeGroupVersion.Group,
			Resource: v1alpha1.ResourceResourceDescriptors,
		}),
	}
}

func (r *Storage) GroupVersionKind(containingGV schema.GroupVersion) schema.GroupVersionKind {
	return v1alpha1.SchemeGroupVersion.WithKind(v1alpha1.ResourceKindResourceDescriptor)
}

func (r *Storage) NamespaceScoped() bool {
	return false
}

// Getter
func (r *Storage) New() runtime.Object {
	return &meta.ResourceDescriptor{}
}

func (r *Storage) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	obj, err := resourcedescriptors.LoadByName(name)
	if err != nil {
		return nil, kerr.NewNotFound(schema.GroupResource{Group: meta.GroupName, Resource: v1alpha1.ResourceKindResourceDescriptor}, name)
	}

	var out meta.ResourceDescriptor
	err = v1alpha1.Convert_v1alpha1_ResourceDescriptor_To_meta_ResourceDescriptor(obj, &out, nil)
	return &out, err
}

// Lister
func (r *Storage) NewList() runtime.Object {
	return &meta.ResourceDescriptorList{}
}

func (r *Storage) List(ctx context.Context, options *metainternalversion.ListOptions) (runtime.Object, error) {
	if options.FieldSelector != nil {
		return nil, kerr.NewBadRequest("fieldSelector is not a supported")
	}

	var names []string
	err := fs.WalkDir(resourcedescriptors.FS(), ".", func(p string, e fs.DirEntry, err error) error {
		if !e.IsDir() {
			names = append(names, path.Join(p, e.Name()))
		}
		return err
	})
	if err != nil {
		return nil, err
	}

	if options.Continue != "" {
		start, err := strconv.Atoi(options.Continue)
		if err != nil {
			return nil, kerr.NewBadRequest(fmt.Sprintf("invalid continue option, err:%v", err))
		}
		if start > len(names) {
			return r.NewList(), nil
		}
		names = names[start:]
	}
	if options.Limit > 0 && int64(len(names)) > options.Limit {
		names = names[:options.Limit]
	}

	items := make([]meta.ResourceDescriptor, 0, len(names))
	for _, filename := range names {
		obj, err := resourcedescriptors.LoadByFile(filename)
		if err != nil {
			return nil, err
		}

		if options.LabelSelector != nil && !options.LabelSelector.Matches(labels.Set(obj.GetLabels())) {
			continue
		}

		var item meta.ResourceDescriptor
		err = v1alpha1.Convert_v1alpha1_ResourceDescriptor_To_meta_ResourceDescriptor(obj, &item, nil)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return &meta.ResourceDescriptorList{Items: items}, nil
}

func (r *Storage) ConvertToTable(ctx context.Context, object runtime.Object, tableOptions runtime.Object) (*metav1.Table, error) {
	return r.convertor.ConvertToTable(ctx, object, tableOptions)
}
