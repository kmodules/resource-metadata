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

// FakeClusterProfiles implements ClusterProfileInterface
type FakeClusterProfiles struct {
	Fake *FakeUiV1alpha1
}

var clusterprofilesResource = v1alpha1.SchemeGroupVersion.WithResource("clusterprofiles")

var clusterprofilesKind = v1alpha1.SchemeGroupVersion.WithKind("ClusterProfile")

// Get takes name of the clusterProfile, and returns the corresponding clusterProfile object, and an error if there is any.
func (c *FakeClusterProfiles) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.ClusterProfile, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(clusterprofilesResource, name), &v1alpha1.ClusterProfile{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ClusterProfile), err
}

// List takes label and field selectors, and returns the list of ClusterProfiles that match those selectors.
func (c *FakeClusterProfiles) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.ClusterProfileList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(clusterprofilesResource, clusterprofilesKind, opts), &v1alpha1.ClusterProfileList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.ClusterProfileList{ListMeta: obj.(*v1alpha1.ClusterProfileList).ListMeta}
	for _, item := range obj.(*v1alpha1.ClusterProfileList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested clusterProfiles.
func (c *FakeClusterProfiles) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(clusterprofilesResource, opts))
}

// Create takes the representation of a clusterProfile and creates it.  Returns the server's representation of the clusterProfile, and an error, if there is any.
func (c *FakeClusterProfiles) Create(ctx context.Context, clusterProfile *v1alpha1.ClusterProfile, opts v1.CreateOptions) (result *v1alpha1.ClusterProfile, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(clusterprofilesResource, clusterProfile), &v1alpha1.ClusterProfile{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ClusterProfile), err
}

// Update takes the representation of a clusterProfile and updates it. Returns the server's representation of the clusterProfile, and an error, if there is any.
func (c *FakeClusterProfiles) Update(ctx context.Context, clusterProfile *v1alpha1.ClusterProfile, opts v1.UpdateOptions) (result *v1alpha1.ClusterProfile, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(clusterprofilesResource, clusterProfile), &v1alpha1.ClusterProfile{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ClusterProfile), err
}

// Delete takes name of the clusterProfile and deletes it. Returns an error if one occurs.
func (c *FakeClusterProfiles) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteActionWithOptions(clusterprofilesResource, name, opts), &v1alpha1.ClusterProfile{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeClusterProfiles) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(clusterprofilesResource, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.ClusterProfileList{})
	return err
}

// Patch applies the patch and returns the patched clusterProfile.
func (c *FakeClusterProfiles) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.ClusterProfile, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(clusterprofilesResource, name, pt, data, subresources...), &v1alpha1.ClusterProfile{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ClusterProfile), err
}
