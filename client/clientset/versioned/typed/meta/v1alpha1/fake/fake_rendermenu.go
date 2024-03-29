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

// FakeRenderMenus implements RenderMenuInterface
type FakeRenderMenus struct {
	Fake *FakeMetaV1alpha1
}

var rendermenusResource = v1alpha1.SchemeGroupVersion.WithResource("rendermenus")

var rendermenusKind = v1alpha1.SchemeGroupVersion.WithKind("RenderMenu")

// Create takes the representation of a renderMenu and creates it.  Returns the server's representation of the renderMenu, and an error, if there is any.
func (c *FakeRenderMenus) Create(ctx context.Context, renderMenu *v1alpha1.RenderMenu, opts v1.CreateOptions) (result *v1alpha1.RenderMenu, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(rendermenusResource, renderMenu), &v1alpha1.RenderMenu{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.RenderMenu), err
}
