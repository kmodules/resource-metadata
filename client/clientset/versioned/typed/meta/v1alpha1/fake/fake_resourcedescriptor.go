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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1alpha1 "kmodules.xyz/resource-metadata/apis/meta/v1alpha1"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeResourceDescriptors implements ResourceDescriptorInterface
type FakeResourceDescriptors struct {
	Fake *FakeMetaV1alpha1
}

var resourcedescriptorsResource = schema.GroupVersionResource{Group: "meta.appscode.com", Version: "v1alpha1", Resource: "resourcedescriptors"}

var resourcedescriptorsKind = schema.GroupVersionKind{Group: "meta.appscode.com", Version: "v1alpha1", Kind: "ResourceDescriptor"}

// Get takes name of the resourceDescriptor, and returns the corresponding resourceDescriptor object, and an error if there is any.
func (c *FakeResourceDescriptors) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.ResourceDescriptor, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(resourcedescriptorsResource, name), &v1alpha1.ResourceDescriptor{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ResourceDescriptor), err
}

// List takes label and field selectors, and returns the list of ResourceDescriptors that match those selectors.
func (c *FakeResourceDescriptors) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.ResourceDescriptorList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(resourcedescriptorsResource, resourcedescriptorsKind, opts), &v1alpha1.ResourceDescriptorList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.ResourceDescriptorList{ListMeta: obj.(*v1alpha1.ResourceDescriptorList).ListMeta}
	for _, item := range obj.(*v1alpha1.ResourceDescriptorList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested resourceDescriptors.
func (c *FakeResourceDescriptors) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(resourcedescriptorsResource, opts))
}

// Create takes the representation of a resourceDescriptor and creates it.  Returns the server's representation of the resourceDescriptor, and an error, if there is any.
func (c *FakeResourceDescriptors) Create(ctx context.Context, resourceDescriptor *v1alpha1.ResourceDescriptor, opts v1.CreateOptions) (result *v1alpha1.ResourceDescriptor, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(resourcedescriptorsResource, resourceDescriptor), &v1alpha1.ResourceDescriptor{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ResourceDescriptor), err
}

// Update takes the representation of a resourceDescriptor and updates it. Returns the server's representation of the resourceDescriptor, and an error, if there is any.
func (c *FakeResourceDescriptors) Update(ctx context.Context, resourceDescriptor *v1alpha1.ResourceDescriptor, opts v1.UpdateOptions) (result *v1alpha1.ResourceDescriptor, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(resourcedescriptorsResource, resourceDescriptor), &v1alpha1.ResourceDescriptor{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ResourceDescriptor), err
}

// Delete takes name of the resourceDescriptor and deletes it. Returns an error if one occurs.
func (c *FakeResourceDescriptors) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(resourcedescriptorsResource, name), &v1alpha1.ResourceDescriptor{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeResourceDescriptors) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(resourcedescriptorsResource, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.ResourceDescriptorList{})
	return err
}

// Patch applies the patch and returns the patched resourceDescriptor.
func (c *FakeResourceDescriptors) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.ResourceDescriptor, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(resourcedescriptorsResource, name, pt, data, subresources...), &v1alpha1.ResourceDescriptor{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ResourceDescriptor), err
}
