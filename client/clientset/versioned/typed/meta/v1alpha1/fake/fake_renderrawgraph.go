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
	testing "k8s.io/client-go/testing"
	v1alpha1 "kmodules.xyz/resource-metadata/apis/meta/v1alpha1"
)

// FakeRenderRawGraphs implements RenderRawGraphInterface
type FakeRenderRawGraphs struct {
	Fake *FakeMetaV1alpha1
}

var renderrawgraphsResource = v1alpha1.SchemeGroupVersion.WithResource("renderrawgraphs")

var renderrawgraphsKind = v1alpha1.SchemeGroupVersion.WithKind("RenderRawGraph")

// Create takes the representation of a renderRawGraph and creates it.  Returns the server's representation of the renderRawGraph, and an error, if there is any.
func (c *FakeRenderRawGraphs) Create(ctx context.Context, renderRawGraph *v1alpha1.RenderRawGraph, opts v1.CreateOptions) (result *v1alpha1.RenderRawGraph, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(renderrawgraphsResource, renderRawGraph), &v1alpha1.RenderRawGraph{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.RenderRawGraph), err
}
