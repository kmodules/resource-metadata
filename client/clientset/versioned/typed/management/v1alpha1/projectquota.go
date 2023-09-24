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

package v1alpha1

import (
	"context"
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
	v1alpha1 "kmodules.xyz/resource-metadata/apis/management/v1alpha1"
	scheme "kmodules.xyz/resource-metadata/client/clientset/versioned/scheme"
)

// ProjectQuotasGetter has a method to return a ProjectQuotaInterface.
// A group's client should implement this interface.
type ProjectQuotasGetter interface {
	ProjectQuotas() ProjectQuotaInterface
}

// ProjectQuotaInterface has methods to work with ProjectQuota resources.
type ProjectQuotaInterface interface {
	Create(ctx context.Context, projectQuota *v1alpha1.ProjectQuota, opts v1.CreateOptions) (*v1alpha1.ProjectQuota, error)
	Update(ctx context.Context, projectQuota *v1alpha1.ProjectQuota, opts v1.UpdateOptions) (*v1alpha1.ProjectQuota, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.ProjectQuota, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.ProjectQuotaList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.ProjectQuota, err error)
	ProjectQuotaExpansion
}

// projectQuotas implements ProjectQuotaInterface
type projectQuotas struct {
	client rest.Interface
}

// newProjectQuotas returns a ProjectQuotas
func newProjectQuotas(c *ManagementV1alpha1Client) *projectQuotas {
	return &projectQuotas{
		client: c.RESTClient(),
	}
}

// Get takes name of the projectQuota, and returns the corresponding projectQuota object, and an error if there is any.
func (c *projectQuotas) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.ProjectQuota, err error) {
	result = &v1alpha1.ProjectQuota{}
	err = c.client.Get().
		Resource("projectquotas").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of ProjectQuotas that match those selectors.
func (c *projectQuotas) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.ProjectQuotaList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.ProjectQuotaList{}
	err = c.client.Get().
		Resource("projectquotas").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested projectQuotas.
func (c *projectQuotas) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("projectquotas").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a projectQuota and creates it.  Returns the server's representation of the projectQuota, and an error, if there is any.
func (c *projectQuotas) Create(ctx context.Context, projectQuota *v1alpha1.ProjectQuota, opts v1.CreateOptions) (result *v1alpha1.ProjectQuota, err error) {
	result = &v1alpha1.ProjectQuota{}
	err = c.client.Post().
		Resource("projectquotas").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(projectQuota).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a projectQuota and updates it. Returns the server's representation of the projectQuota, and an error, if there is any.
func (c *projectQuotas) Update(ctx context.Context, projectQuota *v1alpha1.ProjectQuota, opts v1.UpdateOptions) (result *v1alpha1.ProjectQuota, err error) {
	result = &v1alpha1.ProjectQuota{}
	err = c.client.Put().
		Resource("projectquotas").
		Name(projectQuota.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(projectQuota).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the projectQuota and deletes it. Returns an error if one occurs.
func (c *projectQuotas) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Resource("projectquotas").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *projectQuotas) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("projectquotas").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched projectQuota.
func (c *projectQuotas) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.ProjectQuota, err error) {
	result = &v1alpha1.ProjectQuota{}
	err = c.client.Patch(pt).
		Resource("projectquotas").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
