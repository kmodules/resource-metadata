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

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
	v1alpha1 "kmodules.xyz/resource-metadata/apis/ui/v1alpha1"
)

// FakeResourceEditors implements ResourceEditorInterface
type FakeResourceEditors struct {
	Fake *FakeUiV1alpha1
}

var resourceeditorsResource = v1alpha1.SchemeGroupVersion.WithResource("resourceeditors")

var resourceeditorsKind = v1alpha1.SchemeGroupVersion.WithKind("ResourceEditor")

// Get takes name of the resourceEditor, and returns the corresponding resourceEditor object, and an error if there is any.
func (c *FakeResourceEditors) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.ResourceEditor, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(resourceeditorsResource, name), &v1alpha1.ResourceEditor{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ResourceEditor), err
}

// List takes label and field selectors, and returns the list of ResourceEditors that match those selectors.
func (c *FakeResourceEditors) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.ResourceEditorList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(resourceeditorsResource, resourceeditorsKind, opts), &v1alpha1.ResourceEditorList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.ResourceEditorList{ListMeta: obj.(*v1alpha1.ResourceEditorList).ListMeta}
	for _, item := range obj.(*v1alpha1.ResourceEditorList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested resourceEditors.
func (c *FakeResourceEditors) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(resourceeditorsResource, opts))
}

// Create takes the representation of a resourceEditor and creates it.  Returns the server's representation of the resourceEditor, and an error, if there is any.
func (c *FakeResourceEditors) Create(ctx context.Context, resourceEditor *v1alpha1.ResourceEditor, opts v1.CreateOptions) (result *v1alpha1.ResourceEditor, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(resourceeditorsResource, resourceEditor), &v1alpha1.ResourceEditor{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ResourceEditor), err
}

// Update takes the representation of a resourceEditor and updates it. Returns the server's representation of the resourceEditor, and an error, if there is any.
func (c *FakeResourceEditors) Update(ctx context.Context, resourceEditor *v1alpha1.ResourceEditor, opts v1.UpdateOptions) (result *v1alpha1.ResourceEditor, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(resourceeditorsResource, resourceEditor), &v1alpha1.ResourceEditor{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ResourceEditor), err
}

// Delete takes name of the resourceEditor and deletes it. Returns an error if one occurs.
func (c *FakeResourceEditors) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteActionWithOptions(resourceeditorsResource, name, opts), &v1alpha1.ResourceEditor{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeResourceEditors) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(resourceeditorsResource, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.ResourceEditorList{})
	return err
}

// Patch applies the patch and returns the patched resourceEditor.
func (c *FakeResourceEditors) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.ResourceEditor, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(resourceeditorsResource, name, pt, data, subresources...), &v1alpha1.ResourceEditor{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ResourceEditor), err
}
