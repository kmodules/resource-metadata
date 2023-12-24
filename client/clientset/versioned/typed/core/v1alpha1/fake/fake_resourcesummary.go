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
	testing "k8s.io/client-go/testing"
	v1alpha1 "kmodules.xyz/resource-metadata/apis/core/v1alpha1"
)

// FakeResourceSummaries implements ResourceSummaryInterface
type FakeResourceSummaries struct {
	Fake *FakeCoreV1alpha1
	ns   string
}

var resourcesummariesResource = v1alpha1.SchemeGroupVersion.WithResource("resourcesummaries")

var resourcesummariesKind = v1alpha1.SchemeGroupVersion.WithKind("ResourceSummary")

// Get takes name of the resourceSummary, and returns the corresponding resourceSummary object, and an error if there is any.
func (c *FakeResourceSummaries) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.ResourceSummary, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(resourcesummariesResource, c.ns, name), &v1alpha1.ResourceSummary{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ResourceSummary), err
}

// List takes label and field selectors, and returns the list of ResourceSummaries that match those selectors.
func (c *FakeResourceSummaries) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.ResourceSummaryList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(resourcesummariesResource, resourcesummariesKind, c.ns, opts), &v1alpha1.ResourceSummaryList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.ResourceSummaryList{ListMeta: obj.(*v1alpha1.ResourceSummaryList).ListMeta}
	for _, item := range obj.(*v1alpha1.ResourceSummaryList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}
