package resourcedescriptor

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	kerr "k8s.io/apimachinery/pkg/api/errors"
	metainternalversion "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/registry/rest"
	"kmodules.xyz/resource-metadata/apis/meta"
	"kmodules.xyz/resource-metadata/apis/meta/v1alpha1"
	hub "kmodules.xyz/resource-metadata/hub/v1alpha1"
	"sigs.k8s.io/yaml"
)

type Storage struct {
}

var _ rest.GroupVersionKindProvider = &Storage{}
var _ rest.Scoper = &Storage{}
var _ rest.Getter = &Storage{}
var _ rest.Lister = &Storage{}

func NewStorage() *Storage {
	return &Storage{}
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
	data, err := hub.Asset(strings.Replace(name, "-", "/", 2) + ".yaml")
	if err != nil {
		return nil, kerr.NewNotFound(schema.GroupResource{Group: meta.GroupName, Resource: v1alpha1.ResourceKindResourceDescriptor}, name)
	}

	var obj v1alpha1.ResourceDescriptor
	err = yaml.Unmarshal(data, &obj)
	if err != nil {
		return nil, kerr.NewInternalError(err)
	}

	var out meta.ResourceDescriptor
	err = v1alpha1.Convert_v1alpha1_ResourceDescriptor_To_meta_ResourceDescriptor(&obj, &out, nil)
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

	names := hub.AssetNames()

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
	for _, name := range hub.AssetNames() {
		data, err := hub.Asset(name)
		if err != nil {
			return nil, kerr.NewNotFound(schema.GroupResource{Group: meta.GroupName, Resource: v1alpha1.ResourceKindResourceDescriptor}, name)
		}

		var obj v1alpha1.ResourceDescriptor
		err = yaml.Unmarshal(data, &obj)
		if err != nil {
			return nil, kerr.NewInternalError(err)
		}

		if options.LabelSelector != nil && !options.LabelSelector.Matches(labels.Set(obj.GetLabels())) {
			continue
		}

		var item meta.ResourceDescriptor
		err = v1alpha1.Convert_v1alpha1_ResourceDescriptor_To_meta_ResourceDescriptor(&obj, &item, nil)
		items = append(items, item)
	}

	return &meta.ResourceDescriptorList{Items: items}, nil
}
